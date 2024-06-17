package userUsecase

import (
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	sharedMethodsInterface "github.com/gtkmk/finder_api/core/port/sharedMethods"
)

const checkPointEditPostTransactionNameConst = "editUserInfo"

type UpdateUserUpdateUserInfo struct {
	LoggedUserId      string
	UserDatabase      repositories.UserRepository
	Transaction       port.ConnectionInterface
	RollBackAndReturn sharedMethodsInterface.RollBackAndReturnErrorInterface
	CustomError       port.CustomErrorInterface
}

func NewUpdateUserUpdateUserInfo(
	loggedUserId string,
	userDatabase repositories.UserRepository,
	transaction port.ConnectionInterface,
	rollBackAndReturn sharedMethodsInterface.RollBackAndReturnErrorInterface,
	customError port.CustomErrorInterface,
) *UpdateUserUpdateUserInfo {
	return &UpdateUserUpdateUserInfo{
		LoggedUserId:      loggedUserId,
		UserDatabase:      userDatabase,
		Transaction:       transaction,
		RollBackAndReturn: rollBackAndReturn,
		CustomError:       customError,
	}
}

func (updateUserUpdateUserInfo *UpdateUserUpdateUserInfo) Execute(decodedName string, decodedUserCellphone string) error {
	if err := updateUserUpdateUserInfo.Transaction.SavePoint(checkPointEditPostTransactionNameConst); err != nil {
		return updateUserUpdateUserInfo.CustomError.ThrowError(err.Error())
	}

	err := updateUserUpdateUserInfo.UserDatabase.UpdateUserEmailAndCellphoneNumber(
		updateUserUpdateUserInfo.LoggedUserId,
		decodedName,
		decodedUserCellphone,
	)

	if err != nil {
		if rollbackErr := updateUserUpdateUserInfo.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err := updateUserUpdateUserInfo.Transaction.Commit(); err != nil {
		if rollbackErr := updateUserUpdateUserInfo.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return rollbackErr
		}
		return updateUserUpdateUserInfo.CustomError.ThrowError(err.Error())
	}

	return nil
}

func (updateUserUpdateUserInfo *UpdateUserUpdateUserInfo) rollbackToSavePointAndCommit() error {
	if transactErr := updateUserUpdateUserInfo.Transaction.RollbackTo(checkPointEditPostTransactionNameConst); transactErr != nil {
		return updateUserUpdateUserInfo.CustomError.ThrowError(transactErr.Error())
	}

	if commitErr := updateUserUpdateUserInfo.Transaction.Commit(); commitErr != nil {
		return updateUserUpdateUserInfo.CustomError.ThrowError(commitErr.Error())
	}

	return nil
}
