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

	dsn := fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=utf8mb4&parseTime=true&multiStatements=true",
		user,
		password,
		addr,
		name,
	)
	var err error
	db, err = sql.Open("mysql", dsn)
	log.Print(dsn)
	if err != nil {
		log.Fatalf("main SQL Open: %v", err)
	}

	// driver, _ := mysql.WithInstance(db, &mysql.Config{})
	// m, err := migrate.NewWithDatabaseInstance(
	// 	"file://migrations",
	// 	"mysql",
	// 	driver,
	// )

	// m.Up()

	return db
}
