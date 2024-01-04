package db

import (
	"WebBlogCraft/internal/conifg"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func CreateOB(cfg *conifg.Config) *sql.DB {
	db, err := sql.Open("mysql", cfg.DbUser+":"+cfg.DbPass+"@tcp(127.0.0.1:3306)/users")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		UserID TEXT,
		Post   TEXT,
		dt DATETIME DEFAULT CURRENT_TIMESTAMP,
	)
`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
