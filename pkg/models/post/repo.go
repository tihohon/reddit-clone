package post

import "fmt"

type PostMemory struct {
	posts map[string]Post
}

func (memory *PostMemory) GetPosts(filterFunc func(*Post) (bool, error)) ([]Post, error) {
	result := []Post{}
	for _, post := range memory.posts {
		add, err := filterFunc(&post)
		if err != nil {
			return nil, fmt.Errorf("failed on filter posts: %w", err)
		}
		if add {
			result = append(result, post)
		}

	}
	return result, nil

}

func NewPostMemory() *PostMemory {
	return &PostMemory{posts: map[string]Post{}}
}
