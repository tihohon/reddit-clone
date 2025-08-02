package handlers

import (
	"encoding/json"
	"net/http"
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
		h.Logger.Error("decode request error")
		return // (400)
	}

	userId, err := h.UserRepo.CreateUser(request.Username, request.Password)
	if err != nil {
		h.Logger.Error("error on user creation %e", err)
	}

	session, err := h.SessionsManager.Create(userId, request.Username)
	if err != nil {
		h.Logger.Error("failed to create session: %e", err)
	}

	response := RegisterResponse{
		Token: session.Id,
	}

	bin, _ := json.Marshal(response)
	w.Write(bin)
}
