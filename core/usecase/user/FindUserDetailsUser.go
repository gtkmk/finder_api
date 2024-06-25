package userUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type FindUserUserDetails struct {
	UserDatabase repositories.UserRepository
	CustomError  port.CustomErrorInterface
}

func NewFindUserUserDetails(
	userDatabase repositories.UserRepository,
	customError port.CustomErrorInterface,
) *FindUserUserDetails {
	return &FindUserUserDetails{
		UserDatabase: userDatabase,
		CustomError:  customError,
	}
}

func (findUserUserDetails *FindUserUserDetails) Execute(userId string, loggedUserId string) ([]map[string]interface{}, error) {
	if err := findUserUserDetails.verifyIfUserExists(userId); err != nil {
		return nil, err
	}

	dbProposal, err := findUserUserDetails.UserDatabase.FindCompleteUserInfoByID(userId, loggedUserId)

	if err != nil {
		return nil, findUserUserDetails.CustomError.ThrowError(err.Error())
	}

	if dbProposal == nil {
		return nil, findUserUserDetails.CustomError.ThrowError(helper.UserNotFoundConst)
	}

	return dbProposal, nil
}

func (findUserUserDetails *FindUserUserDetails) verifyIfUserExists(userId string) error {
	user, err := findUserUserDetails.UserDatabase.FindUserById(userId)

	if err != nil {
		return findUserUserDetails.CustomError.ThrowError(err.Error())
	}

	if user == nil {
		return findUserUserDetails.CustomError.ThrowError(helper.UserNotFoundConst)
	}

	return nil
}
