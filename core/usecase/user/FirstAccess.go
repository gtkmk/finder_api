package userUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/domain/userDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
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
	dbUser, err := userFirstAccess.userRepository.FindUserById(userId)

	if err != nil {
		return userFirstAccess.ThrowError(err.Error())
	}

	if dbUser == nil {
		return userFirstAccess.ThrowError(helper.UserNotFoundConst)
	}

	if dbUser.Status == userDomain.UserStatusLoggedConst {
		return userFirstAccess.ThrowError(helper.FirstAccessLinkAlreadyUsedConst)
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
