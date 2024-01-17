package main

import (
	"fmt"
	"net/http"

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

	repositoryUsers := repository.NewRepositoryUsers(db)
	repositorySessions := repository.NewRepositorySessions(db)
	repositoryPosts := repository.NewRepositoryPosts(db)
	repositoryRedis := redis.NewRepositoryRedis(redisClient)
	var redisPostId response.PostsIdRedis
	PostService := service.NewPostService(logger, repositoryUsers, repositorySessions, repositoryPosts, repositoryRedis, cache, redisPostId)
	AuthService := service.NewAuthService(logger, repositoryUsers, repositorySessions, repositoryPosts)
	handler := handler.NewHandler(logger, PostService, AuthService)
	router.HandleFunc("/", handler.MainPageHandler).Methods("GET")
	router.HandleFunc("/posts/add", handler.PostsHandler).Methods("GET", "POST")
	router.HandleFunc("/signup", handler.SignupHandler).Methods("GET", "POST")
	router.HandleFunc("/signin", handler.SigninHandler).Methods("GET", "POST")
	router.HandleFunc("/posts/list", handler.ViewingPostsHandler).Methods("GET", "POST")
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	fmt.Println("Starting server at", cfg.Addr)

	httpserver.NewHttpServer(router)
}
