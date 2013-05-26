package main

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/PuerkitoBio/purell"
	"github.com/PuerkitoBio/renku/cache"
	"github.com/PuerkitoBio/renku/config"
	"github.com/PuerkitoBio/renku/iface"
	"github.com/PuerkitoBio/renku/io"
	"github.com/PuerkitoBio/renku/watcher"
	"github.com/PuerkitoBio/renku/web"
	"github.com/jessevdk/go-flags"
)

func main() {
	_, err := flags.Parse(&config.Settings)
	if err == nil {
		if config.Settings.Output != "" && config.Settings.Output != "stdout" {
			// Open file for logging
			f, err := os.Open(config.Settings.Output)
			if err != nil {
				log.Fatal("error opening specified output file:", err)
			}
			log.SetOutput(f)
			defer f.Close()
		}
		// Set up dependencies
		var postsWatcher iface.Watcher
		if !config.Settings.NoWatch {
			postsWatcher = watcher.New(path.Join(config.Settings.Root, config.Settings.PostsDir))
		}
		web.Reader = io.NewBlogReader(postsWatcher)
		web.CacheHandler = func(h http.Handler) http.Handler {
			return cache.LRUCacheHandler(h, config.Settings.CacheSz, purell.FlagsSafe)
		}

		web.ListenAndServe()
	}
}
