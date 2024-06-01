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
	TranslatedFoundConst = "achado"
	TranslatedLostConst  = "perdido"
)

const (
	PrivacyPublicConst      = "public"
	PrivacyPrivateConst     = "private"
	PrivacyFriendsOnlyConst = "friends_only"
)

const (
	TranslatedPrivacyPublicConst      = "público"
	TranslatedPrivacyPrivateConst     = "privado"
	TranslatedPrivacyFriendsOnlyConst = "apenas amigos"
)

const (
	CategoryDefaultConst = "default"
	CategoryPaidConst    = "paid"
	CategoryAddConst     = "add"
)

const (
	TranslatedCategoryDefaultConst = "padrão"
	TranslatedCategoryPaidConst    = "impulsionado"
	TranslatedCategoryAddConst     = "anúncio"
)

const (
	RewardOptionTrueConst  = "1"
	RewardOptionFalseConst = "0"
)

const (
	FoundOptionTrueConst  = "1"
	FoundOptionFalseConst = "0"
)

const (
	TranslatedRewardOptionTrueConst  = "possui recompensa"
	TranslatedRewardOptionFalseConst = "sem recompensa"
)

const (
	AnimalTypeDogConst   = "dog"
	AnimalTypeCatConst   = "cat"
	AnimalTypeBirdConst  = "bird"
	AnimalTypeOtherConst = "other"
)

const (
	TranslatedAnimalTypeDogConst   = "cachorro"
	TranslatedAnimalTypeCatConst   = "gato"
	TranslatedAnimalTypeBirdConst  = "ave"
	TranslatedAnimalTypeOtherConst = "outro"
)

const (
	AnimalSizeSmallConst  = "small"
	AnimalSizeMediumConst = "medium"
	AnimalSizeBigConst    = "large"
)

const (
	TranslatedAnimalSizeSmallConst  = "pequeno"
	TranslatedAnimalSizeMediumConst = "médio"
	TranslatedAnimalSizeBigConst    = "grande"
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
	Id                   string                   `json:"id"`
	Text                 string                   `json:"text"`
	Media                *documentDomain.Document `json:"media"`
	Location             string                   `json:"location"`
	Reward               bool                     `json:"reward"`
	Privacy              string                   `json:"privacy"`
	SharesCount          int                      `json:"shares_count"`
	Category             string                   `json:"category"`
	LostFound            *string                  `json:"lost_found"`
	AnimalType           *string                  `json:"animal_type"`
	AnimalSize           *string                  `json:"animal_size"`
	Found                bool                     `json:"reward"`
	UpdatedFoundStatusAt *time.Time               `json:"updated_found_status_at"`
	UserId               string                   `json:"user_id"`
	CreatedAt            *time.Time               `json:"created_at"`
	UpdatedAt            *time.Time               `json:"updated_at"`
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
	found bool,
	updatedFoundStatusAt *time.Time,
	userID string,
	createdAt *time.Time,
	updatedAt *time.Time,
) *Post {
	return &Post{
		Id:                   id,
		Text:                 text,
		Media:                media,
		Location:             location,
		Reward:               reward,
		Privacy:              privacy,
		SharesCount:          sharesCount,
		Category:             category,
		LostFound:            lostFound,
		AnimalType:           animalType,
		AnimalSize:           animalSize,
		Found:                found,
		UpdatedFoundStatusAt: updatedFoundStatusAt,
		UserId:               userID,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}
}
