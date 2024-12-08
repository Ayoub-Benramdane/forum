package structs

import (
	"time"
)

type User struct {
	ID       int64  `sqlite:"id" json:"id"`
	Username string `sqlite:"username" json:"username"`
	Password string `sqlite:"password" json:"-"`
}

type Post struct {
	ID        int64     `sqlite:"id" json:"id"`
	Title     string    `sqlite:"title" json:"title"`
	Content   string    `sqlite:"content" json:"content"`
	UserID    int64     `sqlite:"user_id" json:"user_id"`
	CreatedAt time.Time `sqlite:"created_at" json:"created_at"`
}

type Comment struct {
	ID        int64     `sqlite:"id" json:"id"`
	Content   string    `sqlite:"content" json:"content"`
	UserID    int64     `sqlite:"user_id" json:"user_id"`
	PostID    int64     `sqlite:"post_id" json:"post_id"`
	CreatedAt time.Time `sqlite:"created_at" json:"created_at"`
}

type PostLike struct {
	ID        int64     `sqlite:"id" json:"id"`
	UserID    int64     `sqlite:"user_id" json:"user_id"`
	PostID    int64     `sqlite:"post_id" json:"post_id"`
	Like      int64     `sqlite:"like" json:"like"`
	CreatedAt time.Time `sqlite:"created_at" json:"created_at"`
}

type CommentLike struct {
	ID        int64     `sqlite:"id" json:"id"`
	UserID    int64     `sqlite:"user_id" json:"user_id"`
	CommentID int64     `sqlite:"comment_id" json:"comment_id"`
	Like      int64     `sqlite:"like" json:"like"`
	CreatedAt time.Time `sqlite:"created_at" json:"created_at"`
}

type Tag struct {
	ID   int64  `sqlite:"id" json:"id"`
	Name string `sqlite:"name" json:"name"`
}

type PostTag struct {
	PostID    int64     `sqlite:"post_id" json:"post_id"`
	TagID     int64     `sqlite:"tag_id" json:"tag_id"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
