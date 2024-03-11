package userHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/domain/userDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	userUsecase "github.com/gtkmk/finder_api/core/usecase/user"
	"github.com/gtkmk/finder_api/infra/database/repository"
)

type FindUserHandler struct {
	connection   port.ConnectionInterface
	uuid         port.UuidInterface
	userDatabase repositories.UserRepository
	port.CustomErrorInterface
}

type FindUserReturn struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Cpf             string `json:"cpf"`
	CellphoneNumber string `json:"cellphone_number"`
	CreatorId       string `json:"creator_id"`
	Role            string `json:"role"`
	IsActive        bool   `json:"is_active"`
	Permissions     any    `json:"permissions"`
	Unities         any    `json:"unities"`
}

func NewFindUserHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
) port.HandlerInterface {
	return &FindUserHandler{
		connection:           connection,
		uuid:                 uuid,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (findUserRoute *FindUserHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, findUserRoute.connection, findUserRoute.uuid)

	findUserRoute.openTableConnection()

	getUser := userUsecase.NewGetUser(findUserRoute.userDatabase)

	userId := context.Query("user-id")

	if userId == "" {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findUserRoute.ThrowError(helper.UserIdIsRequiredConst),
			routesConstants.BadRequestConst,
		)
		return
	}

	user, err := getUser.Execute(userId, true)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	userReturn := findUserRoute.defineGetUserReturn(user)

	jsonResponse.SendJson("user", userReturn, routesConstants.StatusOk)
}

func (findUserRoute *FindUserHandler) defineGetUserReturn(user *userDomain.User) FindUserReturn {
	var userReturn FindUserReturn

	userReturn.Id = user.Id
	userReturn.Email = user.Email
	userReturn.Name = user.Name
	userReturn.Cpf = user.Cpf
	userReturn.CellphoneNumber = user.CellphoneNumber
	userReturn.IsActive = user.IsActive

	return userReturn
}

func (findUserRoute *FindUserHandler) openTableConnection() {
	findUserRoute.userDatabase = repository.NewUserDatabase(findUserRoute.connection)
}
