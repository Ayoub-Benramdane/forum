package database

import (
	structs "forum/Data"
	"time"
)

func InsertReport(PostID, UserID int64, Descreption string) error {
	_, err := DB.Exec("INSERT INTO reports (user_id, post_id, description, type, created_at) VALUES (?, ?, ?, ?, ?)", UserID, PostID, Descreption, "Pending", time.Now())
	return err
}

func GetAllReports() ([]structs.Reports, error) {
	rows, err := DB.Query("SELECT r.id, r.post_id, r.description, r.created_at, u.username, p.title FROM reports r JOIN users u on u.id = r.user_id JOIN posts p ON p.id = r.post_id ORDER BY r.created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reports []structs.Reports
	for rows.Next() {
		var report structs.Reports
		if err := rows.Scan(&report.ID, &report.PostID, &report.Description, &report.ReportedAt, &report.ReportedBy, &report.Title); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func GetReportsType(Type string) ([]structs.Reports, error) {
	rows, err := DB.Query("SELECT r.id, r.post_id, r.description, r.created_at, u.username, p.title FROM reports r JOIN users u on u.id = r.user_id JOIN posts p ON p.id = r.post_id, WHERE type = ? ORDER BY r.created_at DESC", Type)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reports []structs.Reports
	for rows.Next() {
		var report structs.Reports
		if err := rows.Scan(&report.ID, &report.PostID, &report.Description, &report.ReportedAt, &report.ReportedBy, &report.Title); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func GetPostsReported() ([]structs.Post, error) {
	rows, err := DB.Query("SELECT p.id, p.title, p.content, p.created_at, u.username FROM posts p JOIN users u ON p.user_id = u.id JOIN reports r ON p.id = r.post_id ORDER BY p.created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		var date time.Time
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &date, &post.Author); err != nil {
			return nil, err
		}
		post.CreatedAt = TimeAgo(date)
		post.TotalLikes, err = CountLikes(post.ID)
		if err != nil {
			return nil, err
		}
		post.TotalDislikes, err = CountDislikes(post.ID)
		if err != nil {
			return nil, err
		}
		post.TotalComments, err = CountComments(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories, err = GetCategories(post.ID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	*Posts = posts
	return posts, nil
}
