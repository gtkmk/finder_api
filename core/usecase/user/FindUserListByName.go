package userUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type FindUsersListByName struct {
	UserDatabase repositories.UserRepository
	CustomError  port.CustomErrorInterface
}

func NewFindUsersListByName(
	userDatabase repositories.UserRepository,
	customError port.CustomErrorInterface,
) *FindUsersListByName {
	return &FindUsersListByName{
		UserDatabase: userDatabase,
		CustomError:  customError,
	}
}

func (FindUsersListByName *FindUsersListByName) Execute(nameToSearch string, loggedUserId string) ([]map[string]interface{}, error) {
	dbProposal, err := FindUsersListByName.UserDatabase.FindUsersListByName(nameToSearch, loggedUserId)

	if err != nil {
		return nil, FindUsersListByName.CustomError.ThrowError(err.Error())
	}

	if dbProposal == nil {
		return nil, FindUsersListByName.CustomError.ThrowError(helper.UserNotFoundConst)
	}

	return dbProposal, nil
}
