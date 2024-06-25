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

type FindUsersListByNameHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	contextExtractor port.HttpContextValuesExtractorInterface
	userDatabase     repositories.UserRepository
	customError      port.CustomErrorInterface
}

func NewFindUsersListByNameHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &FindUsersListByNameHandler{
		connection:       connection,
		uuid:             uuid,
		contextExtractor: contextExtractor,
		customError:      customError.NewCustomError(),
	}
}

func (findUserUserDetailsHandler *FindUsersListByNameHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, findUserUserDetailsHandler.connection, findUserUserDetailsHandler.uuid)

	nameToSearch := context.Query("name")

	if nameToSearch == "" {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findUserUserDetailsHandler.customError.ThrowError(helper.FieldIsMandatoryConst, "nome"),
			routesConstants.BadRequestConst,
		)

		return
	}

	loggedUserId, extractErr := findUserUserDetailsHandler.contextExtractor.Extract(context)
	if extractErr != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findUserUserDetailsHandler.customError.ThrowError(extractErr.Error()),
			routesConstants.InternarServerErrorConst,
		)

		return
	}

	findUserUserDetailsHandler.openTableConnection()

	usersInfo, err := userUsecase.NewFindUsersListByName(
		findUserUserDetailsHandler.userDatabase,
		findUserUserDetailsHandler.customError,
	).Execute(nameToSearch, loggedUserId)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findUserUserDetailsHandler.customError.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson("data", usersInfo, routesConstants.StatusOk)
}

func (findUserUserDetailsHandler *FindUsersListByNameHandler) openTableConnection() {
	findUserUserDetailsHandler.userDatabase = repository.NewUserDatabase(findUserUserDetailsHandler.connection)
}
