package session

import (
	"context"
	"redditclone/pkg/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Session struct {
	Id     string
	UserId string
}

type sessKey string

var SessionKey sessKey = "sessionKey"

func NewSession(userId string, username string, key string) (*Session, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"id":       userId,
			"username": username,
		},
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return nil, &models.SignError{}
	}

	return &Session{
		Id:     tokenString,
		UserId: userId,
	}, nil
}

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, &models.NoAuthError{}
	}
	return sess, nil
}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}
