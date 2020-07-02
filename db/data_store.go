package db

import (
	"database/sql"
	"go-api/model"
)

type DataStore struct {
	Closable
	UserProvider    userProvider
	ProjectProvider projectProvider
	ThemeProvider   themeProvider
	DiscussionProvider discussionProvider
	ScopeProvider	scopeProvider
}

type Closable interface {
	Close()
}

type userProvider interface {
	CreateUser(user *model.IDUser) (string, error)
	GetUser(id string) (*model.User, error)
	GetUserProfile(id string) (*model.UserProfile, error)
	UpdateUser(id string, user *model.UserProfile) (*model.UserProfile, error)
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

	CreateProject(project *model.Project) (string, error)
	RemoveProject(id string) error

	// Admin APIs

	UpdateProject(project *model.Project) (*model.Project, error)
	RemoveMember(projectID string, userID string) error
	ChangeAdmin(projectID string, userID string) error
	AddTheme(themeName string, projectID string) error
	RemoveTheme(themeName string, projectID string) error

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

type scopeProvider interface {
	GetCreatorID(token string, projectID string) (sql.NullString, error)
	GetAdminID(token string, projectID string) (sql.NullString, error)
	GetUserID(token string) (sql.NullString, error)
}