package database

import (
	"database/sql"
	"forum/Structs"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func ConnectDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return err
	} else if err = createTables(); err != nil {
		return err
	}
	return nil
}

func createTables() error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL
        )
    `)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS posts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            user_id INTEGER NOT NULL,
            created_at DATETIME NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users(id)
        )
    `)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS comments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            content TEXT NOT NULL,
            user_id INTEGER NOT NULL,
            post_id INTEGER NOT NULL,
            created_at DATETIME NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (post_id) REFERENCES posts(id)
        )
    `)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS post_likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			like INT NOT NULL,
			created_at DATETIME NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (post_id) REFERENCES posts(id)
		)
    `)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS comment_likes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			comment_id INTEGER NOT NULL,
			like INT NOT NULL,
			created_at DATETIME NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (comment_id) REFERENCES comments(id)
		)
    `)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
    		name TEXT NOT NULL UNIQUE
		)
    `)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS post_tags (
			post_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			PRIMARY KEY (post_id, tag_id),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		)
    `)
	return err
}

func CreateNewUser(username, hashedPassword string) (*structs.User, error) {
	result, err := db.Exec(`INSERT INTO users (username, password) VALUES (?, ?)`, username, hashedPassword)
	if err != nil {
		return nil, err
	}
	lastInsertID, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }
	user := &structs.User{
        ID:       int64(lastInsertID),
        Username: username,
        Password: hashedPassword,
    }
	return user, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUserByUsername(username string) (*structs.User, error) {
	var user structs.User
	err := db.QueryRow(`SELECT * FROM users WHERE username = ?`, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreatePost(title, content string, userID int64) error {
	now := time.Now()
	_, err := db.Exec(`
        INSERT INTO posts (title, content, user_id, created_at)
        VALUES (?, ?, ?, ?)
    `, title, content, userID, now)
	return err
}

func GetRecentPosts(limit int) ([]structs.Post, error) {
	rows, err := db.Query(`
        SELECT * FROM posts p
        JOIN users u ON p.user_id = u.id
        WHERE p.is_published = true
        ORDER BY p.created_at DESC
        LIMIT ?
    `, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var p structs.Post
		var username, categoryTitle string
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &p.CreatedAt, &username, &categoryTitle)
		if err != nil {
			return nil, err
		}
		// p.Author = &structs.User{Username: username}
		// p.Category = &structs.Category{Title: categoryTitle}
		posts = append(posts, p)
	}
	return posts, nil
}
