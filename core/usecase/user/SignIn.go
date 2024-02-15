package userUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/helper"
	"os"
	"time"

	"github.com/gtkmk/finder_api/core/domain/credentialsDomain"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/jwtAuth"
	"github.com/gtkmk/finder_api/core/domain/userDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	sharedMethodsInterface "github.com/gtkmk/finder_api/core/port/sharedMethods"
	"github.com/gtkmk/finder_api/infra/envMode"
)

const maxAllowedTimeInMinutesConst = 120 * time.Minute

type SignIn struct {
	userDatabase      repositories.UserRepository
	passwordEncryptor port.EncryptionInterface
	credentials       *credentialsDomain.Credential
	userEvent         sharedMethodsInterface.CreateUserEventInterface
	port.CustomErrorInterface
}

func NewSign(
	userDatabase repositories.UserRepository,
	passwordEncryptor port.EncryptionInterface,
	credentials *credentialsDomain.Credential,
	userEvent sharedMethodsInterface.CreateUserEventInterface,
) *SignIn {
	return &SignIn{
		userDatabase:         userDatabase,
		passwordEncryptor:    passwordEncryptor,
		credentials:          credentials,
		CustomErrorInterface: customError.NewCustomError(),
		userEvent:            userEvent,
	}
}

func (signIn *SignIn) Execute(userIP string, userDevice string) (string, error) {
	user, err := signIn.userDatabase.FindUserByEmail(signIn.credentials.Email)

	if err != nil {
		return "", signIn.ThrowError(err.Error())
	}

	if user == nil {
		return "", signIn.ThrowError(helper.EmailOrPasswordIncorrectConst)
	}

	if err := signIn.isUserAvailableToLogin(user); err != nil {
		return "", err
	}

	checkPassword := signIn.passwordEncryptor.CheckHashedPassword(signIn.credentials.Password, user.Password)

	if !checkPassword {
		return "", signIn.ThrowError(helper.EmailOrPasswordIncorrectConst)
	}

	if err != nil {
		return "", err
	}

	jwtAuthentication := jwtAuth.NewjwtAuth(os.Getenv(envMode.JwtSecretConst))

	token, err := jwtAuthentication.GenerateJWT(map[string]interface{}{
		"i": user.Id,
		"l": user.Group.Layer,
	})

	if err != nil {
		return "", signIn.ThrowError(err.Error())
	}

	if err := signIn.userEvent.SaveLoginEvent(user.Id, userIP, userDevice); err != nil {
		return "", signIn.ThrowError(err.Error())
	}

	return token, nil
}

func (signIn *SignIn) isUserAvailableToLogin(user *userDomain.User) error {
	if user.Status == userDomain.UserStatusExpiredConst {
		return signIn.ThrowError(helper.TokenExpiredGenerateNewEmailConst)
	}

	if user.Status == userDomain.UserStatusPendingConst && isMoreThanTheAllowedTime(user.CreatedAt) {
		if err := signIn.userDatabase.SetUserStatus(user.Id, userDomain.UserStatusExpiredConst); err != nil {
			return signIn.ThrowError(err.Error())
		}

		return signIn.ThrowError(helper.TokenExpiredGenerateNewEmailConst)
	}

	if user.Status == userDomain.UserStatusPendingConst {
		return signIn.ThrowError(helper.EmailSentToChangePasswordConst)
	}

	if !user.IsActive {
		return signIn.ThrowError(helper.UserIsNotActiveConst)
	}

	return nil
}

func isMoreThanTheAllowedTime(createdTime time.Time) bool {
	currentTime := time.Now()
	timeDifference := currentTime.Sub(createdTime)

	return timeDifference > maxAllowedTimeInMinutesConst
}
