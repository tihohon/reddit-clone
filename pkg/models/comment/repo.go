package comment

import (
	"crypto/rand"
	"encoding/hex"
	"redditclone/pkg/models"
	"sync"
	"time"
)

type CommentMemory struct {
	Comments map[string]Comment
	mutex    *sync.RWMutex
}

func (memory *CommentMemory) GetCommentsForPost(postId string) ([]Comment, error) {
	response := []Comment{}
	memory.mutex.RLock()
	for _, comment := range memory.Comments {
		if comment.PostId == postId {
			response = append(response, comment)
		}
	}
	memory.mutex.RUnlock()

	return response, nil

}

func (memory *CommentMemory) CreateComment(postId string, userId string, text string) error {
	randId := make([]byte, 12)
	rand.Read(randId)
	comment := Comment{
		PostId:  postId,
		UserId:  userId,
		Text:    text,
		Created: time.Now().Format(time.RFC3339),
		Id:      hex.EncodeToString(randId),
	}
	memory.mutex.Lock()
	memory.Comments[comment.Id] = comment
	memory.mutex.Unlock()
	return nil

}

func (memory *CommentMemory) DeleteComment(commentId string) error {
	memory.mutex.Lock()
	delete(memory.Comments, commentId)
	memory.mutex.Unlock()
	return nil
}

func (memory *CommentMemory) GetComment(commentId string) (*Comment, error) {
	memory.mutex.RLock()
	comment, ok := memory.Comments[commentId]
	memory.mutex.RUnlock()
	if !ok {
		return nil, &models.InvalidValueError{}
	}
	return &comment, nil

}

func NewCommentMemory() *CommentMemory {
	return &CommentMemory{Comments: map[string]Comment{}, mutex: &sync.RWMutex{}}
}
