package db

import (
	"go-api/model"
)

type DataStore struct {
	Closable
	UserProvider userProvider
	ProjectProvider projectProvider
}

type Closable interface {
	Close()
}

type userProvider interface {
	CreateUser(user *model.User) (string, error)
	GetUser(id string) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	RemoveUser(id string) error
	FollowUser(follow *model.Follows) error
	UnfollowUser(follow *model.Follows) error
	IntrestedProject(up *model.UserProject) error
	UnintrestedProject(up *model.UserProject) error
	JoinProject(up *model.UserProject) error
	QuitProject(up *model.UserProject) error
}

type projectProvider interface {
	CreateProject(project *model.Project) (string, error)
	GetProject(id string) (*model.Project, error)
	UpdateProject(project *model.Project) (*model.Project, error)
	RemoveProject(id string) error
	RemoveMember(projectID string, userID string) error
	ChangeAdmin(projectID string, userID string) error
	GetAdmin(projectID string) error
}
