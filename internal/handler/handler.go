package handler

import (
	"WebBlogCraft/internal/pkg"
	"WebBlogCraft/internal/repository"
	"WebBlogCraft/internal/service"
	"fmt"
	"net/http"
	"time"
)

type Handler struct {
	UsersRepository repository.RepositoryUsers
}

func NewHandler(UsersRepository repository.RepositoryUsers) *Handler {
	return &Handler{
		UsersRepository: UsersRepository,
	}
}

func (h *Handler) PublishHandler(w http.ResponseWriter, r *http.Request) {
	htmlContent := service.HtmlContent("html/index.html")
	fmt.Fprint(w, string(htmlContent))
}

func (h *Handler) SetUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userID := pkg.GenerateUserID()

	cookie := http.Cookie{
		Name:    "userId",
		Value:   userID,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	h.UsersRepository.AddUsers(userID)

	htmlContent := service.HtmlContent("html/getUserID.html")

	fmt.Fprint(w, string(htmlContent))
}
