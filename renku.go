package main

import (
	"log"
	"os"

	"github.com/PuerkitoBio/renku/config"
	"github.com/PuerkitoBio/renku/io"
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
		web.Reader = io.NewBlogReader()
		web.ListenAndServe()
	}
}
