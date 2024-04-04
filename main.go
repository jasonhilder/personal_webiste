package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/jasonhilder/personal_website/utils"
)

//go:embed html
var htmlPages embed.FS

//go:embed html/entries
var entries embed.FS

//go:embed static
var staticFiles embed.FS

type PostLinks struct {
    Posts []utils.Post
}

func main() {
    fs, err := fs.Sub(staticFiles, "static")
    if err != nil {
        log.Fatal(err)
    }
    http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(fs))))

    http.HandleFunc("GET /entries/", listEntries)
    http.HandleFunc("GET /entries/{entry_url}", serveEntry)

	http.HandleFunc("GET /", serveRoot)

    log.Println("Listening on http://localhost:42069")
	err = http.ListenAndServe(":42069", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func renderPage(w http.ResponseWriter, r *http.Request, page string, data any) {
	page = filepath.Clean(page)
	page = strings.TrimPrefix(page, "/")
    
    // @todo move to main function
	view := template.Must(template.ParseFS(
        htmlPages, 
        "html/*.html", 
        "html/partials/*.html",
        "html/entries/*.html",
    ))

    err := view.ExecuteTemplate(w, page, data)
	if err != nil {
		log.Println(err)

		//http.Error(w, err.Error(), 500)
        var i interface{}
        renderPage(w, r, "/404.html", i)
	}
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		r.URL.Path = "index.html"
    }

    var i interface{}
	renderPage(w, r, r.URL.Path, i)
}

func listEntries(w http.ResponseWriter, r *http.Request) {
    // @todo move to main function
    files, err := fs.ReadDir(entries, "html/entries")
    if err != nil {
        log.Println("Error reading directory:", err)
        return
    }

    var postLinks PostLinks
    var posts []utils.Post
    for _, file := range files {
        p := utils.GetPost(file)
        posts = append([]utils.Post{p}, posts...)
    }

    postLinks.Posts = posts
    renderPage(w, r, "/entries.html", postLinks)
}

func serveEntry(w http.ResponseWriter, r *http.Request) {
    entry := r.PathValue("entry_url")

    var i interface{}
    renderPage(w, r, entry, i)
}
