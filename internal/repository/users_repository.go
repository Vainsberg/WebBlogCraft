package repository

import (
	"database/sql"
	"fmt"
)

type RepositoryUsers struct {
	db *sql.DB
}

func NewRepositoryUsers(db *sql.DB) *RepositoryUsers {
	return &RepositoryUsers{db: db}
}

func (r *RepositoryUsers) AddUsers(Userid string) error {
	_, err := r.db.Exec("INSERT INTO users (UserID,dt) VALUES(?,CURRENT_TIMESTAMP())", Userid)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
