package commentHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/commentDomain"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/domain/success"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	commentUsecase "github.com/gtkmk/finder_api/core/usecase/comment"
	"github.com/gtkmk/finder_api/infra/database/repository"
	"github.com/gtkmk/finder_api/infra/requestEntity/commentRequestEntity"
)

type PostCommentHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	ContextExtractor port.HttpContextValuesExtractorInterface
	commentDatabase  repositories.CommentRepository
	postDatabase     repositories.PostRepositoryInterface
	port.CustomErrorInterface
}

func NewCreateCommentHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &PostCommentHandler{
		connection:           connection,
		uuid:                 uuid,
		ContextExtractor:     contextExtractor,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (postCreateCommentHandler *PostCommentHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, postCreateCommentHandler.connection, postCreateCommentHandler.uuid)

	if context.ContentType() != "multipart/form-data" {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			postCreateCommentHandler.ThrowError(helper.ContentTypeErrorConst),
			routesConstants.BadRequestConst,
		)

		return
	}

	loggedUserId, extractErr := postCreateCommentHandler.ContextExtractor.Extract(context)
	if extractErr != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			postCreateCommentHandler.ThrowError(extractErr.Error()),
			routesConstants.InternarServerErrorConst,
		)

		return
	}

	comment, err := postCreateCommentHandler.defineComment(context, loggedUserId)
	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
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

	if err := commentUsecase.NewCreateComment(
		postCreateCommentHandler.commentDatabase,
		postCreateCommentHandler.postDatabase,
		*comment,
		transaction,
	).Execute(); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson(routesConstants.MessageKeyConst, success.SuccessfullyCreatedCommentConst, routesConstants.CreatedConst)
}

func (postCreateCommentHandler *PostCommentHandler) defineComment(context *gin.Context, userId string) (
	decodedComment *commentDomain.Comment,
	err error,
) {
	comment, decodeErr := commentRequestEntity.NewCommentRequest(
		context,
		postCreateCommentHandler.uuid,
		userId,
	)

	if decodeErr != nil {
		return nil, postCreateCommentHandler.ThrowError(decodeErr.Error())
	}

	if validateErr := comment.Validate(context, false); validateErr != nil {
		return nil, postCreateCommentHandler.ThrowError(validateErr.Error())
	}

	return comment.BuildCommentObject(nil)
}

func (postCreateCommentHandler *PostCommentHandler) openTableConnection(transaction port.ConnectionInterface) {
	postCreateCommentHandler.commentDatabase = repository.NewCommentDatabase(transaction)
	postCreateCommentHandler.postDatabase = repository.NewPostDatabase(postCreateCommentHandler.connection)
}
