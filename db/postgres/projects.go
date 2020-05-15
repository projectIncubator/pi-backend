package postgres

import "go-api/model"

func (p PostgresDBStore) CreateProject(project *model.Project) (string, error) {
	sqlStatement :=
		`INSERT INTO projects(id, title, state, /*tags,*/ user_id, created_date, end_date, oneliner, discussion_id, /*members,*/ logo, cover_photo /*media*/) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	var id string
	err := p.database.QueryRow(sqlStatement,
		project.ID,
		project.Title,
		project.State,
	//	project.Tags,
		project.Creator,
		project.CreatedDate,
		project.EndDate,
		project.OneLiner,
		project.Discussion,
	//	project.Members,
		project.Logo,
		project.CoverPhoto,
	//	project.Media,

	).Scan(&id)
	if err != nil {
		return "", err
	}
	if id != project.ID {
		return "", CreateError
	}
	return id, nil
}

func (p PostgresDBStore) GetProject(id string) (*model.Project, error) {
	sqlStatement := `SELECT id, title, state, /*tags,*/ user_id, created_date, end_date, oneliner, discussion_id, /*members,*/ logo, cover_photo /*media*/ FROM projects WHERE id=$1;`
	var project model.Project
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&project.ID,
		&project.Title,
		&project.State,
	//	&project.Tags,
		&project.Creator,
		&project.CreatedDate,
		&project.EndDate,
		&project.OneLiner,
		&project.Discussion,
	//	&project.Members,
		&project.Logo,
		&project.CoverPhoto,
	//	&project.Media,
	)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (p PostgresDBStore) UpdateProject(project *model.Project) (*model.Project, error) {
	sqlStatement :=
		`UPDATE project
				SET title = $2, stage = $3, /*tags = $4,*/ user_id = $4, created_date = $5, end_date = $6, oneliner = $7, discussion_id = $8, logo = $9, cover_photo = $10
				WHERE id = $1
				RETURNING id;`
	var _id string
	err := p.database.QueryRow(sqlStatement,
		project.ID,
		project.Title,
		project.State,
		project.Creator,
		project.CreatedDate,
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
