package httpserver

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/Vainsberg/WebBlogCraft/internal/config"
	"github.com/gorilla/mux"
)

func NewHttpServer(r *mux.Router) {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println("Error creating config:", err.Error())
		return
	}

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: r,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
