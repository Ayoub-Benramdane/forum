package database

import "fmt"

func CheckLike(userID, postID int64) bool {
	var likes int64
	DB.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE user_id = ? AND post_id = ? AND type = ?", userID, postID, like).Scan(&likes)
	return likes > 0
}

func AddLike(userID, postID int64) error {
	if userID == 0 {
		return fmt.Errorf("session closed")
	}
	if err := DeleteDislike(userID, postID); err != nil {
		return err
	}
	_, err := DB.Exec("INSERT INTO post_reactions (user_id, post_id, type) VALUES (?, ?, ?)", userID, postID, like)
	return err
}

func DeleteLike(userID, postID int64) error {
	if userID == 0 {
		return fmt.Errorf("session closed")
	}
	_, err := DB.Exec("DELETE FROM post_reactions WHERE user_id = ? AND post_id = ? AND type = ?", userID, postID, like)
	return err
}

func CountLikes(postID int64) (int64, error) {
	var likes int64
	err := DB.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND type = ?", postID, like).Scan(&likes)
	return likes, err
}

func CheckDislike(userID, postID int64) bool {
	var likes int64
	DB.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE user_id = ? AND post_id = ? AND type = ?", userID, postID, dislike).Scan(&likes)
	return likes > 0
}

func AddDislike(userID, postID int64) error {
	if userID == 0 {
		return fmt.Errorf("session closed")
	}
	if err := DeleteLike(userID, postID); err != nil {
		return err
	}
	_, err := DB.Exec("INSERT INTO post_reactions (user_id, post_id, type) VALUES (?, ?, ?)", userID, postID, dislike)
	return err
}

func DeleteDislike(userID, postID int64) error {
	if userID == 0 {
		return fmt.Errorf("session closed")
	}
	_, err := DB.Exec("DELETE FROM post_reactions WHERE user_id = ? AND post_id = ? AND type = ?", userID, postID, dislike)
	return err
}

func CountDislikes(postID int64) (int64, error) {
	var likes int64
	err := DB.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND type = ?", postID, dislike).Scan(&likes)
	return likes, err
}
