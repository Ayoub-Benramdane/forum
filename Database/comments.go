package database

import (
	structs "forum/Structs"
	"time"
)

func CreateComment(content string, userID, postID int64) error {
	_, err := DB.Exec(`
        INSERT INTO comments (content, user_id, post_id, created_at)
        VALUES (?, ?, ?, ?)
    `, content, userID, postID, time.Now())
	return err
}

func GetAllComments(PostID int64, statut string) ([]structs.Comment, error) {
	rows, err := DB.Query(`
        SELECT c.id, c.content, c.user_id, post_id, c.created_at, u.username
        FROM comments c JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
        ORDER BY c.created_at DESC
    `, PostID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []structs.Comment
	for rows.Next() {
		var c structs.Comment
		err := rows.Scan(&c.ID, &c.Content, &c.UserID, &c.PostID, &c.CreatedAt, &c.Author)
		if err != nil {
			return nil, err
		}
		like, err := CountLikesComment(c.PostID, c.ID)
		if err != nil {
			return nil, err
		}
		dislike, err := CountDislikesComment(c.PostID, c.ID)
		if err != nil {
			return nil, err
		}
		c.TotalLike = like
		c.TotalDislike = dislike
		c.Statut = statut
		comments = append(comments, c)
	}
	return comments, nil
}
