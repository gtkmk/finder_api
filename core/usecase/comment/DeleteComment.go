package postUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type DeleteComment struct {
	CommentDatabase repositories.CommentRepository
	PostDatabase    repositories.PostRepositoryInterface
	Transaction     port.ConnectionInterface
	CommentId       string
	LoggedUserId    string
	port.CustomErrorInterface
}

const checkPointDeleteCommentTransactionNameConst = "editComment"

func NewDeleteComment(
	commentDatabase repositories.CommentRepository,
	transaction port.ConnectionInterface,
	commentId string,
	loggedUserId string,
) *DeleteComment {
	return &DeleteComment{
		CommentDatabase:      commentDatabase,
		Transaction:          transaction,
		CommentId:            commentId,
		LoggedUserId:         loggedUserId,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (deleteComment *DeleteComment) Execute() error {
	if err := deleteComment.verifyIfCommentExistsAndCanBeDeleted(); err != nil {
		return err
	}

	if err := deleteComment.Transaction.SavePoint(checkPointDeleteCommentTransactionNameConst); err != nil {
		return deleteComment.ThrowError(err.Error())
	}

	err := deleteComment.deleteComment()
	if err != nil {
		if rollbackErr := deleteComment.rollbackToSavePointAndCommit(); err != nil {
			return rollbackErr
		}
		return err
	}

	if err := deleteComment.Transaction.Commit(); err != nil {
		if rollbackErr := deleteComment.rollbackToSavePointAndCommit(); err != nil {
			return rollbackErr
		}
		return deleteComment.ThrowError(err.Error())
	}

	return nil
}

func (deleteComment *DeleteComment) verifyIfCommentExistsAndCanBeDeleted() error {
	comment, err := deleteComment.CommentDatabase.FindCommentByID(deleteComment.CommentId)

	if err != nil {
		return deleteComment.ThrowError(err.Error())
	}

	if comment == nil {
		return deleteComment.ThrowError(helper.CommentNotFoundMessageConst)
	}

	if comment.UserId != deleteComment.LoggedUserId {
		return deleteComment.ThrowError(helper.UnauthorizedConst)
	}

	return nil
}

func (deleteComment *DeleteComment) deleteComment() error {
	if err := deleteComment.CommentDatabase.DeleteComment(
		deleteComment.CommentId,
	); err != nil {
		return deleteComment.ThrowError(err.Error())
	}

	return nil
}

func (deleteComment *DeleteComment) rollbackToSavePointAndCommit() error {
	if transactErr := deleteComment.Transaction.RollbackTo(checkPointDeleteCommentTransactionNameConst); transactErr != nil {
		return deleteComment.ThrowError(transactErr.Error())
	}

	if commitErr := deleteComment.Transaction.Commit(); commitErr != nil {
		return deleteComment.ThrowError(commitErr.Error())
	}

	return nil
}
