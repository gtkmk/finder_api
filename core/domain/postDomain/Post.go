package postDomain

import (
	"time"

	"github.com/gtkmk/finder_api/core/domain/documentDomain"
)

const (
	FoundConst = "found"
	LostConst  = "lost"
)

const (
	PrivacyPublicConst      = "public"
	PrivacyPrivateConst     = "private"
	PrivacyFriendsOnlyConst = "friends_only"
)

const (
	CategoryDefaultConst = "default"
	CategoryPaidConst    = "paid"
	CategoryAddConst     = "add"
)

var LostAndFoundStatus = []string{
	FoundConst,
	LostConst,
}

var AcceptedPrivacySettings = []string{
	PrivacyPublicConst,
	PrivacyPrivateConst,
	PrivacyFriendsOnlyConst,
}

var AcceptedCategories = []string{
	CategoryDefaultConst,
	CategoryPaidConst,
	CategoryAddConst,
}

type Post struct {
	Id          string                   `json:"id"`
	Text        string                   `json:"text"`
	Media       *documentDomain.Document `json:"media"`
	Location    string                   `json:"location"`
	Reward      bool                     `json:"reward"`
	Privacy     string                   `json:"privacy"`
	SharesCount int                      `json:"shares_count"`
	Category    string                   `json:"category"`
	LostFound   string                   `json:"lost_found"`
	UserId      string                   `json:"user_id"`
	CreatedAt   *time.Time               `json:"created_at"`
	UpdatedAt   *time.Time               `json:"updated_at"`
}

func NewPost(
	id string,
	text string,
	media *documentDomain.Document,
	location string,
	reward bool,
	privacy string,
	sharesCount int,
	category string,
	lostFound string,
	userID string,
	createdAt *time.Time,
	updatedAt *time.Time,
) *Post {
	return &Post{
		Id:          id,
		Text:        text,
		Media:       media,
		Location:    location,
		Reward:      reward,
		Privacy:     privacy,
		SharesCount: sharesCount,
		Category:    category,
		LostFound:   lostFound,
		UserId:      userID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
