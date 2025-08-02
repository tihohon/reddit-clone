package post

type Post struct {
	Id               string
	Views            int    `json:"views"`
	PostType         int    `json:"type"`
	Title            string `json:"title"`
	Category         string `json:"category"`
	Text             string `json:"text"`
	Created          string `json:"created"`
	UserId           string
}

type PostRepo interface {
	GetPosts(filterFunc func(*Post) (bool, error)) ([]Post, error)
}
