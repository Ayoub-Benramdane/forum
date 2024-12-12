package database

import structs "forum/Structs"

func GetAllCategorys() ([]structs.Category, error) {
	rows, err := DB.Query(`SELECT name FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categorys []structs.Category
	for rows.Next() {
		var c structs.Category
		err := rows.Scan(&c.Name)
		if err != nil {
			return nil, err
		}
		categorys = append(categorys, c)
	}
	return categorys, nil
}

func GetCategories(postID int64) ([]string, error) {
	rows, err := DB.Query(`SELECT name FROM categories c
	JOIN post_category cp ON c.id = cp.category_id
	WHERE cp.post_id = ?`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []string
	for rows.Next() {
		var c structs.Category
		err := rows.Scan(&c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c.Name)
	}
	return categories, nil
}
