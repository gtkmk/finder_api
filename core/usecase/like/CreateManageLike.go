package likeUsecase

import (
	"github.com/gtkmk/finder_api/core/port"

	"github.com/gtkmk/finder_api/core/domain/likeDomain"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

const checkPointManageLikeTransactionNameConst = "manageLike"

type CreateManageLike struct {
	LikeDatabase repositories.LikeRepository
	Like         *likeDomain.Like
	Transaction  port.ConnectionInterface
	CustomError  port.CustomErrorInterface
}

func NewCreateManageLike(
	likeDatabase repositories.LikeRepository,
	like *likeDomain.Like,
	transaction port.ConnectionInterface,
	customError port.CustomErrorInterface,
) *CreateManageLike {
	return &CreateManageLike{
		LikeDatabase: likeDatabase,
		Like:         like,
		Transaction:  transaction,
		CustomError:  customError,
	}
}

func (createManageLike *CreateManageLike) Execute() (int, error) {
	if err := createManageLike.Transaction.SavePoint(checkPointManageLikeTransactionNameConst); err != nil {
		return 0, createManageLike.CustomError.ThrowError(err.Error())
	}

	err := createManageLike.manageLike()
	if err != nil {
		if rollbackErr := createManageLike.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return 0, rollbackErr
		}
		return 0, err
	}

	updatedLikeCount, err := createManageLike.LikeDatabase.FindCurrentLikesCount(createManageLike.Like)
	if err != nil {
		return 0, createManageLike.CustomError.ThrowError(err.Error())
	}

	if err := createManageLike.Transaction.Commit(); err != nil {
		if rollbackErr := createManageLike.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return 0, rollbackErr
		}
		return 0, createManageLike.CustomError.ThrowError(err.Error())
	}

	return updatedLikeCount, nil
}

func (createManageLike *CreateManageLike) manageLike() error {
	alreadyLiked, existingLikeInfo, err := createManageLike.LikeDatabase.ConfirmExistingLike(createManageLike.Like)
	if err != nil {
		return err
	}

	if alreadyLiked {
		return createManageLike.LikeDatabase.RemoveLike(existingLikeInfo.Id)
	}

	return createManageLike.LikeDatabase.CreateLike(createManageLike.Like)
}

func (createManageLike *CreateManageLike) rollbackToSavePointAndCommit() error {
	if transactErr := createManageLike.Transaction.RollbackTo(checkPointManageLikeTransactionNameConst); transactErr != nil {
		return createManageLike.CustomError.ThrowError(transactErr.Error())
	}

	if commitErr := createManageLike.Transaction.Commit(); commitErr != nil {
		return createManageLike.CustomError.ThrowError(commitErr.Error())
	}

	return nil
}
