package database

import (
	"time"

	structs "forum/Data"
)

func CreateNewUser(username, email, hashedPassword string) error {
	_, err := DB.Exec("INSERT INTO users (username, email, password, created_at, status, role, request) VALUES (?, ?, ?, ?, ?, ?, ?)", username, email, hashedPassword, time.Now(), "Disconnected", "user", true)
	return err
}

func GetUserByUsername(username string) (*structs.User, error) {
	var user structs.User
	err := DB.QueryRow("SELECT username, email, password, created_at FROM users WHERE username = ?", username).Scan(&user.Username, &user.Email, &user.Password, &user.CreatedAt)
	return &user, err
}

func GetAllUsers() ([]structs.User, error) {
	var users []structs.User
	rows, err := DB.Query("SELECT id, username, email, role, created_at, status FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user structs.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.Status); err != nil {
			return nil, err
		}
		if user.Role != "admin" {
			users = append(users, user)
		}
	}
	return users, nil
}
