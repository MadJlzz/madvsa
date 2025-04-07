package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	s := chi.URLParam(r, "scanner")
	h.log.Info("using scanner", "scanner", s)

	img := r.URL.Query().Get("img")
	if img == "" {
		return InvalidRequest(map[string]string{
			"img": "needs to be a non blank query parameter",
		})
	}

	requestId := middleware.GetReqID(r.Context())
	h.log.Info("should start here a new container/pod that is running the scan of the given image")
	if err := h.is.Scan(r.Context(), requestId, s, img); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(w, "triggering new scan using %s scanner for image %s\n", s, img)
	return nil
}
