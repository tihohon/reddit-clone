package comment

type Comment struct {
	Created string `json:"created"`
	UserId  string `json:"-"`
	Text    string `json:"body"`
	PostId  string `json:"-"`
	Id      string `json:"id"`
}

type CommentRepo interface {
	GetCommentsForPost(postId string) ([]Comment, error)
	CreateComment(postId string, userId string, text string) error
	DeleteComment(commentId string) error
	GetComment(commentId string) (*Comment, error)
}
