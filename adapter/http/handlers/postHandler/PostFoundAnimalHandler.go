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
)

type FoundAnimalHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	ContextExtractor port.HttpContextValuesExtractorInterface
	postDatabase     repositories.PostRepositoryInterface
	port.CustomErrorInterface
}

func NewPostAnimalFoundHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &FoundAnimalHandler{
		connection:           connection,
		uuid:                 uuid,
		ContextExtractor:     contextExtractor,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (foundAnimalHandler *FoundAnimalHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, foundAnimalHandler.connection, foundAnimalHandler.uuid)

	loggedUserId, extractErr := foundAnimalHandler.ContextExtractor.Extract(context)
	if extractErr != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			foundAnimalHandler.ThrowError(extractErr.Error()),
			routesConstants.InternarServerErrorConst,
		)

		return
	}

	postId, err := foundAnimalHandler.getPostId(context)
	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	foundStatus, err := foundAnimalHandler.getFoundStatus(context)
	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	transaction, err := foundAnimalHandler.connection.BeginTransaction()
	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			foundAnimalHandler.ThrowError(err.Error()),
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	foundAnimalHandler.openTableConnection(transaction)

	if err := postUsecase.NewFoundAnimalPost(
		foundAnimalHandler.postDatabase,
		transaction,
		postId,
		foundStatus,
		loggedUserId,
	).Execute(); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson(routesConstants.MessageKeyConst, success.SuccessfullyUpdatedFoundStatusConst, routesConstants.CreatedConst)
}

func (foundAnimalHandler *FoundAnimalHandler) getPostId(context *gin.Context) (string, error) {
	postId := context.Query("post-id")
	if postId == "" {
		return "", helper.ErrorBuilder(helper.FieldIsMandatoryConst, "post-id")
	}
	return postId, nil
}

func (foundAnimalHandler *FoundAnimalHandler) getFoundStatus(context *gin.Context) (string, error) {
	foundStatus := context.Query("found")

	if foundStatus == "" {
		return "", helper.ErrorBuilder(helper.FieldIsMandatoryConst, "found")
	}
	if foundStatus != postDomain.FoundOptionTrueConst && foundStatus != postDomain.FoundOptionFalseConst {
		return "", helper.ErrorBuilder(helper.FieldNotInAllowedValuesConst, "found")
	}
	return foundStatus, nil
}

func (foundAnimalHandler *FoundAnimalHandler) openTableConnection(transaction port.ConnectionInterface) {
	foundAnimalHandler.postDatabase = repository.NewPostDatabase(transaction)
}
