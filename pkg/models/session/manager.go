package session

import (
	"fmt"
	"net/http"
	"redditclone/pkg/models"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type SessionsManager struct {
	data      map[string]*Session
	mu        *sync.RWMutex
	Logger    *zap.SugaredLogger
	secretKey string
}

func NewSessionsManager(logger *zap.SugaredLogger) *SessionsManager {
	return &SessionsManager{
		data:      make(map[string]*Session, 10),
		mu:        &sync.RWMutex{},
		Logger:    logger,
		secretKey: "secret-key",
	}
}

func (sm *SessionsManager) Check(r *http.Request) (*Session, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, &models.NoAuthError{}
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(sm.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error on token parse: %e", err)
	}

	sm.mu.RLock()
	sess, ok := sm.data[tokenString]
	sm.mu.RUnlock()

	if !ok {
		return nil, &models.NoAuthError{}
	}

	return sess, nil
}

func (sm *SessionsManager) Create(userId string, username string) (*Session, error) {
	sess, err := NewSession(userId, username, sm.secretKey)

	if err != nil {
		return nil, fmt.Errorf("error on session create: %e", err)
	}

	sm.mu.Lock()
	sm.data[sess.Id] = sess
	sm.mu.Unlock()

	return sess, nil
}

func (sm *SessionsManager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	sess, err := SessionFromContext(r.Context())
	if err != nil {
		return err
	}

	sm.mu.Lock()
	delete(sm.data, sess.Id)
	sm.mu.Unlock()

	return nil
}
