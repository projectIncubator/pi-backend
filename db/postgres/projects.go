package postgres

import (
	"go-api/model"
)

// Creator APIs

func (p PostgresDBStore) CreateProject(token string, project *model.Project) (string, error) {
	sqlStatement :=
		`INSERT INTO projects(title, state, creator, start_date, end_date, oneliner, discussion_id, logo, cover_photo ) 
			SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9 FROM users
			WHERE id = $3 AND id_token = $10
			RETURNING id;`
	var id string
	err := p.database.QueryRow(sqlStatement,
		project.Title,
		project.State,
		project.Creator,
		project.StartDate,
		project.EndDate,
		project.OneLiner,
		project.Discussion,
		project.Logo,
		project.CoverPhoto,
		token,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	sqlStatement =
		`INSERT INTO contributing(user_id, project_id, is_admin)
			VALUES ($1, $2, true);`
	p.database.QueryRow(sqlStatement,project.Creator,id)

	return id, nil
}
func (p PostgresDBStore) RemoveProject(token string, id string) error {
	sqlStatement :=
		`DELETE FROM projects
			WHERE id IN (
				SELECT p.id
				FROM projects AS p INNER JOIN
				   (SELECT users.id FROM users WHERE users.id_token = $2) AS u
				   ON p.creator = u.id
				WHERE p.id = $1
			)
			RETURNING id;`
	var _id string
	err := p.database.QueryRow(sqlStatement,
		id,
		token,
	).Scan(&_id)
	if err != nil {
		return err
	}
	if _id != id {
		return CreateError
	}
	return nil
}

// ... + Admins APIs

func (p PostgresDBStore) IsAdmin(token string, projectID string) (bool, error) {

	var isAdmin bool

	sqlStatement :=
		`SELECT CASE
   			WHEN EXISTS (
       			SELECT u.id FROM users AS u, contributing AS c
       			WHERE u.id_token = $1 AND c.project_id = $2
        		AND u.id = c.user_id AND c.is_admin = true
    		) THEN TRUE
   			ELSE FALSE
		END;`

	err := p.database.QueryRow(sqlStatement,
		token,
		projectID,
	).Scan(&isAdmin)
	if err != nil {
		return false, err
	}

	return isAdmin, nil
}

func (p PostgresDBStore) UpdateProject(token string, project *model.Project) (*model.Project, error) {

	isAdmin, err := p.IsAdmin(token, project.ID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, nil // TODO: Scope Error Handling
	}

	sqlStatement :=
		`UPDATE projects
				SET title = $2, state = $3, creator = $4, start_date = $5, end_date = $6, oneliner = $7, logo = $8, cover_photo = $9
				WHERE id = $1
				RETURNING id;`
	var _id string
	err = p.database.QueryRow(sqlStatement,
		project.ID,
		project.Title,
		project.State,
		project.Creator,
		project.StartDate,
		project.EndDate,
		project.OneLiner,
		project.Logo,
		project.CoverPhoto,
	).Scan(&_id)
	if err != nil {
		return nil, err
	}
	if _id != project.ID {
		return nil, CreateError
	}
	return project, nil
}
func (p PostgresDBStore) RemoveMember(token string, projectID string, userID string) error {

	isAdmin, err := p.IsAdmin(token, projectID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return nil // TODO: Scope Error Handling
	}

	sqlStatement :=
		`DELETE FROM contributing
				WHERE project_id = $1 AND user_id = $2
				RETURNING project_id, user_id;`
	var _projectID string
	var _userID string
	err = p.database.QueryRow(sqlStatement, projectID, userID).Scan(&_projectID, &_userID)
	if err != nil {
		return err
	}
	if _projectID != projectID {
		return CreateError
	}
	if _userID != userID {
		return CreateError
	}
	return nil
}
func (p PostgresDBStore) ChangeAdmin(token string, projectID string, userID string) error {

	isAdmin, err := p.IsAdmin(token, projectID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return nil // TODO: Scope Error Handling
	}

	sqlStatement :=
		`UPDATE contributing
				SET is_admin = NOT is_admin 
				WHERE project_id = $1 AND user_id = $2
				RETURNING project_id, user_id;`
	var _projectID string
	var _userID string
	err = p.database.QueryRow(sqlStatement, projectID, userID).Scan(&_projectID, &_userID)
	if err != nil {
		return err
	}
	if _projectID != projectID {
		return CreateError
	}
	if _userID != userID {
		return CreateError
	}
	return nil
}
func (p PostgresDBStore) AddTheme(token string, themeName string, projectID string) error {

	isAdmin, err := p.IsAdmin(token, projectID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return nil // TODO: Scope Error Handling
	}

	sqlStatement := `INSERT INTO project_has_theme(theme_name, project_id) VALUES ($1, $2)
						RETURNING theme_name, project_id`

	var _themeName, _projectID string
	err = p.database.QueryRow(sqlStatement,
		themeName,
		projectID,
	).Scan(&_themeName, &_projectID)

	if err != nil {
		return err
	}
	return nil
}
func (p PostgresDBStore) RemoveTheme(token string, themeName string, projectID string) error {

	isAdmin, err := p.IsAdmin(token, projectID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return nil // TODO: Scope Error Handling
	}

	sqlStatement := `DELETE FROM project_has_theme
						WHERE theme_name = $1 AND project_id = $2
						RETURNING theme_name, project_id`

	var _userID, _themeName string
	err = p.database.QueryRow(sqlStatement,
		themeName,
		projectID,
	).Scan(&_userID, &_themeName)

	if err != nil {
		return err
	}
	return nil
}

// Public APIs

func (p PostgresDBStore) GetProjMembers(id string) ([]model.User, error) {
	var members []model.User
	sqlStatement :=
		`SELECT users.id, users.first_name, users.last_name, users.image, users.profile_id
			FROM users, contributing
			WHERE contributing.project_id = $1 AND users.id = contributing.user_id`
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
		members = append(members, user)
	}
	return members, nil
}
func (p PostgresDBStore) GetProjectStub(id string) (*model.ProjectStub, error) {
	sqlStatement := `SELECT id, title, state, logo FROM projects WHERE id=$1;`
	var projectStub model.ProjectStub
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&projectStub.ID,
		&projectStub.Title,
		&projectStub.State,
		&projectStub.Logo,
		)
	if err != nil {
		return nil, err
	}

	// Fill in the themes array
	sqlStatement = `SELECT themes.name, themes.logo, themes.description
						FROM themes,project_has_theme
						WHERE themes.name = project_has_theme.theme_name AND project_has_theme.project_id = $1;`
	rows, err := p.database.Query(sqlStatement, id)
	for rows.Next() {
		var theme model.Theme
		if err = rows.Scan(
			&theme.Name,
			&theme.Logo,
			&theme.Description,
			); err!= nil {
			return nil, err
		}
		projectStub.Themes = append(projectStub.Themes, theme)
	}

	// Member Count
	sqlStatement = `SELECT COUNT(*) from contributing where project_id = $1`
	row = p.database.QueryRow(sqlStatement, id)
	err = row.Scan (
		&projectStub.MemberCount,
		)
	if err != nil {
		return nil, err
	}
	// Interested Count
	sqlStatement = `SELECT COUNT(*) from interested where project_id = $1`
	row = p.database.QueryRow(sqlStatement, id)
	err = row.Scan (
		&projectStub.InterestedCount,
	)
	if err != nil {
		return nil, err
	}

	return &projectStub, nil
}
func (p PostgresDBStore) GetProject(id string) (*model.Project, error) {
	var project model.Project

	sqlStatement := `SELECT id, title, state, creator, start_date, end_date, oneliner, logo, cover_photo FROM projects WHERE id=$1;`
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&project.ID,
		&project.Title,
		&project.State,
		&project.Creator,
		&project.StartDate,
		&project.EndDate,
		&project.OneLiner,
		&project.Logo,
		&project.CoverPhoto,
	)
	if err != nil {
		return nil, err
	}

	//Fill in the admins array
	sqlStatement = `SELECT users.id, users.first_name, users.last_name, users.image, users.profile_id 
							FROM users, contributing
							WHERE users.id = contributing.user_id AND contributing.project_id=$1 AND contributing.is_admin = true;`
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
		project.Admins = append(project.Admins, user)
	}
	// Fill in the discussion array
	sqlStatement = `SELECT proj_id, disc_num, creator, creation_date, title, text, closed 
						FROM discussions
						WHERE proj_id = $1`
	rows, err = p.database.Query(sqlStatement, id)
	for rows.Next() {
		var discussion model.DiscussionOut
		if err = rows.Scan(
			&discussion.ProjID,
			&discussion.DiscNum,
			&discussion.UserID,
			&discussion.CreatedAt,
			&discussion.Title,
			&discussion.Text,
			&discussion.Closed,
		); err!= nil {
			return nil, err
		}
		project.Discussion = append(project.Discussion, discussion)
	}

	// Fill in the themes array
	sqlStatement = `SELECT themes.name, themes.logo, themes.description
						FROM themes,project_has_theme
						WHERE themes.name = project_has_theme.theme_name AND project_has_theme.project_id = $1;`
	rows, err = p.database.Query(sqlStatement, id)
	for rows.Next() {
		var theme model.Theme
		if err = rows.Scan(
			&theme.Name,
			&theme.Logo,
			&theme.Description,
		); err!= nil {
			return nil, err
		}
		project.Themes = append(project.Themes, theme)
	}

	// Fill in the sidebar array
	sqlStatement = `SELECT module_type, content
						FROM sidebar_modules
						WHERE project_id = $1
						ORDER BY index ASC;`
	rows, err = p.database.Query(sqlStatement, id)
	for rows.Next() {
		var sidebar model.SideBarModule
		if err = rows.Scan(
			&sidebar.Type,
			&sidebar.Content,
		); err!= nil {
			return nil, err
		}
		project.SideBar = append(project.SideBar, sidebar)
	}

	// Member Count
	sqlStatement = `SELECT COUNT(*) from contributing where project_id = $1`
	row = p.database.QueryRow(sqlStatement, id)
	err = row.Scan (
		&project.MemberCount,
	)
	if err != nil {
		return nil, err
	}
	// Interested Count
	sqlStatement = `SELECT COUNT(*) from interested where project_id = $1`
	row = p.database.QueryRow(sqlStatement, id)
	err = row.Scan (
		&project.InterestedCount,
	)
	if err != nil {
		return nil, err
	}


	return &project, nil
}


