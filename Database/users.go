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
	err := DB.QueryRow("SELECT * FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	return &user, err
}
