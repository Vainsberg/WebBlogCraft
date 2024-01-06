package handler

import (
	"WebBlogCraft/internal/pkg"
	"WebBlogCraft/internal/response"
	"WebBlogCraft/internal/service"
	"fmt"
	"io"
	"net"
	"net/http"
	"text/template"
	"time"
)

type Handler struct {
	service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Service: *service,
	}
}

var pageVariables response.Page

func (h *Handler) MainPageHandler(w http.ResponseWriter, r *http.Request) {
	htmlContent := h.Service.HtmlContent("html/mainPage.html")
	fmt.Fprint(w, htmlContent)
}

func (h *Handler) PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Fprintf(w, "Error getting IP address: %v", err)
			return
		}
		fetchedIP := h.Service.FetchUserDataByIP(ip)

		if !fetchedIP {
			htmlContent := h.Service.HtmlContent("html/registration_page.html")
			fmt.Fprint(w, htmlContent)
			return
		}

		resp, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Errorf("Error ReadAll:%s", err)
		}

		pageVariables = pkg.DecodingContentText(resp, pageVariables)

		htmlContent := h.Service.HtmlContent("html/blog.html")

		tmpl, err := template.New("blog").Parse(htmlContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, pageVariables)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
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
	userID := pkg.GenerateUserID()
	cookie := http.Cookie{
		Name:    "userId",
		Value:   userID,
		Expires: time.Now(),
		Path:    "/",
	}

	http.SetCookie(w, &cookie)
	h.UsersRepository.AddUsers(userID, string(ip))

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

		resp, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Errorf("Error ReadAll:%s", err)
		}

		pageVariables = pkg.DecodingName(resp, pageVariables)
		name := pageVariables.UserName[len(pageVariables.UserName)-1]
		text := pkg.RemovingPreposition(name)

		h.UsersRepository.GetSetName(ip, text)
	} else {
		htmlContent := h.Service.HtmlContent("html/set_name_page.html")
		fmt.Fprint(w, htmlContent)
	}
}
