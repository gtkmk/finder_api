package userUsecase

import (
	"github.com/em-cash/simulador.em.cash/core/domain/customError"
	"github.com/em-cash/simulador.em.cash/core/domain/helper/constants"

	"github.com/em-cash/simulador.em.cash/core/domain/userDomain"
	"github.com/em-cash/simulador.em.cash/core/port"
	"github.com/em-cash/simulador.em.cash/core/port/repositories"
)

type UserFirstAccess struct {
	userRepository    repositories.UserRepository
	passwordEncryptor port.EncryptionInterface
	port.CustomErrorInterface
}

func NewUserFirstAccess(
	userRepository repositories.UserRepository,
	passwordEncryptor port.EncryptionInterface,
) *UserFirstAccess {
	return &UserFirstAccess{
		userRepository:       userRepository,
		passwordEncryptor:    passwordEncryptor,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (userFirstAccess *UserFirstAccess) Execute(userId string, password string) error {
	dbUser, err := userFirstAccess.userRepository.FindUserByIdWithoutPermissions(userId)

	if err != nil {
		return userFirstAccess.ThrowError(err.Error())
	}

	if dbUser == nil {
		return userFirstAccess.ThrowError(constants.UserNotFoundConst)
	}

	if dbUser.Status == userDomain.UserStatusLoggedConst {
		return userFirstAccess.ThrowError(constants.FirstAccessLinkAlreadyUsedConst)
	}

	encryptedPassword, err := userFirstAccess.passwordEncryptor.GenerateHashPassword(password)

	if err != nil {
		return userFirstAccess.ThrowError(err.Error())
	}

	if err := userFirstAccess.userRepository.ResetUserPassword(dbUser.Id, encryptedPassword); err != nil {
		return userFirstAccess.ThrowError(err.Error())
	}

	if err := userFirstAccess.userRepository.SetUserStatus(dbUser.Id, userDomain.UserStatusLoggedConst); err != nil {
		return userFirstAccess.ThrowError(err.Error())
	}

	return nil
}
