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
	CreateUser(user *model.UserProfile) (string, error)
	GetUser(id string) (*model.User, error)
	GetUserProfile(id string) (*model.UserProfile, error)
	UpdateUser(user *model.UserProfile) (*model.UserProfile, error)
	RemoveUser(id string) error
	FollowUser(followedID string, followerID string) error
	UnfollowUser(followedID string, followerID string) error
	InterestedProject(userID string, projectID string) error
	UninterestedProject(userID string, projectID string) error
	JoinProject(up *model.UserProject) error
	QuitProject(up *model.UserProject) error
}

type projectProvider interface {
	CreateProject(project *model.Project) (string, error)
	GetProject(id string) (*model.Project, error)
	UpdateProject(project *model.Project) (*model.Project, error)
	RemoveProject(id string) error
}
