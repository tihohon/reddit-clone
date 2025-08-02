package vote

type Vote struct {
	User   string `json:"user"`
	Vote   int    `json:"vote"`
	PostId string
}

type VoteRepo interface {
	GetPostVotes(postId string) ([]Vote, error)
}
