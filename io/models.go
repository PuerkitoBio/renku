package io

import (
	"time"

	"github.com/PuerkitoBio/renku/config"
)

type Server struct {
	Port       int
	Root       string
	Categories []string
	StartTime  time.Time
}

func newServer() *Server {
	return &Server{
		config.Settings.Port,
		config.Settings.Root,
		nil,
		config.Settings.StartTime,
	}
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
	Text string
}

type IndexTemplateData struct {
	Server *Server
	Index  *Index
}

type PostTemplateData struct {
	Server *Server
	Post   *PostDetail
}
