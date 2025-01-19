package database

import (
	structs "forum/Data"
	"time"
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

func GetUserFromToken(token string) (string, error) {
	var Username string
	err := DB.QueryRow("SELECT username FROM session WHERE token = ?", token).Scan(&Username)
	return Username, err
}

func DeleteSession(username string) error {
	_, err := DB.Exec("DELETE FROM session WHERE username = ?", username)
	return err
}

func GenerateToken(user string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(user), bcrypt.DefaultCost)
}
