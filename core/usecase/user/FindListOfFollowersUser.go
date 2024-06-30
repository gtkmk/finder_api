package userUsecase

import (
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type FindUserListOfFollowers struct {
	UserDatabase repositories.UserRepository
	CustomError  port.CustomErrorInterface
}

func NewFindUserListOfFollowers(
	userDatabase repositories.UserRepository,
	customError port.CustomErrorInterface,
) *FindUserListOfFollowers {
	return &FindUserListOfFollowers{
		UserDatabase: userDatabase,
		CustomError:  customError,
	}
}

func (findUserListOfFollowers *FindUserListOfFollowers) Execute(userProfileToSearch string, followOrFollowing string) ([]map[string]interface{}, error) {
	dbFollowInfo, err := findUserListOfFollowers.UserDatabase.FindUsersListByName(userProfileToSearch, userProfileToSearch)

	if err != nil {
		return nil, findUserListOfFollowers.CustomError.ThrowError(err.Error())
	}

	return dbFollowInfo, nil
}
