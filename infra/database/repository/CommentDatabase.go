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

	query := `
		UPDATE comment SET 
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

func (commentDatabase *CommentDatabase) DeleteComment(
	id string,
) error {
	deletedAt, err := datetimeDomain.CreateNow()
	if err != nil {
		return err
	}

	query := `
		UPDATE comment SET 
			deleted_at = ?
		WHERE id = ?`

	var statement interface{}
	return commentDatabase.connection.Raw(
		query,
		statement,
		deletedAt,
		id,
	)
}

func (commentDatabase *CommentDatabase) FindAllComments(postId string, loggedUserId string, offset int64, limit int64) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			comment.id AS comment_id,
			comment.text AS text,
			comment.created_at AS created_at,
			comment.updated_at AS updated_at,
			user.name AS author_name,
			user.user_name AS author_username,
			document.type AS comment_author_avatar_type,
			document.path AS comment_author_avatar,
			document.mime_type AS comment_author_avatar_mime_type,
			COALESCE(likes.likes_count, 0) AS likes,
			CASE
				WHEN user.id = ? THEN true
				ELSE false
			END AS is_own_comment,
			total_records.total_count AS total_records
		FROM
			comment
		INNER JOIN 
			user ON comment.user_id = user.id
		LEFT JOIN 
			document ON document.owner_id = user.id AND document.type = 'profile_picture' AND document.deleted_at IS NULL
		LEFT JOIN (
			SELECT 
				comment_id,
				COUNT(*) AS likes_count
			FROM 
				interaction_likes
			WHERE 
				like_type = 'comment'
			GROUP BY 
				comment_id
		) AS likes ON likes.comment_id = comment.id
		CROSS JOIN (
			SELECT 
				COUNT(*) AS total_count
			FROM 
				comment
			WHERE 
				post_id = ?
		) AS total_records
		WHERE 
			comment.post_id = ?
			AND comment.deleted_at IS NULL
		ORDER BY 
			comment.created_at DESC
		LIMIT 
			?
		OFFSET 
			?
		;
	`

	dbComments, err := commentDatabase.connection.Rows(
		query,
		loggedUserId,
		postId,
		postId,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	return dbComments, nil
}
