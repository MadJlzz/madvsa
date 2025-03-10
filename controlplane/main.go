package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	mr := chi.NewRouter()
	mr.Use(middleware.Logger)
	mr.Use(middleware.Recoverer)

	// API versioned routes
	mr.Route("/api/v1", func(r chi.Router) {

		r.Get("/health", healthHandler)

		r.Route("/scans", func(r chi.Router) {
			r.Post("/trigger", triggerScanHandler)
		})

	})

	if err := http.ListenAndServe(":3000", mr); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "ok")
}

func triggerScanHandler(w http.ResponseWriter, r *http.Request) {
	img := r.URL.Query().Get("img")
	if img == "" {
		http.Error(w, "please provide a valid image name", http.StatusBadRequest)
		return
	}
	log.Println("should start here a new container/pod that is running the scan of the given image")

	cs := NewContainerService("unix:///var/run/docker.sock")
	if err := cs.ScanContainer(r.Context(), img); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprintf(w, "triggering new scan for image %s\n", img)
}
