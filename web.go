package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/PuerkitoBio/ghost/handlers"
	"github.com/PuerkitoBio/ghost/templates"
	_ "github.com/PuerkitoBio/ghost/templates/amber"
)

type logMode int

const (
	lmNormal logMode = iota
	lmQuiet
	lmVerbose

	faviconCacheTTL = 30 * 24 * time.Hour
)

type serverOptions struct {
	Port      int
	Root      string
	Posts     string
	Templates string
	Public    string
	LogMode   logMode
	Watch     bool
}

var (
	webServerOpts serverOptions
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func servePost(w http.ResponseWriter, r *http.Request) {
	if f, err := os.Open(path.Join(webServerOpts.Root, webServerOpts.Posts, r.URL.Path)); err != nil {
		http.NotFound(w, r)
		return
	} else {
		defer f.Close()
	}
	if err := templates.Render("testdata/templates/post.amber", w, nil); err != nil {
		if err == templates.ErrTemplateNotExist {
			http.NotFound(w, r)
			return
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func servePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		serveIndex(w, r)
	} else {
		servePost(w, r)
	}
}

func listenAndServe() {
	// Compile templates
	if err := templates.CompileDir(path.Join(webServerOpts.Root, webServerOpts.Templates)); err != nil {
		log.Fatal("error compiling templates", err)
	}

	mux := http.NewServeMux()
	// TODO : Eventually, will go through cache first
	mux.Handle("/public/", http.StripPrefix("/public/",
		http.FileServer(http.Dir(path.Join(webServerOpts.Root, webServerOpts.Public)))))
	mux.HandleFunc("/", servePage)

	h := handlers.FaviconHandler(
		handlers.PanicHandler(
			handlers.LogHandler(
				mux,
				&handlers.LogOptions{
					Format: handlers.Lshort,
				}),
			nil),
		path.Join(webServerOpts.Root, path.Join(webServerOpts.Root, webServerOpts.Public, "favicon.ico")),
		faviconCacheTTL)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", webServerOpts.Port), h); err != nil {
		log.Fatal("^", err)
	}
}
