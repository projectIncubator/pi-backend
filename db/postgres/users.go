package postgres

import (
	"go-api/model"
	"log"
)

//CREATE TABLE users
//(
//id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
//first_name      TEXT NOT NULL,
//last_name       TEXT NOT NULL,
//email           TEXT DEFAULT '' UNIQUE,
//image           TEXT,
//password        TEXT,
//profile_id      TEXT UNIQUE,
//deactivated     BOOLEAN DEFAULT FALSE,
//banned          BOOLEAN DEFAULT FALSE
//);

func (p PostgresDBStore) CreateUser(user *model.UserProfile) (string, error) {
	sqlStatement :=
		`INSERT INTO users(first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id`
	var id string
	err := p.database.QueryRow(sqlStatement,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (p PostgresDBStore) GetUser(id string) (*model.User, error) {
	sqlStatement := `SELECT id, first_name, last_name, image, profile_id FROM users WHERE id=$1;`
	var user model.User
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Image,
		&user.ProfileID,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}

func (p PostgresDBStore) GetUserProfile(id string) (*model.UserProfile, error) {
	sqlStatement := `SELECT * FROM users WHERE id=$1;`
	var userProfile model.UserProfile
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&userProfile.ID,
		&userProfile.FirstName,
		&userProfile.LastName,
		&userProfile.Email,
		&userProfile.Image,
		&userProfile.Password,
		&userProfile.ProfileID,
		&userProfile.Deactivated,
		&userProfile.Banned,
	)

	if err != nil {
		return nil, err
	}
	//People the user is following
	sqlStatement = `SELECT users.id, users.first_name, users.last_name, users.image, users.profile_id 
						FROM users, follows
						WHERE users.id = follows.followed_id AND follows.follower_id=$1;`

	rows, err := p.database.Query(sqlStatement, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user model.User

		if err := rows.Scan(
				&user.ID,
				&user.FirstName,
				&user.LastName,
				&user.Image,
				&user.ProfileID,
			); err != nil { log.Fatal(err) }

		userProfile.Following = append(userProfile.Following, user)
	}
	log.Println("finish the first part")
	//  followers of the user
	sqlStatement = `SELECT users.id, users.first_name, users.last_name, users.image, users.profile_id 
						FROM users, follows
						WHERE users.id = follows.follower_id AND follows.followed_id=$1;`

	rows, err = p.database.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}
	log.Println("before second loop")
	for rows.Next() {
		var user model.User

		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Image,
			&user.ProfileID,
		); err != nil {
			log.Fatal(err) }

		userProfile.Followers = append(userProfile.Followers, user)
	}
	log.Println("finish the second part")
	//Fill in the interested table
	sqlStatement = `SELECT projects.id, projects.title, projects.state, projects.logo
						FROM projects, intrested
						WHERE intrested.user_id = $1 AND intrested.project_id=projects.id;`

	rows, err = p.database.Query(sqlStatement, id)
	log.Println(err)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var proj model.ProjectStub

		if err := rows.Scan(
			&proj.ID,
			&proj.Title,
			&proj.State,
			&proj.Logo,
		); err != nil { log.Fatal(err) }

		userProfile.Interested = append(userProfile.Interested, proj)
	}

	sqlStatement = `SELECT projects.id, projects.title, projects.state, projects.logo
						FROM projects, contributing
						WHERE contributing.user_id = $1 AND contributing.project_id=projects.id;`

	rows, err = p.database.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var proj model.ProjectStub

		if err := rows.Scan(
			&proj.ID,
			&proj.Title,
			&proj.State,
			&proj.Logo,
		); err != nil { log.Fatal(err) }

		userProfile.Contributing = append(userProfile.Contributing, proj)
	}

	sqlStatement  =  `SELECT projects.id, projects.title, projects.state, projects.logo
						FROM projects
						WHERE projects.user_id = $1;`

	for rows.Next() {
		var proj model.ProjectStub

		if err := rows.Scan(
			&proj.ID,
			&proj.Title,
			&proj.State,
			&proj.Logo,
		); err != nil { log.Fatal(err) }

		userProfile.Created = append(userProfile.Created, proj)
	}

	return &userProfile, nil
}
//TODO: Problem: pq: invalid input syntax for type uuid: "" error when including ProfileID
func (p PostgresDBStore) UpdateUser(user *model.UserProfile) (*model.UserProfile, error) {
	sqlStatement :=
		`UPDATE users
				SET first_name = $2, last_name = $3, email = $4, image = $5, password = $6,/* profile_id = $7,*/ deactivated = $7, banned = $8
				WHERE id = $1
				RETURNING id;`
	var _id string
	err := p.database.QueryRow(sqlStatement,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Image,
		user.Password,
	//	user.ProfileID,
		user.Deactivated,
		user.Banned,
	).Scan(&_id)
	log.Println(err)
	log.Println("updateError")
	if err != nil {
		return nil, err
	}
	if _id != user.ID {
		return nil, CreateError
	}
	return user, nil
}

