package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"redditclone/pkg/helpers"
	"redditclone/pkg/models/comment"
	"redditclone/pkg/models/post"
	"redditclone/pkg/models/session"
	"redditclone/pkg/models/user"
	"redditclone/pkg/models/vote"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type PostHandler struct {
	Logger       *zap.SugaredLogger
	PostRepo     post.PostRepo
	UserRepo     user.UserRepo
	CommentsRepo comment.CommentRepo
	VotesRepo    vote.VoteRepo
}

type CommentWithUser struct {
	comment.Comment
	Author user.User `json:"author"`
}

type PostWithCommentsAndVotes struct {
	post.Post
	Votes            []vote.Vote       `json:"votes"`
	Comments         []CommentWithUser `json:"comments"`
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
	if positive == 0 {
		return 0
	}

	return int(float64(positive) / float64(positive+negative) * 100)
}

func (h *PostHandler) getCommentWithUser(comment comment.Comment) (*CommentWithUser, error) {
	commentWithUser := CommentWithUser{Comment: comment}
	user, err := h.UserRepo.GetUserById(comment.UserId)
	if err != nil {
		return nil, fmt.Errorf("could not find user %e", err)
	}
	commentWithUser.Author = *user
	return &commentWithUser, nil
}

func (h *PostHandler) getPostWithCommentsAndVotes(post post.Post) (*PostWithCommentsAndVotes, error) {
	postWithInfo := PostWithCommentsAndVotes{}
	postWithInfo.Post = post

	comments, commentsErr := h.CommentsRepo.GetCommentsForPost(post.Id)
	if commentsErr != nil {
		return nil, commentsErr
	}
	commentsWithUser := []CommentWithUser{}
	for _, comment := range comments {
		commentWithUser, err := h.getCommentWithUser(comment)
		if err != nil {
			return nil, err
		}
		commentsWithUser = append(commentsWithUser, *commentWithUser)
	}

	postWithInfo.Comments = commentsWithUser

	user, userErr := h.UserRepo.GetUserById(post.UserId)
	if userErr != nil {
		return nil, userErr
	}
	postWithInfo.Author = *user

	votes, votesErr := h.VotesRepo.GetPostVotes(post.Id)
	if votesErr != nil {
		return nil, votesErr
	}
	postWithInfo.Votes = votes
	postWithInfo.Score = calculatePostScore(&votes)
	postWithInfo.UpvotePercentage = calculatePostUpvotePercentage(&votes)
	return &postWithInfo, nil

}

func (h *PostHandler) getPostWithCommentsAndVotesById(postId string) (*PostWithCommentsAndVotes, error) {
	posts, err := h.PostRepo.GetPosts(func(p *post.Post) (bool, error) {
		if p.Id == postId {
			return true, nil
		}
		return false, nil
	})
	if err != nil || len(posts) != 1 {
		h.Logger.Error(err)
		return nil, err
	}
	postWithCommentsAndVotes, err := h.getPostWithCommentsAndVotes(posts[0])
	if err != nil {
		return nil, err
	}
	return postWithCommentsAndVotes, nil

}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	category := vars["category"]
	filter := func(p *post.Post) (bool, error) {
		if p.Category == category {
			return true, nil
		}
		if category == "" {
			return true, nil
		}

		return false, nil
	}
	posts, postErr := h.PostRepo.GetPosts(filter)
	if postErr != nil {
		h.Logger.Error(postErr)
		helpers.WriteBadRequest(w, "Error on fetching posts")
		return
	}

	response := []PostWithCommentsAndVotes{}
	for _, post := range posts {
		postWithCommentsAndVotes, err := h.getPostWithCommentsAndVotes(post)
		if err != nil {
			h.Logger.Error(err)
			helpers.WriteBadRequest(w, "Error on building post")
			return
		}
		response = append(response, *postWithCommentsAndVotes)
	}

	bin, _ := json.Marshal(response)
	w.Write(bin)
}

