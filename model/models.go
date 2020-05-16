package model

import "time"

// TODO: Figure out how to store the password securely

type User struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Image       string `json:"img"`
	Password    string `json:"password"`
	ProfileID   string `json:"profile_id"`
	Deactivated bool   `json:"deactivated"`
	Banned      bool   `json:"banned"`
}

type UserProfile struct {
	ID           string   `json:"id"`
	UserID       string   `json:"user_id"`
	Bio          string   `json:"bio"`
	Following    []string `json:"following"`
	Followers    []string `json:"followers"`
	Interested   []string `json:"interested"`   // These only store the id's of the projects rather than projects to reduce duplicated data
	Contributing []string `json:"contributing"` // ^
	Created      []string `json:"created"`      // ^
	Banned       bool     `json:"banned"`
}

// TODO: Make project state enum not string

type Project struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	State       string    `json:"state"`
	//Tags        []string  `json:"tags"`
	Creator     string    `json:"user_id"`
	CreatedDate time.Time `json:"created_date"`
	EndDate     time.Time `json:"end_date"`
	OneLiner    string    `json:"oneliner"`
	Discussion  string    `json:"discussion_id"`
//	Members     []string  `json:"members"`
	Logo        string    `json:"logo"`
	CoverPhoto  string    `json:"coverphoto"`
//	Media       []string  `json:"media"`
	// Modules []ProjectModule `json:"projectModules"`
}

type Theme struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Discussion struct {
	ID    string   `json:"id"`
	Posts []string `json:"posts"`
}

type Post struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
