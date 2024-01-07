package repository

import (
	"database/sql"
	"fmt"
	"log"
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

func (r *RepositoryUsers) GetSetName(ip string, name string) error {
	stmt, err := r.db.Prepare("UPDATE users_posts SET UserID = ? WHERE UserIP = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, ip)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
