package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Vainsberg/WebBlogCraft/internal/db"
	"github.com/Vainsberg/WebBlogCraft/internal/handler"
	"github.com/Vainsberg/WebBlogCraft/internal/repository"
	"github.com/Vainsberg/WebBlogCraft/internal/service"
	"github.com/Vainsberg/WebBlogCraft/internal/viper"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
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

	logger, err := zap.NewProduction()
	if err != nil {
		panic("Error create logger")
	}

	repository := repository.NewRepositoryUsers(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service, logger)
	r.HandleFunc("/", handler.MainPageHandler).Methods("GET")
	r.HandleFunc("/getuserid", handler.SetUserIDHandler).Methods("GET")
	r.HandleFunc("/posts", handler.PostsHandler).Methods("GET", "POST")
	r.HandleFunc("/getuserid/setname", handler.SetNameHandler).Methods("GET", "POST")
	fmt.Println("Starting server at :8080")

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
