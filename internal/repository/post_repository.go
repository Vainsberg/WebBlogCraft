package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Vainsberg/WebBlogCraft/internal/response"
)

type RepositoryPosts struct {
	db *sql.DB
}

func NewRepositoryPosts(db *sql.DB) *RepositoryPosts {
	return &RepositoryPosts{db: db}
}

func (p *RepositoryPosts) AddContentAndUserId(UsersId, content string) error {
	_, err := p.db.Exec(`
        INSERT INTO Users_posts (Users_id, Content, DtCreate)
        VALUES (?, ?, CURRENT_TIMESTAMP());
    `, UsersId, content)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (p *RepositoryPosts) ContentOutput() (response.Posts, error) {
	rows, err := p.db.Query("SELECT Content FROM Users_posts;")
	if err != nil {
		return response.Posts{}, err
	}
	defer rows.Close()

	Posts := response.Posts{}
	for rows.Next() {
		var item string
		err := rows.Scan(&item)
		if err != nil {
			log.Fatal(err)
		}
		Posts.Content = append(Posts.Content, item)
	}
	return Posts, nil
}

func (p *RepositoryPosts) CalculatePageOffset(offset int) []response.Posts {
	rows, err := p.db.Query("SELECT Content FROM Users_posts LIMIT ? OFFSET ?", 10, offset)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var posts []response.Posts
	for rows.Next() {
		var contentBytes []byte
		err := rows.Scan(&contentBytes)
		if err != nil {
			return nil
		}
		content := string(contentBytes)

		contentSlice := strings.Fields(content)

		post := response.Posts{Content: contentSlice}
		posts = append(posts, post)
	}

	return posts
}

func (p *RepositoryPosts) CountPosts() (float64, error) {
	var count float64

	err := p.db.QueryRow("SELECT COUNT(*) FROM Users_posts").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *RepositoryPosts) GetLastTenPosts() ([]response.PostsRedis, error) {
	rows, err := p.db.Query("SELECT Content FROM Users_Posts ORDER BY DtCreate DESC LIMIT 10;")
	if err != nil {
		return []response.PostsRedis{}, err
	}
	defer rows.Close()

	var posts []response.PostsRedis

	if rows.Next() {
		var content string
		err := rows.Scan(&content)
		if err != nil {
			return []response.PostsRedis{}, err
		}

		contentSlice := strings.Fields(content)

		post := response.PostsRedis{Content: contentSlice}
		posts = append(posts, post)

	}
	return posts, nil

}
