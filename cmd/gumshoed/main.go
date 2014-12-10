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

// Get this flag set working!
var (
	tc         = gumshoe.TrackerConfig{}
	home       = os.Getenv("HOME")
	user       = os.Getenv("USER")
	gopath     = os.Getenv("GOPATH")
	gumshoeSrc = os.Getenv("GUMSHOESRC")
	gcstat     = debug.GCStats{}
)

type gumshoeSignals struct {
	configModified chan bool
	shutdown       chan bool
	// logger          chan Logger
	tcSignal   chan gumshoe.TrackerConfig
	showSignal chan *gumshoe.Shows
}

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
	//signals := new(gumshoeSignals)
	//signals.tcSignal <- tc

	allShows := gumshoe.NewShowsConfig()
	if numShows, err := allShows.LoadShows(tc); err == nil {
		log.Printf("You have %d shows that you are tracking.", numShows)
	}
	//signals.showSignal <- allShows
	log.Println("Starting up gumshoe...")

	// go StartMetrics()
	gumshoe.InitWatcher(tc, allShows)
	if tc.Operations.WatchMethod == "irc" {
		go gumshoe.StartIRC(tc)
	}
	gumshoe.StartHTTPServer(*baseDir, *port)
	log.Println("Exiting gumshoe...")
}
