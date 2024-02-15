package customError

import (
	"github.com/go-errors/errors"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
)

type CustomError struct{}

func NewCustomError() port.CustomErrorInterface {
	return &CustomError{}
}

func (customError *CustomError) ThrowError(errorString string, params ...any) error {
	errorString = helper.DefineMassage(errorString, params...)
	errorWithStackTrace := errors.Errorf(errorString)

	return errors.New(errorWithStackTrace)
}

func (customError *CustomError) IsCustomError(err error) (errorString string, stackTrance string, isCustomError bool) {
	var customErr *errors.Error

	if errors.As(err, &customErr) {
		return err.(*errors.Error).Err.Error(), err.(*errors.Error).ErrorStack(), true
	}

	return err.Error(), "", false
}
