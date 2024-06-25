package repository

import (
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/documentDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	"github.com/gtkmk/finder_api/infra/database/models"
)

type DocumentDatabase struct {
	connection port.ConnectionInterface
}

func NewDocumentDatabase(connection port.ConnectionInterface) repositories.DocumentRepository {
	return &DocumentDatabase{
		connection,
	}
}

func (documentDatabase *DocumentDatabase) CreateMedia(
	document *documentDomain.Document,
	documentPath string,
) error {
	createdAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO document (
			id,
			type,
			path,
			post_id,
			owner_id,
			mime_type,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`

	var statement interface{}

	return documentDatabase.connection.Raw(
		query,
		&statement,
		document.ID,
		document.Type,
		documentPath,
		document.PostId,
		document.OwnerId,
		document.MimeType,
		createdAt,
	)
}

func (documentDatabase *DocumentDatabase) FindCurrentUserMediaByType(userId string, mediaType string) (*documentDomain.Document, error) {
	query := `SELECT * FROM document WHERE owner_id = ? AND type = ? AND deleted_at IS NULL`

	var databaseDocument *models.Document

	if err := documentDatabase.connection.Raw(query, &databaseDocument, userId, mediaType); err != nil {
		return nil, err
	}

	if databaseDocument == nil {
		return nil, nil
	}

	post := documentDomain.NewDocument(
		databaseDocument.ID,
		databaseDocument.Type,
		nil,
		"",
		&databaseDocument.PostId,
		databaseDocument.OwnerId,
		databaseDocument.MimeType,
		"",
	)

	return post, nil
}

func (documentDatabase *DocumentDatabase) DeleteUserMedia(documentId string) error {
	deletedAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	query := `
		UPDATE document SET 
			deleted_at = ?
		WHERE id = ?`

	var statement interface{}

	return documentDatabase.connection.Raw(
		query,
		&statement,
		deletedAt,
		documentId,
	)
}
