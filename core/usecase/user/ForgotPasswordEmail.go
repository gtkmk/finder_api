package userUsecase

import (
	"fmt"
	"os"

	"github.com/em-cash/simulador.em.cash/core/domain/customError"
	"github.com/em-cash/simulador.em.cash/core/domain/helper/constants"
	"github.com/em-cash/simulador.em.cash/core/domain/notificationDomain"
	"github.com/em-cash/simulador.em.cash/core/domain/userDomain"
	"github.com/em-cash/simulador.em.cash/core/port"
	"github.com/em-cash/simulador.em.cash/core/port/repositories"
	"github.com/em-cash/simulador.em.cash/infra/envMode"
)

type ForgotPasswordEmail struct {
	userRepository     repositories.UserRepository
	walletNotification port.NotificationInterface
	port.CustomErrorInterface
}

func NewForgotPasswordEmail(
	userRepository repositories.UserRepository,
	walletNotification port.NotificationInterface,
) *ForgotPasswordEmail {
	return &ForgotPasswordEmail{
		userRepository:       userRepository,
		walletNotification:   walletNotification,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (forgotPasswordEmail *ForgotPasswordEmail) Execute(userEmail string) error {
	dbUser, err := forgotPasswordEmail.userRepository.FindUserByEmail(userEmail)

	if err != nil {
		return forgotPasswordEmail.ThrowError(err.Error())
	}

	if dbUser == nil {
		return forgotPasswordEmail.ThrowError(constants.UserNotFoundConst)
	}

	if err := forgotPasswordEmail.userRepository.UpdateResetPasswordStatus(true, userDomain.UserStatusPendingConst, dbUser.Id); err != nil {
		return forgotPasswordEmail.ThrowError(err.Error())
	}

	url := forgotPasswordEmail.generateUrlResetPassword(dbUser.Id)

	if err := forgotPasswordEmail.sendForgotPasswordEmail(dbUser, url); err != nil {
		return forgotPasswordEmail.ThrowError(err.Error())
	}

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
		"",
	)

	return forgotPasswordEmail.walletNotification.SendNotifications(notification)
}

func (forgotPasswordEmail *ForgotPasswordEmail) generateUrlResetPassword(
	userId string,
) string {
	return fmt.Sprintf("%s/signin?user-id=%s&reset=true", os.Getenv(envMode.PortalUrlConst), userId)
}
