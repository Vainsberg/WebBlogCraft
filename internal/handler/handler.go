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
	fmt.Fprint(w, h.Service.HtmlContent("html/mainPage.html"))
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
			fmt.Fprint(w, h.Service.HtmlContent("html/registration_page.html"))
			return
		}

		contentText := r.FormValue("postContent")

		tmpl, err := template.New("blog").Parse(h.Service.HtmlContent("html/blog.html"))
		if err != nil {
			h.Logger.Error("Error parsing HTML content", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, pkg.AddContentToPosts(contentText, pageVariables))
		if err != nil {
			h.Logger.Error("Error executing template", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		h.Logger.Info("GET request to PostsHandler")
		fmt.Fprint(w, h.Service.HtmlContent("html/blog.html"))
	}
}

func (h *Handler) SetUserIDHandler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		h.Logger.Error("Error getting IP address:", zap.Error(err))
		fmt.Fprintf(w, "Error getting IP address: %v", err)
		return
	}
	fetchedIP := h.Service.FetchUserDataByIP(ip)

	if fetchedIP {
		fmt.Fprint(w, h.Service.HtmlContent("html/registration_confirmation_page.html"))
		return
	}

	http.SetCookie(w, h.Service.CreateUserIDCookie(ip))
	fmt.Fprint(w, h.Service.HtmlContent("html/user_exists.html"))
}

func (h *Handler) SetNameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			h.Logger.Error("Error getting IP address:", zap.Error(err))
			fmt.Fprintf(w, "Error getting IP address: %v", err)
			return
		}
		nameResult := r.FormValue("username")
		http.SetCookie(w, h.Service.SetUserNameAndRepository(nameResult, pageVariables, ip))

	} else {
		fmt.Fprint(w, h.Service.HtmlContent("html/set_name_page.html"))
	}
}
