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

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS Likes (
		Id INT AUTO_INCREMENT PRIMARY KEY,
		Users_id INT,
		Posts_id INT,
		Liked_at DATETIME,
		FOREIGN KEY (Users_id) REFERENCES Users(id),
		FOREIGN KEY (Posts_id) REFERENCES Users_posts(id)
	);
`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS Comments (
		Id INT AUTO_INCREMENT PRIMARY KEY,
		Comment VARCHAR(255),
		Users_id INT,
		Posts_id INT,
		Comment_at DATETIME,
		FOREIGN KEY (Users_id) REFERENCES Users(id),
		FOREIGN KEY (Posts_id) REFERENCES Users_posts(id)
	);
`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS Comments_like (
		Id INT AUTO_INCREMENT PRIMARY KEY,
		Users_id INT,
		Comments_id INT,
		Comment_like_at DATETIME,
		FOREIGN KEY (Users_id) REFERENCES Users(id),
		FOREIGN KEY (Comments_id) REFERENCES Comments(id)
	);
`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS EmailVerifications (
		Id INT AUTO_INCREMENT PRIMARY KEY,
		Users_id INT,
		email VARCHAR(255),
		is_email_verified BOOLEAN DEFAULT FALSE,
		FOREIGN KEY (Users_id) REFERENCES Users(id)
	);
`)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
