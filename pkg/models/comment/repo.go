package comment

type CommentMemory struct {
	Comments map[string]Comment
}

func (memory *CommentMemory) GetCommentsForPost(postId string) ([]Comment, error) {
	response := []Comment{}
	for _, comment := range memory.Comments {
		if comment.postId == postId {
			response = append(response, comment)
		}
	}

	return response, nil

}

func NewCommentMemory() *CommentMemory {
	return &CommentMemory{Comments: map[string]Comment{}}
}
