package main

import (
	"fmt"

	config "github.com/Vainsberg/WebBlogCraft/internal/config"
	"github.com/Vainsberg/WebBlogCraft/internal/db"
	"github.com/Vainsberg/WebBlogCraft/internal/handler"
	httpserver "github.com/Vainsberg/WebBlogCraft/internal/httpServer"
	"github.com/Vainsberg/WebBlogCraft/internal/redis"
	"github.com/Vainsberg/WebBlogCraft/internal/repository"
	"github.com/Vainsberg/WebBlogCraft/internal/response"
	"github.com/Vainsberg/WebBlogCraft/internal/service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	router := mux.NewRouter()

	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println("Error creating config:", err.Error())
		return
	}

	db := db.CreateOB(cfg)
	defer db.Close()
	redisClient := redis.NewRedisClient(cfg)
	cache := redis.CreateRedisCache(cfg)

	logger, err := zap.NewProduction()
	if err != nil {
		panic("Error create logger")
	}
	postsStorage := response.StoragePostsRedis{}

	repositoryUsers := repository.NewRepositoryUsers(db)
	repositorySessions := repository.NewRepositorySessions(db)
	repositoryPosts := repository.NewRepositoryPosts(db)
	repositoryRedis := redis.NewRepositoryRedis(redisClient)
	service := service.NewService(logger, repositoryUsers, repositorySessions, repositoryPosts, repositoryRedis, cache, postsStorage)
	handler := handler.NewHandler(logger, service)
	router.HandleFunc("/", handler.MainPageHandler).Methods("GET")
	router.HandleFunc("/posts", handler.PostsHandler).Methods("GET", "POST")
	router.HandleFunc("/singup", handler.SignupHandler).Methods("GET", "POST")
	router.HandleFunc("/singin", handler.SigninHandler).Methods("GET", "POST")
	router.HandleFunc("/posts/viewing", handler.ViewingPostsHandler).Methods("GET", "POST")
	fmt.Println("Starting server at", cfg.Addr)

	httpserver.NewHttpServer(router)
}
