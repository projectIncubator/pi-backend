package db

import (
	"go-api/model"
)

type DataStore struct {
	Closable
	UserProvider    userProvider
	ProjectProvider projectProvider
	ThemeProvider   themeProvider
}

type Closable interface {
	Close()
}

type userProvider interface {
	CreateUser(user *model.NewUserProfile) (string, error)
	GetUser(id string) (*model.User, error)
	GetUserProfile(id string) (*model.UserProfile, error)
	UpdateUser(user *model.UserProfile) (*model.UserProfile, error)
	RemoveUser(id string) error
	FollowUser(followerID string, followedID string) error
	UnfollowUser(followerID string, followedID string) error
	InterestedProject(userID string, projectID string) error
	UninterestedProject(userID string, projectID string) error
	JoinProject(userID string, projectID string) error
	QuitProject(userID string, projectID string) error
	InterestedTheme(userID string, name string) error
	UninterestedTheme(userID string, name string) error
}

type projectProvider interface {
	CreateProject(project *model.Project) (string, error)
	GetProject(id string) (*model.Project, error)
	UpdateProject(project *model.Project) (*model.Project, error)
	RemoveProject(id string) error
	RemoveMember(projectID string, userID string) error
	ChangeAdmin(projectID string, userID string) error
	AddTheme(themeName string, projectID string) error
	RemoveTheme(themeName string, projectID string) error
}

type themeProvider interface {
	CreateTheme(theme *model.Theme) error
	GetTheme(themeName string) (*model.Theme, error)
	UpdateTheme(theme *model.Theme) (string, error)
	//GetProjectsWithTheme(themeName string) error
	DeleteTheme(themeName string) error
}
