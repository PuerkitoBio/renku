package watcher

import (
	"log"
	"path"

	"github.com/PuerkitoBio/renku/config"
	"github.com/howeyc/fsnotify"
)

var (
	stop chan struct{}
	w    *fsnotify.Watcher
)

func processEvents() {
	for {
		select {
		case ev := <-w.Event:
			log.Println(ev)
		case err := <-w.Error:
			log.Println(err)
		case <-stop:
			return
		}
	}
}

func Start() {
	var err error

	// Error if started twice
	if w != nil {
		log.Print("watcher already started")
		return
	}

	// Create watcher
	w, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("error creating watcher: ", err)
	}

	// Start processing events
	stop = make(chan struct{})
	go processEvents()

	// Watch the posts directory
	if err = w.Watch(path.Join(config.Settings.Root, config.Settings.PostsDir)); err != nil {
		log.Fatalf("error watching %s: %s", config.Settings.PostsDir, err)
	}
}

func Stop() {
	if w != nil {
		w.Close()
		close(stop)
	}
}
