package userUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"

	"github.com/gtkmk/finder_api/core/domain/userDomain"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type GetUser struct {
	userDatabase repositories.UserRepository
	port.CustomErrorInterface
}

func NewGetUser(
	userDatabase repositories.UserRepository,
) *GetUser {
	return &GetUser{
		userDatabase:         userDatabase,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (signIn *GetUser) Execute(userId string, withPermissions bool) (*userDomain.User, error) {
	user, err := signIn.userDatabase.FindUserById(userId)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, signIn.ThrowError(helper.UserNotFoundConst)
	}

	return user, nil
}
