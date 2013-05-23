package config

import (
	"time"
)

var Settings struct {
	Root           string `short:"d" long:"directory" description:"the root directory of the website" default:"./"`
	Port           int    `short:"p" long:"port" description:"the port to use for the web server" default:"9000"`
	Output         string `short:"o" long:"output" description:"output file for logging" default:"stdout"`
	ReqLogOutput   string `short:"r" long:"request-output" description:"output file for request logging" default:"stdout"`
	Verbose        bool   `short:"v" long:"verbose" description:"log everything"`
	Quiet          bool   `short:"q" long:"quiet" description:"don't log anything unless it's important"`
	NoCache        bool   `short:"C" long:"no-cache" description:"disable the response cache"`
	NoPrefillCache bool   `short:"P" long:"no-prefill-cache" description:"don't prefill the response cache"`
	CacheSz        int    `short:"c" long:"cache-size" description:"set the maximum number of items in the response cache"`
	NoWatch        bool   `short:"W" long:"no-watch" description:"disable the file watcher"`
	TemplatesDir   string
	PublicDir      string
	PostsDir       string
	DraftsDir      string
	StartTime      time.Time
}

func init() {
	Settings.TemplatesDir = "templates/"
	Settings.PublicDir = "public/"
	Settings.PostsDir = "posts/"
	Settings.DraftsDir = "drafts/"
}
