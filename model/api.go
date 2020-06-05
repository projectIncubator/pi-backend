// This file contains all data types (structs) coming in from the frontend.
package model

type Follows struct {
	FollowerID 	string
	FollowedID  string
}

type UserProject struct {
	UserID string
	ProjectID string
}

type NewUserProfile struct {
	UserProfile
	IDToken string       `json:"id_token"`
}