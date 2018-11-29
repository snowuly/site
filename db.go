package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "tony:123@/site")

	err := db.Ping()
	if err != nil {
		panic(err)
	}

}
