package commentHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/domain/success"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	commentUsecase "github.com/gtkmk/finder_api/core/usecase/comment"
	"github.com/gtkmk/finder_api/infra/database/repository"
)

type DeleteCommentHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	ContextExtractor port.HttpContextValuesExtractorInterface
	commentDatabase  repositories.CommentRepository
	port.CustomErrorInterface
}

func NewDeleteCommentHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &DeleteCommentHandler{
		connection:           connection,
		uuid:                 uuid,
		ContextExtractor:     contextExtractor,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (postCreateCommentHandler *DeleteCommentHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, postCreateCommentHandler.connection, postCreateCommentHandler.uuid)

	loggedUserId, extractErr := postCreateCommentHandler.ContextExtractor.Extract(context)
	if extractErr != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			postCreateCommentHandler.ThrowError(extractErr.Error()),
			routesConstants.InternarServerErrorConst,
		)

		return
	}

	commentId := context.Query("comment-id")
	if commentId == "" {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			helper.ErrorBuilder(helper.FieldIsMandatoryConst, "comment-id"),
			routesConstants.InternarServerErrorConst,
		)

		return
	}

	transaction, err := postCreateCommentHandler.connection.BeginTransaction()
	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			postCreateCommentHandler.ThrowError(err.Error()),
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	postCreateCommentHandler.openTableConnection(transaction)

	if err := commentUsecase.NewDeleteComment(
		postCreateCommentHandler.commentDatabase,
		transaction,
		commentId,
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

func (postCreateCommentHandler *DeleteCommentHandler) openTableConnection(transaction port.ConnectionInterface) {
	postCreateCommentHandler.commentDatabase = repository.NewCommentDatabase(transaction)
}
