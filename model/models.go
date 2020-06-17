package model

import (
	"time"
)

// TODO: Figure out how to store the password securely

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Image     string `json:"image"`
	ProfileID string `json:"profile_id"`
}

type UserProfile struct {
	User
	Email            string        `json:"email"`
	Password         string        `json:"password"`
	Deactivated      bool          `json:"deactivated"`
	Banned           bool          `json:"banned"`
	Bio              string        `json:"bio"`
	Following        []User        `json:"following"`    //TODO make count
	Followers        []User        `json:"followers"` 	 //TODO make count
	Interested       []ProjectStub `json:"interested"`   // These only store the id's of the projects rather than projects to reduce duplicated data
	Contributing     []ProjectStub `json:"contributing"` // ^
	Created          []ProjectStub `json:"created"`      // ^
	InterestedThemes []Theme       `json:"interested_themes"`
}

// TODO: Make project state enum not string

type ProjectStub struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	State  string  `json:"state"`
	Logo   string  `json:"logo"`
	Themes []Theme `json:"themes"`
	MemberCount int `json:"member_count"`
	InterestedCount int `json:"interested_count"`
}

type Project struct {
	ProjectStub
	//	Tags        []string  `json:"tags"`
	Creator    User    `json:"user_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	OneLiner   string    `json:"oneliner"`
	Discussion []DiscussionOut    `json:"discussion_id"`
	Admins     []User    `json:"admins"`
	CoverPhoto string    `json:"coverphoto"`
	//	Media       []string  `json:"media"`
	//	Modules     []ProjectModule `json:"projectModules"`
}

type Theme struct {
	Name        string `json:"name"`
	Colour      string `json:"colour"` //TODO: remove
	Logo        string `json:"logo"`
	Description string `json:"description"`
}

//TODO: this will be returned with filtered/algorithm-ed list of best project results
type SearchResults struct {
	//BestResults []ProjectStub `json:"results"`
}

type DiscussionOut struct {
	ProjID      string `json:"proj_id"`
	DiscNum     string
	UserID		string `json:"creator"`
	CreatedAt   string `json:"creation_date"`
	Title		string `json:"title"`
	Text		string `json:"text"`
	Closed      bool
	Media 		[]string
}

type Post struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
