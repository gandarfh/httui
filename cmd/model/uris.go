package model

import "github.com/gandarfh/httui/db"

func GetUris() ([]string, error) {
	db, err := db.Conn()
	defer db.Close()

	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT name FROM uris")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var list []string

	for rows.Next() {
		var name string
		err := rows.Scan(&name)

		if err != nil {
			break
		}
		list = append(list, name)
	}

	return list, nil
}

func CreateUri(name string) error {
	db, err := db.Conn()
	defer db.Close()

	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`INSERT INTO uris (name) VALUES (?)`)
	defer stmt.Close()

	if err != nil {
		return err
	}

	stmt.Exec(name)

	return nil
}
