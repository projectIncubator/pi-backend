package db

import (
	"go-api/model"
)

type DataStore struct {
	Closable
	UserProvider userProvider
}

type Closable interface {
	Close()
}

type userProvider interface {
	CreateUser(user *model.User) (string, error)
	GetUser(id string) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	RemoveUser(id string) error
}
