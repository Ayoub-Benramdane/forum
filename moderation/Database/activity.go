package database

import (
	structs "forum/Data"
	"time"
)

func GetData(Id int64) (*structs.Activity, error) {
	var activities structs.Activity
	posts, err := PostsActivity(Id)
	if err != nil {
		return nil, err
	}
	comments, err := CommentsActivity(Id)
	if err != nil {
		return nil, err
	}
	reacted_posts, err := ReactPostsActivity(Id)
	if err != nil {
		return nil, err
	}
	reacted_comments, err := ReactCommentsActivity(Id)
	if err != nil {
		return nil, err
	}
	activities.Posts = posts
	activities.Comments = comments
	activities.ReactedPosts = reacted_posts
	activities.ReactedComments = reacted_comments
	activities.TotalPosts = int64(len(activities.Posts))
	activities.TotalComments = int64(len(activities.Comments))
	activities.TotalLikes = int64(len(activities.ReactedPosts) + len(activities.ReactedComments))
	return &activities, err
}

func PostsActivity(id int64) ([]structs.Post, error) {
	rows, err := DB.Query("SELECT p.id, p.title, p.content, p.created_at, u.username FROM posts p JOIN users u ON p.user_id = u.id  WHERE u.id == ? ORDER BY p.created_at DESC", id)
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
	return posts, nil
}

func CommentsActivity(id int64) ([]structs.Comment, error) {
	rows, err := DB.Query("SELECT c.id, c.post_id, c.content, c.created_at, u.username FROM comments c JOIN users u ON c.user_id = u.id WHERE c.user_id = ? ORDER BY c.created_at DESC", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []structs.Comment
	var date time.Time
	for rows.Next() {
		var comment structs.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.Content, &date, &comment.Author); err != nil {
			return nil, err
		}
		comment.CreatedAt = TimeAgo(date)
		comment.TotalLikes, err = CountLikesComment(comment.PostID, comment.ID)
		if err != nil {
			return nil, err
		}
		comment.TotalDislikes, err = CountDislikesComment(comment.PostID, comment.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func ReactPostsActivity(id int64) ([]structs.Post, error) {
	rows, err := DB.Query(`SELECT p.id, p.title, p.content, p.created_at, u.username FROM post_reactions r JOIN posts p ON r.post_id = p.id JOIN users u ON p.user_id = u.id WHERE r.user_id = ? ORDER BY p.created_at DESC`, id)
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
	return posts, nil
}

func ReactCommentsActivity(id int64) ([]structs.Comment, error) {
	rows, err := DB.Query(`SELECT c.id, c.post_id, c.content, c.created_at, u.username FROM comment_reactions r JOIN comments c ON r.comment_id = c.id JOIN users u ON c.user_id = u.id WHERE r.user_id = ? ORDER BY c.created_at DESC`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []structs.Comment
	var date time.Time
	for rows.Next() {
		var comment structs.Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.Content, &date, &comment.Author); err != nil {
			return nil, err
		}
		comment.CreatedAt = TimeAgo(date)
		comment.TotalLikes, err = CountLikesComment(comment.PostID, comment.ID)
		if err != nil {
			return nil, err
		}
		comment.TotalDislikes, err = CountDislikesComment(comment.PostID, comment.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func GetActivities() ([]structs.RecentActivities, error) {
	rows, err := DB.Query(`SELECT 'post' AS source, title AS name, created_at FROM posts UNION ALL SELECT 'report' AS source, description AS name, created_at FROM reports UNION ALL SELECT 'user' AS source, username AS name, created_at FROM users ORDER BY created_at DESC;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var activities []structs.RecentActivities
	for rows.Next() {
		var activity structs.RecentActivities
		if err := rows.Scan(&activity.Type, &activity.Description, &activity.CreatedAt); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}
