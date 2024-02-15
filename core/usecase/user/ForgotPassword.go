package userUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/domain/userDomain"

	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type ForgotPassword struct {
	userRepository    repositories.UserRepository
	passwordEncryptor port.EncryptionInterface
	port.CustomErrorInterface
}

func NewForgotPassword(
	userRepository repositories.UserRepository,
	passwordEncryptor port.EncryptionInterface,
) *ForgotPassword {
	return &ForgotPassword{
		userRepository:       userRepository,
		passwordEncryptor:    passwordEncryptor,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (forgotPassword *ForgotPassword) Execute(userId string, password string) error {
	dbUser, err := forgotPassword.userRepository.FindUserByIdWithoutPermissions(userId)

	if err != nil {
		return forgotPassword.ThrowError(err.Error())
	}

	if dbUser == nil {
		return forgotPassword.ThrowError(helper.UserNotFoundConst)
	}

	if !dbUser.ResetPassword {
		return forgotPassword.ThrowError(helper.UserDoNotRequestPasswordChangingConst)
	}

	encryptedPassword, err := forgotPassword.passwordEncryptor.GenerateHashPassword(password)

	if err != nil {
		return forgotPassword.ThrowError(err.Error())
	}

	if err := forgotPassword.userRepository.ResetUserPassword(dbUser.Id, encryptedPassword); err != nil {
		return forgotPassword.ThrowError(err.Error())
	}

	if err := forgotPassword.userRepository.UpdateResetPasswordStatus(false, userDomain.UserStatusLoggedConst, dbUser.Id); err != nil {
		return forgotPassword.ThrowError(err.Error())
	}

	return nil
}
