// This file contains all data types (structs) coming in from the frontend.
package model

type Follows struct {
	FollowerID 	string
	FollowedID  string
}
func NewFollows() Follows {
	follows := Follows{}
	return follows
}

type UserProject struct {
	UserID string
	ProjectID string
}
func NewUserProject()  UserProject {
	userProject := UserProject{}
	return userProject
}

type IDUser struct {
	IDToken   string `json:"id_token"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
func NewIDUser() IDUser {
	idUser := IDUser{}
	return idUser
}


type DiscussionIn struct {
	UserID		string `json:"creator"`
	Title		string `json:"title"`
	Text		string `json:"text"`
	Media 		[]string `json:"media"`
}
func NewDiscussionIn() DiscussionIn {
	discussionIn := DiscussionIn{}
	discussionIn.Media = []string{}
	return discussionIn
}
