package main

import (
	"errors"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"log/slog"
	"net/http"
	"os"
)

var socketPath string
var orchestrationMode string

func init() {
	const (
		defaultSocket                 = "unix:///var/run/docker.sock"
		defaultSocketUsage            = "container runtime socket to use when running scans as containers"
		defaultOrchestrationMode      = "container"
		defaultOrchestrationModeUsage = "decide whether scans are running on Kubernetes or via container runtimes as containers"
	)
	flag.StringVar(&socketPath, "socket", defaultSocket, defaultSocketUsage)
	flag.StringVar(&socketPath, "s", defaultSocket, defaultSocketUsage+" (shorthand)")
	flag.StringVar(&orchestrationMode, "orchestration", defaultOrchestrationMode, defaultOrchestrationModeUsage)
}

func main() {
	flag.Parse()
	mr := chi.NewRouter()
	mr.Use(middleware.Logger)
	mr.Use(middleware.Recoverer)

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	var is ImageScanner
	switch orchestrationMode {
	case "container":
		is = NewContainerService(socketPath)
	default:
		panic(errors.New("not implemented yet"))
	}

	// Step 1: parse the service mode to know what kind of orchestrator we run.
	// Step 2:
	sh := &scanHandler{is: is, log: l}

	// API versioned routes
	mr.Route("/api/v1", func(r chi.Router) {

		r.Get("/health", healthHandler)

		r.Route("/scanner", func(r chi.Router) {
			// Not sure that's the best way to support multiple scanner, but for now it's okay.
			r.Post("/{scanner:^(trivy|grype)$}/trigger", Make(sh.triggerScanHandler))
		})

	})

	if err := http.ListenAndServe(":3000", mr); err != nil {
		log.Fatal(err)
	}
}
