package userHandler

import (
	"github.com/gin-gonic/gin"

	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	userUsecase "github.com/gtkmk/finder_api/core/usecase/user"
	"github.com/gtkmk/finder_api/infra/database/repository"
)

type FindUserUserDetailsHandler struct {
	connection   port.ConnectionInterface
	uuid         port.UuidInterface
	userDatabase repositories.UserRepository
	customError  port.CustomErrorInterface
}

func NewFindUserUserDetailsHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
) port.HandlerInterface {
	return &FindUserUserDetailsHandler{
		connection:  connection,
		uuid:        uuid,
		customError: customError.NewCustomError(),
	}
}

func (findUserUserDetailsHandler *FindUserUserDetailsHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, findUserUserDetailsHandler.connection, findUserUserDetailsHandler.uuid)

	userId := context.Query("user-id")

	if userId == "" {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findUserUserDetailsHandler.customError.ThrowError(helper.FieldIsMandatoryConst, "user-id"),
			routesConstants.BadRequestConst,
		)
		return
	}

	findUserUserDetailsHandler.openTableConnection()

	userInfo, err := userUsecase.NewFindUserUserDetails(
		findUserUserDetailsHandler.userDatabase,
		findUserUserDetailsHandler.customError,
	).Execute(userId)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findUserUserDetailsHandler.customError.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson("data", userInfo, routesConstants.StatusOk)
}

func (findUserUserDetailsHandler *FindUserUserDetailsHandler) openTableConnection() {
	findUserUserDetailsHandler.userDatabase = repository.NewUserDatabase(findUserUserDetailsHandler.connection)
}
