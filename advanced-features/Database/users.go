package database

import (
	structs "forum/Data"
	"time"
)

func CreateNewUser(username, email, hashedPassword string) error {
	_, err := DB.Exec("INSERT INTO users (username, email, password, created_at) VALUES (?, ?, ?, ?)", username, email, hashedPassword, time.Now())
	return err
}

func GetUserByUsername(username string) (*structs.User, error) {
	var user structs.User
	err := DB.QueryRow("SELECT username, email, password, created_at FROM users WHERE username = ?", username).Scan(&user.Username, &user.Email, &user.Password, &user.CreatedAt)
	return &user, err
}

func CreateSession(username string, id int64, token string) error {
	_, err := DB.Exec("UPDATE users SET status = ?, token = ? WHERE username = ?", "Connected", token, username)
	return err
}

func GetUserConnected(token string) (*structs.User, error) {
	var user structs.User
	err := DB.QueryRow("SELECT id, username, status FROM users WHERE token = ?", token).Scan(&user.ID, &user.Username, &user.Status)
	return &user, err
}

func DeleteSession(username string) error {
	_, err := DB.Exec("UPDATE users SET status = ?, token = ? WHERE username = ?", "Disconnected", "", username)
	return err
}
