package structs

import (
	"time"
)

type User struct {
	ID       int64  `sqlite:"id" json:"id"`
	Username string `sqlite:"username" json:"username"`
	Email    string `sqlite:"email" json:"email"`
	Password string `sqlite:"password" json:"-"`
	CreatedAt    time.Time `sqlite:"created_at" json:"created_at"`
}

type Session struct {
	ID       int64  `sqlite:"id" json:"id"`
	Username string `sqlite:"username" json:"username"`
	UserID   int64  `sqlite:"user_id" json:"user_id"`
	Statut   string `sqlite:"statut" json:"statut"`
}

type Post struct {
	ID           int64     `sqlite:"id" json:"id"`
	Title        string    `sqlite:"title" json:"title"`
	Content      string    `sqlite:"content" json:"content"`
	UserID       int64     `sqlite:"user_id" json:"user_id"`
	CreatedAt    time.Time `sqlite:"created_at" json:"created_at"`
	Author       string    `sqlite:"author" json:"author"`
	TotalLike    int64     `sqlite:"total_like" json:"total_like"`
	TotalDislike int64     `sqlite:"total_dislike" json:"total_dislike"`
	Statut       string    `sqlite:"statut" json:"statut"`
}

type Comment struct {
	ID           int64     `sqlite:"id" json:"id"`
	Content      string    `sqlite:"content" json:"content"`
	UserID       int64     `sqlite:"user_id" json:"user_id"`
	PostID       int64     `sqlite:"post_id" json:"post_id"`
	CreatedAt    time.Time `sqlite:"created_at" json:"created_at"`
	Author       string    `sqlite:"author" json:"author"`
	TotalLike    int64     `sqlite:"total_like" json:"total_like"`
	TotalDislike int64     `sqlite:"total_dislike" json:"total_dislike"`
	Statut       string    `sqlite:"statut" json:"statut"`
}

type Category struct {
	ID   int64  `sqlite:"id" json:"id"`
	Name string `sqlite:"name" json:"name"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
