package routes

import (
	"fmt"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/documentDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/database/repository"
	"github.com/gtkmk/finder_api/infra/errorPersistence"
	"strings"

	"github.com/gin-gonic/gin"
)

type JsonResponse struct {
	*gin.Context
	Connection    port.ConnectionInterface
	UuidGenerator port.UuidInterface
	CustomError   port.CustomErrorInterface
	port.CustomErrorInterface
}

func NewJsonResponse(
	ctx *gin.Context,
	connection port.ConnectionInterface,
	uuidGenerator port.UuidInterface,
) *JsonResponse {
	return &JsonResponse{
		Context:              ctx,
		Connection:           connection,
		UuidGenerator:        uuidGenerator,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (jsonResponse *JsonResponse) ThrowError(
	key string,
	err error,
	statusCode int,
) {
	errorRepository := repository.NewErrorDatabase(jsonResponse.Connection)
	errorHandler := errorPersistence.NewErrorPersistenceHandler(errorRepository, jsonResponse.UuidGenerator)

	errorString, stack, isCustomError := jsonResponse.IsCustomError(err)
	errorHandler.HandleError(errorString, stack, isCustomError, statusCode)

	jsonResponse.Writer.Header().Set("Content-Type", "application/json")
	jsonResponse.Writer.WriteHeader(statusCode)

	jsonResponse.JSON(statusCode, map[string]string{key: errorString})
}

func (jsonResponse *JsonResponse) SendJson(key string, data interface{}, statusCode int) {
	jsonResponse.Writer.Header().Set("Content-Type", "application/json")
	jsonResponse.Writer.WriteHeader(statusCode)

	jsonResponse.JSON(statusCode, map[string]interface{}{key: data})
}

func (jsonResponse *JsonResponse) SendDocument(document *documentDomain.Document, statusCode int) error {
	data := document.Data

	jsonResponse.Writer.Header().Set("Content-Type", document.MimeType)
	jsonResponse.Writer.Header().Set("Content-Disposition", "attachment; filename="+defineFileName(document))
	jsonResponse.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	if _, err := jsonResponse.Writer.Write(data); err != nil {
		return err
	}
	jsonResponse.Writer.WriteHeader(statusCode)

	return nil
}

func (jsonResponse *JsonResponse) SendRawDocument(data string, statusCode int) error {
	jsonResponse.Writer.Header().Set("Content-Type", "application/pdf")
	jsonResponse.Writer.Header().Set("Content-Disposition", "attachment")

	if _, err := jsonResponse.Writer.Write([]byte(data)); err != nil {
		return err
	}

	jsonResponse.Writer.WriteHeader(statusCode)

	return nil
}

func defineFileName(document *documentDomain.Document) string {
	fileName := document.File.Filename
	lastIndex := strings.LastIndex(fileName, "/")

	if lastIndex != -1 {
		fileName = fileName[lastIndex+1:]
	}

	return fileName
}
