package sharedMethods

import (
	"github.com/gtkmk/finder_api/core/port"
)

type RollBackAndReturnError struct {
	transaction port.ConnectionInterface
	port.CustomErrorInterface
}

func NewRollBackAndReturnError(
	transaction port.ConnectionInterface,
) *RollBackAndReturnError {
	return &RollBackAndReturnError{
		transaction: transaction,
	}
}

func (rollBackAndReturnError *RollBackAndReturnError) RollbackAndReturnError(err error) error {
	if rollbackErr := rollBackAndReturnError.transaction.Rollback(); rollbackErr != nil {
		return rollBackAndReturnError.ThrowError(rollbackErr.Error())
	}

	return err
}
