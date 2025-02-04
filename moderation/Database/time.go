package database

import (
	"fmt"
	"time"
)

func TimeAgo(date time.Time) string {
	timeAgo := time.Since(date)
	if timeAgo.Minutes() < 1 {
		return "Just now"
	} else if timeAgo.Minutes() < 60 {
		return fmt.Sprintf("%d minutes ago", int(timeAgo.Minutes()))
	} else if timeAgo.Minutes() < 60*24 {
		return fmt.Sprintf("%d hours ago", int(timeAgo.Hours()))
	}
	return fmt.Sprintf("%d days ago", int(timeAgo.Hours())/24)
}
