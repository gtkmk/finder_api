package repository

import (
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/followDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	"github.com/gtkmk/finder_api/infra/database/models"
)

type FollowDatabase struct {
	connection port.ConnectionInterface
}

func NewFollowDatabase(connection port.ConnectionInterface) repositories.FollowRepository {
	return &FollowDatabase{
		connection,
	}
}

func (followDatabase *FollowDatabase) ConfirmExistingFollow(followInfo *followDomain.Follow) (bool, *followDomain.Follow, error) {
	query := `
		SELECT
			id,
			follower_id,
			followed_id,
			created_at
		FROM follow
		WHERE follower_id = ? AND followed_id = ?`

	var databaseFollow *models.Follow

	if err := followDatabase.connection.Raw(
		query,
		&databaseFollow,
		followInfo.FollowerId,
		followInfo.FollowedId,
	); err != nil {
		return false, nil, err
	}

	if databaseFollow == nil {
		return false, nil, nil
	}

	follow := followDomain.NewFollow(
		databaseFollow.ID,
		databaseFollow.FollowerId,
		databaseFollow.FollowedId,
		&databaseFollow.CreatedAt.Time,
	)

	return true, follow, nil

}

func (followDatabase *FollowDatabase) CreateFollow(followInfo *followDomain.Follow) error {
	createdAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO follow (
			id,
			follower_id,
			followed_id,
			created_at
		) 
		VALUES (?, ?, ?, ?)`

	var statement interface{}

	if err := followDatabase.connection.Raw(
		query,
		statement,
		followInfo.Id,
		followInfo.FollowerId,
		followInfo.FollowedId,
		createdAt,
	); err != nil {
		return err
	}

	return nil
}

func (followDatabase *FollowDatabase) RemoveFollow(followId string) error {
	query := `DELETE FROM follow WHERE id = ?`

	var statement interface{}

	return followDatabase.connection.Raw(
		query,
		&statement,
		followId,
	)
}
