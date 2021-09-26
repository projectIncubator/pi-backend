package postgres

import (
	"go-api/model"
	"log"
	"strconv"
	"net/http"
)

// Creator APIs

func (p PostgresDBStore) CreateProject(project *model.Project) (string, error) {
	sqlStatement :=
		`INSERT INTO projects(title, state, creator, start_date, end_date, oneliner, discussion_id, logo, cover_photo ) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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
	).Scan(&id)
	if err != nil {
		return "", err
	}

	sqlStatement =
		`INSERT INTO contributing(user_id, project_id, is_admin)
			VALUES ($1, $2, true);`
	p.database.QueryRow(sqlStatement, project.Creator, id)

	return id, nil
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

// ... + Admins APIs

func (p PostgresDBStore) UpdateProject(project *model.Project) (*model.Project, error) {

	sqlStatement :=
		`UPDATE projects
				SET title = $2, state = $3, creator = $4, start_date = $5, end_date = $6, oneliner = $7, logo = $8, cover_photo = $9
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

func (p PostgresDBStore) AddTheme(themeName string, projectID string) error {

	sqlStatement := `INSERT INTO project_has_theme(theme_name, project_id) VALUES ($1, $2)
						RETURNING theme_name, project_id`

	var _themeName, _projectID string
	err := p.database.QueryRow(sqlStatement,
		themeName,
		projectID,
	).Scan(&_themeName, &_projectID)

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

// Public APIs

func (p PostgresDBStore) GetProjects(r *http.Request) ([]model.ProjectStub, error) {
	var projects []model.ProjectStub
	sqlStatement :=	`SELECT id, title, state, logo, start_date, end_date, oneliner, COALESCE(contributors, 0), COALESCE(interested, 0)
					FROM
					projects p
					LEFT JOIN (SELECT project_id, COUNT(*) AS contributors FROM contributing GROUP BY project_id) c
					ON c.project_id = id
					LEFT JOIN (SELECT project_id, COUNT(*) AS interested FROM interested GROUP BY project_id) i
					ON i.project_id = id
					ORDER BY $1 $2
					LIMIT $3
					OFFSET $4;` // TODO replace the limit/offset with a faster query, keyset or something
					//sortBy=name&sort=desc&page=3&perPage=50 TODO create index on anything sortby
					
	sortBy := "name"
	key := r.URL.Query().Get("sortBy")
    if key != "" { sortBy = key } //TODO try catch, sanitize

	sort := "DESC" //TODO

	var err error
	var perPage int
	key = r.URL.Query().Get("perPage")
    if key != "" {
		perPage, err = strconv.Atoi(key)
		if err != nil { perPage = 5 }
	} //TODO try catch, sanitize
	
	var page int
	key = r.URL.Query().Get("page")
    if key != "" {
		page, err = strconv.Atoi(key)
		if err != nil { page = 0 }
	} //TODO try catch, sanitize

	//TODO filter

	rows, err := p.database.Query(sqlStatement, sortBy, sort, perPage, perPage * page)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		projectStub := model.NewProjectStub()
		if err := rows.Scan(
			&projectStub.ID,
			&projectStub.Title,
			&projectStub.State,
			&projectStub.Logo,
			&projectStub.StartDate,
			&projectStub.EndDate,
			&projectStub.OneLiner,
			&projectStub.MemberCount,
			&projectStub.InterestedCount,
		); err != nil {
			return nil, err
		}

		//TODO refactor all this to use a common func for getting stubs
		sqlStatement = `SELECT themes.name, themes.logo, themes.description
					FROM themes, project_has_theme
					WHERE themes.name = project_has_theme.theme_name AND project_has_theme.project_id = $1;`
		rows, _ := p.database.Query(sqlStatement, &projectStub.ID)
		for rows.Next() {
			var theme model.Theme
			if err = rows.Scan(
				&theme.Name,
				&theme.Logo,
				&theme.Description,
			); err != nil {
				return nil, err
			}
			projectStub.Themes = append(projectStub.Themes, theme)
		}

		projects = append(projects, projectStub)
	}

	return projects, nil
}
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
		members = append(members, user)
	}
	return members, nil
}

func (p PostgresDBStore) CreateProjectMedia(projectID string, mediaURL string) error {
	sqlStatement :=
		`INSERT INTO project_has_media(project_id, media) VALUES ($1, $2) RETURNING project_id, media`
	var returnedID, returnedURL string
	err := p.database.QueryRow(sqlStatement,
		projectID,
		mediaURL,
	).Scan(&returnedID, &returnedURL)
	if err != nil {
		return err
	}
	return nil
}

