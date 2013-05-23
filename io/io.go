package io

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/PuerkitoBio/renku/config"
	"github.com/russross/blackfriday"
)

type BlogReader struct {
	posts      map[string]*PostTemplateData
	serverData *Server
}

func NewBlogReader() *BlogReader {
	b := new(BlogReader)
	b.posts = make(map[string]*PostTemplateData)
	// Sync with file system
	b.createServer()
	b.readPosts()
	return b
}

func (ø *BlogReader) getPostData(fi os.FileInfo) (*PostTemplateData, error) {
	f, err := os.Open(path.Join(config.Settings.Root, config.Settings.PostsDir, fi.Name()))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return &PostTemplateData{
		Server: ø.serverData,
		Post: &PostDetail{
			Post: &Post{
				Path: f.Name(),
			},
			Text: string(blackfriday.MarkdownCommon(b)),
		}}, nil
}

func (ø *BlogReader) readPosts() {
	fis, err := ioutil.ReadDir(path.Join(config.Settings.Root, config.Settings.PostsDir))
	if err != nil {
		log.Println("error reading posts: ", err)
	}
	for _, fi := range fis {
		if ext := filepath.Ext(fi.Name()); ext == ".md" || ext == ".markdown" {
			pd, err := ø.getPostData(fi)
			if err != nil {
				log.Printf("error building post data for %s: %s\n", fi.Name(), err)
			} else {
				log.Printf("storing post %s\n", pd.Post.Path)
				ø.posts[pd.Post.Path] = pd
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
