package main

import (
	"github.com/jessevdk/go-flags"
	"log"
)

var (
	opts struct {
		Verbose bool `short:"v" long:"verbose" description:"log everything"`
	}
)

func main() {
	args, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal("^", err)
	}
	log.Print(args)
}
