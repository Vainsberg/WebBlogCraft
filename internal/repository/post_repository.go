package repository

import (
	"database/sql"
	"fmt"
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
