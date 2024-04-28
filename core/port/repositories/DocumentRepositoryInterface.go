package repositories

import "github.com/gtkmk/finder_api/core/domain/documentDomain"

type DocumentRepository interface {
	CreateMedia(document *documentDomain.Document, documentPath string) error
}
