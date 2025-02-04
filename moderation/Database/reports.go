package database

import (
	structs "forum/Data"
	"time"
)

func InsertReport(PostID, UserID int64, Descreption string) error {
	_, err := DB.Exec("INSERT INTO reports (user_id, post_id, description, created_at) VALUES (?, ?, ?, ?)", UserID, PostID, Descreption, time.Now())
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
