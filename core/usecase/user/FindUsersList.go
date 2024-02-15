package userUsecase

import (
	"time"

	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/userDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type GetUsersList struct {
	userDatabase repositories.UserRepository
	port.CustomErrorInterface
}

func NewGetUsersList(
	userDatabase repositories.UserRepository,
) *GetUsersList {
	return &GetUsersList{
		userDatabase:         userDatabase,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (getUsersList *GetUsersList) Execute(
	userId string,
	groupLayer float64,
	orderBy string,
	unityId string,
) ([]*userDomain.User, error) {
	users, err := getUsersList.userDatabase.FindUsersList(unityId, userId, orderBy)
	if err != nil {
		return nil, getUsersList.ThrowError(err.Error())
	}
	return getUsersList.transformsMapToUserDomain(users), nil
}

func (getUsersList *GetUsersList) transformsMapToUserDomain(users []map[string]interface{}) []*userDomain.User {
	var (
		formattedUsers      []*userDomain.User
		userCellphoneNumber string
	)

	uniqueUsers := make(map[string]*userDomain.User)

	for _, user := range users {
		createdAtFormatted, _ := time.Parse(datetimeDomain.DATETIME_FORMAT, user["user_created_at"].(string))

		if user["user_cellphone_number"] != nil {
			userCellphoneNumber = user["user_cellphone_number"].(string)
		}

		formattedUsers = getUsersList.defineUniqueUsersList(
			user,
			uniqueUsers,
			formattedUsers,
			userCellphoneNumber,
			createdAtFormatted,
		)
	}

	return formattedUsers
}

func (getUsersList *GetUsersList) defineUniqueUsersList(
	user map[string]interface{},
	uniqueUsers map[string]*userDomain.User,
	formattedUsers []*userDomain.User,
	userCellphoneNumber string,
	createdAtFormatted time.Time,
) []*userDomain.User {
	uniqueUsers[user["id"].(string)] =
		userDomain.NewUser(
			user["id"].(string),
			user["user_name"].(string),
			user["user_email"].(string),
			"",
			user["user_cpf"].(string),
			userCellphoneNumber,
			user["user_status"].(string),
			user["user_is_active"].(int64) > 0,
			user["user_password_reset"].(int64) > 0,
			createdAtFormatted,
		)

	return append(formattedUsers, uniqueUsers[user["id"].(string)])
}
