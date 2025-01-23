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
	Posts          int64     `sqlite:"posts" json:"posts"`
	Comments       int64     `sqlite:"comments" json:"comments"`
	Likes          int64     `sqlite:"likes" json:"likes"`
	Dislikes       int64     `sqlite:"dislikes" json:"dislikes"`
	RecentActivity *Post     `sqlite:"recent_activity" json:"recent_activity"`
}

type Session struct {
	ID       int64  `sqlite:"id" json:"id"`
	Username string `sqlite:"username" json:"username"`
	UserID   int64  `sqlite:"user_id" json:"user_id"`
	Status   string `sqlite:"status" json:"status"`
}

type Post struct {
	ID            int64    `sqlite:"id" json:"id"`
	Title         string   `sqlite:"title" json:"title"`
	Content       string   `sqlite:"content" json:"content"`
	UserID        int64    `sqlite:"user_id" json:"user_id"`
	CreatedAt     string   `sqlite:"created_at" json:"created_at"`
	FilePath	  string   `sqlite:"file_path" json:"file_path"`
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
