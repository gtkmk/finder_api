package userUsecase

import (
	"encoding/json"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port/sharedMethods"

	"github.com/gtkmk/finder_api/core/domain/notificationDomain"
	"github.com/gtkmk/finder_api/core/domain/userDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type SignUp struct {
	userDatabase        repositories.UserRepository
	transaction         port.ConnectionInterface
	user                *userDomain.User
	uuid                port.UuidInterface
	notificationService port.NotificationInterface
	path                string
	customError         port.CustomErrorInterface
	userEvent           sharedMethods.CreateUserEventInterface
}

func NewCreateUser(
	transaction port.ConnectionInterface,
	userDatabase repositories.UserRepository,
	definedUser *userDomain.User,
	uuid port.UuidInterface,
	notificationService port.NotificationInterface,
	path string,
	customErrorInterface port.CustomErrorInterface,
	userEvent sharedMethods.CreateUserEventInterface,
) *SignUp {
	return &SignUp{
		transaction:         transaction,
		userDatabase:        userDatabase,
		user:                definedUser,
		uuid:                uuid,
		notificationService: notificationService,
		path:                path,
		customError:         customErrorInterface,
		userEvent:           userEvent,
	}
}

func (signUp *SignUp) Execute(userIP string, userDevice string) error {
	if err := signUp.verifyIfUserExistsByCpf(); err != nil {
		if transactionErr := signUp.transaction.Rollback(); transactionErr != nil {
			return signUp.customError.ThrowError(transactionErr.Error())
		}

		return err
	}

	if err := signUp.verifyIfUserExists(); err != nil {
		if transactionErr := signUp.transaction.Rollback(); transactionErr != nil {
			return signUp.customError.ThrowError(transactionErr.Error())
		}

		return err
	}

	if err := signUp.createUser(userIP, userDevice); err != nil {
		if transactionErr := signUp.transaction.Rollback(); transactionErr != nil {
			return signUp.customError.ThrowError(transactionErr.Error())
		}

		return err
	}

	if err := signUp.transaction.Commit(); err != nil {
		return signUp.customError.ThrowError(err.Error())
	}

	_ = signUp.sendSingUpEmail(signUp.user)

	return nil
}

func (signUp *SignUp) verifyIfUserExistsByCpf() error {
	exists := signUp.userDatabase.VerifyIfUserExistsByCpf(signUp.user.Cpf)

	if exists {
		return signUp.customError.ThrowError(helper.UserAlreadyRegisteredConst)
	}

	return nil
}

func (signUp *SignUp) verifyIfUserExists() error {
	dbUser, err := signUp.userDatabase.FindUserByEmail(signUp.user.Email)

	if err != nil {
		return signUp.customError.ThrowError(err.Error())
	}

	if dbUser != nil {
		return signUp.customError.ThrowError(helper.UserAlreadyRegisteredWithEmailConst)
	}

	return nil
}

func (signUp *SignUp) createUser(userIP string, userDevice string) error {
	if err := signUp.userDatabase.CreateUser(signUp.user); err != nil {
		return signUp.customError.ThrowError(err.Error())
	}

	if err := signUp.saveNewUserEvent(userIP, userDevice); err != nil {
		return err
	}

	return nil
}

func (signUp *SignUp) sendSingUpEmail(
	user *userDomain.User,
) error {
	notification := notificationDomain.NewNotification(
		userDomain.FirstAccessConst,
		user.Name,
		user.Email,
		"",
		"",
		signUp.path,
		nil,
	)

	_ = signUp.notificationService.SendNotifications(notification)

	return nil
}

func (signUp *SignUp) saveNewUserEvent(userIP string, userDevice string) error {
	signUp.user.Password = ""
	interfaceToJson, err := json.MarshalIndent(signUp.user, "", "\t")

	if err != nil {
		return err
	}

	if err := signUp.userEvent.SaveNewUserEvent(signUp.user.Id, userIP, userDevice, string(interfaceToJson)); err != nil {
		return err
	}
	return nil
}
