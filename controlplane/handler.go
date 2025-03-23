package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type APIError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error code: %d", e.StatusCode)
}

func NewAPIError(statusCode int, msg any) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        msg,
	}
}

func InvalidRequest(errors map[string]string) APIError {
	return APIError{
		StatusCode: http.StatusBadRequest,
		Msg:        errors,
	}
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func Make(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			var apiErr APIError
			if errors.As(err, &apiErr) {
				_ = writeJSON(w, apiErr.StatusCode, apiErr)
			} else {
				errResp := NewAPIError(http.StatusInternalServerError, "internal server error")
				_ = writeJSON(w, errResp.StatusCode, errResp)
			}
			slog.Error("HTTP API error", "err", err.Error(), "path", r.URL.Path)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintln(w, "ok")
}

type scanHandler struct {
	log *slog.Logger
	is  ImageScanner
}

func (h *scanHandler) triggerScanHandler(w http.ResponseWriter, r *http.Request) error {
	img := r.URL.Query().Get("img")
	if img == "" {
		return InvalidRequest(map[string]string{
			"img": "needs to be a non blank query parameter",
		})
	}

	h.log.Info("should start here a new container/pod that is running the scan of the given image")
	if err := h.is.Scan(r.Context(), img); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(w, "triggering new scan for image %s\n", img)
	return nil
}
