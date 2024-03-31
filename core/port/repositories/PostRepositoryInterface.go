package repositories

import (
	"github.com/gtkmk/finder_api/core/domain/documentDomain"
	"github.com/gtkmk/finder_api/core/domain/postDomain"
)

type PostRepositoryInterface interface {
	CreatePost(post *postDomain.Post) error
	CreatePostMedia(document *documentDomain.Document, documentPath string) error
}
