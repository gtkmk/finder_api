package postUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/documentDomain"
	"github.com/gtkmk/finder_api/core/domain/postDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type CreatePost struct {
	PostDatabase     repositories.PostRepositoryInterface
	DocumentDatabase repositories.DocumentRepository
	FileService      port.FileFactoryInterface
	Post             postDomain.Post
	transaction      port.ConnectionInterface
	dist             string
	port.CustomErrorInterface
}

const checkPointCreatePostTransactionNameConst = "createPost"

func NewCreatePost(
	postDatabase repositories.PostRepositoryInterface,
	documentDatabase repositories.DocumentRepository,
	fileService port.FileFactoryInterface,
	post postDomain.Post,
	transaction port.ConnectionInterface,
	dist string,
) *CreatePost {
	return &CreatePost{
		PostDatabase:         postDatabase,
		DocumentDatabase:     documentDatabase,
		FileService:          fileService,
		Post:                 post,
		transaction:          transaction,
		dist:                 dist,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (createPost *CreatePost) Execute() error {
	if err := createPost.transaction.SavePoint(checkPointCreatePostTransactionNameConst); err != nil {
		return createPost.ThrowError(err.Error())
	}

	err := createPost.persistPostAndMedia()
	if err != nil {
		if rollbackErr := createPost.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err := createPost.transaction.Commit(); err != nil {
		if rollbackErr := createPost.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return rollbackErr
		}
		return createPost.ThrowError(err.Error())
	}

	return nil
}

func (createPost *CreatePost) persistPostAndMedia() error {
	fileService := createPost.FileService.Make(createPost.Post.Media.NewName)

	if err := fileService.SaveFileFromMultipart(
		createPost.Post.Media.File,
		createPost.dist,
	); err != nil {
		return createPost.ThrowError(err.Error())
	}

	if err := createPost.PostDatabase.CreatePost(
		&createPost.Post,
	); err != nil {
		if removeFileErr := fileService.RemoveTempFile(); removeFileErr != nil {
			return createPost.ThrowError(removeFileErr.Error())
		}

		return createPost.ThrowError(err.Error())
	}

	createPost.Post.Media.Type = documentDomain.PostMediaConst

	if err := createPost.DocumentDatabase.CreateMedia(
		createPost.Post.Media,
		fileService.GetFullPath(),
	); err != nil {
		if removeFileErr := fileService.RemoveTempFile(); removeFileErr != nil {
			return createPost.ThrowError(removeFileErr.Error())
		}

		return createPost.ThrowError(err.Error())
	}

	return nil
}

func (createPost *CreatePost) rollbackToSavePointAndCommit() error {
	if transactErr := createPost.transaction.RollbackTo(checkPointCreatePostTransactionNameConst); transactErr != nil {
		return createPost.ThrowError(transactErr.Error())
	}

	if commitErr := createPost.transaction.Commit(); commitErr != nil {
		return createPost.ThrowError(commitErr.Error())
	}

	return nil
}
