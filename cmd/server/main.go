package main

import (
	"WebBlogCraft/internal/db"
	"WebBlogCraft/internal/handler"
	"WebBlogCraft/internal/repository"
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
	handler := handler.NewHandler(repository)

	r.HandleFunc("/", handler.MainPageHandler).Methods("GET")
	r.HandleFunc("/getuserid", handler.SetUserIDHandler).Methods("GET")
	r.HandleFunc("/publish", handler.PublishHandler).Methods("GET", "POST")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