type CreatePostRequest struct {
	Category string `json:"category"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	Text     string `json:"text"`
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var request CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helpers.WriteBadRequest(w, "decode request error")
		h.Logger.Info("decode request error")
		return
	}
	ctx := r.Context()
	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		helpers.WriteBadRequest(w, "auth error")
		h.Logger.Error(err)
		return
	}

	userId := sess.UserId
	post, err := h.PostRepo.CreatePost(request.Type, request.Title, request.Category, request.Text, request.Url, userId)
	if err != nil {
		helpers.WriteBadRequest(w, "failed to crete post")
		h.Logger.Error(err)
		return
	}
	_, err = h.VotesRepo.CreateVote(post.Id, userId, 1)
	if err != nil {
		helpers.WriteBadRequest(w, "failed to crete vote")
		h.Logger.Error(err)
		return
	}

	postWithCommentsAndVotes, err := h.getPostWithCommentsAndVotes(*post)
	if err != nil {
		helpers.WriteBadRequest(w, "failed to crete post")
		h.Logger.Error(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	bin, _ := json.Marshal(postWithCommentsAndVotes)
	w.Write(bin)
}

func (h *PostHandler) GetPostInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId := vars["id"]
	postWithCommentsAndVotes, err := h.getPostWithCommentsAndVotesById(postId)
	if err != nil {
		helpers.WriteBadRequest(w, "post not found")
		return
	}

	w.WriteHeader(http.StatusCreated)
	bin, _ := json.Marshal(postWithCommentsAndVotes)
	w.Write(bin)
}

type CreateCommentRequest struct {
	Comment string `json:"comment"`
}

func (h *PostHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId := vars["id"]

	var request CreateCommentRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helpers.WriteBadRequest(w, "decode request error")
		h.Logger.Info("decode request error")
		return
	}
	ctx := r.Context()
	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		helpers.WriteBadRequest(w, "auth error")
		h.Logger.Error(err)
		return
	}
	h.CommentsRepo.CreateComment(postId, sess.UserId, request.Comment)

	postWithCommentsAndVotes, err := h.getPostWithCommentsAndVotesById(postId)
	if err != nil {
		helpers.WriteBadRequest(w, "post not found")
		return
	}

	w.WriteHeader(http.StatusCreated)
	bin, _ := json.Marshal(postWithCommentsAndVotes)
	w.Write(bin)
}

func (h *PostHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId := vars["id"]
	commentId := vars["commentId"]

	ctx := r.Context()
	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		helpers.WriteBadRequest(w, "auth error")
		h.Logger.Error(err)
		return
	}
	comment, err := h.CommentsRepo.GetComment(commentId)
	if err != nil || comment.PostId != postId || comment.UserId != sess.UserId {
		helpers.WriteBadRequest(w, "comment not found")
		h.Logger.Error(err)
		return
	}
	h.CommentsRepo.DeleteComment(commentId)

	postWithCommentsAndVotes, err := h.getPostWithCommentsAndVotesById(postId)
	if err != nil {
		helpers.WriteBadRequest(w, "post not found")
		return
	}

	w.WriteHeader(http.StatusCreated)
	bin, _ := json.Marshal(postWithCommentsAndVotes)
	w.Write(bin)
}

func (h *PostHandler) UpVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId := vars["id"]
	ctx := r.Context()
	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		helpers.WriteBadRequest(w, "auth error")
		h.Logger.Error(err)
		return
	}
	h.VotesRepo.CreateVote(postId, sess.UserId, 1)
	postWithCommentsAndVotes, err := h.getPostWithCommentsAndVotesById(postId)
	if err != nil {
		helpers.WriteBadRequest(w, "post not found")
		return
	}

	w.WriteHeader(http.StatusCreated)
	bin, _ := json.Marshal(postWithCommentsAndVotes)
	w.Write(bin)
}

func (h *PostHandler) DownVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId := vars["id"]
	ctx := r.Context()
	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		helpers.WriteBadRequest(w, "auth error")
		h.Logger.Error(err)
		return
	}
	h.VotesRepo.CreateVote(postId, sess.UserId, -1)
	postWithCommentsAndVotes, err := h.getPostWithCommentsAndVotesById(postId)
	if err != nil {
		helpers.WriteBadRequest(w, "post not found")
		return
	}

	w.WriteHeader(http.StatusCreated)
	bin, _ := json.Marshal(postWithCommentsAndVotes)
	w.Write(bin)
}

func (h *PostHandler) WithdrawVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId := vars["id"]
	ctx := r.Context()
	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		helpers.WriteBadRequest(w, "auth error")
		h.Logger.Error(err)
		return
	}
	h.VotesRepo.WithdrawVote(postId, sess.UserId)
	postWithCommentsAndVotes, err := h.getPostWithCommentsAndVotesById(postId)
	if err != nil {
		helpers.WriteBadRequest(w, "post not found")
		return
	}

	w.WriteHeader(http.StatusCreated)
	bin, _ := json.Marshal(postWithCommentsAndVotes)
	w.Write(bin)
}

type DeletePostResponse struct {
	Message string `json:"comment"`
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId := vars["id"]
	ctx := r.Context()
	sess, err := session.SessionFromContext(ctx)
	if err != nil {
		helpers.WriteBadRequest(w, "auth error")
		h.Logger.Error(err)
		return
	}
	post, err := h.getPostWithCommentsAndVotesById(postId)
	if err != nil || post.UserId != sess.UserId {
		helpers.WriteBadRequest(w, "post not found")
		h.Logger.Error(err)
		return
	}
	err = h.PostRepo.DeletePost(postId)
	if err != nil {
		helpers.WriteBadRequest(w, "post not found")
		return
	}

	bin, _ := json.Marshal(DeletePostResponse{Message: "success"})
	w.Write(bin)
}

func (h *PostHandler) GetUserPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login := vars["login"]

	users, err := h.UserRepo.GetUsers(func(u *user.User) (bool, error) {
		if u.Username == login {
			return true, nil
		}
		return false, nil
	})
	if err != nil || len(users) != 1 {
		helpers.WriteBadRequest(w, "User not found")
		h.Logger.Error(err)
		return
	}

	filter := func(p *post.Post) (bool, error) {
		if p.UserId == users[0].Id {
			return true, nil
		}
		return false, nil
	}
	posts, postErr := h.PostRepo.GetPosts(filter)

	if postErr != nil {
		helpers.WriteBadRequest(w, "Failed to fetch users posts")
		h.Logger.Error(err)
		return
	}

	response := []PostWithCommentsAndVotes{}
	for _, post := range posts {
		postWithCommentsAndVotes, err := h.getPostWithCommentsAndVotes(post)
		if err != nil {
			h.Logger.Error(err)
			helpers.WriteBadRequest(w, "Error on building post")
			return
		}
		response = append(response, *postWithCommentsAndVotes)
	}

	bin, _ := json.Marshal(response)
	w.Write(bin)
}
