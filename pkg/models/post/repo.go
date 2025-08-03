package post

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"redditclone/pkg/models"
	"sync"
	"time"
)

type PostMemory struct {
	posts map[string]Post
	mutex *sync.RWMutex
}

func (memory *PostMemory) GetPosts(filterFunc func(*Post) (bool, error)) ([]Post, error) {
	result := []Post{}

	memory.mutex.RLock()
	for _, post := range memory.posts {
		add, err := filterFunc(&post)
		if err != nil {
			return nil, fmt.Errorf("failed on filter posts: %e", err)
		}
		if add {
			result = append(result, post)
		}

	}
	memory.mutex.RUnlock()
	return result, nil

}

func (memory *PostMemory) CreatePost(postType string, title string, category string, text string, url string, userId string) (*Post, error) {
	existedTypes := map[string]struct{}{"text": {}, "link": {}}
	if _, ok := existedTypes[postType]; !ok {
		return nil, &models.InvalidValueError{}
	}
	randId := make([]byte, 12)
	rand.Read(randId)
	post := Post{
		Id:       hex.EncodeToString(randId),
		Views:    0,
		PostType: postType,
		Title:    title,
		Category: category,
		Text:     text,
		Url:      url,
		Created:  time.Now().Format(time.RFC3339),
		UserId:   userId,
	}

	memory.mutex.Lock()
	memory.posts[post.Id] = post
	memory.mutex.Unlock()
	return &post, nil
}

func (memory *PostMemory) DeletePost(postId string) error {
	_, ok := memory.posts[postId]
	if !ok {
		return &models.NotFoundError{}
	}
	memory.mutex.Lock()
	delete(memory.posts, postId)
	memory.mutex.Unlock()
	return nil

}

func NewPostMemory() *PostMemory {
	return &PostMemory{posts: map[string]Post{}, mutex: &sync.RWMutex{}}
}
