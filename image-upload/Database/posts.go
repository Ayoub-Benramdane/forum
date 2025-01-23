package database

import (
	"fmt"
	"time"

	structs "forum/Data"
)

func CreatePost(title, content string, categories []string, userID int64, filePath string) error {
	if userID == 0 {
		return fmt.Errorf("session closed")
	}

	result, err := DB.Exec("INSERT INTO posts (title, content, user_id, created_at, file_path) VALUES (?, ?, ?, ?, ?)",
		title, content, userID, time.Now(), filePath)
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

func GetAllPosts() ([]structs.Post, error) {
	rows, err := DB.Query("SELECT p.id, p.title, p.content, p.file_path, p.created_at, u.username FROM posts p JOIN users u ON p.user_id = u.id ORDER BY p.created_at")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		var date time.Time
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.FilePath, &date, &post.Author); err != nil {
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
		posts = append(posts, post)
	}
	SortingPost(posts)
	*Posts = posts
	if len(posts) > 10 {
		return posts[:10], nil
	}
	return posts, nil
}

func GetPostByID(id int64) (*structs.Post, error) {
	post := &structs.Post{}
	var date time.Time
	err := DB.QueryRow("SELECT p.id, p.title, p.user_id, p.content, p.file_path, p.created_at, u.username FROM posts p JOIN users u ON p.user_id = u.id WHERE p.id == ?",
		id).Scan(&post.ID, &post.Title, &post.UserID, &post.Content, &post.FilePath, &date, &post.Author)
	if err != nil {
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
	return post, err
}

func CountPosts() (float64, error) {
	var posts float64
	err := DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&posts)
	return posts, err
}

func CountPostsByCat() (float64, error) {
	var posts float64
	err := DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&posts)
	return posts, err
}
