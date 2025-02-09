package database

func DeletePostId(PostID int64) error {
	_, err := DB.Exec("DELETE FROM posts WHERE id = ?", PostID)
	return err
}

func UpdatePost(title, content string, categories []string, PostID int64, status string) error {
	if status != "" {
		_, err := DB.Exec("UPDATE posts SET status = ? WHERE id = ?", status, PostID)
		return err
	}
	_, err := DB.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ?", title, content, PostID)
	if err != nil {
		return err
	}
	if categories != nil {
		_, err = DB.Exec("DELETE FROM post_category WHERE post_id = ?", PostID)
		if err != nil {
			return err
		}
		for _, category := range categories {
			var catID int64
			err = DB.QueryRow("SELECT id FROM categories WHERE name = ?", category).Scan(&catID)
			if err != nil {
				return err
			}
			_, err = DB.Exec("INSERT INTO post_category (category_id, post_id) VALUES (?, ?)", catID, PostID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
