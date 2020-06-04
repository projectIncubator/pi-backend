package postgres

import (
	"go-api/model"
)

func (p PostgresDBStore) CreateDiscussion(proj_id string, discussion *model.DiscussionIn) (string, error) {
	sqlStatement :=
		`INSERT INTO post(proj_id, creator, title, text) VALUES ($1, $2, $3, $4) RETURNING disc_num`
	var id string
	err := p.database.QueryRow(sqlStatement,
		proj_id,
		discussion.UserID,
		discussion.Title,
		discussion.Text,
	).Scan(&id)

	if err != nil {
		return "", err
	}

	for i, mediaUrl := range discussion.Media {
		sqlStatement :=
			`INSERT INTO discussion_has_media(proj_id, disc_num, media_url) VALUES ($1, $2, $3)`

		err := p.database.QueryRow(sqlStatement,
			proj_id,
			id,
			mediaUrl,
		).Scan()

		if err != nil {
			return "", err
		}
	}

	return id, nil
}

func (p PostgresDBStore) GetDiscussion(proj_id string, discNum string) (model.DiscussionOut, error) {
	sqlStatement := `SELECT * FROM discussions WHERE proj_id=$1 AND disc_num=$2;`

	var discussion model.DiscussionOut
	row := p.database.QueryRow(sqlStatement, proj_id, discNum)
	err := row.Scan(
		&discussion.ProjID,
		&discussion.DiscNum,
		&discussion.UserID,
		&discussion.CreatedAt,
		&discussion.Title,
		&discussion.Text,
		&discussion.Closed,
	)
	if err != nil {
		return discussion, err
	}
	//Fill in members array
	sqlStatement = `SELECT media_url FROM discussion_has_media 
							WHERE proj_id = $1 AND disc_num=$2;`
	rows, err := p.database.Query(sqlStatement, proj_id, discNum)
	if err != nil {
		return discussion, err
	}
	for rows.Next() {
		var url string
		if err := rows.Scan(url); err != nil {
			return discussion, err
		}
		discussion.Media = append(discussion.Media, url)
	}

	return discussion, nil
}
