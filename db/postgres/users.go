package postgres

import (
	"go-api/model"
)

// Private APIs

func (p PostgresDBStore) CreateUser(user *model.IDUser, userInfo *model.UserSessionInfo) error {

	sqlStatement :=
		`INSERT INTO users(id_token, first_name, last_name, email) VALUES ($1, $2, $3, $4) 
			RETURNING id, first_name, last_name, email, image, deactivated, banned, bio`
	err := p.database.QueryRow(sqlStatement,
		user.IDToken,
		user.FirstName,
		user.LastName,
		user.Email,
	).Scan(
		&userInfo.ID,
		&userInfo.FirstName,
		&userInfo.LastName,
		&userInfo.Email,
		&userInfo.Image,
		&userInfo.Deactivated,
		&userInfo.Banned,
		&userInfo.Bio,
	)
	if err != nil {
		return err
	}

	sqlStatement =
		`UPDATE users SET profile_id = $1 WHERE id = $1 RETURNING profile_id`
	err = p.database.QueryRow(sqlStatement,userInfo.ID).Scan(&userInfo.ProfileID)
	if err != nil {
		return err
	}

	return nil
}
func (p PostgresDBStore) LoginUser(user *model.IDUser, userInfo *model.UserSessionInfo) error {
	sqlStatement :=
		`SELECT id, profile_id, first_name, last_name, email, image, deactivated, banned, bio
			FROM users WHERE id_token = $1`
	err := p.database.QueryRow(sqlStatement,
		user.IDToken,
	).Scan(
		&userInfo.ID,
		&userInfo.ProfileID,
		&userInfo.FirstName,
		&userInfo.LastName,
		&userInfo.Email,
		&userInfo.Image,
		&userInfo.Deactivated,
		&userInfo.Banned,
		&userInfo.Bio,
		)
	if err != nil {
		return err
	}

	userInfo.Following, err = p.GetUserFollows(userInfo.ID)
	if err != nil {
		return err
	}
	userInfo.Followers, err = p.GetUserFollowers(userInfo.ID)
	if err != nil {
		return err
	}
	userInfo.Interested, err = p.GetUserInterested(userInfo.ID)
	if err != nil {
		return err
	}
	userInfo.Contributing, err = p.GetUserContributing(userInfo.ID)
	if err != nil {
		return err
	}
	userInfo.Created, err = p.GetUserCreated(userInfo.ID)
	if err != nil {
		return err
	}
	userInfo.InterestedThemes, err = p.GetUserInterestedThemes(userInfo.ID)
	if err != nil {
		return err
	}

	return nil
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

	user := model.NewUser()

	var sqlStatement string

	if IsValidUUID(id) {
		sqlStatement = `SELECT id, first_name, last_name, image, profile_id FROM users WHERE id=$1;`
	} else {
		sqlStatement = `SELECT id, first_name, last_name, image, profile_id FROM users WHERE profile_id=$1;`
	}

	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Image,
		&user.ProfileID,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (p PostgresDBStore) GetUserProfile(id string) (*model.UserProfile, error) {

	userProfile := model.NewUserProfile()
	var sqlStatement string

	if IsValidUUID(id) {
		sqlStatement =
			`SELECT id, first_name, last_name, email, image, profile_id, deactivated, banned FROM users WHERE id=$1;`
	} else {
		sqlStatement =
			`SELECT id, first_name, last_name, email, image, profile_id, deactivated, banned FROM users WHERE profile_id=$1`
	}


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

	//Fill in interested array
	userProfile.Interested, err = p.GetUserInterested(userProfile.ID)
	if err != nil {
		return nil, err
	}
	//Fill in contributing array
	userProfile.Contributing, err = p.GetUserContributing(userProfile.ID)
	if err != nil {
		return nil, err
	}
	//Fill in created array
	userProfile.Created, err = p.GetUserCreated(userProfile.ID)

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
        user := model.NewUser()

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
		user := model.NewUser()

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

// Helpers

func (p PostgresDBStore) GetUserInterested(id string) ([]model.ProjectStub, error) {
	projects := []model.ProjectStub{}

	sqlStatement := `SELECT project_id
						FROM interested
						WHERE interested.user_id = $1;`

	rows, err := p.database.Query(sqlStatement, id)
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

		projects = append(projects, *proj)
	}

	return projects, err
}
func (p PostgresDBStore) GetUserContributing(id string) ([]model.ProjectStub, error) {
	projects := []model.ProjectStub{}

	sqlStatement := `SELECT project_id
						FROM contributing
						WHERE contributing.user_id = $1;`

	rows, err := p.database.Query(sqlStatement, id)
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

		projects = append(projects, *proj)
	}

	return projects, err
}
func (p PostgresDBStore) GetUserCreated(id string) ([]model.ProjectStub, error) {

	projects := []model.ProjectStub{}

	sqlStatement := `SELECT id
						FROM projects
						WHERE projects.creator = $1;`

	rows, err := p.database.Query(sqlStatement, id)
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

		projects = append(projects, *proj)
	}
	return projects, err
}
func (p PostgresDBStore) GetUserInterestedThemes(id string) ([]model.Theme, error) {
	themes := []model.Theme{}

	sqlStatement :=  `SELECT name, logo, description FROM themes, user_interested_theme
						WHERE user_interested_theme.user_id = $1 
						AND user_interested_theme.theme_name = themes.name`

	rows, err := p.database.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		theme := model.NewTheme()

		if err := rows.Scan(
			&theme.Name,
			&theme.Logo,
			&theme.Description,
		); err != nil {
			return nil, err
		}

		themes = append(themes, theme)
	}

	return themes, nil
}