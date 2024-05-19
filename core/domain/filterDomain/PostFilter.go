package filterDomain

const (
	MaxItensPerPageConst = 150
)

type PostFilter struct {
	Page               *int64 `json:"page"`
	LoggedUserId       string
	Neighborhood       *string `json:"neighborhood"`
	LostFound          *string `json:"lostFound"`
	Reward             *string `json:"reward"`
	UserId             *string `json:"user_id"`
	OnlyFollowingPosts *string `json:"only_following_posts"`
	Limit              int64   `json:"limit"`
	OffSet             *int64  `json:"offset"`
}

func NewPostFilter(
	Page *int64,
	LoggedUserId string,
	Neighborhood *string,
	LostFound *string,
	Reward *string,
	UserId *string,
	OnlyFollowingPosts *string,
	Limit int64,
	OffSet *int64,
) *PostFilter {
	return &PostFilter{
		Page,
		LoggedUserId,
		Neighborhood,
		LostFound,
		Reward,
		UserId,
		OnlyFollowingPosts,
		Limit,
		OffSet,
	}
}
