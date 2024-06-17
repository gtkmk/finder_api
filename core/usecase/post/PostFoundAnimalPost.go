package postUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type FoundAnimalPost struct {
	PostDatabase repositories.PostRepositoryInterface
	Transaction  port.ConnectionInterface
	PostId       string
	FoundStatus  string
	LoggedUserId string
	port.CustomErrorInterface
}

const checkPointUpdatePostFoundStatusTransactionNameConst = "updatePostFoundStatus"

func NewFoundAnimalPost(
	postDatabase repositories.PostRepositoryInterface,
	transaction port.ConnectionInterface,
	postId string,
	foundStatus string,
	loggedUserId string,
) *FoundAnimalPost {
	return &FoundAnimalPost{
		PostDatabase:         postDatabase,
		Transaction:          transaction,
		PostId:               postId,
		FoundStatus:          foundStatus,
		LoggedUserId:         loggedUserId,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (foundAnimalPost *FoundAnimalPost) Execute() error {
	if err := foundAnimalPost.verifyIfPostExistsAndCanBeUpdated(); err != nil {
		return err
	}

	if err := foundAnimalPost.Transaction.SavePoint(checkPointUpdatePostFoundStatusTransactionNameConst); err != nil {
		return foundAnimalPost.ThrowError(err.Error())
	}

	err := foundAnimalPost.updatePostFoundStatus()
	if err != nil {
		if rollbackErr := foundAnimalPost.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err := foundAnimalPost.Transaction.Commit(); err != nil {
		if rollbackErr := foundAnimalPost.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return rollbackErr
		}
		return foundAnimalPost.ThrowError(err.Error())
	}

	return nil
}

func (foundAnimalPost *FoundAnimalPost) verifyIfPostExistsAndCanBeUpdated() error {
	post, err := foundAnimalPost.PostDatabase.FindPostByID(foundAnimalPost.PostId)

	if err != nil {
		return foundAnimalPost.ThrowError(err.Error())
	}

	if post == nil {
		return foundAnimalPost.ThrowError(helper.PostNotFoundMessageConst)
	}

	if post.UserId != foundAnimalPost.LoggedUserId {
		return foundAnimalPost.ThrowError(helper.UnauthorizedConst)
	}

	return nil
}

func (foundAnimalPost *FoundAnimalPost) updatePostFoundStatus() error {
	if err := foundAnimalPost.PostDatabase.UpdatePostFoundStatus(
		foundAnimalPost.PostId,
		foundAnimalPost.FoundStatus,
	); err != nil {
		return foundAnimalPost.ThrowError(err.Error())
	}

	return nil
}

func (foundAnimalPost *FoundAnimalPost) rollbackToSavePointAndCommit() error {
	if transactErr := foundAnimalPost.Transaction.RollbackTo(checkPointUpdatePostFoundStatusTransactionNameConst); transactErr != nil {
		return foundAnimalPost.ThrowError(transactErr.Error())
	}

	if commitErr := foundAnimalPost.Transaction.Commit(); commitErr != nil {
		return foundAnimalPost.ThrowError(commitErr.Error())
	}

	return nil
}
