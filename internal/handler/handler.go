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
	Logger  *zap.Logger
	Service *service.Service
}

func NewHandler(logger *zap.Logger, service *service.Service) *Handler {
	return &Handler{
		Logger:  logger,
		Service: service,
	}
}

func (h *Handler) MainPageHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Main page accessed")
	_, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		fmt.Fprint(w, h.Service.HtmlContent("html/main_page_authorization.html"))
		return
	}
	fmt.Fprint(w, h.Service.HtmlContent("html/main_page.html"))
}

func (h *Handler) PostsHandler(w http.ResponseWriter, r *http.Request) {
	var outputContent response.Posts

	if r.Method == "POST" {
		h.Logger.Info("POST request to PostsHandler")
		contentText := r.FormValue("postContent")
		c, err := r.Cookie("session_token")
		if errors.Is(err, http.ErrNoCookie) {
			fmt.Fprint(w, h.Service.HtmlContent("html/authorization_wrong.html"))
			return

		} else if err != nil {
			h.Logger.Error("Error retrieving cookie", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		searchUserName, err := h.Service.SessionsRepository.SearchUserNameSessionCookie(c.Value)
		if err != nil {
			h.Logger.Error("SearchUserNameSessionCookie error: ", zap.Error(err))
		}

		if searchUserName == "" {
			fmt.Fprint(w, h.Service.HtmlContent("html/authorization_wrong.html"))
			return
		}

		h.Service.PublishPostWithSessionUser(searchUserName, contentText)
		h.Service.AddContentToPosts(contentText)

		tmpl := h.Service.ParseHtml("html/blog.tmpl", "blog")

		outputContent = h.Service.PostsRepository.ContentOutput()
		err = tmpl.Execute(w, outputContent)
		if err != nil {
			h.Logger.Error("Error executing template", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	h.Logger.Info("GET request to PostsHandler")
	tmpl := h.Service.ParseHtml("html/blog.tmpl", "blog")
	err := tmpl.Execute(w, outputContent)
	if err != nil {
		h.Logger.Error("tmpl.Execute error:", zap.Error(err))
	}
}

func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userName := r.FormValue("username")
		userPassword := r.FormValue("password")

		searchName := h.Service.CheckUserExistence(userName)

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

		if h.Service.SessionsRepository.SearchAccountInSessions(userName) != "" {
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
	pageStr := r.URL.Query().Get("page")

	page, offset := h.Service.ParsePageAndCalculateOffset(pageStr)
	if page == 1 {
		tmpl := h.Service.ParseHtml("html/viewing_posts_redis.html", "viewing_posts_redis")
		err := tmpl.Execute(w, h.Service.PostsRedis)
		if err != nil {
			h.Logger.Error("tmpl.Execute error:", zap.Error(err))
		}

	} else if page != 1 {
		templateData := h.Service.CreateTemplateData(page, offset)

		tmpl := h.Service.ParseHtml("html/viewing_posts.html", "viewing_posts")
		err := tmpl.Execute(w, templateData)
		if err != nil {
			h.Logger.Error("tmpl.Execute error:", zap.Error(err))
		}
	}
}
