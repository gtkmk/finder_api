package models

import (
	"database/sql"
)

type Comment struct {
	ID        string
	Text      string
	PostId    string
	UserId    string
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}
