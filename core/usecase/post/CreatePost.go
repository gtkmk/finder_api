package postUsecase

import (
	"fmt"

	"github.com/gtkmk/finder_api/core/domain/customError"
	"github.com/gtkmk/finder_api/core/domain/postDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
)

type CreatePost struct {
	PostDatabase repositories.PostRepositoryInterface
	FileService  port.FileFactoryInterface
	Post         postDomain.Post
	transaction  port.ConnectionInterface
	dist         string
	port.CustomErrorInterface
}

const CheckPointTransactionNameConst = "createPost"

func NewCreatePost(
	postDatabase repositories.PostRepositoryInterface,
	fileService port.FileFactoryInterface,
	post postDomain.Post,
	transaction port.ConnectionInterface,
	dist string,
) *CreatePost {
	return &CreatePost{
		PostDatabase:         postDatabase,
		FileService:          fileService,
		Post:                 post,
		transaction:          transaction,
		dist:                 dist,
		CustomErrorInterface: customError.NewCustomError(),
	}
}

func (createPost *CreatePost) Execute() error {
	fmt.Print("***************************************************************")
	fmt.Print(createPost.Post)
	if err := createPost.transaction.SavePoint(CheckPointTransactionNameConst); err != nil {
		return createPost.ThrowError(err.Error())
	}

	err := createPost.persistPostAndMedia()
	if err != nil {
		if rollbackErr := createPost.rollbackToSavePointAndCommit(); err != nil {
			return rollbackErr
		}
		return err
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

	if err := createPost.PostDatabase.CreatePostMedia(
		createPost.Post.Media,
		fileService.GetFullPath(),
	); err != nil {
		if removeFileErr := fileService.RemoveTempFile(); removeFileErr != nil {
			return createPost.ThrowError(removeFileErr.Error())
		}

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

	return nil
}

func (createPost *CreatePost) rollbackToSavePointAndCommit() error {
	if transactErr := createPost.transaction.RollbackTo(CheckPointTransactionNameConst); transactErr != nil {
		return createPost.ThrowError(transactErr.Error())
	}

	if commitErr := createPost.transaction.Commit(); commitErr != nil {
		return createPost.ThrowError(commitErr.Error())
	}

	return nil
}
