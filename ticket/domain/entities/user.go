package entities

import "github.com/google/uuid"

type User struct {
	Id    string
	Name  string
	Email string
}

func NewUser(name string, email string) *User {
	id := uuid.NewString()
	return &User{
		Id:    id,
		Name:  name,
		Email: email,
	}
}
