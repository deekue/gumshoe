package webui

import (
	"log"
	"net/http"
	"path"
	"time"
)

func GumshoeHandlers(base_dir string) http.Handler {
	gumshoe_handlers := http.NewServeMux()
	gumshoe_handlers.Handle("/", http.FileServer(http.Dir(path.Join(base_dir, "html"))))
	return gumshoe_handlers
}

func StartHttpServer(base_dir string, port string) {
	s := &http.Server{
		Addr:           "127.0.0.1:" + port,
		Handler:        GumshoeHandlers(base_dir),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Starting up webserver...")
	log.Fatal(s.ListenAndServe())
}
