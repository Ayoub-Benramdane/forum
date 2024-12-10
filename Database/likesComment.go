package database

func CheckLikeComment(userID, postID, commentID int64) bool {
	var likes int64
	DB.QueryRow(`
        SELECT COUNT(*) FROM comment_likes WHERE user_id = ? AND post_id = ? AND comment_id = ?
    `, userID, postID, commentID).Scan(&likes)
	return likes > 0
}

func CheckDislikeComment(userID, postID, commentID int64) bool {
	var likes int64
	DB.QueryRow(`
        SELECT COUNT(*) FROM comment_dislikes WHERE user_id = ? AND post_id = ? AND comment_id = ?
    `, userID, postID, commentID).Scan(&likes)
	return likes > 0
}

func AddLikeComment(userID, postID, commentID int64) error {
	DeleteDislikeComment(userID, postID, commentID)
	_, err := DB.Exec(`
        INSERT INTO comment_likes (user_id, post_id, comment_id)
        VALUES (?, ?, ?)
    `, userID, postID, commentID)
	return err
}

func DeleteLikeComment(userID, postID, commentID int64) error {
	_, err := DB.Exec(`
        DELETE FROM comment_likes WHERE user_id = ? AND post_id = ? AND comment_id = ?
    `, userID, postID, commentID)
	return err
}

func AddDislikeComment(userID, postID, commentID int64) error {
	DeleteLikeComment(userID, postID, commentID)
	_, err := DB.Exec(`
        INSERT INTO comment_dislikes (user_id, post_id, comment_id)
        VALUES (?, ?, ?)
    `, userID, postID, commentID)
	return err
}

func DeleteDislikeComment(userID, postID, commentID int64) error {
	_, err := DB.Exec(`
        DELETE FROM comment_dislikes WHERE user_id = ? AND post_id = ? AND comment_id = ?
    `, userID, postID, commentID)
	return err
}

func CountLikesComment(postID, commentID int64) (int64, error) {
	var likes int64
	err := DB.QueryRow(`
        SELECT COUNT(*) FROM comment_likes WHERE post_id = ? AND comment_id = ?
    `, postID, commentID).Scan(&likes)
	if err != nil {
		return 0, err
	}
	return likes, nil
}

func CountDislikesComment(postID, commentID int64) (int64, error) {
	var likes int64
	err := DB.QueryRow(`
        SELECT COUNT(*) FROM comment_dislikes WHERE post_id = ? AND comment_id = ?
    `, postID, commentID).Scan(&likes)
	if err != nil {
		return 0, err
	}
	return likes, nil
}
