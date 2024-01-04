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

func (r *RepositoryUsers) AddUsers(Userid, Userip string) error {
	_, err := r.db.Exec("INSERT INTO users_posts (UserID,UserIP,dt) VALUES(?,?,CURRENT_TIMESTAMP())", Userid, Userip)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *RepositoryUsers) GetIpAdress(ip string) (string, error) {
	var userIP string

	row := r.db.QueryRow("SELECT UserIP FROM users_posts WHERE UserIP = ?", ip)
	if err := row.Scan(&userIP); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return userIP, nil
}
