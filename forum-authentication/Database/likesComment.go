package database

import "fmt"

var like = "like"
var dislike = "dislike"

func CheckLikeComment(userID, postID, commentID int64) bool {
	var likes int64
	DB.QueryRow("SELECT COUNT(*) FROM comment_reactions WHERE user_id = ? AND post_id = ? AND comment_id = ? AND type = ?", userID, postID, commentID, like).Scan(&likes)
	return likes > 0
}

func AddLikeComment(userID, postID, commentID int64) error {
	if userID == 0 {
		return fmt.Errorf("session closed")
	}
	if err := DeleteDislikeComment(userID, postID, commentID); err != nil {
		return err
	}
	_, err := DB.Exec("INSERT INTO comment_reactions (user_id, post_id, comment_id, type) VALUES (?, ?, ?, ?)", userID, postID, commentID, like)
	return err
}

func DeleteLikeComment(userID, postID, commentID int64) error {
	if userID == 0 {
		return fmt.Errorf("session closed")
	}
	_, err := DB.Exec("DELETE FROM comment_reactions WHERE user_id = ? AND post_id = ? AND comment_id = ? AND type = ?", userID, postID, commentID, like)
	return err
}

func CountLikesComment(postID, commentID int64) (int64, error) {
	var likes int64
	err := DB.QueryRow("SELECT COUNT(*) FROM comment_reactions WHERE post_id = ? AND comment_id = ? AND type = ?", postID, commentID, like).Scan(&likes)
	return likes, err
}

func CheckDislikeComment(userID, postID, commentID int64) bool {
	var dislikes int64
	DB.QueryRow("SELECT COUNT(*) FROM comment_reactions WHERE user_id = ? AND post_id = ? AND comment_id = ? AND type = ?", userID, postID, commentID, dislike).Scan(&dislikes)
	return dislikes > 0
}

func AddDislikeComment(userID, postID, commentID int64) error {
	if userID == 0 {
		return fmt.Errorf("session closed")
	}
	if err := DeleteLikeComment(userID, postID, commentID); err != nil {
		return err
	}
	_, err := DB.Exec("INSERT INTO comment_reactions (user_id, post_id, comment_id, type) VALUES (?, ?, ?, ?)", userID, postID, commentID, dislike)
	return err
}

func DeleteDislikeComment(userID, postID, commentID int64) error {
	if userID == 0 {
		return fmt.Errorf("session closed")
	}
	_, err := DB.Exec("DELETE FROM comment_reactions WHERE user_id = ? AND post_id = ? AND comment_id = ? AND type = ?", userID, postID, commentID, dislike)
	return err
}

func CountDislikesComment(postID, commentID int64) (int64, error) {
	var likes int64
	err := DB.QueryRow("SELECT COUNT(*) FROM comment_reactions WHERE post_id = ? AND comment_id = ? AND type = ?", postID, commentID, dislike).Scan(&likes)
	return likes, err
}
