package db

import (
	"go-api/model"
)

type DataStore struct {
	Closable
	UserProvider    userProvider
	ProjectProvider projectProvider
	ThemeProvider   themeProvider
	DiscussionProvider discussionProvider
}

type Closable interface {
	Close()
}

type userProvider interface {
	CreateUser(user *model.IDUser) (string, error)
	GetUser(id string) (*model.User, error)
	GetUserProfile(id string) (*model.UserProfile, error)
	UpdateUser(token string, user *model.UserProfile) (*model.UserProfile, error)
	RemoveUser(id string) error
	GetUserFollowers(id string) ([]model.User, error)
	GetUserFollows(id string) ([]model.User, error)
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

	// Creator APIs

	CreateProject(token string, project *model.Project) (string, error)
	RemoveProject(token string, id string) error

	// Admin APIs

	IsAdmin(token string, projectID string) (bool, error)

	UpdateProject(token string, project *model.Project) (*model.Project, error)
	RemoveMember(token string, projectID string, userID string) error
	ChangeAdmin(token string, projectID string, userID string) error
	AddTheme(token string, themeName string, projectID string) error
	RemoveTheme(token string, themeName string, projectID string) error

	// Public APIs

	GetProject(id string) (*model.Project, error)
	GetProjectStub(id string) (*model.ProjectStub, error)
	GetProjMembers(id string) ([]model.User, error)
}

type themeProvider interface {
	CreateTheme(theme *model.Theme) error
	GetTheme(themeName string) (*model.Theme, error)
	UpdateTheme(theme *model.Theme) (string, error)
	//GetProjectsWithTheme(themeName string) error
	DeleteTheme(themeName string) error
}

type discussionProvider interface {
	CreateDiscussion(proj_id string, discussion *model.DiscussionIn) (string, error)
	GetDiscussion(proj_id string, discNum string) (model.DiscussionOut, error)
}