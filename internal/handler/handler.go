package handler

import (
	"fmt"
	"html/template"
	"net"
	"net/http"

	"github.com/Vainsberg/WebBlogCraft/internal/pkg"
	"github.com/Vainsberg/WebBlogCraft/internal/response"
	"github.com/Vainsberg/WebBlogCraft/internal/service"
	"go.uber.org/zap"
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
	htmlContent := h.Service.HtmlContent("html/mainPage.html")
	fmt.Fprint(w, htmlContent)
	h.Logger.Info("Main page accessed")
}

func (h *Handler) PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.Logger.Info("Post request to PostsHandler")
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			h.Logger.Error("Error getting IP address", zap.Error(err))
			fmt.Fprintf(w, "Error getting IP address: %v", err)
			return

		}
		fetchedIP := h.Service.FetchUserDataByIP(ip)

		if !fetchedIP {
			h.Logger.Info("User IP not found, redirecting to registration page")
			htmlContent := h.Service.HtmlContent("html/registration_page.html")
			fmt.Fprint(w, htmlContent)
			return
		}
		contentText := r.FormValue("postContent")
		pageVariables = pkg.AddContentToPosts(contentText, pageVariables)

		htmlContent := h.Service.HtmlContent("html/blog.html")

		tmpl, err := template.New("blog").Parse(htmlContent)
		if err != nil {
			h.Logger.Error("Error parsing HTML content", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, pageVariables)
		if err != nil {
			h.Logger.Error("Error executing template", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		h.Logger.Info("GET request to PostsHandler")
		htmlContent := h.Service.HtmlContent("html/blog.html")
		fmt.Fprint(w, htmlContent)
	}
}

func (h *Handler) SetUserIDHandler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Fprintf(w, "Error getting IP address: %v", err)
		return
	}
	fetchedIP := h.Service.FetchUserDataByIP(ip)

	if fetchedIP {
		htmlContent := h.Service.HtmlContent("html/registration_confirmation_page.html")
		fmt.Fprint(w, htmlContent)
		return
	}
	cookie, userID := h.Service.CreateUserIDCookie()
	http.SetCookie(w, &cookie)
	h.Service.UsersRepository.AddUsers(userID, ip)

	htmlContent := h.Service.HtmlContent("html/user_exists.html")
	fmt.Fprint(w, htmlContent)
}

func (h *Handler) SetNameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Fprintf(w, "Error getting IP address: %v", err)
			return
		}
		nameResult := r.FormValue("username")

		h.Service.SetUserNameAndRepository(nameResult, pageVariables, ip)
	} else {
		htmlContent := h.Service.HtmlContent("html/set_name_page.html")
		fmt.Fprint(w, htmlContent)
	}
}
