package models

import (
	"database/sql"
)

type Post struct {
	ID                   string
	Text                 string
	Location             string
	Reward               bool
	LostFound            string
	Privacy              string
	SharesCount          int
	Category             string
	AnimalType           string
	AnimalSize           string
	Found                bool
	UpdatedFoundStatusAt sql.NullTime
	UserId               string
	CreatedAt            sql.NullTime
	UpdatedAt            sql.NullTime
}
