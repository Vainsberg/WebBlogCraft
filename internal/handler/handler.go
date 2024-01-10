package handler

import (
	"fmt"
	"net/http"

	"github.com/Vainsberg/WebBlogCraft/internal/pkg"
	"github.com/Vainsberg/WebBlogCraft/internal/response"
	"github.com/Vainsberg/WebBlogCraft/internal/service"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Service *service.Service
	Logger  *zap.Logger
}

func NewHandler(service *service.Service, logger *zap.Logger) *Handler {
	return &Handler{
		Service: service,
		Logger:  logger,
	}
}

var pageVariables response.Page

func (h *Handler) MainPageHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Main page accessed")
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
		fmt.Fprint(w, h.Service.HtmlContent("html/session_cookie.html"))
			return
		}
		h.Logger.Error("Error:", zap.Error(err))
	}

	if !h.Service.UsersRepository.SearchSessionCookie(c.Value) {
		fmt.Fprint(w, h.Service.HtmlContent("html/session_cookie.html"))
		return
	}
	if !h.Service.UsersRepository.CheckingTimeforCookie(c.Value) {
		h.Service.UsersRepository.DeleteSessionCookie(c.Value)
		fmt.Fprint(w, h.Service.HtmlContent("html/session_expiration.html"))
		return
	}
	fmt.Fprint(w, h.Service.HtmlContent("html/main_page.html"))
}

func (h *Handler) PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.Logger.Info("POST request to PostsHandler")
		contentText := r.FormValue("postContent")

		tmpl := h.Service.ParseHtml("html/blog.tmpl", "blog")
		err := tmpl.Execute(w, pkg.AddContentToPosts(contentText, pageVariables))
		if err != nil {
			h.Logger.Error("Error executing template", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	h.Logger.Info("GET request to PostsHandler")
	tmpl := h.Service.ParseHtml("html/blog.tmpl", "blog")
	err := tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userName := r.FormValue("username")
		userPassword := r.FormValue("password")
		h.Service.AddUserWithHashedPassword(userName, userPassword)

		fmt.Fprint(w, h.Service.HtmlContent("html/registration_confirmation_page.html"))
		return
	}
	fmt.Fprint(w, h.Service.HtmlContent("html/singup.html"))
}

func (h *Handler) SinginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userName := r.FormValue("username")
		userPassword := r.FormValue("password")
		searchPassword := h.Service.SearchPassword(userName)

		if err := bcrypt.CompareHashAndPassword([]byte(searchPassword), []byte(userPassword)); err != nil {
			fmt.Fprint(w, h.Service.HtmlContent("html/singin_wrong.html"))
			return
		}
		http.SetCookie(w, h.Service.CreateSessionCookie(userName))
		fmt.Fprint(w, h.Service.HtmlContent("html/authorization.html"))
		return
	}
	fmt.Fprint(w, h.Service.HtmlContent("html/singin.html"))
}
