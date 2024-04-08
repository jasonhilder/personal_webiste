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

var htmlTemplates *template.Template
var htmlEntries []fs.DirEntry

func main() {
    loadHtmlFiles()

    fs, err := fs.Sub(staticFiles, "static")
    if err != nil {
        log.Fatal(err)
    }
    http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(fs))))

    http.HandleFunc("GET /entries/", listEntries)
    http.HandleFunc("GET /entries/{entry_url}", serveEntry)

    http.HandleFunc("GET /galleries/", listGalleries)

    http.HandleFunc("GET /books/", readingList)

	http.HandleFunc("GET /", serveRoot)

    log.Println("Listening on http://localhost:42069")
	err = http.ListenAndServe(":42069", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func loadHtmlFiles() {
	htmlTemplates = template.Must(template.ParseFS(
        htmlPages, 
        "html/*.html", 
        "html/partials/*.html",
        "html/entries/*.html",
    ))

    files, err := fs.ReadDir(entries, "html/entries")
    if err != nil {
        log.Println("Error reading directory:", err)
        return
    }
    htmlEntries = files
}

func renderPage(w http.ResponseWriter, r *http.Request, page string, data any) {
	page = filepath.Clean(page)
	page = strings.TrimPrefix(page, "/")
    
    err := htmlTemplates.ExecuteTemplate(w, page, data)
	if err != nil {
		log.Println(err)

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
    var postLinks PostLinks
    var posts []utils.Post
    for _, file := range htmlEntries {
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

func listGalleries(w http.ResponseWriter, r *http.Request) {
    var i interface{}
    renderPage(w, r, "galleries.html",i) 
}

func serveGallery() {
    // @todo
}

func readingList(w http.ResponseWriter, r *http.Request) {
    var i interface{}
    renderPage(w, r, "reading_list.html", i)
}
