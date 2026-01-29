package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func Connect(url string) *DB {
	db, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	return &DB{db}
}
