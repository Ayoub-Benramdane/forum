package database

import (
	"time"

	structs "forum/Data"
)

func CreateComment(content string, userID, postID int64) error {
	_, err := DB.Exec("INSERT INTO comments (content, user_id, post_id, created_at) VALUES (?, ?, ?, ?)", content, userID, postID, time.Now())
	return err
}

func GetAllComments(PostID int64, status string) ([]structs.Comment, error) {
	rows, err := DB.Query("SELECT c.id, c.user_id, c.content, c.created_at, u.username FROM comments c JOIN users u ON c.user_id = u.id WHERE c.post_id = ? ORDER BY c.created_at DESC", PostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []structs.Comment
	for rows.Next() {
		var comment structs.Comment
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.Author)
		if err != nil {
			return nil, err
		}
		likes, err := CountLikesComment(PostID, comment.ID)
		if err != nil {
			return nil, err
		}
		dislikes, err := CountDislikesComment(PostID, comment.ID)
		if err != nil {
			return nil, err
		}
		comment.PostID = PostID
		comment.TotalLikes = likes
		comment.TotalDislikes = dislikes
		comment.Status = status
		comments = append(comments, comment)
	}
	return comments, nil
}

func CountComments(postID int64) (int64, error) {
	var comments int64
	err := DB.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", postID).Scan(&comments)
	if err != nil {
		return 0, err
	}
	return comments, nil
}
