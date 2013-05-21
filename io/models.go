package io

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/PuerkitoBio/renku/config"
)

type Server struct {
	Port       int
	Root       string
	Categories []string
	StartTime  time.Time
}

type Index struct {
	Path     string
	Category string
	Posts    []Post
}

type Post struct {
	Path    string
	Title   string
	Author  string
	Lead    string
	PubTime time.Time
	ModTime time.Time
}

type PostDetail struct {
	*Post
	Text []byte
}

type IndexTemplateData struct {
	Server *Server
	Index  *Index
}

type PostTemplateData struct {
	Server *Server
	Post   *PostDetail
}

var (
	startTime  = time.Now()
	serverData *Server
)

func ensureServerCreated() {
	if serverData == nil {
		serverData = &Server{
			config.Settings.Port,
			config.Settings.Root,
			nil,
			startTime,
		}
	}
}

func GetPostData(postPath string) (*PostTemplateData, error) {
	ensureServerCreated()

	if f, err := os.Open(postPath); err != nil {
		return nil, err
	} else {
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		return &PostTemplateData{
			serverData,
			&PostDetail{
				Post: &Post{
					Path: postPath,
				},
				Text: b,
			},
		}, nil
	}
}
