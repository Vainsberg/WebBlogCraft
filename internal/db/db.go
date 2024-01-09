package db

import (
	"database/sql"
	"log"

	config "github.com/Vainsberg/WebBlogCraft/internal/config"
)

func CreateOB(cfg *config.Сonfigurations) *sql.DB {
	db, err := sql.Open("mysql", cfg.DbUser+":"+cfg.DbPass+"@tcp("+cfg.DbIp+":"+cfg.DbPort+")/user_posts_db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users_posts (
		id INT PRIMARY KEY AUTO_INCREMENT,
		UserName TEXT,
		Сontent TEXT,
		dt DATETIME DEFAULT CURRENT_TIMESTAMP
	);
`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INT PRIMARY KEY AUTO_INCREMENT,
		UserName VARCHAR(255),
		UserPassword TEXT,
		dt DATETIME DEFAULT CURRENT_TIMESTAMP
	);
`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS sessions (
		Session_id VARCHAR(40) PRIMARY KEY,
		UserName VARCHAR(255) NOT NULL,
		Expiry TIMESTAMP NOT NULL
	);
`)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
