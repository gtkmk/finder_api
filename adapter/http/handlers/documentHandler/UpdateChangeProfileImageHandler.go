package documentHandler

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/documentDomain"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/domain/success"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	documentUsecase "github.com/gtkmk/finder_api/core/usecase/document"
	"github.com/gtkmk/finder_api/core/usecase/sharedMethods"
	"github.com/gtkmk/finder_api/infra/database/repository"
	"github.com/gtkmk/finder_api/infra/envMode"
	"github.com/gtkmk/finder_api/infra/file"
	documentRequestEntity "github.com/gtkmk/finder_api/infra/requestEntity/documentRequestEntity"
)

type UpdateDocumentChangeProfileImageHandler struct {
	connection       port.ConnectionInterface
	uuid             port.UuidInterface
	contextExtractor port.HttpContextValuesExtractorInterface
	documentDatabase repositories.DocumentRepository
	customError      port.CustomErrorInterface
}

func NewUpdateDocumentChangeProfileImageHandler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
	contextExtractor port.HttpContextValuesExtractorInterface,
) port.HandlerInterface {
	return &UpdateDocumentChangeProfileImageHandler{
		connection:       connection,
		uuid:             uuid,
		contextExtractor: contextExtractor,
		customError:      customError.NewCustomError(),
	}
}

func (updateDocumentChangeProfileImageHandler *UpdateDocumentChangeProfileImageHandler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, updateDocumentChangeProfileImageHandler.connection, updateDocumentChangeProfileImageHandler.uuid)

	if context.ContentType() != "multipart/form-data" {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			updateDocumentChangeProfileImageHandler.customError.ThrowError(helper.ContentTypeErrorConst),
			routesConstants.BadRequestConst,
		)

		return
	}

	loggedUserId, extractErr := updateDocumentChangeProfileImageHandler.contextExtractor.Extract(context)
	if extractErr != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			updateDocumentChangeProfileImageHandler.customError.ThrowError(extractErr.Error()),
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	decodedDocument, documentType, err := updateDocumentChangeProfileImageHandler.defineDocument(context, loggedUserId)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)

		return
	}

	transaction, err := updateDocumentChangeProfileImageHandler.connection.BeginTransaction()

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			updateDocumentChangeProfileImageHandler.customError.ThrowError(err.Error()),
			routesConstants.InternarServerErrorConst,
		)
		return
	}

	updateDocumentChangeProfileImageHandler.openTableConnection(transaction)

	rollBackAndReturn := sharedMethods.NewRollBackAndReturnError(transaction)

	if err := documentUsecase.NewUpdateDocumentChangeProfileImage(
		updateDocumentChangeProfileImageHandler.documentDatabase,
		file.NewFileFactory(),
		decodedDocument,
		os.Getenv(envMode.TempDirConst),
		transaction,
		rollBackAndReturn,
		updateDocumentChangeProfileImageHandler.customError,
	).Execute(documentType); err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			err,
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson(routesConstants.DataKeyConst, success.SuccessfullyUpdatedProfileConst, routesConstants.StatusOk)
}

func (updateDocumentChangeProfileImageHandler *UpdateDocumentChangeProfileImageHandler) defineDocument(context *gin.Context, userId string) (
	decodedPost *documentDomain.Document,
	mediaType string,
	err error,
) {
	documen, decodeErr := documentRequestEntity.NewUpdateDocumentRequest(
		context,
		updateDocumentChangeProfileImageHandler.uuid,
		userId,
	)

	if decodeErr != nil {
		return nil, "", updateDocumentChangeProfileImageHandler.customError.ThrowError(decodeErr.Error())
	}

	if validateErr := documen.Validate(context); validateErr != nil {
		return nil, "", updateDocumentChangeProfileImageHandler.customError.ThrowError(validateErr.Error())
	}

	return documen.BuildDocumentObjectAndType(context)
}

func (updateDocumentChangeProfileImageHandler *UpdateDocumentChangeProfileImageHandler) openTableConnection(transaction port.ConnectionInterface) {
	updateDocumentChangeProfileImageHandler.documentDatabase = repository.NewDocumentDatabase(transaction)
}
