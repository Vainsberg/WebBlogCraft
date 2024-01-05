package handler

import (
	"WebBlogCraft/internal/pkg"
	"WebBlogCraft/internal/repository"
	"WebBlogCraft/internal/response"
	"WebBlogCraft/internal/service"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"text/template"
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
func (h *Handler) MainPageHandler(w http.ResponseWriter, r *http.Request) {
	htmlContent := service.HtmlContent("html/mainPage.html")
	fmt.Fprint(w, string(htmlContent))
}

func (h *Handler) PublishHandler(w http.ResponseWriter, r *http.Request) {
	var pageVariables response.Page
	// htmlContent := service.HtmlContent("html/blog.html")
	// fmt.Fprint(w, string(htmlContent))

	if r.Body != nil {
		resp, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Errorf("read body error: %s, %s", resp, err)
		}
		err = json.Unmarshal(resp, &pageVariables)
		if err != nil {
			fmt.Errorf("json.Unmarshel error: %s", err)
		}

		tmpl, err := template.New("blog").Parse(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>BlOG Content</title>
			</head>
			<body>
			
				<h1>BlOG Content</h1>
				<form action="/publish" method="post">
					<label for="postContent">Введите ваш пост:</label><br>
					<textarea id="postContent" name="postContent" rows="4" cols="50"></textarea><br>
					<input type="submit" value="Отправить">
				</form>
				<h3>Лента постов:</h3>
				<ul>
					{{range .Posts}}
						<li>{{.}}</li>
					{{end}}
				</ul>
			</body>
			</html>
		`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, pageVariables)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
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
