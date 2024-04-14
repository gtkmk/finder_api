package postUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/commentDomain"
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type CreateComment struct {
	CommentDatabase repositories.CommentRepository
	PostDatabase    repositories.PostRepositoryInterface
	FileService     port.FileFactoryInterface
	Comment         commentDomain.Comment
	Transaction     port.ConnectionInterface
	port.CustomErrorInterface
}

const checkPointCreateCommentTransactionNameConst = "createComment"

func NewCreateComment(
	commentDatabase repositories.CommentRepository,
	postDatabase repositories.PostRepositoryInterface,
	comment commentDomain.Comment,
	transaction port.ConnectionInterface,
) *CreateComment {
	return &CreateComment{
		CommentDatabase:      commentDatabase,
		PostDatabase:         postDatabase,
		Comment:              comment,
		Transaction:          transaction,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (createComment *CreateComment) Execute() error {
	if err := createComment.verifyIfPostExists(); err != nil {
		return err
	}

	if err := createComment.Transaction.SavePoint(checkPointCreateCommentTransactionNameConst); err != nil {
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

func (createComment *CreateComment) verifyIfPostExists() error {
	proposal, err := createComment.PostDatabase.FindPostByID(createComment.Comment.PostId)

	if err != nil {
		return createComment.ThrowError(err.Error())
	}

	if proposal == nil {
		return createComment.ThrowError(helper.PostNotFoundMessageConst)
	}

	return nil
}

func (createComment *CreateComment) persistComment() error {
	if err := createComment.CommentDatabase.CreateComment(
		&createComment.Comment,
	); err != nil {
		return createComment.ThrowError(err.Error())
	}

	return nil
}

func (createComment *CreateComment) rollbackToSavePointAndCommit() error {
	if transactErr := createComment.Transaction.RollbackTo(checkPointCreateCommentTransactionNameConst); transactErr != nil {
		return createComment.ThrowError(transactErr.Error())
	}

	if commitErr := createComment.Transaction.Commit(); commitErr != nil {
		return createComment.ThrowError(commitErr.Error())
	}

	return nil
}
