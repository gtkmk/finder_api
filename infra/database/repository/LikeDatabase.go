package repository

import (
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/likeDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	"github.com/gtkmk/finder_api/infra/database/models"
)

type LikeDatabase struct {
	connection port.ConnectionInterface
}

func NewLikeDatabase(connection port.ConnectionInterface) repositories.LikeRepository {
	return &LikeDatabase{
		connection,
	}
}

func (likeDatabase *LikeDatabase) ConfirmExistingLike(likeInfo *likeDomain.Like) (bool, *likeDomain.Like, error) {
	query := `
		SELECT 
			id,
			like_type,
			post_id,
			comment_id,
			user_id,
			created_at
		FROM interaction_likes
		WHERE like_type = ? AND (post_id = ? OR comment_id = ?)`

	var databaseLike *models.Like

	if err := likeDatabase.connection.Raw(
		query,
		&databaseLike,
		likeInfo.LikeType,
		likeInfo.PostId,
		likeInfo.CommentId,
	); err != nil {
		return false, nil, err
	}

	if databaseLike == nil {
		return false, nil, nil
	}

	unityUser := likeDomain.NewLike(
		databaseLike.ID,
		databaseLike.LikeType,
		&databaseLike.PostId,
		&databaseLike.CommentId,
		databaseLike.UserId,
		&databaseLike.CreatedAt.Time,
	)

	return true, unityUser, nil
}

func (likeDatabase *LikeDatabase) CreateLike(likeInfo *likeDomain.Like) error {
	createdAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO interaction_likes (
			id,
			like_type,
			post_id,
			comment_id,
			user_id,
			created_at
		) 
		VALUES (?, ?, ?, ?, ?, ?)`

	var statement interface{}

	if err := likeDatabase.connection.Raw(
		query,
		statement,
		likeInfo.Id,
		likeInfo.LikeType,
		likeInfo.PostId,
		likeInfo.CommentId,
		likeInfo.UserId,
		createdAt,
	); err != nil {
		return err
	}

	return nil
}

func (likeDatabase *LikeDatabase) RemoveLike(likeId string) error {
	query := `DELETE FROM interaction_likes WHERE id = ?`

	var statement interface{}

	return likeDatabase.connection.Raw(
		query,
		&statement,
		likeId,
	)
}

func (likeDatabase *LikeDatabase) FindCurrentLikesCount(likeInfo *likeDomain.Like) (int, error) {
	var amount int

	query := `SELECT COUNT(*) FROM interaction_likes WHERE like_type = ? AND (post_id = ? OR comment_id = ?)`

	if err := likeDatabase.connection.Raw(query, &amount, likeInfo.LikeType, likeInfo.PostId, likeInfo.CommentId); err != nil {
		return 0, err
	}

	return amount, nil

}
