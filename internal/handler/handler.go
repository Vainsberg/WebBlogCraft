package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Vainsberg/WebBlogCraft/internal/response"
	"github.com/Vainsberg/WebBlogCraft/internal/service"
	"github.com/gorilla/mux"
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
	c, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		fmt.Fprint(w, h.PostService.HtmlContent("html/main_page_authorization.html"))
		return
	}
	if !h.PostService.SearchVerifEmail(c.Value) {
		fmt.Fprint(w, h.PostService.HtmlContent("html/main_page_no_verif_email.html"))
		return
	}
	fmt.Fprint(w, h.PostService.HtmlContent("html/main_page.html"))
}

func (h *Handler) PostsHandler(w http.ResponseWriter, r *http.Request) {
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

		if searchUsersId == 0 {
			fmt.Fprint(w, h.PostService.HtmlContent("html/authorization_wrong.html"))
			return
		}

		h.PostService.PublishPostWithSessionUser(searchUsersId, contentText)

		postsLastTen, err := h.PostService.PostsRepository.GetLastTenPosts()
		if err != nil {
			h.Logger.Error("GetLastTenPosts error: ", zap.Error(err))
		}

		if !h.PostService.SearchVerifEmail(c.Value) {
			tmpl := h.PostService.ParseHtml("html/no_verif_email.html", "blog")
			err = tmpl.Execute(w, postsLastTen)
			if err != nil {
				h.Logger.Error("tmpl.Execute error:", zap.Error(err))
			}
			return
		}

		tmpl := h.PostService.ParseHtml("html/blog.html", "blog")
		err = tmpl.Execute(w, postsLastTen)
		if err != nil {
			h.Logger.Error("tmpl.Execute error:", zap.Error(err))
		}
		return
	}

	h.Logger.Info("GET request to PostsHandler")

	postsLastTen, err := h.PostService.PostsRepository.GetLastTenPosts()
	if err != nil {
		h.Logger.Error("GetLastTenPosts error: ", zap.Error(err))
	}

	cookie, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		tmpl := h.PostService.ParseHtml("html/blog_no_authorization.html", "blog")
		err = tmpl.Execute(w, postsLastTen)
		if err != nil {
			h.Logger.Error("tmpl.Execute error:", zap.Error(err))
		}

	} else if !h.PostService.SearchVerifEmail(cookie.Value) {
		tmpl := h.PostService.ParseHtml("html/blog_no_verif_email.html", "blog")
		err = tmpl.Execute(w, postsLastTen)
		if err != nil {
			h.Logger.Error("tmpl.Execute error:", zap.Error(err))
		}

	} else {
		tmpl := h.PostService.ParseHtml("html/blog.html", "blog")
		err = tmpl.Execute(w, postsLastTen)
		if err != nil {
			h.Logger.Error("tmpl.Execute error:", zap.Error(err))
		}
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

		h.PostService.SessionsRepository.DeleteExpiredSessions()

		fmt.Fprint(w, h.PostService.HtmlContent("html/authorization.html"))
		return
	}
	fmt.Fprint(w, h.PostService.HtmlContent("html/signin.html"))
}

func (h *Handler) ViewingPostsHandler(w http.ResponseWriter, r *http.Request) {
	var posts response.ResponseData
	pageStr := r.URL.Query().Get("page")
	_, err := r.Cookie("session_token")

	page, offset := h.PostService.ParsePageAndCalculateOffset(pageStr)

	if page == 1 {
		posts = h.PostService.GenerateTemplateDataPostsRedis(page)
	} else {
		posts = h.PostService.GenerateTemplateDataPosts(page, offset)
	}

	if errors.Is(err, http.ErrNoCookie) {
		tmpl := h.PostService.ParseHtml("html/viewing_posts_no_auth.html", "viewing_posts")
		err = tmpl.Execute(w, posts)
		if err != nil {
			h.Logger.Error("tmpl.Execute error:", zap.Error(err))
		}

	} else {
		tmpl := h.PostService.ParseHtml("html/viewing_posts.html", "viewing_posts")
		err = tmpl.Execute(w, posts)
		if err != nil {
			h.Logger.Error("tmpl.Execute error:", zap.Error(err))
		}
	}
}

func (h *Handler) AddLikeToPostHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("AddLikeToPostHandler accessed")

	cookie, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	postIdStr := vars["postId"]

	updatedLikesCount, err := h.PostService.ProcessLikeAction(cookie.Value, postIdStr)
	if err != nil {
		h.Logger.Error("ProcessLikeAction error:", zap.Error(err))
		http.Error(w, "Error processing like", http.StatusInternalServerError)
		return
	}

	responseNewLikes := response.LikeResponse{
		NewLikesCount: updatedLikesCount,
	}

	jsonResponseLikes, err := json.Marshal(responseNewLikes)
	if err != nil {
		h.Logger.Error("JSON marshaling error:", zap.Error(err))
		http.Error(w, "JSON marshaling error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponseLikes)
}

func (h *Handler) SignOutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		return
	}
	h.PostService.SessionsRepository.DeleteSessionCookieAccount(cookie.Value)

	http.SetCookie(w, h.AuthService.DeleteSessionCookie())

	fmt.Fprint(w, h.PostService.HtmlContent("html/signout.html"))
}

func (h *Handler) AddCommentToPostHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("AddCommentToPostHandler accessed")

	cookie, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	postIdStr := vars["postId"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var cmt response.CommentInput
	err = json.Unmarshal(body, &cmt)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	commentId, username, err := h.PostService.AddUserCommentToPostAndSearchUserName(cookie.Value, postIdStr, cmt.Comment)
	if err != nil {
		http.Error(w, "Error AddUserCommentToPost", http.StatusInternalServerError)
		return
	}

	response := response.CommentResponse{
		CommentID: commentId,
		Comment:   cmt.Comment,
		UserName:  username,
		Likes:     0,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *Handler) LikeToCommentHandler(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("LikeToCommentHandler accessed")

	cookie, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	commentIdStr := vars["commentId"]

	updatedLikesCount, err := h.PostService.LikeActionToComment(cookie.Value, commentIdStr)

	responseNewLikes := response.LikeResponse{
		NewLikesCount: updatedLikesCount,
	}

	jsonResponseLikes, err := json.Marshal(responseNewLikes)
	if err != nil {
		h.Logger.Error("JSON marshaling error:", zap.Error(err))
		http.Error(w, "JSON marshaling error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponseLikes)
}

func (h *Handler) EmailVerificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		cookie, err := r.Cookie("session_token")
		if errors.Is(err, http.ErrNoCookie) {
			return
		}
		email := r.FormValue("email")

		h.PostService.PublishCodeToRabbitMQ(cookie.Value, email)
		http.Redirect(w, r, "/verify-email/code", http.StatusSeeOther)
		return
	}
	fmt.Fprint(w, h.PostService.HtmlContent("html/email_verification.html"))
}

func (h *Handler) EmailCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		c, err := r.Cookie("session_token")
		if errors.Is(err, http.ErrNoCookie) {
			return
		}
		code := r.FormValue("code")

		email := h.PostService.SearchEmailAdr(c.Value)
		SearchCode := h.PostService.GetCodeToRedis(email)

		if code != SearchCode {
			fmt.Fprint(w, h.PostService.HtmlContent("html/email_code_wrong.html"))
			return
		}
		h.PostService.ProcessVerifiedEmail(email)
		fmt.Fprint(w, h.PostService.HtmlContent("html/email_verif.html"))
		return
	}
	fmt.Fprint(w, h.PostService.HtmlContent("html/email_code.html"))
}
