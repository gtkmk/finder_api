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

type FindUserListOfFollowersHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	contextExtractor port.HttpContextValuesExtractorInterface
	userDatabase     repositories.UserRepository
	customError      port.CustomErrorInterface
}

func NewFindUserListOfFollowersHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &FindUserListOfFollowersHandler{
		connection:       connection,
		uuid:             uuid,
		contextExtractor: contextExtractor,
		customError:      customError.NewCustomError(),
	}
}

func (findUserListOfFollowersHandler *FindUserListOfFollowersHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, findUserListOfFollowersHandler.connection, findUserListOfFollowersHandler.uuid)

	userProfileToSearch, followOrFollowing, err := findUserListOfFollowersHandler.extractAndValidateQueryParams(context)
	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findUserListOfFollowersHandler.customError.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	loggedUserId, extractErr := findUserListOfFollowersHandler.contextExtractor.Extract(context)
	if extractErr != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findUserListOfFollowersHandler.customError.ThrowError(extractErr.Error()),
			routesConstants.InternarServerErrorConst,
		)

		return
	}

	if userProfileToSearch == "" {
		userProfileToSearch = loggedUserId
	}

	findUserListOfFollowersHandler.openTableConnection()

	dbFollowInfo, err := userUsecase.NewFindUserListOfFollowers(
		findUserListOfFollowersHandler.userDatabase,
		findUserListOfFollowersHandler.customError,
	).Execute(userProfileToSearch, followOrFollowing)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findUserListOfFollowersHandler.customError.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson(routesConstants.DataKeyConst, map[string]interface{}{"info": dbFollowInfo}, routesConstants.StatusOk)
}

func (findUserListOfFollowersHandler *FindUserListOfFollowersHandler) extractAndValidateQueryParams(context *gin.Context) (string, string, error) {
	userProfileToSearch := context.Query("user-id")
	followOrFollowing := context.Query("type")

	if followOrFollowing == "" {
		return "", "", findUserListOfFollowersHandler.customError.ThrowError(helper.FieldIsMandatoryConst, "type")
	}

	if followOrFollowing != "following" && followOrFollowing != "followers" {
		return "", "", findUserListOfFollowersHandler.customError.ThrowError(helper.FieldNotInAllowedValuesConst, "type")
	}

	return userProfileToSearch, followOrFollowing, nil
}

func (findUserListOfFollowersHandler *FindUserListOfFollowersHandler) openTableConnection() {
	findUserListOfFollowersHandler.userDatabase = repository.NewUserDatabase(findUserListOfFollowersHandler.connection)
}
