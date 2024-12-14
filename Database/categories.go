package database

import structs "forum/Structs"

func GetAllCategorys() ([]structs.Category, error) {
	rows, err := DB.Query("SELECT name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []structs.Category
	for rows.Next() {
		var category structs.Category
		err := rows.Scan(&category.Name)
		if err != nil {
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
		err := rows.Scan(&category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category.Name)
	}
	return categories, nil
}
