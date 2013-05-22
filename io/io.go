package io

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/renku/config"
	"github.com/russross/blackfriday"
)

// TODO : As-is, no caching of data nor response, ~500-700 TPS with Siege/OSX

type BlogReader struct {
	posts      map[string]*PostDetail
	serverData *Server
}

func NewBlogReader() *BlogReader {
	b := new(BlogReader)
	b.posts = make(map[string]*PostDetail)
	// Sync with file system
	b.createServer()
	b.readPosts()
	return b
}

func (ø *BlogReader) getPostDetail(fi os.FileInfo) (*PostDetail, error) {
	f, err := os.Open(fi.Name())
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return &PostDetail{
		Post: &Post{
			Path: fi.Name(),
		},
		Text: string(blackfriday.MarkdownCommon(b)),
	}, nil
}

func (ø *BlogReader) readPosts() {
	fis, err := ioutil.ReadDir(filepath.Join(config.Settings.Root, config.Settings.PostsDir))
	if err != nil {
		log.Println("error reading posts: ", err)
	}
	for _, fi := range fis {
		if ext := filepath.Ext(fi.Name()); ext == ".md" || ext == ".markdown" {
			pd, err := ø.getPostDetail(fi)
			if err != nil {
				log.Printf("error building post detail for %s: %s\n", fi.Name(), err)
			} else {
				ø.posts[fi.Name()] = pd
			}
		}
	}
}

func (ø *BlogReader) createServer() {
	ø.serverData = &Server{
		config.Settings.Port,
		config.Settings.Root,
		nil,
		config.Settings.StartTime,
	}
}

func (ø *BlogReader) GetPost(postPath string) (interface{}, error) {
	if pd, ok := ø.posts[postPath]; ok {
		return pd, nil
	}

	log.Printf("post data not cached for %s\n", postPath)
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
