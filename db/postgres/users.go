package postgres

import (
	"go-api/model"
	"log"
)

// Private APIs

func (p PostgresDBStore) CreateUser(user *model.IDUser) (string, error) {
	sqlStatement :=
		`INSERT INTO users(id_token, first_name, last_name, email) VALUES ($1, $2, $3, $4) RETURNING id`
	var id string
	err := p.database.QueryRow(sqlStatement,
		user.IDToken,
		user.FirstName,
		user.LastName,
		user.Email,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	sqlStatement =
		`UPDATE users SET profile_id = $1 WHERE id = $1 RETURNING id`
	err = p.database.QueryRow(sqlStatement,id,).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
func (p PostgresDBStore) LoginUser(user *model.IDUser) (string, error) {
	sqlStatement :=
		`SELECT id FROM users WHERE id_token = $1`
	var id string
	err := p.database.QueryRow(sqlStatement,
		user.IDToken,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}


//TODO: Problem: pq: invalid input syntax for type uuid: "" error when including ProfileID
func (p PostgresDBStore) UpdateUser(id string, user *model.UserProfile) (*model.UserProfile, error) {
	sqlStatement :=
		`UPDATE users
				SET first_name = $2, last_name = $3, email = $4, image = $5, profile_id = $6, deactivated = $7, banned = $8
				WHERE id = $1
				RETURNING id;`
	var _id string
	err := p.database.QueryRow(sqlStatement,
		id,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Image,
		user.ProfileID,
		user.Deactivated,
		user.Banned,
	).Scan(&_id)
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

func (p PostgresDBStore) FollowUser(followerID string, followedID string) error {
	sqlStatement := `INSERT INTO follows(follower_id, followed_id) 
						VALUES ($1, $2)
						RETURNING follower_id, followed_id`

	var _followedID, _followerID string
	err := p.database.QueryRow(sqlStatement,
		followerID,
		followedID,
	).Scan(&_followerID, &_followedID)
	if err != nil {
		return err
	}
	return nil
}
func (p PostgresDBStore) UnfollowUser(followerID string, followedID string) error {

	sqlStatement := `DELETE FROM follows
						WHERE follower_id = $1 AND followed_id = $2
						RETURNING follower_id, followed_id`

	var _followedID, _followerID string
	err := p.database.QueryRow(sqlStatement,
		followerID,
		followedID,
	).Scan(&_followerID, &_followedID)
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresDBStore) InterestedProject(userID string, projectID string) error {

	sqlStatement := `INSERT INTO interested(user_id, project_id) 
						VALUES ($1,$2)
						RETURNING user_id, project_id;`

	var _userID, _projectID string
	err := p.database.QueryRow(sqlStatement,
		userID,
		projectID,
	).Scan(&_userID, &_projectID)

	if err != nil {
		return err
	}
	return nil
}
func (p PostgresDBStore) UninterestedProject(userID string, projectID string) error {

	sqlStatement := `DELETE FROM interested
						WHERE user_id = $1 AND project_id = $2
						RETURNING user_id, project_id`

	var _userID, _projectID string
	err := p.database.QueryRow(sqlStatement,
		userID,
		projectID,
	).Scan(&_userID, &_projectID)

	if err != nil {
		return err
	}
	return nil
}

func (p PostgresDBStore) JoinProject(userID string, projectID string) error {

	sqlStatement := `INSERT INTO contributing(user_id, project_id) 
						VALUES ($1, $2)`

	var _userID, _projectID string
	err := p.database.QueryRow(sqlStatement,
		userID,
		projectID,
	).Scan(&_userID, &_projectID)

	if err != nil {
		return err
	}
	return nil
}
func (p PostgresDBStore) QuitProject(userID string, projectID string) error {

	sqlStatement := `DELETE FROM contributing
						WHERE user_id = $1 AND project_id = $2
						RETURNING user_id, project_id`

	var _userID, _projectID string
	err := p.database.QueryRow(sqlStatement,
		userID,
		projectID,
	).Scan(&_userID, &_projectID)

	if err != nil {
		return err
	}
	return nil
}

func (p PostgresDBStore) InterestedTheme(userID string, themeName string) error {

	sqlStatement := `INSERT INTO user_interested_theme(user_id, theme_name) 
						VALUES ($1, $2)
						RETURNING user_id, project_id`

	var _userID, _themeName string
	err := p.database.QueryRow(sqlStatement,
		userID,
		themeName,
	).Scan(&_userID, &_themeName)

	if err != nil {
		return err
	}
	return nil
}
func (p PostgresDBStore) UninterestedTheme(userID string, themeName string) error {

	sqlStatement := `DELETE FROM user_interested_theme
						WHERE user_id = $1 AND theme_name = $2
						RETURNING user_id, project_id`

	var _userID, _themeName string
	err := p.database.QueryRow(sqlStatement,
		userID,
		themeName,
	).Scan(&_userID, &_themeName)

	if err != nil {
		return err
	}
	return nil
}

// Public APIs

func (p PostgresDBStore) GetUser(id string) (*model.User, error) {
	sqlStatement :=
		`SELECT id, first_name, last_name, image, profile_id FROM users WHERE id=$1 OR profile_id=$1;`
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

	sqlStatement :=
		`SELECT id, first_name, last_name, email, image, profile_id, deactivated, banned FROM users WHERE id=$1 OR profile_id=$1;`
	var userProfile model.UserProfile
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&userProfile.ID,
		&userProfile.FirstName,
		&userProfile.LastName,
		&userProfile.Email,
		&userProfile.Image,
		&userProfile.ProfileID,
		&userProfile.Deactivated,
		&userProfile.Banned,
	)

	if err != nil {
		return nil, err
	}
	//People the user is following
	sqlStatement =
		`SELECT COUNT(*) 
			FROM users, follows
			WHERE users.id = follows.followed_id AND follows.follower_id=$1;`

	row = p.database.QueryRow(sqlStatement, userProfile.ID)
	err = row.Scan(&userProfile.FollowingCount)
	if err != nil {
		return nil, err
	}

	//  followers of the user
	sqlStatement =
		`SELECT COUNT(*)
			FROM users, follows
			WHERE users.id = follows.follower_id AND follows.followed_id=$1;`

	row = p.database.QueryRow(sqlStatement, userProfile.ID)
	err = row.Scan(&userProfile.FollowersCount)
	if err != nil {
		return nil, err
	}


	//Fill in the interested table
	sqlStatement = `SELECT project_id
						FROM interested
						WHERE interested.user_id = $1;`

	rows, err := p.database.Query(sqlStatement, userProfile.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var proj_id string

		if err := rows.Scan(
			&proj_id,
		); err != nil {
			return nil, err
		}

		proj, err := p.GetProjectStub(proj_id)
		if err != nil {
			return nil, err
		}

		userProfile.Interested = append(userProfile.Interested, *proj)
	}

	sqlStatement = `SELECT project_id
						FROM contributing
						WHERE contributing.user_id = $1;`

	rows, err = p.database.Query(sqlStatement, userProfile.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var proj_id string

		if err := rows.Scan(
			&proj_id,
		); err != nil {
			return nil, err
		}

		proj, err := p.GetProjectStub(proj_id)
		if err != nil {
			return nil, err
		}

		userProfile.Contributing = append(userProfile.Contributing, *proj)
	}

	sqlStatement = `SELECT id
						FROM projects
						WHERE projects.creator = $1;`

	rows, err = p.database.Query(sqlStatement, userProfile.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var proj_id string

		if err := rows.Scan(
			&proj_id,
		); err != nil {
			return nil, err
		}

		proj, err := p.GetProjectStub(proj_id)
		if err != nil {
			return nil, err
		}

		userProfile.Created = append(userProfile.Created, *proj)
	}

	return &userProfile, nil
}
func (p PostgresDBStore) GetUserFollowers(id string) ([]model.User, error) {
	sqlStatement := `SELECT users.id, users.first_name, users.last_name, users.image, users.profile_id
						FROM users, follows
						WHERE users.id = follows.follower_id AND follows.followed_id=$1;`

	var followers []model.User
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
        ); err != nil {
            return nil, err
        }

        followers = append(followers, user)
    }
    return followers, nil
}
func (p PostgresDBStore) GetUserFollows(id string) ([]model.User, error) {
	sqlStatement := `SELECT users.id, users.first_name, users.last_name, users.image, users.profile_id
						FROM users, follows
						WHERE users.id = follows.followed_id AND follows.follower_id=$1;`

	var follows []model.User
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
		); err != nil {
			return nil, err
		}

		follows = append(follows, user)
	}
	return follows, nil
}

