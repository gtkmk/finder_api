package postHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/domain/success"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	postUsecase "github.com/gtkmk/finder_api/core/usecase/post"
	"github.com/gtkmk/finder_api/infra/database/repository"
)

type DeletePostHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	ContextExtractor port.HttpContextValuesExtractorInterface
	postDatabase     repositories.PostRepositoryInterface
	port.CustomErrorInterface
}

func NewDeletePostHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &DeletePostHandler{
		connection:           connection,
		uuid:                 uuid,
		ContextExtractor:     contextExtractor,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (deletePostHandler *DeletePostHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, deletePostHandler.connection, deletePostHandler.uuid)

	loggedUserId, extractErr := deletePostHandler.ContextExtractor.Extract(context)
	if extractErr != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			deletePostHandler.ThrowError(extractErr.Error()),
			routesConstants.InternarServerErrorConst,
		)

		return
	}

	postId := context.Query("post-id")
	if postId == "" {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			helper.ErrorBuilder(helper.FieldIsMandatoryConst, "post-id"),
			routesConstants.InternarServerErrorConst,
		)

		return
	}

	transaction, err := deletePostHandler.connection.BeginTransaction()
	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			deletePostHandler.ThrowError(err.Error()),
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	deletePostHandler.openTableConnection(transaction)

	if err := postUsecase.NewDeletePost(
		deletePostHandler.postDatabase,
		transaction,
		postId,
		loggedUserId,
	).Execute(); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson(routesConstants.MessageKeyConst, success.SuccessfullyDeletedCommentConst, routesConstants.CreatedConst)
}

func (deletePostHandler *DeletePostHandler) openTableConnection(transaction port.ConnectionInterface) {
	deletePostHandler.postDatabase = repository.NewPostDatabase(transaction)
}
