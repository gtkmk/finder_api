package port

type ErrorPersistenceHandlerInterface interface {
	HandleError(message string, stack string, isCustomError bool, statusCode int)
}
