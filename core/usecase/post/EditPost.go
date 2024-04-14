package postUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/domain/postDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type EditPost struct {
	PostDatabase repositories.PostRepositoryInterface
	FileService  port.FileFactoryInterface
	Post         postDomain.Post
	Transaction  port.ConnectionInterface
	FoundPost    *postDomain.Post
	port.CustomErrorInterface
}

const checkPointEditPostTransactionNameConst = "editPost"

func NewEditPost(
	postDatabase repositories.PostRepositoryInterface,
	fileService port.FileFactoryInterface,
	post postDomain.Post,
	transaction port.ConnectionInterface,
) *EditPost {
	return &EditPost{
		PostDatabase:         postDatabase,
		FileService:          fileService,
		Post:                 post,
		Transaction:          transaction,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (editPost *EditPost) Execute() error {
	if err := editPost.verifyIfPostExists(); err != nil {
		return err
	}

	if err := editPost.Transaction.SavePoint(checkPointEditPostTransactionNameConst); err != nil {
		return editPost.ThrowError(err.Error())
	}

	err := editPost.persistPostEdition()
	if err != nil {
		if rollbackErr := editPost.rollbackToSavePointAndCommit(); err != nil {
			return rollbackErr
		}
		return err
	}

	if err := editPost.Transaction.Commit(); err != nil {
		if rollbackErr := editPost.rollbackToSavePointAndCommit(); err != nil {
			return rollbackErr
		}
		return editPost.ThrowError(err.Error())
	}

	return nil
}

func (editPost *EditPost) verifyIfPostExists() error {
	proposal, err := editPost.PostDatabase.FindPostByID(editPost.Post.Id)

	if err != nil {
		return editPost.ThrowError(err.Error())
	}

	if proposal == nil {
		return editPost.ThrowError(helper.PostNotFoundMessageConst)
	}

	editPost.FoundPost = proposal

	return nil
}

func (editPost *EditPost) persistPostEdition() error {
	if err := editPost.PostDatabase.EditPost(
		&editPost.Post,
	); err != nil {
		return editPost.ThrowError(err.Error())
	}

	return nil
}

func (editPost *EditPost) rollbackToSavePointAndCommit() error {
	if transactErr := editPost.Transaction.RollbackTo(checkPointEditPostTransactionNameConst); transactErr != nil {
		return editPost.ThrowError(transactErr.Error())
	}

	if commitErr := editPost.Transaction.Commit(); commitErr != nil {
		return editPost.ThrowError(commitErr.Error())
	}

	return nil
}
