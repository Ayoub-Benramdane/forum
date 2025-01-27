package database

import (
	"fmt"
	"strings"
	"time"

	structs "forum/Data"
)

func CreateNotification(content, Type string, user_id, post_id, comment_id int64, title, author string) error {
	if comment_id == -1 {
		_, err := DB.Exec("INSERT INTO notifications (content, user_id, post_id, title, type, notif_by, created_at, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", content, user_id, post_id, title, Type, author, time.Now(), "Unread")
		return err
	}
	_, err := DB.Exec("INSERT INTO notifications (content, user_id, post_id, title, type, notif_by, created_at, status, comment_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", content, user_id, post_id, title, Type, author, time.Now(), "Unread", comment_id)
	fmt.Println(user_id, post_id, comment_id, err)
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
		if err := rows.Scan(&notification.ID, &notification.Content, &notification.UserID, &notification.PostID, &notification.Title, &notification.Type, &notification.Author, &date, &notification.Status, &notification.CommentID); err != nil {
			if strings.Contains(err.Error(), "converting NULL to int64") {
				notification.CommentID = -1
			} else {
				return nil, err
			}
		}
		notification.CreatedAt = TimeAgo(date)
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func DeleteNotification(content string, post_id, comment_id int64, author string) error {
	if comment_id == -1 {
		_, err := DB.Exec("DELETE FROM notifications WHERE content = ? AND post_id = ? AND notif_by = ?", content, post_id, author)
		return err
	}
	_, err := DB.Exec("DELETE FROM notifications WHERE content = ? AND post_id = ? AND comment_id = ? AND notif_by = ?", content, post_id, comment_id, author)
	return err
}
