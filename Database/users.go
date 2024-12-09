package database

import structs "forum/Structs"

func CreateNewUser(username, hashedPassword string) error {
	_, err := DB.Exec(`INSERT INTO users (username, password) VALUES (?, ?)`, username, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByUsername(username string) (*structs.User, error) {
	var user structs.User
	err := DB.QueryRow(`SELECT * FROM users WHERE username = ?`, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
