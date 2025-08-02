package handlers

import (
	"encoding/json"
	"net/http"
	"redditclone/pkg/models/comment"
	"redditclone/pkg/models/post"
	"redditclone/pkg/models/user"
	"redditclone/pkg/models/vote"

	"go.uber.org/zap"
)

type PostHandler struct {
	Logger       *zap.SugaredLogger
	PostRepo     post.PostRepo
	UserRepo     user.UserRepo
	CommentsRepo comment.CommentRepo
	VotesRepo    vote.VoteRepo
}

type PostWithCommentsAndVotes struct {
	post.Post
	Votes            []vote.Vote       `json:"votes"`
	Comments         []comment.Comment `json:"comments"`
	Author           user.User         `json:"author"`
	UpvotePercentage int               `json:"upvotePercentage"`
	Score            int               `json:"score"`
}

func calculatePostScore(votes *[]vote.Vote) int {
	sum := 0
	for _, vote := range *votes {
		sum += vote.Vote
	}
	return sum
}

func calculatePostUpvotePercentage(votes *[]vote.Vote) int {
	positive := 0
	negative := 0
	for _, vote := range *votes {
		if vote.Vote > 0 {
			positive++
		}
		if vote.Vote < 0 {
			negative++
		}
	}

	return positive / (positive + negative)
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {

	filter := func(*post.Post) (bool, error) { return true, nil }
	posts, postErr := h.PostRepo.GetPosts(filter)
	if postErr != nil {
		h.Logger.Error(postErr)
		return
	}

	response := []PostWithCommentsAndVotes{}
	for _, post := range posts {
		postWithInfo := PostWithCommentsAndVotes{}
		postWithInfo.Post = post

		comments, commentsErr := h.CommentsRepo.GetCommentsForPost(post.Id)
		if commentsErr != nil {
			h.Logger.Error(commentsErr)
		}
		postWithInfo.Comments = comments

		user, userErr := h.UserRepo.GetUserById(post.UserId)
		if userErr != nil {
			h.Logger.Error(userErr)
		}
		postWithInfo.Author = user

		votes, votesErr := h.VotesRepo.GetPostVotes(post.UserId)
		if votesErr != nil {
			h.Logger.Error(votesErr)
		}
		postWithInfo.Votes = votes
		postWithInfo.Score = calculatePostScore(&votes)
		postWithInfo.UpvotePercentage = calculatePostUpvotePercentage(&votes)
	}

	bin, _ := json.Marshal(response)
	w.Write(bin)
}