func (p PostgresDBStore) RemoveUser(id string) error {
	sqlStatement :=
		`UPDATE users
			SET deactivated = TRUE
			WHERE id = $1
			RETURNING id;`
	var _id string
	err := p.database.QueryRow(sqlStatement,
		id,
	).Scan(&_id)
	if err != nil {
		return err
	}
	if _id != id {
		return CreateError
	}
	return nil
}

func (p PostgresDBStore) FollowUser(follow *model.Follows) error {
	log.Println("we are in followUser")
	sqlStatement := `INSERT INTO follows(followed_id, follower_id) VALUES ($1, $2)
						RETURNING followed_id, follower_id`
	//var _follow *model.Follows
	var followed_id, follower_id string
	log.Println("we are in followUser2")
	err := p.database.QueryRow(sqlStatement,
		follow.FollowedID,
		follow.FollowerID,
	).Scan(&followed_id, &follower_id)
		/*Scan(&_follow.FollowedID,&_follow.FollowerID)*/
	log.Println("we are in followUser3")
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresDBStore) UnfollowUser(follow *model.Follows) error {
	sqlStatement := `DELETE FROM follows 
						WHERE followed_id = $1 AND follower_id = $2
						RETURNING followed_id, follower_id`
	var _follow *model.Follows
	err := p.database.QueryRow(sqlStatement,
		follow.FollowedID,
		follow.FollowerID,
	).Scan(&_follow.FollowedID,&_follow.FollowerID)
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresDBStore) IntrestedProject(up *model.UserProject) error {
	sqlStatement := `INSERT INTO intrested(user_id, project_id) VALUES ($1, $2) 
						RETURNING user_id, project_id`
	var _up *model.UserProject
	err := p.database.QueryRow(sqlStatement,
		up.UserID,
		up.ProjectID,
		).Scan(&_up.UserID,&_up.ProjectID)

	if err != nil {
		return err
	}
	if _up != up {
		return CreateError
	}
	return nil
}

func (p PostgresDBStore) UnintrestedProject(up *model.UserProject) error {
	sqlStatement := `DELETE FROM intrested 
						WHERE user_id = $1 AND project_id = $2
						RETURNING user_id, project_id`
	var _up *model.UserProject
	err := p.database.QueryRow(sqlStatement,
		up.UserID,
		up.ProjectID,
	).Scan(&_up.UserID,&_up.ProjectID)
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresDBStore) JoinProject(up *model.UserProject) error {
	sqlStatement := `INSERT INTO contributing(user_id, project_id) VALUES ($1, $2) 
						RETURNING user_id, project_id`
	var _up *model.UserProject
	err := p.database.QueryRow(sqlStatement,
		up.UserID,
		up.ProjectID,
	).Scan(&_up.UserID,&_up.ProjectID)

	if err != nil {
		return err
	}
	if _up != up {
		return CreateError
	}
	return nil
}

func (p PostgresDBStore) QuitProject(up *model.UserProject) error {
	sqlStatement := `DELETE FROM contributing
						WHERE user_id = $1 AND project_id = $2
						RETURNING user_id, project_id`
	var _up *model.UserProject
	err := p.database.QueryRow(sqlStatement,
		up.UserID,
		up.ProjectID,
	).Scan(&_up.UserID,&_up.ProjectID)
	if err != nil {
		return err
	}
	return nil
}
