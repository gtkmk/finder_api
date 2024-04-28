package models

import "database/sql"

type Follow struct {
	ID         string
	FollowerId string
	FollowedId string
	CreatedAt  sql.NullTime
}
