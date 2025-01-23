package database

import (
	"database/sql"
	structs "forum/Data"
	"time"
)

var Posts = &structs.PostsShowing

func GetFilterPosts(user *structs.Session, categories []string) ([]structs.Post, error) {
	var posts []structs.Post
	var rows *sql.Rows
	var err error
	for _, Category := range categories {
		switch Category {
		case "All":
			return GetAllPosts()
		case "MyPosts":
			rows, err = SelectPost(user.UserID)
		case "MyLikes":
			rows, err = SelectLike(user.UserID)
		default:
			rows, err = SelectCategory(Category)
		}
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var post structs.Post
			var date time.Time
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &date, &post.Author); err != nil {
				return nil, err
			}
			post.CreatedAt = TimeAgo(date)
			post.TotalLikes, err = CountLikes(post.ID)
			if err != nil {
				return nil, err
			}
			post.TotalDislikes, err = CountDislikes(post.ID)
			if err != nil {
				return nil, err
			}
			post.TotalComments, err = CountComments(post.ID)
			if err != nil {
				return nil, err
			}
			post.Categories, err = GetCategories(post.ID)
			if err != nil {
				return nil, err
			}
			if NotExist(post.ID, posts) {
				posts = append(posts, post)
			}
		}
	}
	SortingPost(posts)
	*Posts = posts
	if len(posts) > 10 {
		return posts[:10], nil
	}
	return posts, nil
}

func SelectCategory(Category string) (*sql.Rows, error) {
	rows, err := DB.Query(`
			SELECT p.id, p.title, p.content, p.created_at, u.username
			FROM posts p JOIN post_category pc ON p.id = pc.post_id
			JOIN users u ON p.user_id = u.id
			JOIN categories c ON c.id = pc.category_id
			WHERE c.name = ?
			ORDER BY p.created_at DESC
		`, Category)
	return rows, err
}

func SelectPost(UserID int64) (*sql.Rows, error) {
	rows, err := DB.Query(`
			SELECT p.id, p.title, p.content, p.created_at, u.username
			FROM posts p JOIN users u ON p.user_id = u.id
			WHERE p.user_id = ?
			ORDER BY p.created_at DESC
		`, UserID)
	return rows, err
}

func SelectLike(UserID int64) (*sql.Rows, error) {
	rows, err := DB.Query(`
			SELECT p.id, p.title, p.content, p.created_at, u.username
			FROM posts p JOIN users u ON p.user_id = u.id
			JOIN post_reactions r ON r.post_id = p.id
			WHERE r.user_id = ?
			ORDER BY p.created_at DESC
		`, UserID)
	return rows, err
}

func NotExist(PostID int64, Posts []structs.Post) bool {
	for _, post := range Posts {
		if post.ID == PostID {
			return false
		}
	}
	return true
}
