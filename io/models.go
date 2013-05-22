package io

import (
	"time"
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
