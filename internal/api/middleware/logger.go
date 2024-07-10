package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type loggerWriter struct {
	http.ResponseWriter
	code int
}

func wrapWriter(w http.ResponseWriter) *loggerWriter {
	return &loggerWriter{ResponseWriter: w}
}

func (l *loggerWriter) WriteHeader(code int) {
	l.code = code
	l.ResponseWriter.WriteHeader(code)
}

type Logger struct {
	logger *slog.Logger
}

func NewLoggerMiddleware(logger *slog.Logger) Logger {
	return Logger{logger: logger}
}

func (l Logger) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			lw := wrapWriter(w)
			timeBefore := time.Now()
			next.ServeHTTP(lw, r)
			l.logger.Info(
				"Incoming request",
				"PATH", r.URL.Path,
				"METHOD", r.Method,
				"STATUS", lw.code,
				"DURATION", fmt.Sprintf("%dms", time.Since(timeBefore).Milliseconds()),
			)
		},
	)
}
