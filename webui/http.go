package webui

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-martini/martini"
)

func getShows() string {
	return "getShows"
}

func newShow() string {
	return "newShow"
}

func updateShow() string {
	return "updateShow"
}

func deleteShow() string {
	return "deleteShow"
}

func getConfigs() string {
	return "getConfigs"
}

func newConfig() string {
	return "newConfig"
}

func updateConfig() string {
	return "updateConfig"
}

func deleteConfig() string {
	return "deleteConfig"
}

func getQueueItems() string {
	return "getQueueItems"
}

func newQueueItem() string {
	return "newQueueItem"
}

func deleteQueueItem() string {
	return "deleteQueueItem"
}

// StartHTTPServer start a HTTP server for configuration and monitoring
func StartHTTPServer(port string) {
	var hostString = fmt.Sprintf(":%s", port)
	var m = martini.Classic()

	static := martini.Static("www", martini.StaticOptions{Fallback: "/index.html", Exclude: "/api"})
	m.NotFound(static, http.NotFound)

	m.Group("/api/show", func(r martini.Router) {
		r.Get("/:id", getShows)
		r.Post("/new", newShow)
		r.Put("/update/:id", updateShow)
		r.Delete("/delete/:id", deleteShow)
	})

	m.Group("/api/config", func(r martini.Router) {
		r.Get("/:id", getConfigs)
		r.Post("/new", newConfig)
		r.Put("/update/:id", updateConfig)
		r.Delete("/delete/:id", deleteConfig)
	})

	m.Group("/api/queue", func(r martini.Router) {
		r.Get("/:id", getQueueItems)
		r.Post("/new", newQueueItem)
		r.Delete("/delete/:id", deleteQueueItem)
	})

	log.Println("Starting up webserver...")
	m.RunOnAddr(hostString)
}
