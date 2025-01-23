package database

import (
	"database/sql"

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
