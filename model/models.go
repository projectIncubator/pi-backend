package model

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

// TODO: Figure out how to store the password securely

type NullTime struct {
	time.Time
	Valid bool // Valid is true if Time is not NULL
}
// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}
// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

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
	Deactivated      bool          `json:"deactivated"`
	Banned           bool          `json:"banned"`
	Bio              sql.NullString`json:"bio"`
	Links			 sql.NullString`json:"links"`
	FollowingCount   int           `json:"following_count"`
	FollowersCount   int           `json:"followers_count"`
	Interested       []ProjectStub `json:"interested"`   // These only store the id's of the projects rather than projects to reduce duplicated data
	Contributing     []ProjectStub `json:"contributing"` // ^
	Created          []ProjectStub `json:"created"`      // ^
	InterestedThemes []Theme       `json:"interested_themes"`
}

// TODO: Make project state enum not string

type ProjectStub struct {
	ID    			string  `json:"id"`
	Title 			string  `json:"title"`
	State 			string  `json:"state"`
	Logo  			string  `json:"logo"`
	Themes 			[]Theme `json:"themes"`
	MemberCount		int `json:"member_count"`
	InterestedCount int `json:"interested_count" `
}

type Project struct {
	ProjectStub
	Creator    string    		`json:"user_id"`
	StartDate  time.Time 		`json:"start_date"`
	EndDate    NullTime 		`json:"end_date"`
	OneLiner   string    		`json:"oneliner"`
	Discussion []DiscussionOut  `json:"discussion_id"`
	Admins     []User    		`json:"admins"`
	CoverPhoto string   		`json:"coverphoto"`
	SideBar	   []SideBarModule  `json: "sidebar"`
	//	Media       []string  `json:"media"`
	//	Modules     []ProjectModule `json:"projectModules"`
}

type SideBarModule struct {
	Type        string `json:"type"`
	Content 	string `json:"content"`
}

type Theme struct {
	Name        string `json:"name"`
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