func (p PostgresDBStore) GetProjectStub(id string) (*model.ProjectStub, error) {
	sqlStatement :=	`SELECT id, title, state, logo, start_date, end_date, oneliner, COALESCE(contributors, 0), COALESCE(interested, 0)
					FROM
					projects p
					LEFT JOIN (SELECT project_id, COUNT(*) AS contributors FROM contributing GROUP BY project_id) c
					ON c.project_id = id
					LEFT JOIN (SELECT project_id, COUNT(*) AS interested FROM interested GROUP BY project_id) i
					ON i.project_id = id
					WHERE id=$1;`
	rows, err := p.database.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}
	projectStub := model.NewProjectStub()
	for rows.Next() {
		if err := rows.Scan(
			&projectStub.ID,
			&projectStub.Title,
			&projectStub.State,
			&projectStub.Logo,
			&projectStub.StartDate,
			&projectStub.EndDate,
			&projectStub.OneLiner,
			&projectStub.MemberCount,
			&projectStub.InterestedCount,
		); err != nil {
			return nil, err
		}

		// Fill in the themes array
		sqlStatement = `SELECT themes.name, themes.logo, themes.description
							FROM themes,project_has_theme
							WHERE themes.name = project_has_theme.theme_name AND project_has_theme.project_id = $1;`
		rows, _ = p.database.Query(sqlStatement, id)
		for rows.Next() {
			var theme model.Theme
			if err = rows.Scan(
				&theme.Name,
				&theme.Logo,
				&theme.Description,
			); err != nil {
				return nil, err
			}
			projectStub.Themes = append(projectStub.Themes, theme)
		}
	}

	return &projectStub, nil
}

func (p PostgresDBStore) GetProject(id string) (*model.Project, error) {
	project := model.NewProject()

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
	rows, _ = p.database.Query(sqlStatement, id)
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
		); err != nil {
			return nil, err
		}
		project.Discussion = append(project.Discussion, discussion)
	}
	// Fill in the themes array
	sqlStatement = `SELECT themes.name, themes.logo, themes.description
						FROM themes,project_has_theme
						WHERE themes.name = project_has_theme.theme_name AND project_has_theme.project_id = $1;`
	rows, _ = p.database.Query(sqlStatement, id)
	for rows.Next() {

		var theme model.Theme
		if err = rows.Scan(
			&theme.Name,
			&theme.Logo,
			&theme.Description,
		); err != nil {
			return nil, err
		}

		project.Themes = append(project.Themes, theme)
	}
	// Fill in the sidebar array
	sqlStatement = `SELECT module_type, content
						FROM sidebar_modules
						WHERE project_id = $1
						ORDER BY index ASC;`
	rows, _ = p.database.Query(sqlStatement, id)
	for rows.Next() {
		var sidebar model.SideBarModule
		if err = rows.Scan(
			&sidebar.Type,
			&sidebar.Content,
		); err != nil {
			return nil, err
		}
		project.SideBar = append(project.SideBar, sidebar)
	}

	// Member Count
	sqlStatement = `SELECT COUNT(*) from contributing where project_id = $1`
	row = p.database.QueryRow(sqlStatement, id)
	err = row.Scan(
		&project.MemberCount,
	)
	if err != nil {
		return nil, err
	}
	// Interested Count
	sqlStatement = `SELECT COUNT(*) from interested where project_id = $1`
	row = p.database.QueryRow(sqlStatement, id)
	err = row.Scan(
		&project.InterestedCount,
	)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (p PostgresDBStore) UpdateCoverPhoto(projectID string, coverURL string) (string, error) {
	sqlStatement :=
		`UPDATE projects
				SET cover_photo = $2
				WHERE id = $1
				RETURNING id;`
	var _id string
	err := p.database.QueryRow(sqlStatement,
		projectID,
		coverURL,
	).Scan(&_id)
	if err != nil {
		return "", err
	}
	if _id != projectID {
		return "", CreateError
	}
	return projectID, nil
}

func (p PostgresDBStore) UpdateLogo(projectID string, logo string) (string, error) {
	sqlStatement :=
		`UPDATE projects
				SET logo = $2
				WHERE id = $1
				RETURNING id;`
	var _id string
	err := p.database.QueryRow(sqlStatement,
		projectID,
		logo,
	).Scan(&_id)
	if err != nil {
		return "", err
	}
	if _id != projectID {
		return "", CreateError
	}
	return projectID, nil
}

func (p PostgresDBStore) ChangeAdmin(projectID string, userID string) error {
	//Make sure the user is not changing the only admin to non-admin
	_sqlStatement := `SELECT COUNT(*) FROM contributing WHERE is_admin = true;`
	var _count int
	_err := p.database.QueryRow(_sqlStatement).Scan(&_count)
	if _err != nil {
		return _err
	}
	if _count == 1 {
		_sqlStatement := `SELECT user_id FROM contributing WHERE is_admin = true;`
		var _userID string
		_err = p.database.QueryRow(_sqlStatement).Scan(&_userID)
		if _err != nil {
			return _err
		}
		if _userID == userID {
			log.Printf("App.ToggleProejct - there must be at least one admin")
			return CreateError
		}
	}
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

func (p PostgresDBStore) CheckAdmin(projectID string, userID string) bool {
	type _UserID struct {
		ID string `json:"id"`
	}
	var userTK _UserID
	userTK.ID = ""
	sqlStatement := `SELECT user_id FROM contributing WHERE project_id = $1 AND user_id = $2 AND is_admin = true`

	row := p.database.QueryRow(
		sqlStatement,
		projectID,
		userID,
	)
	err := row.Scan(
		&userTK.ID,
	)
	if err != nil {
		return false
	}
	if userTK.ID == "" {
		return false
	} else {
		return true
	}

}
