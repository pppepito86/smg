package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func OpenConnection() error {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/smg?parseTime=true")
	return err
}

func getConnection() *sql.DB {
	return db
}

func Close() {
	db.Close()
}
