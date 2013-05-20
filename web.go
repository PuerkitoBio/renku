package main

import (
	"fmt"
	"log"
	"net/http"
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
	Port    int
	Root    string
	LogMode logMode
	Watch   bool
}

func servePage(w http.ResponseWriter, r *http.Request) {
	if err := templates.Render("testdata/templates/post.amber", w, nil); err != nil {
		if err == templates.ErrTemplateNotExist {
			http.NotFound(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func listenAndServe(opts serverOptions) {
	// Compile templates
	if err := templates.CompileDir(path.Join(opts.Root, "templates/")); err != nil {
		log.Fatal("error compiling templates", err)
	}

	mux := http.NewServeMux()
	// TODO : Eventually, will go through cache first
	mux.Handle("/public/", http.StripPrefix("/public/",
		http.FileServer(http.Dir(path.Join(opts.Root, "public/")))))
	mux.HandleFunc("/", servePage)

	h := handlers.FaviconHandler(
		handlers.PanicHandler(
			handlers.LogHandler(
				mux,
				&handlers.LogOptions{
					Format: handlers.Lshort,
				}),
			nil),
		path.Join(opts.Root, "public/favicon.ico"),
		faviconCacheTTL)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", opts.Port), h); err != nil {
		log.Fatal("^", err)
	}
}
