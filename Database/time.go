package database

import (
	"fmt"
	structs "forum/Data"
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

func SortingPost(Posts []structs.Post) []structs.Post {
	for i := 0; i < len(Posts); i++ {
		for j := i + 1; j < len(Posts); j++ {
			if Posts[j].ID > Posts[i].ID {
				Posts[i], Posts[j] = Posts[j], Posts[i]
			}
		}
	}
	return Posts
}
