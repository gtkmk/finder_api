package likerequestentity

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/domain/likeDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

const (
	PostFieldConst          = "a publicação"
	CommentFieldConst       = "o comentário"
	UserFieldConst          = "o usuário"
	LikeTypeFieldConst      = "o tipo da publicação"
	PostOrCommentFieldConst = "a publicação ou comentário"
)

const (
	PostFieldNameConst     = "post_id"
	CommentFieldNameConst  = "comment_id"
	UserFieldNameConst     = "user_id"
	LikeTypeFieldNameConst = "like_type"
)

const (
	MaxTextLengthConst = 800
	MaxLocationLength  = 255
)

type LikeRequest struct {
	uuid        port.UuidInterface
	LikeType    string  `form:"like_type" json:"like_type"`
	PostId      *string `form:"post_id" json:"post_id"`
	CommentId   *string `form:"comment_id" json:"comment_id"`
	UserId      string
	Like        *likeDomain.Like
	customError port.CustomErrorInterface
}

func NewLikeRequest(
	context *gin.Context,
	uuid port.UuidInterface,
	userId string,
	customError port.CustomErrorInterface,
) (*LikeRequest, error) {
	likeRequest := &LikeRequest{
		uuid:        uuid,
		UserId:      userId,
		customError: customError,
	}

	if err := context.ShouldBind(likeRequest); err != nil {
		return nil, err
	}

	return likeRequest, nil
}

func (likeRequest *LikeRequest) Validate(context *gin.Context) error {
	if err := requestEntityFieldsValidation.ValidateFieldInArray(
		likeRequest.LikeType,
		LikeTypeFieldConst,
		likeDomain.LikeAcceptedTypes,
	); err != nil {
		return err
	}

	if likeRequest.PostId == nil && likeRequest.CommentId == nil {
		return likeRequest.customError.ThrowError(helper.InformFieldConst, PostOrCommentFieldConst)
	}

	if likeRequest.PostId == nil {
		if likeRequest.LikeType != likeDomain.CommentTypeConst {
			return likeRequest.customError.ThrowError(helper.InvalidLikeTypeConst)
		}

		if err := requestEntityFieldsValidation.IsValidUUID(
			CommentFieldConst,
			*likeRequest.CommentId,
		); err != nil {
			return err
		}
	}

	if likeRequest.CommentId == nil {
		if likeRequest.LikeType != likeDomain.PostTypeConst {
			return likeRequest.customError.ThrowError(helper.InvalidLikeTypeConst)
		}

		if err := requestEntityFieldsValidation.IsValidUUID(
			PostFieldConst,
			*likeRequest.PostId,
		); err != nil {
			return err
		}
	}

	if likeRequest.CommentId != nil && likeRequest.PostId != nil {
		return likeRequest.customError.ThrowError(helper.InvalidLikeRequestConst)
	}

	return nil
}

func (likeRequest *LikeRequest) BuildLikeObject() (*likeDomain.Like, error) {
	dateTime, err := datetimeDomain.CreateNow()
	if err != nil {
		return nil, err
	}

	return likeDomain.NewLike(
		likeRequest.uuid.GenerateUuid(),
		likeRequest.LikeType,
		likeRequest.PostId,
		likeRequest.CommentId,
		likeRequest.UserId,
		&dateTime,
	), nil
}
