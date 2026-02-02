package db

import (
	"database/sql"
	"server/schema"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
	*schema.Queries
}

func Connect(url string) *DB {
	db, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	queries := schema.New(db)

	return &DB{db, queries}
}
