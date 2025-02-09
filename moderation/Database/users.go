package database

import (
	structs "forum/Data"
	"time"
)

func CreateNewUser(username, email, hashedPassword string) error {
	_, err := DB.Exec("INSERT INTO users (username, email, password, created_at, role) VALUES (?, ?, ?, ?, ?)", username, email, hashedPassword, time.Now(), "user")
	return err
}

func GetUserByUsername(username string) (*structs.User, error) {
	var user structs.User
	err := DB.QueryRow("SELECT username, email, password, created_at FROM users WHERE username = ?", username).Scan(&user.Username, &user.Email, &user.Password, &user.CreatedAt)
	return &user, err
}

func GetAllUsers() ([]structs.User, error) {
	rows, err := DB.Query("SELECT id, username, email, role, created_at, status FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []structs.User
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

func UpdateRole(role string, id int64) error {
	_, err := DB.Exec("UPDATE users SET role = ? WHERE id = ?", role, id)
	return err
}
