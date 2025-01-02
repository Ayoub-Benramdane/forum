package database

import (
	"log"
	"time"

	structs "forum/Data"

	"golang.org/x/crypto/bcrypt"
)

func CreateSession(username string, id int64, token string) error {
	_, err := DB.Exec("INSERT INTO session (username, user_id,  status, token,created_at) VALUES (?, ?, ?, ?, ?)", username, id, "Connected", token, time.Now())
	return err
}

func GetUserConnected(token string) *structs.Session {
	var session structs.Session
	err := DB.QueryRow("SELECT id, username, user_id, status FROM session WHERE token = ?", token).Scan(&session.ID, &session.Username, &session.UserID, &session.Status)
	if err != nil {
		return nil
	}
	return &session
}

func DeleteSession() error {
	_, err := DB.Exec("DELETE FROM session")
	return err
}

func GenerateToken(email string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}
