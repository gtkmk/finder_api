package repositories

import "github.com/gtkmk/finder_api/core/domain/likeDomain"

type LikeRepository interface {
	ConfirmExistingLike(likeInfo *likeDomain.Like) (bool, *likeDomain.Like, error)
	CreateLike(likeInfo *likeDomain.Like) error
	RemoveLike(likeId string) error
	FindCurrentLikesCount(likeInfo *likeDomain.Like) (int, error)
}
