package database

import (
	"fmt"
	structs "forum/Data"
	"time"
)

func GetInfoUser(UserID int64) (*structs.User, error) {
	var user structs.User
	err := DB.QueryRow("SELECT * FROM users WHERE id = ?", UserID).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	user.Posts, err = CountPostsUser(UserID)
	if err != nil {
		return nil, err
	}
	user.Comments, err = CountCommentsUser(UserID)
	if err != nil {
		return nil, err
	}
	user.Likes, user.Dislikes, err = CountLikesUser(UserID)
	if err != nil {
		return nil, err
	}
	user.RecentActivity, err = LastPost(UserID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CountPostsUser(UserID int64) (int64, error) {
	var posts int64
	err := DB.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", UserID).Scan(&posts)
	if err != nil {
		return 0, err
	}
	return posts, nil
}

func CountCommentsUser(UserID int64) (int64, error) {
	var comments int64
	err := DB.QueryRow("SELECT COUNT(*) FROM comments WHERE user_id = ?", UserID).Scan(&comments)
	if err != nil {
		return 0, err
	}
	return comments, nil
}

func CountLikesUser(UserID int64) (int64, int64, error) {
	var likes, dislikes int64
	err := DB.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE user_id = ? AND type = ?", UserID, "like").Scan(&likes)
	if err != nil {
		return 0, 0, err
	}
	err = DB.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE user_id = ? AND type = ?", UserID, "dislike").Scan(&dislikes)
	if err != nil {
		return 0, 0, err
	}
	return likes, dislikes, nil
}

func LastPost(UserID int64) (*structs.Post, error) {
	rows, err := DB.Query("SELECT title, content, created_at FROM posts WHERE user_id = ? ORDER BY created_at DESC", UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []structs.Post
	var date time.Time
	for rows.Next() {
		var post structs.Post
		err := rows.Scan(&post.Title, &post.Content, &date)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = TimeAgo(date)
		posts = append(posts, post)
	}
	if posts != nil {
		return &posts[0], nil
	}
	return nil, err
}

func UpdateInfo(userID int64, username, email string) error {
	if userID == 0 {
		return fmt.Errorf("session closed")
	}
	_, err := DB.Exec("UPDATE users SET username = ?, email = ? WHERE id = ?", username, email, userID)
	if err != nil {
		return err
	}
	_, err = DB.Exec("UPDATE session SET username = ? WHERE user_id = ?", username, userID)
	return err
}

func UpdatePass(userID int64, password string) error {
	if userID == 0 {
		return fmt.Errorf("session closed")
	}
	_, err := DB.Exec("UPDATE users SET  password = ? WHERE id = ?", password, userID)
	return err
}
