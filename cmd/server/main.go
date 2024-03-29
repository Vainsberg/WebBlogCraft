package main

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/Vainsberg/WebBlogCraft/internal/config"
	"github.com/Vainsberg/WebBlogCraft/internal/db"
	"github.com/Vainsberg/WebBlogCraft/internal/handler"
	httpserver "github.com/Vainsberg/WebBlogCraft/internal/httpServer"
	"github.com/Vainsberg/WebBlogCraft/internal/rabbitmq"
	"github.com/Vainsberg/WebBlogCraft/internal/redis"
	"github.com/Vainsberg/WebBlogCraft/internal/repository"
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

	conn, err := rabbitmq.ConnectToRabbitMQ(*cfg)
	if err != nil {
		panic("Error create connection Rabbit MQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel. Error: %s", err)
	}

	repositoryUsers := repository.NewRepositoryUsers(db)
	repositorySessions := repository.NewRepositorySessions(db)
	repositoryPosts := repository.NewRepositoryPosts(db)
	repositoryLikes := repository.NewRepositoryLikes(db)
	repositoryComments := repository.NewRepositoryComments(db)
	repositoryEmail := repository.NewRepositoryEmail(db)
	repositoryRedis := redis.NewRepositoryRedis(redisClient)
	repositoryRabbitMQ := rabbitmq.NewRepositoryRabbitMQ(ch, conn)
	PostService := service.NewPostService(logger, repositoryUsers, repositorySessions, repositoryPosts, repositoryLikes, repositoryComments, repositoryEmail, repositoryRabbitMQ, repositoryRedis, cache)
	AuthService := service.NewAuthService(logger, repositoryUsers, repositorySessions, repositoryPosts)
	handler := handler.NewHandler(logger, PostService, AuthService)
	router.HandleFunc("/", handler.MainPageHandler).Methods("GET")
	router.HandleFunc("/posts/add", handler.PostsHandler).Methods("GET", "POST")
	router.HandleFunc("/signup", handler.SignupHandler).Methods("GET", "POST")
	router.HandleFunc("/signin", handler.SigninHandler).Methods("GET", "POST")
	router.HandleFunc("/posts/list", handler.ViewingPostsHandler).Methods("GET", "POST")
	router.HandleFunc("/posts/{postId}/like", handler.AddLikeToPostHandler).Methods("POST")
	router.HandleFunc("/signout", handler.SignOutHandler).Methods("POST", "GET")
	router.HandleFunc("/posts/{postId}/comment", handler.AddCommentToPostHandler).Methods("POST")
	router.HandleFunc("/posts/{commentId}/comment/like", handler.LikeToCommentHandler).Methods("GET")
	router.HandleFunc("/verify-email", handler.EmailVerificationsHandler).Methods("GET", "POST")
	router.HandleFunc("/verify-email/code", handler.EmailCodeHandler).Methods("GET", "POST")
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	fmt.Println("Starting server at", cfg.Addr)

	httpserver.NewHttpServer(router)

}
