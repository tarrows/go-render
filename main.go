package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

const templatesDir = "templates"
const assetsDir = "assets"

func main() {
	addr := flag.String("addr", ":2983", "The address of the application")

	assetHandler := http.FileServer(http.Dir(assetsDir))
	http.Handle("/favicon.ico", assetHandler)
	http.Handle("/assets", http.StripPrefix("/assets/", assetHandler))
	http.HandleFunc("/docs/", docsHandler)
	http.HandleFunc("/", indexHandler)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DocsHandler] %s %s", r.Method, r.URL)

	//     segs[0] / segs[1] / segs[2]
	// localhost:8080/docs/foo
	segs := strings.Split(r.URL.Path, "/")

	if len(segs) < 2 {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	target := filepath.Join(templatesDir, "index.html")
	t := template.Must(template.New("index.html").ParseFiles(target))
	_ = t.ExecuteTemplate(w, "index.html", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[IndexHandler] %s %s", r.Method, r.URL)

	target := filepath.Join(templatesDir, "index.html")
	t := template.Must(template.New("index.html").ParseFiles(target))
	_ = t.ExecuteTemplate(w, "index.html", nil)
}
