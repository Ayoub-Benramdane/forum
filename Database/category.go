package database

import structs "forum/Structs"

func GetAllCategorys() ([]structs.Category, error) {
	rows, err := DB.Query(`SELECT * FROM post_category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categorys []structs.Category
	for rows.Next() {
		var c structs.Category
		err:= rows.Scan(&c.ID, &c.Name, &c.PostID)
		if err != nil {
			return nil, err
		}
		categorys = append(categorys, c)
	}
	return categorys, nil
}