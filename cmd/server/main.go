package main

import (
	"WebBlogCraft/internal/conifg"
	"WebBlogCraft/internal/db"
	"WebBlogCraft/internal/handler"
	"WebBlogCraft/internal/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	cfg, err := conifg.NewConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	db := db.CreateOB(cfg)
	defer db.Close()

	repository := repository.NewRepositoryUsers(db)
	handler := handler.NewHandler(*repository)

	r.HandleFunc("/", handler.PublishHandler).Methods("GET")
	r.HandleFunc("/getuserid", handler.SetUserIDHandler).Methods("GET")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
