package postHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/domain/postDomain"
	"github.com/gtkmk/finder_api/core/domain/success"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	postUsecase "github.com/gtkmk/finder_api/core/usecase/post"
	"github.com/gtkmk/finder_api/infra/database/repository"
	"github.com/gtkmk/finder_api/infra/file"
	"github.com/gtkmk/finder_api/infra/requestEntity/postRequestEntity"
)

type PostEditPostHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	ContextExtractor port.HttpContextValuesExtractorInterface
	postDatabase     repositories.PostRepositoryInterface
	port.CustomErrorInterface
}

func NewEditPostHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &PostEditPostHandler{
		connection:           connection,
		uuid:                 uuid,
		ContextExtractor:     contextExtractor,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (postCreatePostHandler *PostEditPostHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, postCreatePostHandler.connection, postCreatePostHandler.uuid)

	if context.ContentType() != "multipart/form-data" {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			postCreatePostHandler.ThrowError(helper.ContentTypeErrorConst),
			routesConstants.BadRequestConst,
		)

		return
	}

	loggedUserId, extractErr := postCreatePostHandler.ContextExtractor.Extract(context)
	if extractErr != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			postCreatePostHandler.ThrowError(extractErr.Error()),
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

	post, err := postCreatePostHandler.definePost(context, loggedUserId, postId)
	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	transaction, err := postCreatePostHandler.connection.BeginTransaction()
	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			postCreatePostHandler.ThrowError(err.Error()),
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	postCreatePostHandler.openTableConnection(transaction)

	if err := postUsecase.NewEditPost(
		postCreatePostHandler.postDatabase,
		file.NewFileFactory(),
		*post,
		transaction,
	).Execute(); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson(routesConstants.MessageKeyConst, success.SuccessfullyEditedPostConst, routesConstants.CreatedConst)
}

func (postCreatePostHandler *PostEditPostHandler) definePost(context *gin.Context, userId string, postId string) (
	decodedPost *postDomain.Post,
	err error,
) {
	post, decodeErr := postRequestEntity.NewPostRequest(
		context,
		postCreatePostHandler.uuid,
		userId,
	)

	if decodeErr != nil {
		return nil, postCreatePostHandler.ThrowError(decodeErr.Error())
	}

	if validateErr := post.Validate(context, true); validateErr != nil {
		return nil, postCreatePostHandler.ThrowError(validateErr.Error())
	}

	return post.BuildPostObject(&postId)
}

func (postCreatePostHandler *PostEditPostHandler) openTableConnection(transaction port.ConnectionInterface) {
	postCreatePostHandler.postDatabase = repository.NewPostDatabase(transaction)
}
