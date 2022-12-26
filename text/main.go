package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:13306)/fileserver")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	test(db)
}

func test(db *sql.DB) {
	stmt, _ := db.Prepare("insert into test(id, age) values (?,?)")
	stmt.Exec(12, "18")
}
