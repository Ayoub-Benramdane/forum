package database

import (
	structs "forum/Data"
	"time"
)

func GetData(Id int64) (*structs.All, error) {
	var all structs.All
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
	all.Posts = posts
	all.Comments = comments
	all.ReactedPosts = reacted_posts
	all.ReactedComments = reacted_comments
	all.TotalPosts = int64(len(all.Posts))
	all.TotalComments = int64(len(all.Comments))
	all.TotalLikes = int64(len(all.ReactedPosts)+len(reacted_comments))
	return &all, err
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

func ReactPostsActivity(id int64) ([]structs.Post, error) {
	rows, err := DB.Query("SELECT p.id, p.title, p.content, p.created_at, u.username FROM posts p JOIN users u ON p.user_id = u.id JOIN post_reactions r ON r.user_id = u.id  WHERE u.id == ? ORDER BY p.created_at DESC", id)
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
	rows, err := DB.Query("SELECT c.id, c.post_id, c.content, c.created_at, u.username FROM comments c JOIN users u ON c.user_id = u.id JOIN comment_reactions r ON r.user_id = u.id WHERE c.user_id = ? ORDER BY c.created_at DESC", id)
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
