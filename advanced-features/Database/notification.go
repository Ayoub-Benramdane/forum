package database

import (
	"fmt"
	structs "forum/Data"
	"time"
)

func CreateNotification(content, Type string, user_id, post_id, comment_id int64, title, author string) error {
	_, err := DB.Exec("INSERT INTO notifications (content, user_id, post_id, comment_id, title, type, notif_by, created_at, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", content, user_id, post_id, comment_id, title, Type, author, time.Now(), "Unread")
	fmt.Println(err)
	return err
}

func GetNotification(id int64) ([]structs.Notification, error) {
	var notifications []structs.Notification
	rows, err := DB.Query("SELECT * FROM notifications WHERE user_id = ? ORDER BY created_at DESC", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var date time.Time
	for rows.Next() {
		var notification structs.Notification
		if rows.Scan(&notification.ID, &notification.Content, &notification.UserID, &notification.PostID, &notification.CommentID, &notification.Title, &notification.Type, &notification.Author, &date, &notification.Status) != nil {
			return nil, err
		}
		notification.CreatedAt = TimeAgo(date)
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func DeleteNotification(content string, post_id, comment_id int64, author string) error {
	_, err := DB.Exec("DELETE FROM notifications WHERE content = ? AND post_id = ? AND comment_id = ? AND notif_by = ?", content, post_id, comment_id, author)
	fmt.Println(err)
	return err
}
