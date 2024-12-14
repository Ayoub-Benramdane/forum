package database

func CheckLike(userID, postID int64) bool {
	var likes int64
	DB.QueryRow("SELECT COUNT(*) FROM post_likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&likes)
	return likes > 0
}

func AddLike(userID, postID int64) error {
	if err := DeleteDislike(userID, postID); err != nil {
		return err
	}
	_, err := DB.Exec("INSERT INTO post_likes (user_id, post_id) VALUES (?, ?)", userID, postID)
	return err
}

func DeleteLike(userID, postID int64) error {
	_, err := DB.Exec("DELETE FROM post_likes WHERE user_id = ? AND post_id = ?", userID, postID)
	return err
}

func CountLikes(postID int64) (int64, error) {
	var likes int64
	err := DB.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ?", postID).Scan(&likes)
	if err != nil {
		return 0, err
	}
	return likes, nil
}

func CheckDislike(userID, postID int64) bool {
	var likes int64
	DB.QueryRow("SELECT COUNT(*) FROM post_dislikes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&likes)
	return likes > 0
}

func AddDislike(userID, postID int64) error {
	if err := DeleteLike(userID, postID); err != nil {
		return err
	}
	_, err := DB.Exec("INSERT INTO post_dislikes (user_id, post_id) VALUES (?, ?)", userID, postID)
	return err
}

func DeleteDislike(userID, postID int64) error {
	_, err := DB.Exec("DELETE FROM post_dislikes WHERE user_id = ? AND post_id = ?", userID, postID)
	return err
}

func CountDislikes(postID int64) (int64, error) {
	var likes int64
	err := DB.QueryRow("SELECT COUNT(*) FROM post_dislikes WHERE post_id = ?", postID).Scan(&likes)
	if err != nil {
		return 0, err
	}
	return likes, nil
}
