package database

import (
	structs "forum/Data"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin() error {
	hashedPassword, errCrepting := bcrypt.GenerateFromPassword([]byte("Aa@00000"), bcrypt.DefaultCost)
	if errCrepting != nil {
		return errCrepting
	}
	if _, err := CheckAdmin(); err != nil {
		_, err := DB.Exec("INSERT INTO users (username, email, password, created_at, type, status, role) VALUES (?, ?, ?, ?, ?, ?, ?)", "molchi", "molchi@gmail.com", hashedPassword, time.Now(), "", "", "admin")
		return err
	}
	return nil
}

func CheckAdmin() (*structs.User, error) {
	var user structs.User
	err := DB.QueryRow("SELECT username FROM users WHERE role = ?", "admin").Scan(&user.Username)
	return &user, err
}
