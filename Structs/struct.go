package structs

import (
    "time"
)

type User struct {
    ID        int64     `sqlite:"id" json:"id"`
    Username  string    `sqlite:"username" json:"username"`
    Password  string    `sqlite:"password" json:"-"`
    CreatedAt time.Time `sqlite:"created_at" json:"created_at"`
    UpdatedAt time.Time `sqlite:"updated_at" json:"updated_at"`
    LastLogin time.Time `sqlite:"last_login" json:"last_login"`
    IsActive  bool      `sqlite:"is_active" json:"is_active"`
}

type Category struct {
    ID          int64     `sqlite:"id" json:"id"`
    Title       string    `sqlite:"title" json:"title"`
    Description string    `sqlite:"description" json:"description"`
    Slug        string    `sqlite:"slug" json:"slug"`
    CreatedAt   time.Time `sqlite:"created_at" json:"created_at"`
    UpdatedAt   time.Time `sqlite:"updated_at" json:"updated_at"`
}

type Post struct {
    ID          int64     `sqlite:"id" json:"id"`
    Title       string    `sqlite:"title" json:"title"`
    Content     string    `sqlite:"content" json:"content"`
    UserID      int64     `sqlite:"user_id" json:"user_id"`
    CategoryID  int64     `sqlite:"category_id" json:"category_id"`
    CreatedAt   time.Time `sqlite:"created_at" json:"created_at"`
    UpdatedAt   time.Time `sqlite:"updated_at" json:"updated_at"`
    Views       int       `sqlite:"views" json:"views"`
    IsPublished bool      `sqlite:"is_published" json:"is_published"`
    Author    *User     `sqlite:"-" json:"author,omitempty"`
    Category  *Category `sqlite:"-" json:"category,omitempty"`
}

type Comment struct {
    ID        int64     `sqlite:"id" json:"id"`
    Content   string    `sqlite:"content" json:"content"`
    UserID    int64     `sqlite:"user_id" json:"user_id"`
    PostID    int64     `sqlite:"post_id" json:"post_id"`
    CreatedAt time.Time `sqlite:"created_at" json:"created_at"`
    UpdatedAt time.Time `sqlite:"updated_at" json:"updated_at"`
    IsEdited  bool      `sqlite:"is_edited" json:"is_edited"`
    Author    *User `sqlite:"-" json:"author,omitempty"`
}

type Error struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}