package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "ok")
}

type scanHandler struct {
	log *slog.Logger
	is  ImageScanner
}

func (h *scanHandler) triggerScanHandler(w http.ResponseWriter, r *http.Request) {
	img := r.URL.Query().Get("img")
	if img == "" {
		http.Error(w, "please provide a valid image name", http.StatusBadRequest)
		return
	}

	h.log.Info("should start here a new container/pod that is running the scan of the given image")

	if err := h.is.Scan(r.Context(), img); err != nil {
		h.log.Error("something bad happened during image scan", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprintf(w, "triggering new scan for image %s\n", img)
}
