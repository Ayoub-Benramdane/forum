package database

import (
	structs "forum/Data"
	"time"
)

func CreateCategories() error {
	if CheckCategory() == nil {
		categories := []string{"Sport", "General", "Tech", "Gaming", "Movies", "Music", "Health", "Travel", "Food", "Fashion", "Education", "Science", "Art", "Finance", "Lifestyle", "History"}
		for _, category := range categories {
			_, err := DB.Exec("INSERT INTO categories (name, created_at) VALUES (?, ?)", category, time.Now())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateCategory(category string) error {
	_, err := DB.Exec("INSERT INTO categories (name, created_at) VALUES (?, ?)", category, time.Now())
	return err
}

func UpdateCategory(id int64, category string) error {
	_, err := DB.Exec("UPDATE categories SET name = ? WHERE id = ?", category, id)
	return err
}

func CheckCategory() *structs.Category {
	var cat structs.Category
	err := DB.QueryRow("SELECT name FROM categories").Scan(&cat.Name)
	if err != nil {
		return nil
	}
	return &cat
}

func GetAllCategories() ([]structs.Category, error) {
	rows, err := DB.Query("SELECT id, name, created_at FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []structs.Category
	for rows.Next() {
		var category structs.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func GetCategories(postID int64) ([]string, error) {
	rows, err := DB.Query("SELECT c.id, c.name FROM categories c JOIN post_category cp ON c.id = cp.category_id WHERE cp.post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []string
	for rows.Next() {
		var category structs.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category.Name)
	}
	return categories, nil
}

func PostsCategory(categories *[]structs.Category) error {
	for i, category := range *categories {
		var count int64
		err := DB.QueryRow("SELECT COUNT(*) FROM posts p JOIN post_category pc ON pc.post_id = p.id JOIN categories c ON c.id = pc.category_id WHERE c.name = ?", category.Name).Scan(&count)
		if err != nil {
			return err
		}
		(*categories)[i].PostCount = count
	}
	return nil
}
