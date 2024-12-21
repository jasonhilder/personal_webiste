package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jasonhilder/personal_website/internal/utils"
)

//go:embed html
var htmlPages embed.FS

//go:embed static
var staticFiles embed.FS

var htmlTemplates *template.Template
var htmlEntries []fs.DirEntry

func main() {
    _, isDebug := os.LookupEnv("DEBUG")
    if !isDebug {
        f, err := os.OpenFile("/var/log/personal_website.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
        if err != nil {
            log.Fatalf("error opening file: %v", err)
        }
        defer f.Close()
        log.SetOutput(f)
    }
    
    loadHtmlFiles()

    fs, err := fs.Sub(staticFiles, "static")
    if err != nil {
        log.Fatal(err)
    }
    http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(fs))))

    http.HandleFunc("GET /music/", InitSpotify(GetSpotifyInfo))

    http.HandleFunc("GET /gists/", listGists)

    http.HandleFunc("GET /gists/{gist_id}", serveGist)

    http.HandleFunc("GET /books/", bookList)

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

func RenderPage(w http.ResponseWriter, r *http.Request, page string, data any) {
    page = filepath.Clean(page)
    page = strings.TrimPrefix(page, "/")

    err := htmlTemplates.ExecuteTemplate(w, page, data)
    if err != nil {
        log.Println(err)

        IPAddress := r.Header.Get("X-Real-Ip")
        if IPAddress == "" {
            IPAddress = r.Header.Get("X-Forwarded-For")
        }
        if IPAddress == "" {
            IPAddress = r.RemoteAddr
        }
        log.Printf("404 - Route Not Found: %s, IP: %s\n", r.URL.Path, IPAddress)
        var i interface{}
        RenderPage(w, r, "404.html", i)
    }
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		r.URL.Path = "index.html"
    }

    var i interface{}
	RenderPage(w, r, r.URL.Path, i)
}

func listGists(w http.ResponseWriter, r *http.Request) {
    var i interface{}
    RenderPage(w, r, "gists.html", i)
}

func serveGist(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("gist_id")

    i := utils.GetGistId(id)
	RenderPage(w, r, "detail.html", i)
}

func bookList(w http.ResponseWriter, r *http.Request) {
    var i interface{}
    RenderPage(w, r, "book_list.html", i)
}
