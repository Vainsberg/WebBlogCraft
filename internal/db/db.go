package db

import (
	"database/sql"
	"log"

	"github.com/Vainsberg/WebBlogCraft/internal/viper"
)

func CreateOB(cfg *viper.Сonfigurations) *sql.DB {
	db, err := sql.Open("mysql", cfg.DbUser+":"+cfg.DbPass+"@tcp(127.0.0.1:3306)/user_posts_db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users_posts (
		id INT PRIMARY KEY AUTO_INCREMENT,
		UserID TEXT,
		UserIP TEXT,
		Сontent TEXT,
		dt DATETIME DEFAULT CURRENT_TIMESTAMP
	);
`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
