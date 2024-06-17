package repositories

import "github.com/gtkmk/finder_api/core/domain/documentDomain"

type DocumentRepository interface {
	CreateMedia(document *documentDomain.Document, documentPath string) error
	FindCurrentUserMediaByType(userId string, mediaType string) (*documentDomain.Document, error)
	DeleteUserMedia(documentId string) error
}
