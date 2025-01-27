package structs

import (
	"time"
)

type User struct {
	ID             int64     `sqlite:"id" json:"id"`
	Username       string    `sqlite:"username" json:"username"`
	Email          string    `sqlite:"email" json:"email"`
	Password       string    `sqlite:"password" json:"-"`
	CreatedAt      time.Time `sqlite:"created_at" json:"created_at"`
	Status         string    `sqlite:"status" json:"status"`
	ConnectedAt    time.Time `sqlite:"connected_at" json:"connected_at"`
	Posts          int64     `sqlite:"posts" json:"posts"`
	Comments       int64     `sqlite:"comments" json:"comments"`
	Likes          int64     `sqlite:"likes" json:"likes"`
	Dislikes       int64     `sqlite:"dislikes" json:"dislikes"`
	RecentActivity *Post     `sqlite:"recent_activity" json:"recent_activity"`
}

type Post struct {
	ID            int64    `sqlite:"id" json:"id"`
	Title         string   `sqlite:"title" json:"title"`
	Content       string   `sqlite:"content" json:"content"`
	UserID        int64    `sqlite:"user_id" json:"user_id"`
	CreatedAt     string   `sqlite:"created_at" json:"created_at"`
	Author        string   `sqlite:"author" json:"author"`
	TotalLikes    int64    `sqlite:"total_likes" json:"total_likes"`
	TotalDislikes int64    `sqlite:"total_dislikes" json:"total_dislikes"`
	TotalComments int64    `sqlite:"total_comments" json:"total_comments"`
	Categories    []string `sqlite:"categories" json:"categories"`
}

var PostsShowing []Post

type Comment struct {
	ID            int64  `sqlite:"id" json:"id"`
	Content       string `sqlite:"content" json:"content"`
	UserID        int64  `sqlite:"user_id" json:"user_id"`
	PostID        int64  `sqlite:"post_id" json:"post_id"`
	CreatedAt     string `sqlite:"created_at" json:"created_at"`
	Author        string `sqlite:"author" json:"author"`
	TotalLikes    int64  `sqlite:"total_likes" json:"total_likes"`
	TotalDislikes int64  `sqlite:"total_dislikes" json:"total_dislikes"`
}

type Notification struct {
	ID        int64  `sqlite:"id" json:"id"`
	Content   string `sqlite:"content" json:"content"`
	UserID    int64  `sqlite:"user_id" json:"user_id"`
	PostID    int64  `sqlite:"post_id" json:"post_id"`
	CommentID int64  `sqlite:"comment_id" json:"comment_id"`
	Title     string `sqlite:"title" json:"title"`
	Type      string `sqlite:"type" json:"type"`
	Author    string `sqlite:"author" json:"author"`
	CreatedAt string `sqlite:"created_at" json:"created_at"`
	Status    string `sqlite:"status" json:"status"`
}

type Activity struct {
	Posts           []Post
	ReactedPosts    []Post
	Comments        []Comment
	ReactedComments []Comment
	TotalPosts      int64
	TotalComments   int64
	TotalLikes      int64
}

type Category struct {
	ID   int64  `sqlite:"id" json:"id"`
	Name string `sqlite:"name" json:"name"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Page    string `json:"page"`
	Path    string `json:"path"`
}
