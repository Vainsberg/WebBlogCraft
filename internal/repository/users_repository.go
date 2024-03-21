package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type UserRepository interface {
	AddUserWithHashedPassword(userName, userPassword string) error
}

type RepositoryUsers struct {
	db *sql.DB
}

func NewRepositoryUsers(db *sql.DB) *RepositoryUsers {
	return &RepositoryUsers{db: db}
}

func (r *RepositoryUsers) AddPasswordAndUserName(userName, password string) error {
	_, err := r.db.Exec("INSERT INTO Users (UserName,UserPassword,DtCreate) VALUES(?,?,CURRENT_TIMESTAMP())", userName, password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *RepositoryUsers) SearchPasswordAndUserName(userName string) (string, error) {
	var UserPassword string
	row := r.db.QueryRow("SELECT UserPassword FROM Users WHERE UserName = ?", userName)
	if err := row.Scan(&UserPassword); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return UserPassword, nil
}

func (r *RepositoryUsers) CheckingPresenceUser(username string) (string, error) {
	var searchUserName string
	row := r.db.QueryRow("SELECT UserName FROM Users WHERE UserName = ?;", username)
	if err := row.Scan(&searchUserName); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return searchUserName, nil
}

func (r *RepositoryUsers) SearchUserName(UserID int) (string, error) {
	var searchUserName string
	row := r.db.QueryRow("SELECT UserName FROM Users WHERE Id = ?;", UserID)
	if err := row.Scan(&searchUserName); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return searchUserName, nil
}
