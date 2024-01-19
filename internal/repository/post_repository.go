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

func (p *RepositoryPosts) ContentOutput() (*response.Posts, error) {
	rows, err := p.db.Query("SELECT Content FROM Users_posts;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Posts := &response.Posts{}
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
	rows, err := p.db.Query("SELECT Content,Id FROM Users_posts LIMIT ? OFFSET ?", 10, offset)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var posts []response.Posts
	for rows.Next() {
		var content, id string
		err := rows.Scan(&content, &id)
		if err != nil {
			return nil
		}

		contentSlice := strings.Fields(content)
		contentIdSlice := strings.Fields(id)

		post := response.Posts{Content: contentSlice, PostId: contentIdSlice}
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
	rows, err := p.db.Query("SELECT Content,Id FROM Users_Posts ORDER BY DtCreate DESC LIMIT 10;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []response.PostsRedis

	for rows.Next() {
		var content, id string
		err := rows.Scan(&content, &id)
		if err != nil {
			return nil, err
		}

		contentSlice := strings.Fields(content)
		postIdSlice := strings.Fields(id)

		post := response.PostsRedis{Content: contentSlice, PostId: postIdSlice}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return []response.PostsRedis{}, err
	}
	return posts, nil
}
