package commentUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/commentDomain"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type EditComment struct {
	CommentDatabase repositories.CommentRepository
	Comment         commentDomain.Comment
	Transaction     port.ConnectionInterface
	port.CustomErrorInterface
}

const checkPointEditCommentTransactionNameConst = "editComment"

func NewEditComment(
	commentDatabase repositories.CommentRepository,
	comment commentDomain.Comment,
	transaction port.ConnectionInterface,
) *EditComment {
	return &EditComment{
		CommentDatabase:      commentDatabase,
		Comment:              comment,
		Transaction:          transaction,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (createComment *EditComment) Execute() error {
	if err := createComment.verifyIfCommentExists(); err != nil {
		return err
	}

	if err := createComment.Transaction.SavePoint(checkPointEditCommentTransactionNameConst); err != nil {
		return createComment.ThrowError(err.Error())
	}

	err := createComment.persistComment()
	if err != nil {
		if rollbackErr := createComment.rollbackToSavePointAndCommit(); err != nil {
			return rollbackErr
		}
		return err
	}

	if err := createComment.Transaction.Commit(); err != nil {
		if rollbackErr := createComment.rollbackToSavePointAndCommit(); err != nil {
			return rollbackErr
		}
		return createComment.ThrowError(err.Error())
	}

	return nil
}

func (createComment *EditComment) verifyIfCommentExists() error {
	proposal, err := createComment.CommentDatabase.FindCommentByID(createComment.Comment.Id)

	if err != nil {
		return createComment.ThrowError(err.Error())
	}

	if proposal == nil {
		return createComment.ThrowError(helper.CommentNotFoundMessageConst)
	}

	return nil
}

func (createComment *EditComment) persistComment() error {
	if err := createComment.CommentDatabase.EditComment(
		&createComment.Comment,
	); err != nil {
		return createComment.ThrowError(err.Error())
	}

	return nil
}

func (createComment *EditComment) rollbackToSavePointAndCommit() error {
	if transactErr := createComment.Transaction.RollbackTo(checkPointEditCommentTransactionNameConst); transactErr != nil {
		return createComment.ThrowError(transactErr.Error())
	}

	if commitErr := createComment.Transaction.Commit(); commitErr != nil {
		return createComment.ThrowError(commitErr.Error())
	}

	return nil
}
