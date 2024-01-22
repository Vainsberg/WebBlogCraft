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
	rows, err := p.db.Query("SELECT Users_posts.Id, Users_posts.Content, COUNT(Likes.Id) FROM Users_posts "+
		"LEFT JOIN Likes ON Users_posts.Id = Likes.Posts_id "+
		"GROUP BY Users_posts.Id, Users_posts.Content "+
		"LIMIT ? OFFSET ?;", 10, offset)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var posts []response.Posts
	for rows.Next() {
		var id, content string
		var likes int
		err := rows.Scan(&id, &content, &likes)
		if err != nil {
			return nil
		}

		posts = append(posts, response.Posts{Content: content, PostId: id, Likes: likes})
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

func (p *RepositoryPosts) GetLastTenPosts() ([]response.Posts, error) {
	rows, err := p.db.Query("SELECT Users_posts.Id, Users_posts.Content, COUNT(Likes.Id) FROM Users_posts " +
		"LEFT JOIN Likes ON Users_posts.Id = Likes.Posts_id " +
		"GROUP BY Users_posts.Id, Users_posts.Content " +
		"ORDER BY DtCreate DESC " +
		"LIMIT 10;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []response.Posts

	for rows.Next() {
		var id, content string
		var likes int
		err := rows.Scan(&id, &content, &likes)
		if err != nil {
			return nil, err
		}

		post := response.Posts{Content: content, PostId: id, Likes: likes}
		posts = append(posts, post)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
