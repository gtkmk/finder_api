package repository

import (
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/documentDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
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
