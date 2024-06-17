package documentUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/documentDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	sharedMethodsInterface "github.com/gtkmk/finder_api/core/port/sharedMethods"
)

const updateUserImagesTransactionNameConst = "createPost"

type UpdateDocumentChangeProfileImage struct {
	DocumentDatabase  repositories.DocumentRepository
	Document          *documentDomain.Document
	Dist              string
	FileService       port.FileFactoryInterface
	Transaction       port.ConnectionInterface
	RollBackAndReturn sharedMethodsInterface.RollBackAndReturnErrorInterface
	CustomError       port.CustomErrorInterface
}

func NewUpdateDocumentChangeProfileImage(
	documentDatabase repositories.DocumentRepository,
	fileService port.FileFactoryInterface,
	document *documentDomain.Document,
	dist string,
	transaction port.ConnectionInterface,
	rollBackAndReturn sharedMethodsInterface.RollBackAndReturnErrorInterface,
	customError port.CustomErrorInterface,
) *UpdateDocumentChangeProfileImage {
	return &UpdateDocumentChangeProfileImage{
		DocumentDatabase:  documentDatabase,
		FileService:       fileService,
		Document:          document,
		Dist:              dist,
		Transaction:       transaction,
		RollBackAndReturn: rollBackAndReturn,
		CustomError:       customError,
	}
}

func (updateDocumentChangeProfileImage *UpdateDocumentChangeProfileImage) Execute(documentType string) error {
	foundExistingMedia, err := updateDocumentChangeProfileImage.DocumentDatabase.FindCurrentUserMediaByType(
		updateDocumentChangeProfileImage.Document.OwnerId,
		documentType,
	)

	if err != nil {
		return updateDocumentChangeProfileImage.CustomError.ThrowError(err.Error())
	}

	if err := updateDocumentChangeProfileImage.Transaction.SavePoint(updateUserImagesTransactionNameConst); err != nil {
		return updateDocumentChangeProfileImage.CustomError.ThrowError(err.Error())
	}

	err = updateDocumentChangeProfileImage.persistNewMedia()
	if err != nil {
		if rollbackErr := updateDocumentChangeProfileImage.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	err = updateDocumentChangeProfileImage.deletePreviousUserMedia(foundExistingMedia.ID)
	if err != nil {
		if rollbackErr := updateDocumentChangeProfileImage.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	if err = updateDocumentChangeProfileImage.Transaction.Commit(); err != nil {
		if rollbackErr := updateDocumentChangeProfileImage.rollbackToSavePointAndCommit(); rollbackErr != nil {
			return rollbackErr
		}
		return updateDocumentChangeProfileImage.CustomError.ThrowError(err.Error())
	}

	return nil
}

func (updateDocumentChangeProfileImage *UpdateDocumentChangeProfileImage) persistNewMedia() error {
	fileService := updateDocumentChangeProfileImage.FileService.Make(updateDocumentChangeProfileImage.Document.NewName)

	if err := fileService.SaveFileFromMultipart(
		updateDocumentChangeProfileImage.Document.File,
		updateDocumentChangeProfileImage.Dist,
	); err != nil {
		return updateDocumentChangeProfileImage.CustomError.ThrowError(err.Error())
	}

	if err := updateDocumentChangeProfileImage.DocumentDatabase.CreateMedia(
		updateDocumentChangeProfileImage.Document,
		fileService.GetFullPath(),
	); err != nil {
		if removeFileErr := fileService.RemoveTempFile(); removeFileErr != nil {
			return updateDocumentChangeProfileImage.CustomError.ThrowError(removeFileErr.Error())
		}

		return updateDocumentChangeProfileImage.CustomError.ThrowError(err.Error())
	}

	return nil
}

func (updateDocumentChangeProfileImage *UpdateDocumentChangeProfileImage) deletePreviousUserMedia(foundExistingMediaId string) error {
	fileService := updateDocumentChangeProfileImage.FileService.Make(updateDocumentChangeProfileImage.Document.NewName)

	if err := updateDocumentChangeProfileImage.DocumentDatabase.DeleteUserMedia(
		foundExistingMediaId,
	); err != nil {
		if removeFileErr := fileService.RemoveTempFile(); removeFileErr != nil {
			return updateDocumentChangeProfileImage.CustomError.ThrowError(removeFileErr.Error())
		}

		return updateDocumentChangeProfileImage.CustomError.ThrowError(err.Error())
	}

	return nil
}

func (updateDocumentChangeProfileImage *UpdateDocumentChangeProfileImage) rollbackToSavePointAndCommit() error {
	if transactErr := updateDocumentChangeProfileImage.Transaction.RollbackTo(updateUserImagesTransactionNameConst); transactErr != nil {
		return updateDocumentChangeProfileImage.CustomError.ThrowError(transactErr.Error())
	}

	if commitErr := updateDocumentChangeProfileImage.Transaction.Commit(); commitErr != nil {
		return updateDocumentChangeProfileImage.CustomError.ThrowError(commitErr.Error())
	}

	return nil
}
