package postgres

import "go-api/model"

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

func (p PostgresDBStore) CreateUser(user *model.User) (string, error) {
	sqlStatement :=
		`INSERT INTO users(first_name, last_name, email, image, password, profile_id, deactivated, banned) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	var id string
	err := p.database.QueryRow(sqlStatement,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Image,
		user.Password,
		user.ProfileID,
		user.Deactivated,
		user.Banned,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	//if id != user.ID {
	//	return "", CreateError
	//}
	return id, nil
}

func (p PostgresDBStore) GetUser(id string) (*model.User, error) {
	sqlStatement := `SELECT id, first_name, last_name, email, image, password, profile_id, deactivated, banned FROM users WHERE id=$1;`
	var user model.User
	row := p.database.QueryRow(sqlStatement, id)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Image,
		&user.Password,
		&user.ProfileID,
		&user.Deactivated,
		&user.Banned,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p PostgresDBStore) UpdateUser(user *model.User) (*model.User, error) {
	sqlStatement :=
		`UPDATE users
				SET first_name = $2, last_name = $3, email = $4, image = $5, password = $6, profile_id = $7, deactivated = $8, banned = $9
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
		`DELETE FROM users
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