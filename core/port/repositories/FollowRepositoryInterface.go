package repositories

import "github.com/gtkmk/finder_api/core/domain/followDomain"

type FollowRepository interface {
	ConfirmExistingFollow(followInfo *followDomain.Follow) (bool, *followDomain.Follow, error)
	CreateFollow(followInfo *followDomain.Follow) error
	RemoveFollow(followId string) error
}
