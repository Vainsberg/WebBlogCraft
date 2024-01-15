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
	CREATE TABLE IF NOT EXISTS Users (
		Id INT PRIMARY KEY AUTO_INCREMENT,
		UserName VARCHAR(255) NOT NULL,
		UserPassword VARCHAR(255),
		DtCreate DATETIME DEFAULT CURRENT_TIMESTAMP
	);
`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS Users_posts (
		Id INT PRIMARY KEY AUTO_INCREMENT,
		Users_id INT NOT NULL,
		Content VARCHAR(255),
		DtCreate DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (Users_id) REFERENCES Users(Id)
	);
`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS Sessions (
		Session_id VARCHAR(40) PRIMARY KEY,
		Users_id INT NOT NULL,
		Expiry TIMESTAMP NOT NULL,
		FOREIGN KEY (Users_id) REFERENCES Users(Id)
	);
`)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
