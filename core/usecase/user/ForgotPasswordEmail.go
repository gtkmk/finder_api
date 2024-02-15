package userUsecase

import (
	"fmt"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"os"

	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/notificationDomain"
	"github.com/gtkmk/finder_api/core/domain/userDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	"github.com/gtkmk/finder_api/infra/envMode"
)

type ForgotPasswordEmail struct {
	userDatabase       repositories.UserRepository
	walletNotification port.NotificationInterface
	port.CustomErrorInterface
}

func NewForgotPasswordEmail(
	userDatabase repositories.UserRepository,
	walletNotification port.NotificationInterface,
) *ForgotPasswordEmail {
	return &ForgotPasswordEmail{
		userDatabase:         userDatabase,
		walletNotification:   walletNotification,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (forgotPasswordEmail *ForgotPasswordEmail) Execute(userEmail string) error {
	dbUser, err := forgotPasswordEmail.userDatabase.FindUserByEmail(userEmail)

	if err != nil {
		return forgotPasswordEmail.ThrowError(err.Error())
	}

	if dbUser == nil {
		return forgotPasswordEmail.ThrowError(helper.UserNotFoundConst)
	}

	if err := forgotPasswordEmail.userDatabase.UpdateResetPasswordStatus(true, userDomain.UserStatusPendingConst, dbUser.Id); err != nil {
		return forgotPasswordEmail.ThrowError(err.Error())
	}

	url := forgotPasswordEmail.generateUrlResetPassword(dbUser.Id)

	_ = forgotPasswordEmail.sendForgotPasswordEmail(dbUser, url)

	return nil
}

func (forgotPasswordEmail *ForgotPasswordEmail) sendForgotPasswordEmail(
	user *userDomain.User,
	url string,
) error {
	notification := notificationDomain.NewNotification(
		userDomain.PasswordResetConst,
		user.Name,
		user.Email,
		"",
		"",
		url,
		nil,
	)

	return forgotPasswordEmail.walletNotification.SendNotifications(notification)
}

func (forgotPasswordEmail *ForgotPasswordEmail) generateUrlResetPassword(
	userId string,
) string {
	return fmt.Sprintf("%s/signin?user-id=%s&reset=true", os.Getenv(envMode.FrontUrlConst), userId)
}
