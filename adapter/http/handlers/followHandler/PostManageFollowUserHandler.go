package followHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/followDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	followUsecase "github.com/gtkmk/finder_api/core/usecase/follow"
	"github.com/gtkmk/finder_api/infra/database/repository"
	followrequestentity "github.com/gtkmk/finder_api/infra/requestEntity/followRequestEntity"
)

type PostFollowUserHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	followDatabase   repositories.FollowRepository
	contextExtractor port.HttpContextValuesExtractorInterface
	customError      port.CustomErrorInterface
}

func NewCreateManageFollowUserHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &PostFollowUserHandler{
		connection:       connection,
		uuid:             uuid,
		contextExtractor: contextExtractor,
		customError:      customError.NewCustomError(),
	}
}

func (postFollowUserHandler *PostFollowUserHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, postFollowUserHandler.connection, postFollowUserHandler.uuid)

	loggedUserId, extractErr := postFollowUserHandler.contextExtractor.Extract(context)
	if extractErr != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			postFollowUserHandler.customError.ThrowError(extractErr.Error()),
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	decodedFollow, err := postFollowUserHandler.defineFollow(context, loggedUserId)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)

		return
	}

	transaction, err := postFollowUserHandler.connection.BeginTransaction()

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			postFollowUserHandler.customError.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)

		return
	}

	postFollowUserHandler.openTableConnection(transaction)

	following, err := followUsecase.NewCreateManageFollowUser(
		postFollowUserHandler.followDatabase,
		decodedFollow,
		transaction,
		postFollowUserHandler.customError,
	).Execute()

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)

		return
	}

	jsonResponse.SendJson("following", following, routesConstants.StatusOk)
}

func (postFollowUserHandler *PostFollowUserHandler) defineFollow(context *gin.Context, userId string) (
	*followDomain.Follow,
	error,
) {
	decodedFollow, decodeErr := followrequestentity.NewFollowRequest(
		context,
		postFollowUserHandler.uuid,
		userId,
		postFollowUserHandler.customError,
	)

	if decodeErr != nil {
		return nil, postFollowUserHandler.customError.ThrowError(decodeErr.Error())
	}

	if err := decodedFollow.Validate(context); err != nil {
		return nil, postFollowUserHandler.customError.ThrowError(err.Error())
	}

	return decodedFollow.BuildFollowObject()
}

func (postFollowUserHandler *PostFollowUserHandler) openTableConnection(transaction port.ConnectionInterface) {
	postFollowUserHandler.followDatabase = repository.NewFollowDatabase(transaction)
}
