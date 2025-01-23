package database

import (
	"time"

	structs "forum/Data"
)

func CreateComment(content string, userID, postID int64) (int64, error) {
	result, errInsert := DB.Exec("INSERT INTO comments (content, user_id, post_id, created_at) VALUES (?, ?, ?, ?)", content, userID, postID, time.Now())
	if errInsert != nil {
		return 0, errInsert
	}
	lastID, err := result.LastInsertId()
	return lastID, err
}

func GetComment(commentID int64) (int64, error) {
	var userID int64
	err := DB.QueryRow("SELECT user_id FROM comments WHERE id = ?", commentID).Scan(&userID)
	return userID, err
}

func GetAllComments(PostID int64) ([]structs.Comment, error) {
	rows, err := DB.Query("SELECT c.id, c.user_id, c.content, c.created_at, u.username FROM comments c JOIN users u ON c.user_id = u.id WHERE c.post_id = ? ORDER BY c.created_at DESC", PostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []structs.Comment
	var date time.Time
	for rows.Next() {
		var comment structs.Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.Content, &date, &comment.Author); err != nil {
			return nil, err
		}
		comment.CreatedAt = TimeAgo(date)
		comment.TotalLikes, err = CountLikesComment(PostID, comment.ID)
		if err != nil {
			return nil, err
		}
		comment.TotalDislikes, err = CountDislikesComment(PostID, comment.ID)
		if err != nil {
			return nil, err
		}
		comment.PostID = PostID
		comments = append(comments, comment)
	}
	return comments, nil
}

func CountComments(postID int64) (int64, error) {
	var comments int64
	err := DB.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", postID).Scan(&comments)
	return comments, err
}
