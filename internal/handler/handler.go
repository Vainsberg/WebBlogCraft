package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Vainsberg/WebBlogCraft/internal/service"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Logger      *zap.Logger
	PostService *service.PostService
	AuthService *service.AuthService
}

func NewHandler(logger *zap.Logger, PostService *service.PostService, AuthService *service.AuthService) *Handler {
	return &Handler{
		Logger:      logger,
		PostService: PostService,
		AuthService: AuthService,
	}
}

func (h *Handler) MainPageHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Main page accessed")
	_, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		fmt.Fprint(w, h.PostService.HtmlContent("html/main_page_authorization.html"))
		return
	}
	fmt.Fprint(w, h.PostService.HtmlContent("html/main_page.html"))
}

func (h *Handler) PostsHandler(w http.ResponseWriter, r *http.Request) {
	postsLastTen, err := h.PostService.PostsRepository.GetLastTenPosts()
	if err != nil {
		h.Logger.Error("GetLastTenPosts error: ", zap.Error(err))
	}

	if r.Method == "POST" {
		h.Logger.Info("POST request to PostsHandler")
		contentText := r.FormValue("postContent")
		c, err := r.Cookie("session_token")
		if errors.Is(err, http.ErrNoCookie) {
			fmt.Fprint(w, h.PostService.HtmlContent("html/authorization_wrong.html"))
			return

		} else if err != nil {
			h.Logger.Error("Error retrieving cookie", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.PostService.ClearRedisCache()

		searchUsersId, err := h.AuthService.SessionsRepository.SearchUsersIdSessionCookie(c.Value)
		if err != nil {
			h.Logger.Error("SearchUsersIdSessionCookie error: ", zap.Error(err))
		}

		if searchUsersId == "" {
			fmt.Fprint(w, h.PostService.HtmlContent("html/authorization_wrong.html"))
			return
		}

		h.PostService.PublishPostWithSessionUser(searchUsersId, contentText)

		postsLastTen, err := h.PostService.PostsRepository.GetLastTenPosts()
		if err != nil {
			h.Logger.Error("GetLastTenPosts error: ", zap.Error(err))
		}

		tmpl := h.PostService.ParseHtml("html/blog.tmpl", "blog")
		err = tmpl.Execute(w, postsLastTen)
		if err != nil {
			h.Logger.Error("tmpl.Execute error:", zap.Error(err))
		}
		return
	}

	h.Logger.Info("GET request to PostsHandler")
	tmpl := h.PostService.ParseHtml("html/blog.tmpl", "blog")
	err = tmpl.Execute(w, postsLastTen)
	if err != nil {
		h.Logger.Error("tmpl.Execute error:", zap.Error(err))
	}
}

func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userName := r.FormValue("username")
		userPassword := r.FormValue("password")

		searchName := h.AuthService.CheckUserExistence(userName)

		if searchName != userName {
			h.AuthService.AddUserWithHashedPassword(userName, userPassword)
			fmt.Fprint(w, h.PostService.HtmlContent("html/registration_confirmation_page.html"))
			return
		}

		fmt.Fprint(w, h.PostService.HtmlContent("html/signup_wrong.html"))
		return
	}
	fmt.Fprint(w, h.PostService.HtmlContent("html/signup.html"))
}

func (h *Handler) SigninHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		userName := r.FormValue("username")
		userPassword := r.FormValue("password")

		_, err := r.Cookie("session_token")
		if !errors.Is(err, http.ErrNoCookie) {
			fmt.Fprint(w, h.PostService.HtmlContent("html/signin_error.html"))
			return
		}

		searchPassword := h.AuthService.SearchPassword(userName)
		if err := bcrypt.CompareHashAndPassword([]byte(searchPassword), []byte(userPassword)); err != nil {
			fmt.Fprint(w, h.PostService.HtmlContent("html/signin_wrong.html"))
			return
		}

		http.SetCookie(w, h.AuthService.CreateSessionCookie(userName))
		fmt.Fprint(w, h.PostService.HtmlContent("html/authorization.html"))
		return
	}
	fmt.Fprint(w, h.PostService.HtmlContent("html/signin.html"))
}

func (h *Handler) ViewingPostsHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")

	page, offset := h.PostService.ParsePageAndCalculateOffset(pageStr)
	templateData := h.PostService.CreateTemplateData(page, offset)

	if page == 1 {
		contentRedis := h.PostService.AddContentToRedis()

		contentRedis.Template = templateData
		tmpl := h.PostService.ParseHtml("html/viewing_posts_redis.html", "viewing_posts_redis")
		err := tmpl.Execute(w, contentRedis)
		if err != nil {
			h.Logger.Error("tmpl.Execute error:", zap.Error(err))
		}

	} else if page != 1 {
		tmpl := h.PostService.ParseHtml("html/viewing_posts.html", "viewing_posts")
		err := tmpl.Execute(w, templateData)
		if err != nil {
			h.Logger.Error("tmpl.Execute error:", zap.Error(err))
		}
	}
}

//func (h *Handler) AddLikeToPostHandler(w http.ResponseWriter, r *http.Request) {

////}
