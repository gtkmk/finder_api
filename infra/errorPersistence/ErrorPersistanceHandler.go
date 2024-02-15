package errorPersistence

import (
	"fmt"
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type ErrorPersistenceHandler struct {
	errorRepository repositories.ErrorRepositoryInterface
	uuidGenerator   port.UuidInterface
}

func NewErrorPersistenceHandler(
	errorRepository repositories.ErrorRepositoryInterface,
	uuidGenerator port.UuidInterface,
) port.ErrorPersistenceHandlerInterface {
	return &ErrorPersistenceHandler{
		errorRepository,
		uuidGenerator,
	}
}

func (errorPersistenceHandler ErrorPersistenceHandler) HandleError(message string, stack string, isCustomError bool, statusCode int) {
	errorMessage := fmt.Sprintf("Time: %s | StatusCode: %d | Message: %s", datetimeDomain.CreateFormattedNow(), statusCode, message)

	errorId := errorPersistenceHandler.uuidGenerator.GenerateUuid()
	createAt := datetimeDomain.CreateFormattedNow()
	if !isCustomError {
		errorPersistenceHandler.errorRepository.SaveErrorWithoutStack(errorId, errorMessage, createAt)
		return
	}
	errorPersistenceHandler.errorRepository.SaveErrorWithStack(errorId, errorMessage, stack, createAt)
}
