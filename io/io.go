package io

import (
	"io/ioutil"
	"os"
	"sync"

	"github.com/PuerkitoBio/renku/config"
	"github.com/russross/blackfriday"
)

// TODO : As-is, no caching of data nor response, ~500-700 TPS with Siege/OSX

type BlogReader struct {
	mu         sync.RWMutex
	serverData *Server
}

func (ø *BlogReader) ensureServerCreated() {
	ø.mu.RLock()
	if ø.serverData == nil {
		ø.mu.RUnlock()
		ø.mu.Lock()
		ø.serverData = &Server{
			config.Settings.Port,
			config.Settings.Root,
			nil,
			config.Settings.StartTime,
		}
		ø.mu.Unlock()
	} else {
		ø.mu.RUnlock()
	}
}

func (ø *BlogReader) GetPost(postPath string) (interface{}, error) {
	ø.ensureServerCreated()

	if f, err := os.Open(postPath); err != nil {
		return nil, err
	} else {
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}

		// No need to lock serverData here, is necessarily present, and won't be
		// written by another thread.
		return &PostTemplateData{
			ø.serverData,
			&PostDetail{
				Post: &Post{
					Path: postPath,
				},
				Text: string(blackfriday.MarkdownCommon(b)),
			},
		}, nil
	}
}

func (ø *BlogReader) GetIndex() (interface{}, error) {
	return nil, nil
}
