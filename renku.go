package main

import (
	"log"

	"github.com/jessevdk/go-flags"
)

var (
	// The various command-line flags and options
	opts struct {
		Root           string `short:"d" long:"directory" description:"the root directory of the website" default:"./"`
		Port           int    `short:"p" long:"port" description:"the port to use for the web server" default:"9000"`
		Output         string `short:"o" long:"output" description:"output file for logging" default:"stdout"`
		ReqLogOutput   string `short:"r" long:"request-output" description:"output file for request logging" default:"stdout"`
		Verbose        bool   `short:"v" long:"verbose" description:"log everything"`
		Quiet          bool   `short:"q" long:"quiet" description:"don't log anything unless it's important"`
		NoCache        bool   `short:"C" long:"no-cache" description:"disable the response cache"`
		NoPrefillCache bool   `short:"P" long:"no-prefill-cache" description:"don't prefill the response cache"`
		CacheCap       int    `short:"c" long:"cache-capacity" description:"set the maximum number of items in the response cache"`
		NoWatch        bool   `short:"W" long:"no-watch" description:"disable the file watcher"`
	}
)

func main() {
	_, err := flags.Parse(&opts)
	if err == nil {
		log.Printf("%+v", opts)
	}
}
