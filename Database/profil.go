package database

import (
	structs "forum/Structs"
)

func GetInfoUser(UserID int64) (*structs.User, error) {
	var user structs.User
	err := DB.QueryRow(`SELECT * FROM users WHERE id = ?`, UserID).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateInfo(UserID int64, username, email string) error {
	_, err := DB.Exec(`UPDATE users SET username = ?, email = ? WHERE id = ?`, username, email, UserID)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`UPDATE session SET username = ? WHERE user_id = ?`, username, UserID)
	return err
}

func UpdatePass(UserID int64, password string) error {
	_, err := DB.Exec(`UPDATE users SET  password = ? WHERE id = ?`, password, UserID)
	return err
}
