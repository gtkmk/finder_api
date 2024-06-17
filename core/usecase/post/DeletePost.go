package postUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type DeletePost struct {
	PostDatabase repositories.PostRepositoryInterface
	Transaction  port.ConnectionInterface
	PostId       string
	LoggedUserId string
	port.CustomErrorInterface
}

const checkPointDeletePostTransactionNameConst = "deletePost"

func NewDeletePost(
	postDatabase repositories.PostRepositoryInterface,
	transaction port.ConnectionInterface,
	postId string,
	loggedUserId string,
) *DeletePost {
	return &DeletePost{
		PostDatabase:         postDatabase,
		Transaction:          transaction,
		PostId:               postId,
		LoggedUserId:         loggedUserId,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (deleteComment *DeletePost) Execute() error {
	if err := deleteComment.verifyIfPostExistsAndCanBeDeleted(); err != nil {
		return err
	}

	if err := deleteComment.Transaction.SavePoint(checkPointDeletePostTransactionNameConst); err != nil {
		return deleteComment.ThrowError(err.Error())
	}

	err := deleteComment.deletePost()
	if err != nil {
		if rollbackErr := deleteComment.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err := deleteComment.Transaction.Commit(); err != nil {
		if rollbackErr := deleteComment.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return rollbackErr
		}
		return deleteComment.ThrowError(err.Error())
	}

	return nil
}

func (deleteComment *DeletePost) verifyIfPostExistsAndCanBeDeleted() error {
	post, err := deleteComment.PostDatabase.FindPostByID(deleteComment.PostId)

	if err != nil {
		return deleteComment.ThrowError(err.Error())
	}

	if post == nil {
		return deleteComment.ThrowError(helper.PostNotFoundMessageConst)
	}

	if post.UserId != deleteComment.LoggedUserId {
		return deleteComment.ThrowError(helper.UnauthorizedConst)
	}

	return nil
}

func (deleteComment *DeletePost) deletePost() error {
	if err := deleteComment.PostDatabase.DeletePost(
		deleteComment.PostId,
	); err != nil {
		return deleteComment.ThrowError(err.Error())
	}

	return nil
}

func (deleteComment *DeletePost) rollbackToSavePointAndCommit() error {
	if transactErr := deleteComment.Transaction.RollbackTo(checkPointDeletePostTransactionNameConst); transactErr != nil {
		return deleteComment.ThrowError(transactErr.Error())
	}

	if commitErr := deleteComment.Transaction.Commit(); commitErr != nil {
		return deleteComment.ThrowError(commitErr.Error())
	}

	return nil
}
