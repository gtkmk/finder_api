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

func (editPost *EditPost) Execute(loggedUserId string) error {
	if err := editPost.verifyIfPostExists(); err != nil {
		return err
	}

	if err := editPost.verifyIfPostBelongsToUser(loggedUserId); err != nil {
		return err
	}

	if err := editPost.Transaction.SavePoint(checkPointEditPostTransactionNameConst); err != nil {
		return editPost.ThrowError(err.Error())
	}

	err := editPost.persistPostEditing()
	if err != nil {
		if rollbackErr := editPost.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err := editPost.Transaction.Commit(); err != nil {
		if rollbackErr := editPost.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return rollbackErr
		}
		return editPost.ThrowError(err.Error())
	}

	return nil
}

func (editPost *EditPost) verifyIfPostExists() error {
	post, err := editPost.PostDatabase.FindPostByID(editPost.Post.Id)

	if err != nil {
		return editPost.ThrowError(err.Error())
	}

	if post == nil {
		return editPost.ThrowError(helper.PostNotFoundMessageConst)
	}

	editPost.FoundPost = post

	return nil
}

func (editPost *EditPost) verifyIfPostBelongsToUser(loggedUserId string) error {
	if editPost.FoundPost.UserId != loggedUserId {
		return editPost.ThrowError(helper.CannotEditPostNotOwnedConst)
	}

	return nil
}

func (editPost *EditPost) persistPostEditing() error {
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
