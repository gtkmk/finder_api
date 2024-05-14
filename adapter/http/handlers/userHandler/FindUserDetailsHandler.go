package userHandler

import (
	"github.com/gin-gonic/gin"

	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	userUsecase "github.com/gtkmk/finder_api/core/usecase/user"
	"github.com/gtkmk/finder_api/infra/database/repository"
)

type FindUserUserDetailsHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	contextExtractor port.HttpContextValuesExtractorInterface
	userDatabase     repositories.UserRepository
	customError      port.CustomErrorInterface
}

func NewFindUserUserDetailsHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &FindUserUserDetailsHandler{
		connection:       connection,
		uuid:             uuid,
		contextExtractor: contextExtractor,
		customError:      customError.NewCustomError(),
	}
}

func (findUserUserDetailsHandler *FindUserUserDetailsHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, findUserUserDetailsHandler.connection, findUserUserDetailsHandler.uuid)

	userId := context.Query("user-id")

	loggedUserId, extractErr := findUserUserDetailsHandler.contextExtractor.Extract(context)
	if extractErr != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findUserUserDetailsHandler.customError.ThrowError(extractErr.Error()),
			routesConstants.InternarServerErrorConst,
		)

		return
	}

	if userId == "" {
		userId = loggedUserId
	}

	findUserUserDetailsHandler.openTableConnection()

	userInfo, err := userUsecase.NewFindUserUserDetails(
		findUserUserDetailsHandler.userDatabase,
		findUserUserDetailsHandler.customError,
	).Execute(userId, loggedUserId)

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
