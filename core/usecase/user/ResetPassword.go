package userUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type ResetPassword struct {
	userRepository    repositories.UserRepository
	passwordEncryptor port.EncryptionInterface
	port.CustomErrorInterface
}

func NewResetPassword(
	userRepository repositories.UserRepository,
	passwordEncryptor port.EncryptionInterface,
) *ResetPassword {
	return &ResetPassword{
		userRepository:       userRepository,
		passwordEncryptor:    passwordEncryptor,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (resetPassword *ResetPassword) Execute(userId string, oldPassword string, newPassword string) error {
	dbUser, err := resetPassword.userRepository.FindUserByIdWithoutPermissions(userId)

	if err != nil {
		return resetPassword.ThrowError(err.Error())
	}

	if dbUser == nil {
		return resetPassword.ThrowError(helper.UserNotFoundConst)
	}

	if !resetPassword.verifyOldPassword(oldPassword, dbUser.Password) {
		return resetPassword.ThrowError(helper.OldPasswordIncorrectConst)
	}

	if resetPassword.verifyOldPassword(newPassword, dbUser.Password) {
		return resetPassword.ThrowError(helper.PasswordCannotBeSameAsOldPasswordConst)
	}

	encryptedPassword, err := resetPassword.passwordEncryptor.GenerateHashPassword(newPassword)

	if err != nil {
		return resetPassword.ThrowError(err.Error())
	}

	if err := resetPassword.userRepository.ResetUserPassword(dbUser.Id, encryptedPassword); err != nil {
		return resetPassword.ThrowError(err.Error())
	}

	return nil
}

func (resetPassword *ResetPassword) verifyOldPassword(oldPassword string, dbUserPassword string) bool {
	return resetPassword.passwordEncryptor.CheckHashedPassword(oldPassword, dbUserPassword)
}
