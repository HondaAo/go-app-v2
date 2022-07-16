package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func NewDB() *sql.DB {
	addr := os.Getenv("MYSQL_HOST")
	name := os.Getenv("MYSQL_DB")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	log.Println(addr)
	log.Println(name)
	log.Println(user)
	log.Println(password)

	dsn := fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=utf8mb4&parseTime=true",
		user,
		password,
		addr,
		name,
	)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("main SQL Open: %v", err)
	}

	return db
}
