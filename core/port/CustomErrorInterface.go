package port

type CustomErrorInterface interface {
	ThrowError(errorString string, params ...any) error
	IsCustomError(err error) (errorString string, stackTrance string, isCustomError bool)
}
