package postHandler

import (
	"os"

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
	"github.com/gtkmk/finder_api/infra/envMode"
	"github.com/gtkmk/finder_api/infra/file"
	"github.com/gtkmk/finder_api/infra/requestEntity/postRequestEntity"
)

type PostCreatePostHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	ContextExtractor port.HttpContextValuesExtractorInterface
	postDatabase     repositories.PostRepositoryInterface
	port.CustomErrorInterface
}

func NewCreatePostHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &PostCreatePostHandler{
		connection:           connection,
		uuid:                 uuid,
		ContextExtractor:     contextExtractor,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (postCreatePostHandler *PostCreatePostHandler) Handle(context *gin.Context) {
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

	post, err := postCreatePostHandler.definePost(context, loggedUserId)
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

	if err := postUsecase.NewCreatePost(
		postCreatePostHandler.postDatabase,
		file.NewFileFactory(),
		*post,
		transaction,
		os.Getenv(envMode.TempDirConst),
	).Execute(); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson(routesConstants.MessageKeyConst, success.SuccessfullyCreatedPostConst, routesConstants.CreatedConst)
}

func (postCreatePostHandler *PostCreatePostHandler) definePost(context *gin.Context, userId string) (
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

	if validateErr := post.Validate(context); validateErr != nil {
		return nil, postCreatePostHandler.ThrowError(validateErr.Error())
	}

	if iterateErr := post.IterateIntoFiles(context); iterateErr != nil {
		return nil, postCreatePostHandler.ThrowError(iterateErr.Error())
	}

	return post.BuildPostObject()
}

func (postCreatePostHandler *PostCreatePostHandler) openTableConnection(transaction port.ConnectionInterface) {
	postCreatePostHandler.postDatabase = repository.NewPostDatabase(transaction)
}
