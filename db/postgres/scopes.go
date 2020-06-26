package postgres

import (
	"database/sql"
)

func (p PostgresDBStore) GetCreatorID(token string, projectID string) (sql.NullString, error) {

	var creatorID sql.NullString

	sqlStatement :=
		`SELECT projects.creator FROM users, projects
		WHERE projects.creator = users.id AND users.id_token = $1 AND projects.id = $2;`

	err := p.database.QueryRow(sqlStatement,
		token,
		projectID,
		).Scan(&creatorID)
	if err != nil {
		return creatorID, err
	}

	return creatorID, nil
}

func (p PostgresDBStore) GetAdminID(token string, projectID string) (sql.NullString, error) {

	var adminID sql.NullString

	sqlStatement :=
		`SELECT u.id FROM users AS u, contributing AS c
       		WHERE u.id_token = $1 AND c.project_id = $2
        	AND u.id = c.user_id AND c.is_admin = TRUE`

	err := p.database.QueryRow(sqlStatement,
		token,
		projectID,
	).Scan(&adminID)
	if err != nil {
		return adminID, err
	}

	return adminID, nil
}

func (p PostgresDBStore) GetUserID(token string) (sql.NullString, error) {

	var userID sql.NullString

	sqlStatement :=
		`SELECT id FROM users WHERE id_token = $1`

	err := p.database.QueryRow(sqlStatement,
		token,
	).Scan(&userID)
	if err != nil {
		return userID, err
	}

	return userID, nil
}