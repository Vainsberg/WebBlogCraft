package repository

import (
	"database/sql"
	"fmt"
)

type RepositoryEmail struct {
	db *sql.DB
}

func NewRepositoryEmail(db *sql.DB) *RepositoryEmail {
	return &RepositoryEmail{db: db}
}

func (e *RepositoryEmail) AddEmailAndUserId(UsersId int, content string) error {
	_, err := e.db.Exec(`
        INSERT INTO EmailVerifications (Users_id, Email)
        VALUES (?, ?);
    `, UsersId, content)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
