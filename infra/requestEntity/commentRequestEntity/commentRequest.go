package commentRequestEntity

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/core/domain/commentDomain"
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

const (
	TextFieldConst = "o texto do comentário"
	PostFieldConst = "a publicação"
)

const (
	TextFieldNameConst = "text"
	PostFieldNameConst = "post_id"
)

const (
	MaxTextLengthConst = 800
)

type CommentRequest struct {
	uuid    port.UuidInterface
	PostId  string `form:"post_id" json:"post_id"`
	UserId  string
	Text    string `form:"text" json:"text"`
	Comment *commentDomain.Comment
}

func NewCommentRequest(context *gin.Context, uuid port.UuidInterface, userId string) (*CommentRequest, error) {
	postRequest := &CommentRequest{
		uuid:   uuid,
		UserId: userId,
	}

	if err := context.ShouldBind(postRequest); err != nil {
		return nil, err
	}

	return postRequest, nil
}

func (postRequest *CommentRequest) Validate(context *gin.Context, edition bool) error {
	if err := postRequest.validatePostFields(edition); err != nil {
		return err
	}

	return nil
}

func (postRequest *CommentRequest) validatePostFields(edition bool) error {
	if err := requestEntityFieldsValidation.ValidateField(
		postRequest.Text,
		TextFieldConst,
		MaxTextLengthConst,
	); err != nil {
		return err
	}

	if !edition {
		if err := requestEntityFieldsValidation.IsValidUUID(
			PostFieldConst,
			postRequest.PostId,
		); err != nil {
			return err
		}
	}

	return nil
}

func (postRequest *CommentRequest) BuildCommentObject(commentId *string) (*commentDomain.Comment, error) {
	dateTime, err := datetimeDomain.CreateNow()
	if err != nil {
		return nil, err
	}

	if commentId == nil {
		generatedId := postRequest.uuid.GenerateUuid()
		commentId = &generatedId
	}

	return commentDomain.NewComment(
		*commentId,
		postRequest.Text,
		postRequest.PostId,
		postRequest.UserId,
		&dateTime,
		nil,
	), nil
}
