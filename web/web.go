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
)

const (
	faviconCacheTTL = 30 * 24 * time.Hour // 30 days
)

type BlogReader interface {
	GetPost(string) (interface{}, error)
	GetIndex() (interface{}, error)
}

var (
	// Dependencies, injected by the executable (renku.go)
	Reader       BlogReader
	CacheHandler func(http.Handler) http.Handler

	pubDir string
	pstDir string
	tplDir string
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func servePost(w http.ResponseWriter, r *http.Request) {
	if data, err := Reader.GetPost(path.Join(pstDir, r.URL.Path)); err != nil {
		log.Print("!", err)
		http.NotFound(w, r)
	} else {
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
	var h http.Handler

	// Store common directories, for convenience
	pstDir = path.Join(config.Settings.Root, config.Settings.PostsDir)
	pubDir = path.Join(config.Settings.Root, config.Settings.PublicDir)
	tplDir = path.Join(config.Settings.Root, config.Settings.TemplatesDir)

	// Compile templates
	if err := templates.CompileDir(tplDir); err != nil {
		log.Fatal("error compiling templates", err)
	}

	// Handle paths
	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(pubDir))))
	mux.HandleFunc("/", servePage)
	h = mux

	// If the cache is enabled, hook it up
	if !config.Settings.NoCache {
		if CacheHandler == nil {
			log.Fatal("cache handler is nil")
		}
		h = CacheHandler(mux)
	}

	// Setup handlers chain
	h = handlers.FaviconHandler(
		handlers.PanicHandler(
			handlers.LogHandler(
				h,
				&handlers.LogOptions{
					Format: handlers.Lshort,
				}),
			nil),
		path.Join(pubDir, "favicon.ico"),
		faviconCacheTTL)

	// Start listening
	config.Settings.StartTime = time.Now()
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Settings.Port), h); err != nil {
		log.Fatal("^", err)
	}
}
