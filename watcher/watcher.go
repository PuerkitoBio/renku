package watcher

import (
	"log"

	"github.com/PuerkitoBio/renku/iface"
	"github.com/howeyc/fsnotify"
)

type Watcher struct {
	ev   chan iface.FileEvent
	dir  string
	w    *fsnotify.Watcher
	stop chan struct{}
}

func New(dir string) *Watcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("could not create watcher: %s", err)
	}
	return &Watcher{
		ev:   make(chan iface.FileEvent),
		dir:  dir,
		w:    w,
		stop: make(chan struct{}),
	}
}

func (ø *Watcher) Event() <-chan iface.FileEvent {
	return ø.ev
}

func (ø *Watcher) Dir() string {
	return ø.dir
}

type FileEvent struct {
	*fsnotify.FileEvent
}

func (ø *FileEvent) Name() string {
	return ø.FileEvent.Name
}

func (ø *Watcher) processEvents() {
	for {
		select {
		case ev := <-ø.w.Event:
			ø.ev <- &FileEvent{ev}
			log.Println(ev)
		case err := <-ø.w.Error:
			log.Println(err)
		case <-ø.stop:
			return
		}
	}
}

func (ø *Watcher) Start() {
	// Start processing events
	go ø.processEvents()

	// Watch the directory
	if err := ø.w.Watch(ø.dir); err != nil {
		log.Fatalf("error watching %s: %s", ø.dir, err)
	}
}

func (ø *Watcher) Stop() {
	ø.w.Close()
	close(ø.stop)
}
