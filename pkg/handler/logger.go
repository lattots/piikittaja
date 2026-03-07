package handler

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func Log(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{w, http.StatusOK}

		handlerFunc(rec, r)

		duration := time.Since(start).Milliseconds()
		statusColor := getStatusColor(rec.statusCode)

		log.Printf("%s%s%s %s -> %s%d - %s%s %stook %d ms%s\n",
			grey, r.Method, reset,
			r.URL.Path,
			statusColor, rec.statusCode,
			http.StatusText(rec.statusCode), reset,
			grey, duration, reset,
		)
	}
}

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	grey   = "\033[90m"
)

func getStatusColor(code int) string {
	switch {
	case code >= 500:
		return red
	case code >= 400:
		return yellow
	case code >= 300:
		return blue
	case code >= 200:
		return green
	default:
		return reset
	}
}
