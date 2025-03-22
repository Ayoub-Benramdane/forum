package database

import (
	structs "forum/Data"
)

func CreateRequest(id int64) error {
	_, err := DB.Exec("UPDATE users SET request = ? WHERE id = ?", false, id)
	return err
}

func GetRequests() ([]structs.Requests, error) {
	var requests []structs.Requests
	rows, err := DB.Query(`SELECT id, username, role FROM users WHERE request = ?`, false)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var request structs.Requests
		var role string
		if err := rows.Scan(&request.ID, &request.Username, &role); err != nil {
			return nil, err
		}
		if role == "user" {
			requests = append(requests, request)
		}
	}
	return requests, nil
}
