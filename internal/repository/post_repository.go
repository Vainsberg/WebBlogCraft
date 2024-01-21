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

func (p *RepositoryPosts) AddContentAndUserId(UsersId int, content string) error {
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
		Posts.Content = item
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

		posts = append(posts, response.Posts{Content: content, PostId: id})
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

func (p *RepositoryPosts) GetLastTenPostsAndPostsId() ([]response.PostsRedis, response.PostsRedis, error) {
	rows, err := p.db.Query("SELECT Content,Id FROM Users_posts ORDER BY DtCreate DESC LIMIT 10;")
	if err != nil {
		return nil, response.PostsRedis{}, err
	}
	defer rows.Close()

	var posts []response.PostsRedis
	var postsId response.PostsRedis

	for rows.Next() {
		var content, id string
		err := rows.Scan(&content, &id)
		if err != nil {
			return nil, response.PostsRedis{}, err
		}

		contentSlice := strings.Fields(content)

		post := response.PostsRedis{Content: contentSlice}

		posts = append(posts, post)
		postsId.PostId = append(postsId.PostId, id)

	}

	if err := rows.Err(); err != nil {
		return nil, response.PostsRedis{}, err
	}
	return posts, postsId, nil
}
