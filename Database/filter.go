package database

import (
	"database/sql"
	structs "forum/Data"
	"math"
)

func GetFilterPosts(user *structs.Session, categories []string) ([]structs.Post, error) {
	var posts []structs.Post
	var rows *sql.Rows
	var err error
	for _, Category := range categories {
		switch Category {
		case "All":
			return GetAllPosts(user.Status, 20, 1)
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
			totalPosts, err := CountPosts()
			if err != nil {
				return nil, err
			}
			for i := int64(1); i <= int64(math.Ceil(totalPosts/20)); i++ {
				post.TotalPosts = append(post.TotalPosts, i)
			}
			post.Status = user.Status
			if NotExist(post.ID, posts) {
				posts = append(posts, post)
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
		for j := i + 1; j < len(Posts); j++ {
			if Posts[j].CreatedAt.Before(Posts[i].CreatedAt) {
				Posts[i], Posts[j] = Posts[j], Posts[i]
			}
		}
	}
	return Posts
}
