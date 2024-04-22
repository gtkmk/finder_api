package repositories

import (
	"github.com/gtkmk/finder_api/core/domain/commentDomain"
)

type CommentRepository interface {
	CreateComment(commentInfo *commentDomain.Comment) error
	FindCommentByID(id string) (*commentDomain.Comment, error)
	EditComment(comment *commentDomain.Comment) error
	DeleteComment(id string) error
	FindAllComments(postId string, offset int64, limit int64) ([]map[string]interface{}, error)
}
