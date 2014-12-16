package gumshoe

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// The primary structure holding the config data.

type IMDBConfig struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Uid  int    `json:"uid"`
}

type IRCChannel struct {
	Nick         string `json:"nick"`
	Key          string `json:"key"`
	Server       string `json:"server"`
	InviteCmd    string `json:"invite_cmd"`
	WatchChannel string `json:"watch_channel"`
	KeepAlive    int    `json:"keep_alive"`
	PingFreq     int    `json:"ping_frequency"`
	IRCPort      int    `json:"irc_port"`
	Timeout      int    `json:"timeout"`
}

type RSSChannel struct {
	FeedURI      string `json:"feed"`
	Passkey      string `json:"passkey"`
	Uid          string `json:"rss_uid"`
	RssTtl       int    `json:"rss_ttl"`
	UseServerTtl bool   `json:"use_server_ttl"`
}

type Operations struct {
	EnableLogging bool            `json:"enable_logging"`
	EnableWeb     bool            `json:"enable_web"`
	HttpPort      string          `json:"http_port"`
	LogLevel      string          `json:"log_level"`
	UseIMDB       bool            `json:"use_imdb_watchlist"`
	WatchMethods  map[string]bool `json:"watch_methods"`
}

type TrackerConfig struct {
	Cookiejar    map[string]string `json:"cookiejar"`
	Files        map[string]string `json:"file_options"`
	IMDB         IMDBConfig
	IRC          IRCChannel
	Operations   Operations
	RSS          RSSChannel
	Tracker      map[string]interface{} `json:"tracker"`
	LastModified int                    `json:"last_modified"`
}

func NewTrackerConfig() *TrackerConfig {
	return &TrackerConfig{}
}

func (tc *TrackerConfig) LoadGumshoeConfig(cfgFile string) error {
	if err := tc.ProcessGumshoeJSON(cfgFile); err != nil {
		log.Println("Error with config file ", cfgFile, ": ", err)
		log.Println("Using basic template.")
		return tc.ProcessGumshoeJSON("config/gumshoe_config.json")
	}
	return nil
}

func (tc *TrackerConfig) ProcessGumshoeJSON(cfgJson string) error {
	if cfgBuf, err := ioutil.ReadFile(cfgJson); err != nil {
		return err
	} else {
		if err := json.Unmarshal(cfgBuf, &tc); err != nil {
			return err
		}
	}
	// config <- true
	return nil
}

func (tc *TrackerConfig) WriteGumshoeConfig(update []byte) error {
	err := json.Unmarshal(update, &tc)
	if err == nil {
		var gCfg []byte
		gCfg, err := json.MarshalIndent(&tc, "", "  ")
		if err == nil {
			return ioutil.WriteFile(tc.Files["base_dir"]+"gumshoe_config.json",
				gCfg, 0655)
		}
		return err
	}
	return err
}
