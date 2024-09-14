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

    http.HandleFunc("GET /entries/{gist_id}", serveEntry)

    http.HandleFunc("GET /galleries/", listGalleries)

    http.HandleFunc("GET /reading_list/", readingList)

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
    ))    
}

func renderPage(w http.ResponseWriter, r *http.Request, page string, data any) {
	page = filepath.Clean(page)
	page = strings.TrimPrefix(page, "/")
    
    err := htmlTemplates.ExecuteTemplate(w, page, data)
	if err != nil {
		log.Println(err)

        var i interface{}
        renderPage(w, r, "404.html", i)
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
    var i interface{}
    renderPage(w, r, "entries.html", i)
}

func serveEntry(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("gist_id")

    i := utils.GetGistId(id);
	renderPage(w, r, "detail.html", i)
}

func listGalleries(w http.ResponseWriter, r *http.Request) {
    var i interface{}
    // @todo
    renderPage(w, r, "galleries.html",i) 
}

func serveGallery() {
    // @todo
}

func readingList(w http.ResponseWriter, r *http.Request) {
    var i interface{}
    renderPage(w, r, "reading_list.html", i)
}
