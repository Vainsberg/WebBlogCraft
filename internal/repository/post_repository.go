package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Vainsberg/WebBlogCraft/internal/response"
)

type RepositoryPosts struct {
	db *sql.DB
}

func NewRepositoryPosts(db *sql.DB) *RepositoryPosts {
	return &RepositoryPosts{db: db}
}

func (p *RepositoryPosts) AddContentAndUserName(username, content string) error {
	_, err := p.db.Exec("INSERT INTO users_posts (UserName,Content,dt) VALUES(?,?,CURRENT_TIMESTAMP())", username, content)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (p *RepositoryPosts) ContentOutput() response.Posts {
	rows, err := p.db.Query("SELECT Content FROM users_posts;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	Posts := response.Posts{}
	for rows.Next() {
		var item string
		err := rows.Scan(&item)
		if err != nil {
			log.Fatal(err)
		}
		Posts.Posts = append(Posts.Posts, item)
	}
	return Posts

}
