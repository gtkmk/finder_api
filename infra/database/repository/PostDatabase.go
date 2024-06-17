package repository

import (
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/filterDomain"
	"github.com/gtkmk/finder_api/core/domain/postDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	"github.com/gtkmk/finder_api/infra/database/models"
)

type PostDatabase struct {
	connection port.ConnectionInterface
}

func NewPostDatabase(connection port.ConnectionInterface) repositories.PostRepositoryInterface {
	return &PostDatabase{
		connection: connection,
	}
}

func (postDatabase *PostDatabase) CreatePost(post *postDomain.Post) error {
	createdAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO post (
			id,
			text,
			location,
			reward,
			lost_found,
			privacy,
			shares_count,
			category,
			animal_type,
			animal_size,
			user_id,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var statement interface{}

	if err := postDatabase.connection.Raw(
		query,
		statement,
		post.Id,
		post.Text,
		post.Location,
		post.Reward,
		post.LostFound,
		post.Privacy,
		post.SharesCount,
		post.Category,
		post.AnimalType,
		post.AnimalSize,
		post.UserId,
		createdAt,
	); err != nil {
		return err
	}

	return nil
}

func (postDatabase *PostDatabase) FindAllPosts(
	filter *filterDomain.PostFilter,
) ([]map[string]interface{}, error) {
	query := `CALL find_paginated_posts (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	dbProposals, err := postDatabase.connection.Rows(
		query,
		filter.LoggedUserId,
		filter.LostFound,
		filter.Neighborhood,
		filter.Reward,
		filter.UserId,
		filter.OnlyFollowingPosts,
		filter.SpecificPost,
		filter.AnimalType,
		filter.AnimalSize,
		filter.Limit,
		filter.OffSet,
	)

	if err != nil {
		return nil, err
	}

	return dbProposals, nil
}

func (postDatabase *PostDatabase) FindPostByID(id string) (*postDomain.Post, error) {
	query := `SELECT * FROM post WHERE id = ? AND deleted_at IS NULL`

	var databasePost *models.Post

	if err := postDatabase.connection.Raw(query, &databasePost, id); err != nil {
		return nil, err
	}

	if databasePost == nil {
		return nil, nil
	}

	post := postDomain.NewPost(
		databasePost.ID,
		databasePost.Text,
		nil,
		databasePost.Location,
		databasePost.Reward,
		databasePost.Privacy,
		databasePost.SharesCount,
		databasePost.Category,
		&databasePost.LostFound,
		&databasePost.AnimalType,
		&databasePost.AnimalSize,
		databasePost.Found,
		&databasePost.UpdatedFoundStatusAt.Time,
		databasePost.UserId,
		&databasePost.CreatedAt.Time,
		&databasePost.UpdatedAt.Time,
	)

	return post, nil
}

func (postDatabase *PostDatabase) EditPost(post *postDomain.Post) error {
	updatedAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	var dbPost *models.Post

	query := `UPDATE post SET 
		text = ?,
		location = ?,
		reward = ?,
		lost_found = ?,
		privacy = ?,
		animal_type = ?,
		animal_size = ?,
		updated_at = ?
	WHERE id = ?`

	if err := postDatabase.connection.Raw(
		query,
		&dbPost,
		post.Text,
		post.Location,
		post.Reward,
		post.LostFound,
		post.Privacy,
		post.AnimalType,
		post.AnimalSize,
		updatedAt,
		post.Id,
	); err != nil {
		return err
	}

	return nil
}

func (postDatabase *PostDatabase) DeletePost(id string) error {
	deletedAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	query := `
		UPDATE post SET 
			deleted_at = ?
		WHERE id = ?`

	var statement interface{}
	return postDatabase.connection.Raw(
		query,
		statement,
		deletedAt,
		id,
	)
}

func (postDatabase *PostDatabase) UpdatePostFoundStatus(id string, postFoundStatus string) error {
	updatedAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	var dbPost *models.Post

	query := `
		UPDATE post SET 
			found = ?,
			updated_found_status_at = ?
		WHERE id = ?`

	if err := postDatabase.connection.Raw(
		query,
		&dbPost,
		postFoundStatus,
		updatedAt,
		id,
	); err != nil {
		return err
	}

	return nil
}
