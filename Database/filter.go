package database

import (
	"database/sql"
	structs "forum/Structs"
)

func GetFilterPosts(user *structs.Session, categories []string) ([]structs.Post, error) {
	var posts []structs.Post
	var rows *sql.Rows
	var err error
	for _, Category := range categories {
		switch Category {
		case "All": return GetAllPosts(user.Status)
		case "MyPosts": rows, err = SelectPost(user.UserID)
		case "MyLikes": rows, err = SelectLike(user.UserID)
		default: rows, err = SelectCategory(Category)
		}
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var p structs.Post
			err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.Author)
			if err != nil {
				return nil, err
			}
			likes, err := CountLikes(p.ID)
			if err != nil {
				return nil, err
			}
			dislikes, err := CountDislikes(p.ID)
			if err != nil {
				return nil, err
			}
			comments, err := CountComments(p.ID)
			if err != nil {
				return nil, err
			}
			categories, err := GetCategories(p.ID)
			if err != nil {
				return nil, err
			}
			p.TotalLikes = likes
			p.TotalDislikes = dislikes
			p.TotalComments = comments
			p.Categories = categories
			p.Status = user.Status
			if NotExist(p.ID, posts) {
				posts = append(posts, p)
			}
		}
	}
	return SortingPost(posts), nil
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
			JOIN post_likes l ON l.post_id = p.id
			WHERE l.user_id = ?
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

func SortingPost(Posts []structs.Post) []structs.Post {
	for i := 0; i < len(Posts); i++ {
		for j := i+1; j < len(Posts); j++ {
			if Posts[j].CreatedAt.Before(Posts[i].CreatedAt) {
				Posts[i], Posts[j] = Posts[j],Posts[i]
			}
		}
	}
	return Posts
}
