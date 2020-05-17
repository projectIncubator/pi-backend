package model

import "time"

// TODO: Figure out how to store the password securely

type User struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Image        string `json:"img"`
}

type UserProfile struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Image        string `json:"img"`
	Password     string `json:"password"`
	ProfileURL   string `json:"profile_url"`
	Deactivated  bool   `json:"deactivated"`
	Banned       bool   `json:"banned"`
	Bio          string  `json:"bio"`
	Following    []User `json:"following"`
	Followers    []User `json:"followers"`
	Interested   []User `json:"interested"`   // These only store the id's of the projects rather than projects to reduce duplicated data
	Contributing []User `json:"contributing"` // ^
	Created      []User `json:"created"`      // ^
}


// TODO: Make project state enum not string
type Project struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	State       string    `json:"state"`
//	Tags        []string  `json:"tags"`
	Creator     string    `json:"user_id"`
	CreatedDate time.Time `json:"created_date"`
	EndDate     time.Time `json:"end_date"`
	OneLiner    string    `json:"oneliner"`
	Discussion  string    `json:"discussion_id"`
	Members     []string  `json:"members"`
	Logo        string    `json:"logo"`
	CoverPhoto  string    `json:"coverphoto"`
//	Media       []string  `json:"media"`
//	Modules     []ProjectModule `json:"projectModules"`
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
