package sharedMethods

type RollBackAndReturnErrorInterface interface {
	RollbackAndReturnError(err error) error
}
