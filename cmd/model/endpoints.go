package model

import "github.com/gandarfh/httui/db"

type result map[string]string

func GetEndpoints(uri string) ([]result, error) {
	db, err := db.Conn()
	defer db.Close()

	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM endpoints INNER JOIN uris on uris.name = endpoints.uri WHERE uri=?", uri)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var list []result

	for rows.Next() {
		var method, path, headers, body string

		err := rows.Scan(method, path, headers, body)

		if err != nil {
			break
		}
		list = append(list, result{method: method, path: path, headers: headers, body: body})
	}

	return list, nil
}

func CreateEndpoint(values result) error {
	db, err := db.Conn()
	defer db.Close()

	if err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO endpoints (method, path, uri) values(?,?,?)")
	defer stmt.Close()

	if err != nil {
		return err
	}

	stmt.Exec(values["method"], values["path"], values["uri"])

	return nil
}
