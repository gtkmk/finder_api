package repository

import (
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/documentDomain"
	"github.com/gtkmk/finder_api/core/domain/postDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
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
			user_id,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		post.UserId,
		createdAt,
	); err != nil {
		return err
	}

	return nil
}

func (postDatabase *PostDatabase) CreatePostMedia(
	document *documentDomain.Document,
	documentPath string,
) error {
	createdAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO document (
			id,
			type,
			path,
			post_id,
			owner_id,
			mime_type,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var statement interface{}

	return postDatabase.connection.Raw(
		query,
		&statement,
		document.ID,
		document.Type,
		documentPath,
		document.PostId,
		document.OwnerId,
		document.MimeType,
		createdAt,
	)
}
