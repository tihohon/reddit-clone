package post

type Post struct {
	Id       string `json:"id"`
	Views    int    `json:"views"`
	PostType string `json:"type"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Text     string `json:"text"`
	Url      string `json:"url"`
	Created  string `json:"created"`
	UserId   string `json:"-"`
}

type PostRepo interface {
	GetPosts(filterFunc func(*Post) (bool, error)) ([]Post, error)
	CreatePost(postType string, title string, category string, text string, url string, userId string) (*Post, error)
	DeletePost(postId string) error
}
