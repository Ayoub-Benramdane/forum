package database

import (
	"time"

	structs "forum/Data"
)

func CreatePost(title, content string, categories []string, userID int64) error {
	result, err := DB.Exec("INSERT INTO posts (title, content, user_id, created_at) VALUES (?, ?, ?, ?)", title, content, userID, time.Now())
	if err != nil {
		return err
	}
	postID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	var catID int64
	for _, category := range categories {
		err = DB.QueryRow("SELECT id FROM categories WHERE name = ?", category).Scan(&catID)
		if err != nil {
			return err
		}
		_, err = DB.Exec("INSERT INTO post_category (category_id, post_id) VALUES (?, ?)", catID, postID)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAllPosts(status string) ([]structs.Post, error) {
	rows, err := DB.Query(" SELECT p.id, p.title, p.content, p.created_at, u.username FROM posts p JOIN users u ON p.user_id = u.id ORDER BY p.created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []structs.Post
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
		p.Status = status
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostByID(id int64) (*structs.Post, error) {
	post := &structs.Post{}
	err := DB.QueryRow("SELECT p.id, p.title, p.user_id, p.content, p.created_at, u.username FROM posts p JOIN users u ON p.user_id = u.id WHERE p.id == ?",
	id).Scan(&post.ID, &post.Title, &post.UserID, &post.Content, &post.CreatedAt, &post.Author)
	if err != nil {
		return nil, err
	}
	likes, err := CountLikes(post.ID)
	if err != nil {
		return nil, err
	}
	dislikes, err := CountDislikes(post.ID)
	if err != nil {
		return nil, err
	}
	comments, err := CountComments(post.ID)
	if err != nil {
		return nil, err
	}
	categories, err := GetCategories(post.ID)
	if err != nil {
		return nil, err
	}
	post.TotalLikes = likes
	post.TotalDislikes = dislikes
	post.TotalComments = comments
	post.Categories = categories
	return post, nil
}
