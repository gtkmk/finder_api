package postRequestEntity

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gtkmk/finder_api/core/domain/filterDomain"
	"github.com/gtkmk/finder_api/core/domain/helper"
	"github.com/gtkmk/finder_api/core/domain/postDomain"
	"github.com/gtkmk/finder_api/core/port"
	"github.com/gtkmk/finder_api/core/port/sharedMethods"
	"github.com/gtkmk/finder_api/infra/requestEntityFieldsValidation"
)

type ListPostsRequest struct {
	uuid                 port.UuidInterface
	LoggedUserId         string
	Page                 *int64  `form:"page"`
	Neighborhood         *string `form:"neighborhood"`
	LostFound            *string `form:"lostFound"`
	Reward               *string `form:"reward"`
	UserId               *string `form:"user_id"`
	OnlyFollowingPosts   *string `form:"only_following_posts"`
	SpecificPost         *string `form:"specific_post"`
	AnimalType           *string `form:"animal_type"`
	AnimalSize           *string `form:"animal_size"`
	checkForSqlInjection sharedMethods.CheckForSqlInjectionInterface
}

const (
	TrueValueConst  = "1"
	FalseValueConst = "0"
)

const (
	RequestRewardFieldNameConst             = "recompensa"
	RequestOnlyFollowingPostsFieldNameConst = "apenas posts de conhecidos"
	RequestOrdenationFieldNameConst         = "ordenação por campo"
	RequestOrdenationTypeFieldNameConst     = "ordenação por tipo"
)

const (
	UserIdFieldConst       = "o usuário"
	SpecificPostFieldConst = "o post especifico"
)

func NewListPostsRequest(
	context *gin.Context,
	uuid port.UuidInterface,
	checkForSqlInjection sharedMethods.CheckForSqlInjectionInterface,
	loggedUserId string,
) (*ListPostsRequest, error) {
	listPostsRequest := &ListPostsRequest{
		checkForSqlInjection: checkForSqlInjection,
		LoggedUserId:         loggedUserId,
	}
	if err := context.ShouldBind(listPostsRequest); err != nil {
		return nil, err
	}

	return listPostsRequest, nil
}

func (listPostsRequest *ListPostsRequest) ValidatePostsFilterFields(context *gin.Context) error {
	if err := listPostsRequest.verifyIfPageIsValid(listPostsRequest.Page); err != nil {
		return err
	}

	if listPostsRequest.LostFound != nil {
		if err := listPostsRequest.verifyIfLostFoundIsValid(*listPostsRequest.LostFound); err != nil {
			return err
		}
	}

	if listPostsRequest.Neighborhood != nil {
		if err := listPostsRequest.checkForSqlInjection.CheckForSqlInjection(*listPostsRequest.Neighborhood); err != nil {
			return err
		}
	}

	if listPostsRequest.Reward != nil {
		if err := listPostsRequest.verifyIfRewardIsValid(listPostsRequest.Reward); err != nil {
			return err
		}
	}

	if listPostsRequest.OnlyFollowingPosts != nil {
		if err := listPostsRequest.verifyIfOnlyFollowingPostsIsValid(listPostsRequest.OnlyFollowingPosts); err != nil {
			return err
		}
	}

	if listPostsRequest.UserId != nil {
		if err := requestEntityFieldsValidation.IsValidUUID(UserIdFieldConst, *listPostsRequest.UserId); err != nil {
			return err
		}
	}

	if listPostsRequest.SpecificPost != nil {
		if err := requestEntityFieldsValidation.IsValidUUID(SpecificPostFieldConst, *listPostsRequest.SpecificPost); err != nil {
			return err
		}
	}

	if listPostsRequest.AnimalType != nil {
		if err := listPostsRequest.verifyIfAnimalTypeIsValid(*listPostsRequest.AnimalType); err != nil {
			return err
		}
	}

	if listPostsRequest.AnimalSize != nil {
		if err := listPostsRequest.verifyIfAnimalSizeIsValid(*listPostsRequest.AnimalSize); err != nil {
			return err
		}
	}

	return nil
}

func (listPostsRequest *ListPostsRequest) verifyIfPageIsValid(page *int64) error {
	if page == nil || *page <= 0 {
		return fmt.Errorf(
			helper.FieldIsMandatoryAndMustToBeGreaterThanZeroConst,
			"página",
		)
	}
	return nil
}

func (listPostsRequest *ListPostsRequest) verifyIfLostFoundIsValid(status string) error {
	validOptions := map[string]struct{}{
		postDomain.FoundConst: {},
		postDomain.LostConst:  {},
	}

	if _, ok := validOptions[status]; !ok {
		return fmt.Errorf(helper.PostLostFoundStatusNotRecognizedConst)
	}

	return nil
}

func (listPostsRequest *ListPostsRequest) verifyIfRewardIsValid(reward *string) error {
	switch *reward {
	case postDomain.RewardOptionTrueConst,
		postDomain.RewardOptionFalseConst:
		return nil
	default:
		return fmt.Errorf(helper.PostRewardNotRecognizedMessageConst)
	}
}

func (listPostsRequest *ListPostsRequest) verifyIfOnlyFollowingPostsIsValid(onlyFollowingPosts *string) error {
	switch *onlyFollowingPosts {
	case TrueValueConst,
		FalseValueConst:
		return nil
	default:
		return fmt.Errorf(helper.OptionNotRecognizedMessageConst, RequestOnlyFollowingPostsFieldNameConst)
	}
}

func (listPostsRequest *ListPostsRequest) verifyIfAnimalTypeIsValid(anmalType string) error {
	validOptions := map[string]struct{}{
		postDomain.AnimalTypeDogConst:   {},
		postDomain.AnimalTypeCatConst:   {},
		postDomain.AnimalTypeBirdConst:  {},
		postDomain.AnimalTypeOtherConst: {},
	}

	if _, ok := validOptions[anmalType]; !ok {
		return fmt.Errorf(helper.PostAnimalTypeNotRecognizedConst)
	}

	return nil
}

func (listPostsRequest *ListPostsRequest) verifyIfAnimalSizeIsValid(animalSize string) error {
	validOptions := map[string]struct{}{
		postDomain.AnimalSizeSmallConst:  {},
		postDomain.AnimalSizeMediumConst: {},
		postDomain.AnimalSizeBigConst:    {},
	}

	if _, ok := validOptions[animalSize]; !ok {
		return fmt.Errorf(helper.PostAnimalSizeNotRecognizedConst)
	}

	return nil
}

func (listPostsRequest *ListPostsRequest) ConvertProposalFiltersIntoFilterDomain() *filterDomain.PostFilter {
	filters := filterDomain.NewPostFilter(
		listPostsRequest.Page,
		listPostsRequest.LoggedUserId,
		listPostsRequest.Neighborhood,
		listPostsRequest.LostFound,
		listPostsRequest.Reward,
		listPostsRequest.UserId,
		listPostsRequest.OnlyFollowingPosts,
		listPostsRequest.SpecificPost,
		listPostsRequest.AnimalType,
		listPostsRequest.AnimalSize,
		filterDomain.MaxItensPerPageConst,
		nil,
	)

	return filters
}
