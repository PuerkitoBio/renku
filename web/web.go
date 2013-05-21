package web

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/PuerkitoBio/ghost/handlers"
	"github.com/PuerkitoBio/ghost/templates"
	_ "github.com/PuerkitoBio/ghost/templates/amber"
	"github.com/PuerkitoBio/renku/config"
	"github.com/PuerkitoBio/renku/io"
)

const (
	faviconCacheTTL = 30 * 24 * time.Hour
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func servePost(w http.ResponseWriter, r *http.Request) {
	if data, err := io.GetPostData(path.Join(config.Settings.Root,
		config.Settings.PostsDir, r.URL.Path)); err != nil {

		log.Print("!", err)
		http.NotFound(w, r)
	} else {
		log.Printf("? %#v", data)
		if err := templates.Render("post.amber", w, data); err != nil {
			log.Print("!", err)
			if err == templates.ErrTemplateNotExist {
				http.NotFound(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
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

func ListenAndServe() {
	// Compile templates
	if err := templates.CompileDir(path.Join(config.Settings.Root, config.Settings.TemplatesDir)); err != nil {
		log.Fatal("error compiling templates", err)
	}

	mux := http.NewServeMux()
	// TODO : Eventually, will go through cache first
	mux.Handle("/public/", http.StripPrefix("/public/",
		http.FileServer(http.Dir(path.Join(config.Settings.Root, config.Settings.PublicDir)))))
	mux.HandleFunc("/", servePage)

	h := handlers.FaviconHandler(
		handlers.PanicHandler(
			handlers.LogHandler(
				mux,
				&handlers.LogOptions{
					Format: handlers.Lshort,
				}),
			nil),
		path.Join(config.Settings.Root, config.Settings.PublicDir, "favicon.ico"),
		faviconCacheTTL)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Settings.Port), h); err != nil {
		log.Fatal("^", err)
	}
}
