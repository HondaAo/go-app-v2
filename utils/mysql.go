package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var db *sql.DB

func NewDB() *sql.DB {
	addr := os.Getenv("DB_ADDR")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
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
