package vote

type Vote struct {
	User   string `json:"user"`
	Vote   int    `json:"vote"`
	PostId string `json:"-"`
	VoteId string `json:"-"`
}

type VoteRepo interface {
	GetPostVotes(postId string) ([]Vote, error)
	CreateVote(postId string, userId string, voteVal int) (*Vote, error)
	WithdrawVote(postId string, userId string) error
}
