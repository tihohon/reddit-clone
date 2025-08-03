package user

type User struct {
	Username string `json:"username"`
	Id       string `json:"id"`
	Password string `json:"-"`
}

type UserRepo interface {
	GetUserById(id string) (*User, error)
	CreateUser(username string, password string) (userId string, err error)
	GetUsers(filterFunc func(*User) (bool, error)) ([]User, error)
}
