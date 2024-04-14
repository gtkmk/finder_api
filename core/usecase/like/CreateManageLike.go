package likeUsecase

import (
	"github.com/gtkmk/finder_api/core/port"

	"github.com/gtkmk/finder_api/core/domain/likeDomain"
	"github.com/gtkmk/finder_api/core/port/repositories"
	sharedMethodsInterface "github.com/gtkmk/finder_api/core/port/sharedMethods"
)

const checkPointManageLikeTransactionNameConst = "manageLike"

type CreateManageLike struct {
	likeDatabase      repositories.LikeRepository
	like              *likeDomain.Like
	transaction       port.ConnectionInterface
	rollBackAndReturn sharedMethodsInterface.RollBackAndReturnErrorInterface
	customError       port.CustomErrorInterface
}

func NewCreateManageLike(
	likeDatabase repositories.LikeRepository,
	like *likeDomain.Like,
	transaction port.ConnectionInterface,
	rollBackAndReturn sharedMethodsInterface.RollBackAndReturnErrorInterface,
	customError port.CustomErrorInterface,
) *CreateManageLike {
	return &CreateManageLike{
		likeDatabase:      likeDatabase,
		like:              like,
		transaction:       transaction,
		rollBackAndReturn: rollBackAndReturn,
		customError:       customError,
	}
}

func (createManageLike *CreateManageLike) Execute() (int, error) {
	if err := createManageLike.transaction.SavePoint(checkPointManageLikeTransactionNameConst); err != nil {
		return 0, createManageLike.customError.ThrowError(err.Error())
	}

	err := createManageLike.manageLike()
	if err != nil {
		if rollbackErr := createManageLike.rollbackToSavePointAndCommit(); err != nil {
			return 0, rollbackErr
		}
		return 0, err
	}

	updatedLikeCount, err := createManageLike.likeDatabase.FindCurrentLikesCount(createManageLike.like)
	if err != nil {
		return 0, createManageLike.customError.ThrowError(err.Error())
	}

	if err := createManageLike.transaction.Commit(); err != nil {
		if rollbackErr := createManageLike.rollbackToSavePointAndCommit(); err != nil {
			return 0, rollbackErr
		}
		return 0, createManageLike.customError.ThrowError(err.Error())
	}

	return updatedLikeCount, nil
}

func (createManageLike *CreateManageLike) manageLike() error {
	alreadyLiked, existingLikeInfo, err := createManageLike.likeDatabase.ConfirmExistingLike(createManageLike.like)
	if err != nil {
		return err
	}

	if alreadyLiked {
		return createManageLike.likeDatabase.RemoveLike(existingLikeInfo.Id)
	}

	return createManageLike.likeDatabase.CreateLike(createManageLike.like)
}

func (createManageLike *CreateManageLike) rollbackToSavePointAndCommit() error {
	if transactErr := createManageLike.transaction.RollbackTo(checkPointManageLikeTransactionNameConst); transactErr != nil {
		return createManageLike.customError.ThrowError(transactErr.Error())
	}

	if commitErr := createManageLike.transaction.Commit(); commitErr != nil {
		return createManageLike.customError.ThrowError(commitErr.Error())
	}

	return nil
}
