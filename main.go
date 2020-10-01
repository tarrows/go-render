package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/russross/blackfriday.v2"
)

const templatesDir = "templates"
const assetsDir = "assets"
const docsDir = "docs"

func main() {
	addr := flag.String("addr", ":2983", "The address of the application")

	assetHandler := http.FileServer(http.Dir(assetsDir))
	http.Handle("/favicon.ico", assetHandler)
	http.Handle("/assets", http.StripPrefix("/assets/", assetHandler))
	http.HandleFunc("/docs/", docsHandler)
	http.HandleFunc("/", indexHandler)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("[ListenAndServe]", err)
	}
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

type pageObject struct {
	Title string
	Body  template.HTML
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DocsHandler] %s %s", r.Method, r.URL)

	//     segs[0] / segs[1] / segs[2]
	// localhost:8080/docs/foo
	segs := strings.Split(r.URL.Path, "/")

	if len(segs) < 2 {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	p := filepath.Join(docsDir, segs[2])

	if !exists(p) {
		http.NotFound(w, r)
		return
	}

	file, err := os.Open(p)
	defer file.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	content, err := ioutil.ReadAll(file)

	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.HrefTargetBlank,
	})

	output := blackfriday.Run(
		content,
		blackfriday.WithExtensions(blackfriday.CommonExtensions),
		blackfriday.WithRenderer(renderer),
	)
	html := bluemonday.UGCPolicy().SanitizeBytes(output)

	po := pageObject{
		Title: strings.Replace(segs[2], ".md", "", 1),
		Body:  template.HTML(html),
	}

	target := filepath.Join(templatesDir, "markdown.html")
	t := template.Must(template.New("markdown.html").ParseFiles(target))
	_ = t.ExecuteTemplate(w, "markdown.html", po)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[IndexHandler] %s %s", r.Method, r.URL)

	target := filepath.Join(templatesDir, "index.html")
	t := template.Must(template.New("index.html").ParseFiles(target))
	_ = t.ExecuteTemplate(w, "index.html", nil)
}
