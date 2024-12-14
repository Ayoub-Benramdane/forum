package database

import (
	"database/sql"
	structs "forum/Data"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDatabase() error {
	var err error
	DB, err = sql.Open("sqlite3", "./Data/forum.DB")
	if err != nil {
		return err
	} else if err = CreateTables(); err != nil {
		return err
	} else if err = CreateCategoryies(); err != nil {
		return err
	}
	return nil
}

func CreateCategoryies() error {
	if cat := CheckCategory(); cat == nil {
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
