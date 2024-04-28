package followrequestentity

import (
	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/core/domain/datetimeDomain"
	"github.com/gtkmk/finder_api/core/domain/followDomain"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

const (
	FollowerIdFieldConst = "o usuário a ser seguido"
)

type FollowRequest struct {
	uuid        port.UuidInterface
	FollowerId  string `form:"follower_id" json:"follower_id"`
	UserId      string
	Follow      *followDomain.Follow
	CustomError port.CustomErrorInterface
}

func NewFollowRequest(
	context *gin.Context,
	uuid port.UuidInterface,
	userId string,
	customError port.CustomErrorInterface,
) (*FollowRequest, error) {
	followRequest := &FollowRequest{
		uuid:        uuid,
		UserId:      userId,
		CustomError: customError,
	}

	if err := context.ShouldBind(followRequest); err != nil {
		return nil, err
	}

	return followRequest, nil
}

func (followRequest *FollowRequest) Validate(context *gin.Context) error {
	if err := requestEntityFieldsValidation.IsValidUUID(
		FollowerIdFieldConst,
		followRequest.FollowerId,
	); err != nil {
		return err
	}

	if followRequest.FollowerId == followRequest.UserId {
		return followRequest.CustomError.ThrowError(
			helper.YouCannotFollowYourselfConst,
		)
	}

	return nil
}

func (followRequest *FollowRequest) BuildFollowObject() (*followDomain.Follow, error) {
	dateTime, err := datetimeDomain.CreateNow()
	if err != nil {
		return nil, err
	}

	return followDomain.NewFollow(
		followRequest.uuid.GenerateUuid(),
		followRequest.FollowerId,
		followRequest.UserId,
		&dateTime,
	), nil
}
