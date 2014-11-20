package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/deekue/gumshoe/config"
	"github.com/deekue/gumshoe/irc"
	"github.com/deekue/gumshoe/watcher"
	"github.com/deekue/gumshoe/webui"
)

// HTTP Server Flags
var port = flag.String("p", "http",
	"Which port do we serve requests from. 0 allows the system to decide.")
var base_dir = flag.String("base_dir",
	filepath.Join(os.Getenv("HOME"), ".local", "gumshoe"),
	"Base path for the HTTP server's files.")

// Base Config Stuff
var config_file = flag.String("c",
	filepath.Join(os.Getenv("HOME"), ".gumshoe", "config.json"),
	"Location of the configuration file.")

// Get this flag set working!
var (
	tc         = config.TrackerConfig{}
	home       = os.Getenv("HOME")
	user       = os.Getenv("USER")
	gopath     = os.Getenv("GOPATH")
	gumshoeSrc = os.Getenv("GUMSHOESRC")
	gcstat     = debug.GCStats{}
)

type GumshoeSignals struct {
	config_modified chan bool
	shutdown        chan bool
	// logger          chan Logger
	tcSignal   chan config.TrackerConfig
	showSignal chan *config.Shows
}

func main() {
	flag.Parse()
	if gumshoeSrc == "" {
		gumshoeSrc = "/home/ryan/gocode/src/gumshoe"
	}
	if err := tc.LoadGumshoeConfig(*config_file); err != nil {
		log.Fatal(err)
	}
	if tc.Operations.HttpPort != *port && tc.Operations.HttpPort != "" {
		if err := flag.Set("p", tc.Operations.HttpPort); err != nil {
			log.Println(err)
		}
	}
	signals := new(GumshoeSignals)
	signals.tcSignal <- tc

	allShows := config.NewShowsConfig()
	if numShows, err := allShows.LoadShows(tc); err == nil {
		log.Printf("You have %d shows that you are tracking.", numShows)
	}
	signals.showSignal <- allShows
	log.Println("Starting up gumshoe...")

	// go StartMetrics()
	watcher.InitWatcher(tc, allShows)
	go webui.StartHttpServer(gumshoeSrc, *port)
	go irc.StartIRC(tc)
}
