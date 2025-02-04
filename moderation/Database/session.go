package database

import structs "forum/Data"

func CreateSession(username string, id int64, token string) error {
	_, err := DB.Exec("UPDATE users SET status = ?, token = ? WHERE username = ?", "Connected", token, username)
	return err
}

func GetUserConnected(token string) (*structs.User, error) {
	var user structs.User
	err := DB.QueryRow("SELECT id, username, status, role FROM users WHERE token = ?", token).Scan(&user.ID, &user.Username, &user.Status, &user.Role)
	return &user, err
}

func DeleteSession(username string) error {
	_, err := DB.Exec("UPDATE users SET status = ?, token = ? WHERE username = ?", "Disconnected", "", username)
	return err
}
