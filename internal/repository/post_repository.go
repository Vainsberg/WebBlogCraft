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

func (p *RepositoryPosts) ContentOutput() (*response.Post, error) {
	rows, err := p.db.Query("SELECT Content FROM Users_posts;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Post := &response.Post{}
	for rows.Next() {
		var item string
		err := rows.Scan(&item)
		if err != nil {
			log.Fatal(err)
		}
		Post.Content = item
	}
	return Post, nil
}

func (p *RepositoryPosts) CalculatePageOffset(offset int) ([]response.Post, error) {
	rows, err := p.db.Query("SELECT Users_posts.Id,Users_posts.Content,COUNT(Likes.Id),Users.UserName FROM Users_posts "+
		"LEFT JOIN Likes ON Users_posts.Id = Likes.Posts_id "+
		"LEFT JOIN Users ON Users_posts.Users_id = Users.Id "+
		"GROUP BY Users_posts.Id,Users.UserName, Users_posts.Content "+
		"ORDER BY Users.DtCreate DESC "+
		"LIMIT ? OFFSET ?;", 10, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []response.Post
	for rows.Next() {
		var id, content, user string
		var likes int
		err := rows.Scan(&id, &content, &likes, &user)
		if err != nil {
			return nil, err
		}

		posts = append(posts, response.Post{Content: content, PostId: id, UserName: user, Likes: likes})
	}
	return posts, nil
}

func (p *RepositoryPosts) CountPosts() (float64, error) {
	var count float64

	err := p.db.QueryRow("SELECT COUNT(*) FROM Users_posts").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *RepositoryPosts) GetLastTenPosts() ([]response.Post, error) {
	rows, err := p.db.Query("SELECT Users_posts.Id,Users_posts.Content,COUNT(Likes.Id),Users.UserName FROM Users_posts " +
		"LEFT JOIN Likes ON Users_posts.Id = Likes.Posts_id " +
		"LEFT JOIN Users ON Users_posts.Users_id = Users.Id " +
		"GROUP BY Users_posts.Id,Users.UserName, Users_posts.Content " +
		"ORDER BY Users_posts.DtCreate DESC " +
		"LIMIT 10;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []response.Post

	for rows.Next() {
		var id, content, user string
		var likes int
		err := rows.Scan(&id, &content, &likes, &user)
		if err != nil {
			return nil, err
		}

		post := response.Post{Content: content, PostId: id, UserName: user, Likes: likes}
		posts = append(posts, post)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
