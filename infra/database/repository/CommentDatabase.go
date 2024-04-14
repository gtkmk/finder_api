package repository

import (
	"github.com/gtkmk/finder_api/core/domain/commentDomain"
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/repositories"
	"github.com/gtkmk/finder_api/infra/database/models"
)

type CommentDatabase struct {
	connection port.ConnectionInterface
}

func NewCommentDatabase(connection port.ConnectionInterface) repositories.CommentRepository {
	return &CommentDatabase{
		connection,
	}
}

func (commentDatabase *CommentDatabase) CreateComment(commentInfo *commentDomain.Comment) error {
	createdAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO comment (
			id,
			post_id ,
			user_id ,
			text,
			created_at
		) 
		VALUES (?, ?, ?, ?, ?)`

	var statement interface{}

	if err := commentDatabase.connection.Raw(
		query,
		statement,
		commentInfo.Id,
		commentInfo.PostId,
		commentInfo.UserId,
		commentInfo.Text,
		createdAt,
	); err != nil {
		return err
	}

	return nil
}

func (commentDatabase *CommentDatabase) FindCommentByID(id string) (*commentDomain.Comment, error) {
	query := `SELECT * FROM comment WHERE id = ? AND deleted_at IS NULL`

	var databaseComment *models.Comment

	if err := commentDatabase.connection.Raw(query, &databaseComment, id); err != nil {
		return nil, err
	}

	if databaseComment == nil {
		return nil, nil
	}

	post := commentDomain.NewComment(
		databaseComment.ID,
		databaseComment.Text,
		databaseComment.PostId,
		databaseComment.UserId,
		&databaseComment.CreatedAt.Time,
		&databaseComment.UpdatedAt.Time,
	)

	return post, nil
}

func (commentDatabase *CommentDatabase) EditComment(comment *commentDomain.Comment) error {
	updatedAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	var dbComment *models.Comment

	query := `UPDATE comment SET 
		text = ?,
		updated_at = ?
	WHERE id = ?`

	if err := commentDatabase.connection.Raw(
		query,
		&dbComment,
		comment.Text,
		updatedAt,
		comment.Id,
	); err != nil {
		return err
	}

	return nil
}
