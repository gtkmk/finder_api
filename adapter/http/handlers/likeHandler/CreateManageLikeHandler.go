package likeHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/likeDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	likeUsecase "github.com/gtkmk/finder_api/core/usecase/like"
	"github.com/gtkmk/finder_api/core/usecase/sharedMethods"
	"github.com/gtkmk/finder_api/infra/database/repository"
	likerequestentity "github.com/gtkmk/finder_api/infra/requestEntity/likeRequestEntity"
)

type CreateManageLikeHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	likeDatabase     repositories.LikeRepository
	contextExtractor port.HttpContextValuesExtractorInterface
	customError      port.CustomErrorInterface
}

func NewCreateLikeCreateLikeHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &CreateManageLikeHandler{
		connection:       connection,
		uuid:             uuid,
		contextExtractor: contextExtractor,
		customError:      customError.NewCustomError(),
	}
}

func (createManageLikeHandler *CreateManageLikeHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, createManageLikeHandler.connection, createManageLikeHandler.uuid)

	loggedUserId, extractErr := createManageLikeHandler.contextExtractor.Extract(context)
	if extractErr != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			createManageLikeHandler.customError.ThrowError(extractErr.Error()),
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	decodedLike, err := createManageLikeHandler.defineLike(context, loggedUserId)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)

		return
	}

	transaction, err := createManageLikeHandler.connection.BeginTransaction()

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			createManageLikeHandler.customError.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)

		return
	}

	createManageLikeHandler.openTableConnection(transaction)

	rollBackAndReturn := sharedMethods.NewRollBackAndReturnError(transaction)

	likesCount, err := likeUsecase.NewCreateManageLike(
		createManageLikeHandler.likeDatabase,
		decodedLike,
		transaction,
		rollBackAndReturn,
		createManageLikeHandler.customError,
	).Execute()

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)

		return
	}

	jsonResponse.SendJson("likesCount", likesCount, routesConstants.StatusOk)
}

func (createManageLikeHandler *CreateManageLikeHandler) defineLike(context *gin.Context, userId string) (
	*likeDomain.Like,
	error,
) {
	decodedLike, decodeErr := likerequestentity.NewLikeRequest(
		context,
		createManageLikeHandler.uuid,
		userId,
	)

	if decodeErr != nil {
		return nil, createManageLikeHandler.customError.ThrowError(decodeErr.Error())
	}

	if err := decodedLike.Validate(context); err != nil {
		return nil, createManageLikeHandler.customError.ThrowError(err.Error())
	}

	return decodedLike.BuildLikeObject()
}

func (createManageLikeHandler *CreateManageLikeHandler) openTableConnection(transaction port.ConnectionInterface) {
	createManageLikeHandler.likeDatabase = repository.NewLikeDatabase(transaction)
}
