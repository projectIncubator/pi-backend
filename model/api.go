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

type DiscussionIn struct {
	UserID		string `json:"creator"`
	Title		string `json:"title"`
	Text		string `json:"text"`
	Media 		[]string
}