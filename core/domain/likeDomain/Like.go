package likeDomain

import "time"

type Like struct {
	Id        string     `json:"id"`
	LikeType  string     `json:"type"`
	PostId    *string    `json:"post_id"`
	CommentId *string    `json:"comment_id"`
	UserId    string     `json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
}

const (
	PostTypeConst    = "post"
	CommentTypeConst = "comment"
)

var LikeAcceptedTypes = []string{
	PostTypeConst,
	CommentTypeConst,
}

func NewLike(
	id string,
	likeType string,
	postId *string,
	commentId *string,
	userID string,
	createdAt *time.Time,
) *Like {
	return &Like{
		Id:        id,
		LikeType:  likeType,
		PostId:    postId,
		CommentId: commentId,
		UserId:    userID,
		CreatedAt: createdAt,
	}
}
