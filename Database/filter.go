package database

import (
	structs "forum/Structs"
)

func GetAllCategorys() ([]structs.Category, error) {
	rows, err := DB.Query(`SELECT * FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categorys []structs.Category
	for rows.Next() {
		var c structs.Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categorys = append(categorys, c)
	}
	return categorys, nil
}

func GetFilterPosts(statut, Category string) ([]structs.Post, error) {
	rows, err := DB.Query(`
        SELECT p.id, p.title, p.content, p.user_id, p.created_at, u.username
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
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.UserID, &p.CreatedAt, &p.Author)
		if err != nil {
			return nil, err
		}
		like, err := CountLikes(p.ID)
		if err != nil {
			return nil, err
		}
		dislike, err := CountDislikes(p.ID)
		if err != nil {
			return nil, err
		}
		p.TotalLike = like
		p.TotalDislike = dislike
		p.Statut = statut
		posts = append(posts, p)
	}
	return posts, nil
}
