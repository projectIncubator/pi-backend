package postgres

import (
	"go-api/model"
)

func (p PostgresDBStore) CreateProject(project *model.Project) (string, error) {
	sqlStatement :=
		`INSERT INTO projects(title, state, user_id, start_date, end_date, oneliner, discussion_id, logo, coverphoto ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
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
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
func (p PostgresDBStore) GetProject(id string) (*model.Project, error) {
	sqlStatement := `SELECT id, title, state, user_id, start_date, end_date, oneliner, discussion_id, logo, coverphoto FROM projects WHERE id=$1;`
	var project model.Project
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&project.ID,
		&project.Title,
		&project.State,
		&project.Creator,
		&project.StartDate,
		&project.EndDate,
		&project.OneLiner,
		&project.Discussion,
		&project.Logo,
		&project.CoverPhoto,
	)
	if err != nil {
		return nil, err
	}
	//Fill in members array
	sqlStatement = `SELECT users.id, users.first_name, users.last_name, users.image, users.profile_id 
							FROM users, contributing
							WHERE users.id = contributing.user_id AND contributing.project_id=$1;`
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
		project.Members = append(project.Members, user)
	}

	//Fill in the admins array
	sqlStatement = `SELECT users.id, users.first_name, users.last_name, users.image, users.profile_id 
							FROM users, contributing
							WHERE users.id = contributing.user_id AND contributing.project_id=$1 AND contributing.is_admin = true;`
	rows, err = p.database.Query(sqlStatement, id)
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

	return &project, nil
}

func (p PostgresDBStore) UpdateProject(project *model.Project) (*model.Project, error) {
	sqlStatement :=
		`UPDATE projects
				SET title = $2, state = $3, user_id = $4, start_date = $5, end_date = $6, oneliner = $7, discussion_id = $8, logo = $9, coverphoto = $10
				WHERE id = $1
				RETURNING id;`
	var _id string
	err := p.database.QueryRow(sqlStatement,
		project.ID,
		project.Title,
		project.State,
		project.Creator,
		project.StartDate,
		project.EndDate,
		project.OneLiner,
		project.Discussion,
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

func (p PostgresDBStore) RemoveProject(id string) error {
	sqlStatement :=
		`DELETE FROM projects
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

func (p PostgresDBStore) RemoveMember(projectID string, userID string) error {
	sqlStatement :=
		`DELETE FROM contributing
				WHERE project_id = $1 AND user_id = $2
				RETURNING project_id, user_id;`
	var _projectID string
	var _userID string
	err := p.database.QueryRow(sqlStatement, projectID, userID).Scan(&_projectID, &_userID)
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

func (p PostgresDBStore) ChangeAdmin(projectID string, userID string) error {
	sqlStatement :=
		`UPDATE contributing
				SET is_admin = NOT is_admin 
				WHERE project_id = $1 AND user_id = $2
				RETURNING project_id, user_id;`
	var _projectID string
	var _userID string
	err := p.database.QueryRow(sqlStatement, projectID, userID).Scan(&_projectID, &_userID)
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

func (p PostgresDBStore) AddTheme(themeName string, projectID string) error {
	sqlStatement := `INSERT INTO project_has_theme(theme_name, project_id) VALUES ($1, $2)
						RETURNING theme_name, project_id`

	var _themeName, _projectID string
	err := p.database.QueryRow(sqlStatement,
		themeName,
		projectID,
	).Scan(_themeName, _projectID)

	if err != nil {
		return err
	}
	return nil
}

func (p PostgresDBStore) RemoveTheme(themeName string, projectID string) error {
	sqlStatement := `DELETE FROM project_has_theme
						WHERE theme_name = $1 AND project_id = $2
						RETURNING theme_name, project_id`

	var _userID, _themeName string
	err := p.database.QueryRow(sqlStatement,
		themeName,
		projectID,
	).Scan(&_userID, &_themeName)

	if err != nil {
		return err
	}
	return nil
}
