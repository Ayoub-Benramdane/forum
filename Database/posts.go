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

func GetAllPosts(status string, size, page int64) ([]structs.Post, error) {
	rows, err := DB.Query("SELECT p.id, p.title, p.content, p.created_at, u.username FROM posts p JOIN users u ON p.user_id = u.id ORDER BY p.created_at DESC LIMIT ? OFFSET ?", size, size*page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Author)
		if err != nil {
			return nil, err
		}
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
		post.Status = status
		posts = append(posts, post)
	}
	return posts, nil
}

func CountPosts() (float64, error) {
	var posts float64
	err := DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&posts)
	if err != nil {
		return 0, err
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
	return post, nil
}
