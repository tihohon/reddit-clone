package user

import (
	"crypto/rand"
	"fmt"
	"redditclone/pkg/models"
)

type UserMemory struct {
	Users map[string]User
}

func (memory *UserMemory) GetUserById(id string) (User, error) {
	user, ok := memory.Users[id]
	if !ok {
		return User{}, fmt.Errorf("error on retrieving user: %e", models.NotFoundError{})
	}
	return user, nil

}

func (memory *UserMemory) CreateUser(username string, password string) (userId string, err error) {
	randId := make([]byte, 16)
	rand.Read(randId)

	//todo check exist

	newUser := User{
		Username: username,
		Password: password,
		Id:       string(randId),
	}
	memory.Users[string(randId)] = newUser
	return string(randId), nil

}

func NewUserMemory() *UserMemory {
	return &UserMemory{Users: map[string]User{}}
}
