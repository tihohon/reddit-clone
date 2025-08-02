package main

import (
	"html/template"
	"net/http"
	"redditclone/pkg/handlers"
	"redditclone/pkg/middleware"
	"redditclone/pkg/models/comment"
	"redditclone/pkg/models/post"
	"redditclone/pkg/models/session"
	"redditclone/pkg/models/user"
	"redditclone/pkg/models/vote"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	postMemory := post.NewPostMemory()
	userMemory := user.NewUserMemory()
	commentMemory := comment.NewCommentMemory()
	voteMemory := vote.NewVoteMemory()

	sm := session.NewSessionsManager(logger)

	postHandler := handlers.PostHandler{
		Logger:       logger,
		PostRepo:     postMemory,
		UserRepo:     userMemory,
		CommentsRepo: commentMemory,
		VotesRepo:    voteMemory,
	}
	registerHandler := handlers.RegisterHandler{
		Logger:          logger,
		UserRepo:        userMemory,
		SessionsManager: sm,
	}
	tmpl := template.Must(template.ParseFiles("./static/html/index.html"))

	r := mux.NewRouter()
	r.HandleFunc("/api/posts/", postHandler.GetPosts).Methods("GET")
	r.HandleFunc("/api/register", registerHandler.RegisterPost).Methods("POST")

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))

	// Обработка корневого пути
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})
	mux := middleware.Auth(sm, r)

	http.ListenAndServe(":3000", mux)
}
