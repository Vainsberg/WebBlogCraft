package handler

import (
	"errors"
	"fmt"
	"net/http"

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

var Posts response.StoragePosts

func (h *Handler) MainPageHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Main page accessed")
	c, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		h.Logger.Error("Error:", zap.Error(err))
		return
	}

	if !h.Service.SessionsRepository.SearchSessionCookie(c.Value) {
		fmt.Fprint(w, h.Service.HtmlContent("html/session_cookie.html"))
		return
	}
	if h.Service.SessionsRepository.CheckingTimeforCookie(c.Value) {
		h.Service.SessionsRepository.DeleteSessionCookie(c.Value)
		fmt.Fprint(w, h.Service.HtmlContent("html/session_expiration.html"))
		return
	}
	fmt.Fprint(w, h.Service.HtmlContent("html/main_page.html"))
}

func (h *Handler) PostsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		h.Logger.Info("POST request to PostsHandler")
		contentText := r.FormValue("postContent")
		c, err := r.Cookie("session_token")
		if errors.Is(err, http.ErrNoCookie) {
			h.Logger.Error("Error:", zap.Error(err))
			return
		}
		h.Service.PublishPostWithSessionUser(c.Value, contentText)

		tmpl := h.Service.ParseHtml("html/blog.tmpl", "blog")

		Posts = h.Service.AddContentToPosts(contentText, Posts)

		err = tmpl.Execute(w, Posts)
		if err != nil {
			h.Logger.Error("Error executing template", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	h.Logger.Info("GET request to PostsHandler")
	tmpl := h.Service.ParseHtml("html/blog.tmpl", "blog")
	err := tmpl.Execute(w, Posts)
	if err != nil {
		h.Logger.Error("tmpl.Execute error:", zap.Error(err))
	}
}

func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userName := r.FormValue("username")
		userPassword := r.FormValue("password")

		searchName, err := h.Service.UsersRepository.CheckingPresenceUser(userName)
		if err != nil {
			h.Logger.Error("CheckingPresenceUser error: ", zap.Error(err))
		}

		if searchName != userName {
			h.Service.AddUserWithHashedPassword(userName, userPassword)
			fmt.Fprint(w, h.Service.HtmlContent("html/registration_confirmation_page.html"))
			return
		}
		fmt.Fprint(w, h.Service.HtmlContent("html/signup_wrong.html"))
		return
	}
	fmt.Fprint(w, h.Service.HtmlContent("html/signup.html"))
}

func (h *Handler) SigninHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userName := r.FormValue("username")
		userPassword := r.FormValue("password")

		if h.Service.SessionsRepository.SearchAccountInSessions(userName) {
			fmt.Fprint(w, h.Service.HtmlContent("html/signin_error.html"))
			return
		}
		searchPassword := h.Service.SearchPassword(userName)

		if err := bcrypt.CompareHashAndPassword([]byte(searchPassword), []byte(userPassword)); err != nil {
			fmt.Fprint(w, h.Service.HtmlContent("html/signin_wrong.html"))
			return
		}
		http.SetCookie(w, h.Service.CreateSessionCookie(userName))
		fmt.Fprint(w, h.Service.HtmlContent("html/authorization.html"))
		return
	}
	fmt.Fprint(w, h.Service.HtmlContent("html/signin.html"))
}

func (h *Handler) ViewingPostsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := h.Service.ParseHtml("html/viewing_posts.html", "viewing_posts")
	err := tmpl.Execute(w, Posts)
	if err != nil {
		h.Logger.Error("tmpl.Execute error:", zap.Error(err))
	}
}
