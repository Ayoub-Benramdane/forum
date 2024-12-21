package database

import (
	structs "forum/Data"
)

func DeleteCommentId(PostID, CommentID int64) error {
	_, err := DB.Exec("DELETE FROM comments WHERE id = ? AND post_id = ?", CommentID, PostID)
	return err
}

func GetCommentByID(id_post, id_comment int64) (*structs.Comment, error) {
	comment := &structs.Comment{}
	err := DB.QueryRow("SELECT content FROM comments WHERE id == ? AND post_id = ?", id_comment, id_post).Scan(&comment.Content)
	comment.ID = id_comment
	comment.PostID = id_post
	return comment, err
}

func UpdateComment(content string, id_comment, id_post int64) error {
	_, err := DB.Exec("UPDATE comments SET content = ? WHERE id = ? AND post_id = ?", content, id_comment, id_post)
	return err
}
