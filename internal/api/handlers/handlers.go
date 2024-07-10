package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type j = map[string]any

func response(message j, code int, w http.ResponseWriter) error {
	w.WriteHeader(code)
	if message == nil {
		return nil
	}
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}
	if _, err = w.Write(b); err != nil {
		return err
	}
	return nil
}

func logError(logger *slog.Logger, r *http.Request, err error) {
	logger.Error(
		"Error while responding",
		"URL", r.URL,
		"METHOD", r.Method,
		"error", err,
	)
}
