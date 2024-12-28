package database

import (
	structs "forum/Data"
)

func CreateSession(username string, id int64) error {
	_, err := DB.Exec("INSERT INTO session (username, user_id, statut) VALUES (?, ?, ?)", username, id, "Connected")
	return err
}

func GetUserConnected() *structs.Session {
	var session structs.Session
	err := DB.QueryRow("SELECT * FROM session").Scan(&session.ID, &session.Username, &session.UserID, &session.Status)
	if err != nil {
		return nil
	}
	return &session
}

func DeleteSession() error {
	_, err := DB.Exec("DELETE FROM session")
	return err
}
