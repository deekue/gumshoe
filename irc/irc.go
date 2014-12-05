package irc

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/deekue/gumshoe/config"
	"github.com/deekue/gumshoe/watcher"
	"github.com/thoj/go-ircevent"
)

var ircClient *irc.Connection

//var IRCEnabled = make(chan bool)
var announceLine, episodePattern *regexp.Regexp

func init() {
	// Metrics

	// don't immediately launch the IRC client
	//IRCEnabled <- false

	// TODO(ryan): make this configurable
	announceLine = regexp.MustCompile("BitMeTV-IRC2RSS: (?P<title>.*?) : (?P<url>.*)")
	episodePattern = regexp.MustCompile("^([\\w\\d\\s.]+)[. ](?:s(\\d{1,2})e(\\d{1,2})|(\\d)x?(\\d{2})|Star.Wars)([. ])")
}

// should this be refactored so that it can reconnect on config changes instead of diconnect and
// connect. TODO(ryan)
func connectToTrackerIRC(tc config.TrackerConfig, ircClient *irc.Connection) {
	// Give the connection the configured defaults
	//ircClient.KeepAlive = time.Duration(tc.IRC.KeepAlive) * time.Minute
	//ircClient.Timeout = time.Duration(tc.IRC.Timeout) * time.Minute
	//ircClient.PingFreq = time.Duration(tc.IRC.PingFreq) * time.Minute
	ircClient.Password = tc.IRC.Key
	ircClient.AddCallback("invite", func(e *irc.Event) {
		if strings.Index(e.Raw, tc.IRC.WatchChannel) != -1 {
			ircClient.Join(tc.IRC.WatchChannel)
		}
	})
	ircClient.AddCallback("public", matchAnnounce)
	var server = fmt.Sprintf("%s:%d", tc.IRC.Server, int(tc.IRC.IRCPort))
	log.Println("Connecting to IRC server...")
	ircClient.Connect(server)
	time.Sleep(60)
	ircClient.SendRawf(tc.IRC.InviteCmd, tc.IRC.Nick, tc.IRC.Key)
}

func matchAnnounce(e *irc.Event) {
	aMatch := announceLine.FindStringSubmatch(e.Raw)
	log.Println(aMatch)
	if aMatch != nil {
		eMatch := episodePattern.FindStringSubmatch(aMatch[1])
		if eMatch != nil {
			err := watcher.IsNewEpisode(eMatch)
			if err == nil {
				log.Println("This is where we would pick up the new episode.")
				// go RetrieveEpisode(aMatch[2])
				return
			}
			log.Println(err)
		}
	}
}

//func EnableIRC(tc config.TrackerConfig, irc_client *irc.Connection) {
//	for {
//		run := <-irc_enabled
//		if run {
//			log.Println("starting up IRC client.")
//			ConnectToTrackerIRC(tc, irc_client)
//		}
//	}
//}

//func DisableIRC(irc_client *irc.Connection) {
//	for {
//		run := <-irc_enabled
//		if !run {
//			log.Println("stopping the IRC client.")
//			irc_client.Disconnect()
//		}
//	}
//}

// StartIRC kick off the IRC client
func StartIRC(tc config.TrackerConfig) {
	ircClient := irc.IRC(tc.IRC.Nick, tc.IRC.Nick)
	// go WatchIRCConfig(signals)
	// go UpdateLog()
	//	for {
	//		go EnableIRC(tc, irc_client)
	//		go DisableIRC(irc_client)
	//		// go WatchIRCConfig(signals)
	//		// go UpdateLog()
	//		if tc.Operations.WatchMethod == "irc" {
	//			irc_enabled <- true
	//		} else {
	//			irc_enabled <- false
	//		}
	//	}
	connectToTrackerIRC(tc, ircClient)
}
