package database

import structs "forum/Data"

func CreateCategoryies() error {
	if CheckCategory() == nil {
		categories := []string{"Sport", "General", "Tech", "Gaming", "Movies", "Music", "Health", "Travel", "Food", "Fashion", "Education", "Science", "Art", "Finance", "Lifestyle", "History"}
		for _, category := range categories {
			_, err := DB.Exec("INSERT INTO categories (name) VALUES (?)", category)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CheckCategory() *structs.Category {
	var cat structs.Category
	err := DB.QueryRow("SELECT * FROM categories").Scan(&cat.ID, &cat.Name)
	if err != nil {
		return nil
	}
	return &cat
}

func GetAllCategorys() ([]structs.Category, error) {
	rows, err := DB.Query("SELECT name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []structs.Category
	for rows.Next() {
		var category structs.Category
		if err := rows.Scan(&category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func GetCategories(postID int64) ([]string, error) {
	rows, err := DB.Query("SELECT name FROM categories c JOIN post_category cp ON c.id = cp.category_id WHERE cp.post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []string
	for rows.Next() {
		var category structs.Category
		if err := rows.Scan(&category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category.Name)
	}
	return categories, nil
}
