package user

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"redditclone/pkg/models"
	"sync"
)

type UserMemory struct {
	Users map[string]User
	mutex *sync.RWMutex
}

func (memory *UserMemory) GetUserById(id string) (*User, error) {
	memory.mutex.RLock()
	user, ok := memory.Users[id]
	memory.mutex.RUnlock()
	if !ok {
		return nil, fmt.Errorf("error on retrieving user: %e", models.NotFoundError{})
	}
	return &user, nil

}

func (memory *UserMemory) CreateUser(username string, password string) (userId string, err error) {
	randId := make([]byte, 12)
	rand.Read(randId)

	//todo check exist

	newUser := User{
		Username: username,
		Password: password,
		Id:       hex.EncodeToString(randId),
	}
	memory.mutex.Lock()
	memory.Users[newUser.Id] = newUser
	memory.mutex.Unlock()
	return newUser.Id, nil

}

func (memory *UserMemory) GetUsers(filterFunc func(*User) (bool, error)) ([]User, error) {
	result := []User{}

	memory.mutex.RLock()
	for _, post := range memory.Users {
		add, err := filterFunc(&post)
		if err != nil {
			return nil, fmt.Errorf("failed on filter users: %e", err)
		}
		if add {
			result = append(result, post)
		}

	}
	memory.mutex.RUnlock()
	return result, nil

}

func NewUserMemory() *UserMemory {
	return &UserMemory{Users: map[string]User{}, mutex: &sync.RWMutex{}}
}
