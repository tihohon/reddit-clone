package vote

import (
	"fmt"
	"redditclone/pkg/models"
	"sync"
)

type VoteMemory struct {
	Votes map[string]Vote
	mutex *sync.RWMutex
}

func (memory *VoteMemory) GetPostVotes(postId string) ([]Vote, error) {
	response := []Vote{}

	memory.mutex.RLock()
	for _, vote := range memory.Votes {
		if postId == vote.PostId {
			response = append(response, vote)
		}
	}
	memory.mutex.RUnlock()

	return response, nil

}

func (memory *VoteMemory) CreateVote(postId string, userId string, voteVal int) (*Vote, error) {
	voteId := fmt.Sprintf("%v_%v", postId, userId)

	if voteVal > 1 || voteVal < -1 {
		return nil, &models.InvalidValueError{}
	}
	vote := Vote{
		PostId: postId,
		User:   userId,
		Vote:   voteVal,
		VoteId: voteId,
	}

	memory.mutex.Lock()
	memory.Votes[vote.VoteId] = vote
	memory.mutex.Unlock()

	return &vote, nil

}

func (memory *VoteMemory) WithdrawVote(postId string, userId string) error {
	voteId := fmt.Sprintf("%v_%v", postId, userId)

	memory.mutex.Lock()
	vote, ok := memory.Votes[voteId]
	if !ok || vote.User != userId {
		memory.mutex.Unlock()
		return &models.NotFoundError{}
	}
	delete(memory.Votes, voteId)
	memory.mutex.Unlock()

	return nil
}

func NewVoteMemory() *VoteMemory {
	return &VoteMemory{Votes: map[string]Vote{}, mutex: &sync.RWMutex{}}
}
