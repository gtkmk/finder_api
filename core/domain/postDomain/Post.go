package postDomain

type Post struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Content        string   `json:"content"`
	Author         string   `json:"author"`
	CreatedAt      string   `json:"created_at"`
	LastModifiedAt string   `json:"last_modified_at"`
	Category       string   `json:"category"`
	Tags           []string `json:"tags"`
	URL            string   `json:"url"`
	Visibility     bool     `json:"visibility"`
	Comments       []string `json:"comments"`
	Likes          int      `json:"likes"`
	Media          []string `json:"media"`
	Location       string   `json:"location"`
	Status         string   `json:"status"`
	Missing        bool     `json:"missing"`
	Reward         bool     `json:"reward"`
}

func NewPost(
	id string,
	title string,
	content string,
	author string,
	createdAt string,
	lastModifiedAt string,
	category string,
	tags []string,
	url string,
	visibility bool,
	comments []string,
	likes int,
	media []string,
	location string,
	status string,
	reward bool,
) *Post {
	return &Post{
		ID:             id,
		Title:          title,
		Content:        content,
		Author:         author,
		CreatedAt:      createdAt,
		LastModifiedAt: lastModifiedAt,
		Category:       category,
		Tags:           tags,
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
