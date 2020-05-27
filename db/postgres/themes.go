package postgres

import (
	"go-api/model"
	"log"
)

func (p PostgresDBStore) CreateTheme(theme *model.Theme) error {
	sqlStatement :=
		`INSERT INTO themes(name, colour, logo, description) VALUES ($1, $2, $3, $4) RETURNING name`

	var name string
	err := p.database.QueryRow(sqlStatement,
		theme.Name,
		theme.Colour,
		theme.Logo,
		theme.Description,
	).Scan(&name)

	if err != nil {
		return err
	}

	return nil
}

func (p PostgresDBStore) GetTheme(themeName string) (*model.Theme, error) {
	sqlStatement := `SELECT name, colour, logo, description FROM themes WHERE name=$1;`

	var theme model.Theme
	row := p.database.QueryRow(sqlStatement, themeName)
	err := row.Scan(
		&theme.Name,
		&theme.Colour,
		&theme.Logo,
		&theme.Description,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &theme, nil
}

func (p PostgresDBStore) UpdateTheme(theme *model.Theme) (string, error) {
	sqlStatement :=
		`UPDATE themes
				SET name = $1, colour = $2, logo = $3, description = $4
				WHERE name = $1
				RETURNING name;`
	var _name string
	err := p.database.QueryRow(sqlStatement,
		theme.Name,
		theme.Colour,
		theme.Logo,
		theme.Description,
	).Scan(&_name)

	if err != nil {
		return "", err
	}
	if _name != theme.Name {
		return "", CreateError
	}
	return _name, nil
}

//TODO: add GetProjectsWithTheme function
func (p PostgresDBStore) GetProjectsWithTheme(themeName string) error {
	// Get all projects by theme
	return nil
}

func (p PostgresDBStore) DeleteTheme(themeName string) error {
	sqlStatement :=
		`DELETE FROM themes
				WHERE name = $1
				RETURNING name;`

	var _name string
	err := p.database.QueryRow(sqlStatement,
		themeName,
	).Scan(&_name)
	if err != nil {
		return err
	}
	if _name != themeName {
		return CreateError
	}
	return nil
}
