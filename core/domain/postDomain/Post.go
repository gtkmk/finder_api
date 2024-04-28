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

const (
	RewardOptionTrueConst  = "1"
	RewardOptionFalseConst = "0"
)

const (
	AnimalTypeDogConst   = "cachorro"
	AnimalTypeCatConst   = "gato"
	AnimalTypeBirdConst  = "ave"
	AnimalTypeOtherConst = "outro"
)

const (
	AnimalSizeSmallConst  = "pequeno"
	AnimalSizeMediumConst = "médio"
	AnimalSizeBigConst    = "grande"
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

var AcceptedAnimalTypes = []string{
	AnimalTypeDogConst,
	AnimalTypeCatConst,
	AnimalTypeBirdConst,
	AnimalTypeOtherConst,
}

var AcceptedAnimalSizes = []string{
	AnimalSizeSmallConst,
	AnimalSizeMediumConst,
	AnimalSizeBigConst,
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
	LostFound   *string                  `json:"lost_found"`
	AnimalType  *string                  `json:"animal_type"`
	AnimalSize  *string                  `json:"animal_size"`
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
	lostFound *string,
	animalType *string,
	animalSize *string,
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
		AnimalType:  animalType,
		AnimalSize:  animalSize,
		UserId:      userID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
