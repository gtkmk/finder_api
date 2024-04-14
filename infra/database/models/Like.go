package models

import (
	"database/sql"
)

type Like struct {
	ID        string
	LikeType  string
	PostId    string
	CommentId string
	UserId    string
	CreatedAt sql.NullTime
}
