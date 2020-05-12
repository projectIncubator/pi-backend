package db

import (
	"go-api/model"
)

type DataStore struct {
	Closable
	// TODO: Will have providers such as the following
	// UserProvider		userProvider
}

type Closable interface {
	Close()
}
// TODO: Define the interfaces here for all the functions that can be called in Postgres
// Example

type userProvider interface {
	// Function(parameters) (output1, output2)
	GetOneUser(id string) (*model.User, error)
	GetAllUsers() ([]*model.User, error)
	CreateUser(user *model.User) error
}