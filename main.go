package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed html
var htmlPages embed.FS

//go:embed static
var staticFiles embed.FS

func main() {
    fs, err := fs.Sub(staticFiles, "static")
    if err != nil {
        log.Fatal(err)
    }
    http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(fs))))

    http.HandleFunc("GET /entries/", serveEntry)
    http.HandleFunc("GET /entries/{entry_url}", serveEntry)

	http.HandleFunc("GET /", serveRoot)

    log.Println("Listening on http://localhost:42069")
	err = http.ListenAndServe(":42069", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func renderPage(w http.ResponseWriter, _ *http.Request, page string) {
	page = filepath.Clean(page)
	page = strings.TrimPrefix(page, "/")
    
	view := template.Must(template.ParseFS(
        htmlPages, 
        "html/*.html", 
        "html/partials/*.html",
        "html/entries/*.html",
    ))

    err := view.ExecuteTemplate(w, page, "")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		r.URL.Path = "index.html"
    }

	renderPage(w, r, r.URL.Path)
}

func serveEntry(w http.ResponseWriter, r *http.Request) {
    entry := r.PathValue("entry_url")

    if entry == "" {
        log.Println("entries root")
        renderPage(w, r, "/entries.html")
    } else {
        page := entry + ".html"
        log.Printf("handling entry with id=%v\n", entry)
        log.Printf("actual page: %s\n", page)
        renderPage(w, r, page)
    }
}
