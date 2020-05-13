package postgres

import "go-api/model"

func (p PostgresDBStore) CreateProject(project *model.Project) (string, error) {
	sqlStatement :=
		`INSERT INTO projects(id, title, state, tags, user_id, created_date, end_date, oneliner, discussion_id, members, logo, cover_photo, media) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`
	var id string
	err := p.database.QueryRow(sqlStatement,
		projects.ID,
		projects.FirstName,
		projects.LastName,
		projects.Email,
		projects.Image,
		projects.Password,
		projects.ProfileID,
		projects.Deactivated,
		projects.Banned,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	if id != user.ID {
		return "", CreateError
	}
	return id, nil
}
