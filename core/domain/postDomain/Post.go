package postDomain

type Post struct {
	ID             string   `json:"id"`
	Text           string   `json:"text"`
	Author         string   `json:"author"`
	CreatedAt      string   `json:"created_at"`
	LastModifiedAt string   `json:"last_modified_at"`
	Category       string   `json:"category"`
	URL            string   `json:"url"`
	Visibility     string   `json:"visibility"`
	Comments       []string `json:"comments"`
	Likes          int      `json:"likes"`
	Media          []string `json:"media"`
	Location       string   `json:"location"`
	Status         string   `json:"status"`
	Reward         int      `json:"reward"`
}

func NewPost(
	id string,
	text string,
	author string,
	createdAt string,
	lastModifiedAt string,
	category string,
	url string,
	visibility string,
	comments []string,
	likes int,
	media []string,
	location string,
	status string,
	reward int,
) *Post {
	return &Post{
		ID:             id,
		Text:           text,
		Author:         author,
		CreatedAt:      createdAt,
		LastModifiedAt: lastModifiedAt,
		Category:       category,
		URL:            url,
		Visibility:     visibility,
		Comments:       comments,
		Likes:          likes,
		Media:          media,
		Location:       location,
		Status:         status,
		Reward:         reward,
	}
}
