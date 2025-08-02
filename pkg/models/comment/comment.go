package comment

type Comment struct {
	Created string `json:"created"`
	Author  string `json:"author"`
	Text    string `json:"body"`
	postId  string
}

type CommentRepo interface {
	GetCommentsForPost(postId string) ([]Comment, error)
}
