package database

import (
	structs "forum/Structs"
)

func GetFilterPosts(status, Category string) ([]structs.Post, error) {
	rows, err := DB.Query(`
        SELECT p.title, p.content, p.user_id, p.created_at, u.username
        FROM posts p JOIN post_category pc ON p.id = pc.post_id
		JOIN users u ON p.user_id = u.id
		JOIN categories c ON c.id = pc.category_id
		WHERE c.name = ?
        ORDER BY p.created_at DESC
    `, Category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []structs.Post
	for rows.Next() {
		var p structs.Post
		err := rows.Scan(&p.Title, &p.Content, &p.UserID, &p.CreatedAt, &p.Author)
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
