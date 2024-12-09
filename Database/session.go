package database

import (
	structs "forum/Structs"
)

func CreateSession(username string, id int64) error {
	_, err := DB.Exec(`INSERT INTO session (username, user_id, statut) VALUES (?, ?, ?)`, username, id, "Logout")
	if err != nil {
		return err
	}
	return nil
}

func GetUserConnected() *structs.Session {
	var session structs.Session
	err := DB.QueryRow(`SELECT * FROM session`).Scan(&session.ID, &session.Username, &session.UserID, &session.Statut)
	if err != nil {
		return nil
	}
	return &session
}

func DeleteSession() error {
	_, err := DB.Exec(`DELETE FROM session`)
	if err != nil {
		return err
	}
	return nil
}
