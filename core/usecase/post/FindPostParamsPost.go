package postUsecase

import (
	"github.com/gtkmk/finder_api/core/domain/postDomain"
)

type FindPostPostParams struct{}

func NewFindPostPostParams() *FindPostPostParams {
	return &FindPostPostParams{}
}

func (findPostPostParams *FindPostPostParams) Execute() map[string]interface{} {
	translatedLostAndFoundStatus := map[string]interface{}{
		postDomain.TranslatedFoundConst: postDomain.FoundConst,
		postDomain.TranslatedLostConst:  postDomain.LostConst,
	}

	// TODO: Return with this when this functionality is to be implemented
	// translatedAcceptedPrivacySettings := map[string]interface{}{
	// 	postDomain.TranslatedPrivacyPublicConst:      postDomain.PrivacyPublicConst,
	// 	postDomain.TranslatedPrivacyPrivateConst:     postDomain.PrivacyPrivateConst,
	// 	postDomain.TranslatedPrivacyFriendsOnlyConst: postDomain.PrivacyFriendsOnlyConst,
	// }

	translatedAcceptedCategories := map[string]interface{}{
		postDomain.TranslatedCategoryDefaultConst: postDomain.CategoryDefaultConst,
		postDomain.TranslatedCategoryPaidConst:    postDomain.CategoryPaidConst,
		postDomain.TranslatedCategoryAddConst:     postDomain.CategoryAddConst,
	}

	translatedAcceptedAnimalTypes := map[string]interface{}{
		postDomain.TranslatedAnimalTypeDogConst:   postDomain.AnimalTypeDogConst,
		postDomain.TranslatedAnimalTypeCatConst:   postDomain.AnimalTypeCatConst,
		postDomain.TranslatedAnimalTypeBirdConst:  postDomain.AnimalTypeBirdConst,
		postDomain.TranslatedAnimalTypeOtherConst: postDomain.AnimalTypeOtherConst,
	}

	translatedAcceptedAnimalSizes := map[string]interface{}{
		postDomain.TranslatedAnimalSizeSmallConst:  postDomain.AnimalSizeSmallConst,
		postDomain.TranslatedAnimalSizeMediumConst: postDomain.AnimalSizeMediumConst,
		postDomain.TranslatedAnimalSizeBigConst:    postDomain.AnimalSizeBigConst,
	}

	return map[string]interface{}{
		"LostAndFoundStatus":  translatedLostAndFoundStatus,
		"AcceptedCategories":  translatedAcceptedCategories,
		"AcceptedAnimalTypes": translatedAcceptedAnimalTypes,
		"AcceptedAnimalSizes": translatedAcceptedAnimalSizes,
	}
}
