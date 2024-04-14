package repositories

import (
	"github.com/gtkmk/finder_api/core/domain/commentDomain"
)

type CommentRepository interface {
	CreateComment(commentInfo *commentDomain.Comment) error
	FindCommentByID(id string) (*commentDomain.Comment, error)
	EditComment(comment *commentDomain.Comment) error
	DeleteComment(id string) error
}
