package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/deekue/gumshoe/gumshoe"
)

// HTTP Server Flags
var port = flag.String("p", "http",
	"Which port do we serve requests from. 0 allows the system to decide.")
var baseDir = flag.String("d",
	filepath.Join(os.Getenv("HOME"), ".local", "gumshoe"),
	"Base path for gumshoe.")

// Base Config Stuff
var configFile = flag.String("c",
	filepath.Join(os.Getenv("HOME"), ".gumshoe", "config.json"),
	"Location of the configuration file.")

// TODO Get this flag set working!
var (
	tc         = gumshoe.TrackerConfig{}
	home       = os.Getenv("HOME")
	user       = os.Getenv("USER")
	gopath     = os.Getenv("GOPATH")
	gumshoeSrc = os.Getenv("GUMSHOESRC")
	gcstat     = debug.GCStats{}
)

func main() {
	flag.Parse()
	if err := tc.LoadGumshoeConfig(*configFile); err != nil {
		log.Fatal(err)
	}
	if tc.Operations.HttpPort != *port && tc.Operations.HttpPort != "" {
		if err := flag.Set("p", tc.Operations.HttpPort); err != nil {
			log.Println(err)
		}
	}

	log.Println("Starting up gumshoe...")
	gumshoe.InitShowDb(*baseDir)

	//DEBUG
	gumshoe.LoadTestData()

	// start enabled watchers
	for k, v := range tc.Operations.WatchMethods {
		if v {
			switch k {
			case "rss":
				log.Println("starting RSS watcher")
			case "irc":
				log.Println("starting IRC watcher")
				go gumshoe.StartIRC(tc)
			case "log":
				log.Println("starting log file watcher")
			}
		}
	}

	gumshoe.StartHTTPServer(*baseDir, *port)
	log.Println("Exiting gumshoe...")
}
