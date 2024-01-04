package handler

import (
	"WebBlogCraft/internal/pkg"
	"WebBlogCraft/internal/repository"
	"WebBlogCraft/internal/service"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Handler struct {
	UsersRepository *repository.RepositoryUsers
}

func NewHandler(UsersRepository *repository.RepositoryUsers) *Handler {
	return &Handler{
		UsersRepository: UsersRepository,
	}
}

func (h *Handler) PublishHandler(w http.ResponseWriter, r *http.Request) {
	htmlContent := service.HtmlContent("html/blog.html")
	fmt.Fprint(w, string(htmlContent))
}

func (h *Handler) SetUserIDHandler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Fprintf(w, "Error getting IP address: %v", err)
		return
	}

	userIP, err := h.UsersRepository.GetIpAdress(ip)
	if err != nil {
		fmt.Println("Error:", err)
	}

	if ip == userIP {
		htmlContent := service.HtmlContent("html/UserIP.html")
		fmt.Fprint(w, string(htmlContent))
		return
	}

	userID := pkg.GenerateUserID()

	cookie := http.Cookie{
		Name:    "userId",
		Value:   userID,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	h.UsersRepository.AddUsers(userID, string(ip))

	htmlContent := service.HtmlContent("html/getUserID.html")
	fmt.Fprint(w, string(htmlContent))
}
