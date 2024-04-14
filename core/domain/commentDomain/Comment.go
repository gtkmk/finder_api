package commentDomain

import (
	"time"
)

type Comment struct {
	Id        string     `json:"id"`
	Text      string     `json:"text"`
	PostId    string     `json:"post_id"`
	UserId    string     `json:"user_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewComment(
	id string,
	text string,
	postId string,
	userID string,
	createdAt *time.Time,
	updatedAt *time.Time,
) *Comment {
	return &Comment{
		Id:        id,
		Text:      text,
		PostId:    postId,
		UserId:    userID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
