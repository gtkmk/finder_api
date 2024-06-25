package repositories

import (
	"github.com/gtkmk/finder_api/core/domain/filterDomain"
	"github.com/gtkmk/finder_api/core/domain/postDomain"
)

type PostRepositoryInterface interface {
	CreatePost(post *postDomain.Post) error
	FindAllPosts(filter *filterDomain.PostFilter) ([]map[string]interface{}, error)
	FindPostByID(id string) (*postDomain.Post, error)
	EditPost(post *postDomain.Post) error
	DeletePost(id string) error
	UpdatePostFoundStatus(id string, postFoundStatus bool) error
}
