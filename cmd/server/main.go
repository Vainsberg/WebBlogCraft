package main

import (
	"WebBlogCraft/internal/db"
	"WebBlogCraft/internal/handler"
	"WebBlogCraft/internal/repository"
	"WebBlogCraft/internal/service"
	"WebBlogCraft/internal/viper"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	cfg, err := viper.NewConfig()
	if err != nil {
		fmt.Println("Error creating config:", err.Error())
		return
	}

	db := db.CreateOB(cfg)
	defer db.Close()

	repository := repository.NewRepositoryUsers(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)
	r.HandleFunc("/", handler.MainPageHandler).Methods("GET")
	r.HandleFunc("/getuserid", handler.SetUserIDHandler).Methods("GET")
	r.HandleFunc("/posts", handler.PostsHandler).Methods("GET", "POST")
	r.HandleFunc("/getuserid/setname", handler.SetNameHandler).Methods("GET", "POST")
	fmt.Println("Starting server at :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
