package vote

type VoteMemory struct {
	Votes map[string]Vote
}

func (memory *VoteMemory) GetPostVotes(postId string) ([]Vote, error) {
	response := []Vote{}
	for _, vote := range memory.Votes {
		if postId == vote.PostId {
			response = append(response, vote)
		}
	}
	return response, nil

}

func NewVoteMemory() *VoteMemory {
	return &VoteMemory{Votes: map[string]Vote{}}
}
