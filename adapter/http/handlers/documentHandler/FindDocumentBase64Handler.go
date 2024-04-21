package documentHandler

import (
	"github.com/gin-gonic/gin"

	"github.com/gtkmk/finder_api/adapter/http/routes"
	"github.com/gtkmk/finder_api/adapter/http/routesConstants"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	documentUsecase "github.com/gtkmk/finder_api/core/usecase/document"
	"github.com/gtkmk/finder_api/infra/file"
)

type FindDocumentImageBase64Handler struct {
	connection  port.ConnectionInterface
	uuid        port.UuidInterface
	customError port.CustomErrorInterface
}

func NewFindDocumentImageBase64Handler(
	connection port.ConnectionInterface,
	uuid port.UuidInterface,
) port.HandlerInterface {
	return &FindDocumentImageBase64Handler{
		connection:  connection,
		uuid:        uuid,
		customError: customError.NewCustomError(),
	}
}

func (findDocumentImageBase64Handler *FindDocumentImageBase64Handler) Handle(context *gin.Context) {
	jsonResponse := routes.NewJsonResponse(context, findDocumentImageBase64Handler.connection, findDocumentImageBase64Handler.uuid)

	documentPath := context.Query("document-path")

	if documentPath == "" {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findDocumentImageBase64Handler.customError.ThrowError(helper.FieldIsMandatoryConst, "document-path"),
			routesConstants.BadRequestConst,
		)
		return
	}

	document, err := documentUsecase.NewFindDocumentImageBase64(
		file.NewFileFactory(),
		findDocumentImageBase64Handler.customError,
	).Execute(documentPath)

	if err != nil {
		jsonResponse.ThrowError(
			routesConstants.MessageKeyConst,
			findDocumentImageBase64Handler.customError.ThrowError(err.Error()),
			routesConstants.BadRequestConst,
		)
		return
	}

	jsonResponse.SendJson(routesConstants.DataKeyConst, document, routesConstants.StatusOk)
}
