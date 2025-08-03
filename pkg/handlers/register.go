package handlers

import (
	"encoding/json"
	"net/http"
	"redditclone/pkg/helpers"
	"redditclone/pkg/models/session"
	"redditclone/pkg/models/user"

	"go.uber.org/zap"
)

type RegisterHandler struct {
	Logger          *zap.SugaredLogger
	UserRepo        user.UserRepo
	SessionsManager *session.SessionsManager
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *RegisterHandler) RegisterPost(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helpers.WriteBadRequest(w, "decode request error")
		h.Logger.Info("decode request error")
		return
	}

	userId, err := h.UserRepo.CreateUser(request.Username, request.Password)
	if err != nil {
		h.Logger.Error("error on user creation %e", err)
		helpers.WriteBadRequest(w, "failed to create user")
		return
	}

	session, err := h.SessionsManager.Create(userId, request.Username)
	if err != nil {
		h.Logger.Error("failed to create session: %e", err)
		helpers.WriteBadRequest(w, "failed to create session")
		return
	}

	response := RegisterResponse{
		Token: session.Id,
	}

	bin, _ := json.Marshal(response)
	w.Write(bin)
}

func (h *RegisterHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helpers.WriteBadRequest(w, "decode request error")
		h.Logger.Info("decode request error")
		return
	}

	users, err := h.UserRepo.GetUsers(func(u *user.User) (bool, error) {
		if u.Password == request.Password && u.Username == u.Password {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		h.Logger.Error("error on finding user %e", err)
		helpers.WriteBadRequest(w, "failed to find user")
		return
	}
	if len(users) != 1 {
		h.Logger.Error("Cold not find user %e", err)
		helpers.WriteBadRequest(w, "failed to find user")
		return
	}
	user := users[0]

	session, err := h.SessionsManager.Create(user.Id, request.Username)
	if err != nil {
		h.Logger.Error("failed to create session: %e", err)
		helpers.WriteBadRequest(w, "failed to create session")
		return
	}

	response := RegisterResponse{
		Token: session.Id,
	}

	bin, _ := json.Marshal(response)
	w.Write(bin)
}
