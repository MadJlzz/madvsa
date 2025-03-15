package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	mr := chi.NewRouter()
	mr.Use(middleware.Logger)
	mr.Use(middleware.Recoverer)

	cs := NewContainerService("unix:///var/run/docker.sock")
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	sh := &scanHandler{is: cs, log: l}

	// API versioned routes
	mr.Route("/api/v1", func(r chi.Router) {

		r.Get("/health", healthHandler)

		r.Route("/scans", func(r chi.Router) {
			r.Post("/trigger", sh.triggerScanHandler)
		})

	})

	if err := http.ListenAndServe(":3000", mr); err != nil {
		log.Fatal(err)
	}
}
