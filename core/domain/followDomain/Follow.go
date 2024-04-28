package followDomain

import "time"

type Follow struct {
	Id         string     `json:"id"`
	FollowerId string     `json:"follower_id"`
	FollowedId string     `json:"followed_id"`
	CreatedAt  *time.Time `json:"created_at"`
}

func NewFollow(
	id string,
	followerId string,
	followedId string,
	createdAt *time.Time,
) *Follow {
	return &Follow{
		Id:         id,
		FollowerId: followerId,
		FollowedId: followedId,
		CreatedAt:  createdAt,
	}
}
