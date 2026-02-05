package db

import (
	"database/sql"
	"server/schema"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
	*schema.Queries
}

func Connect(url string) *DB {

	var dsn string = url

	if strings.Contains(url, "?") {
		dsn += "&"
	} else {
		dsn += "?"
	}

	dsn += "parseTime=true&charset=utf8mb4&loc=UTC"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	queries := schema.New(db)

	return &DB{db, queries}
}
