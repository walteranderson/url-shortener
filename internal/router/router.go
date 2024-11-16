package router

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/walteranderson/url-shortener/internal/database"
)

var indexTemplate = template.Must(template.ParseFiles("web/index.html"))
var toastTemplate = template.Must(template.ParseFiles("web/_toast-message.html"))

type Router struct {
	repo *database.Repository
}

func NewRouter(repo *database.Repository) http.Handler {
	r := &Router{repo: repo}
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /", r.homePageHandler)
	mux.HandleFunc("GET /{code}", r.redirectHandler)
	mux.HandleFunc("POST /create", r.createHandler)
	return mux
}

func (r *Router) homePageHandler(w http.ResponseWriter, req *http.Request) {
	err := indexTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *Router) createHandler(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	url := req.Form.Get("url")
	link, err := r.repo.CreateLink(url)

	data := struct {
		Success bool
		*database.Link
	}{
		Success: true,
		Link:    link,
	}

	err = toastTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *Router) redirectHandler(w http.ResponseWriter, req *http.Request) {
	code := req.PathValue("code")
	if code == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	link, err := r.repo.GetLink(code)
	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Not Found"}`))
		return
	} else if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	log.Printf("Redirecting: %s to %s", code, link.Url)
	http.Redirect(w, req, link.Url, http.StatusFound)
}
