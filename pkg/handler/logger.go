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

		log.Printf("%s -> %d - %s: Response took %d ms\n",
			r.URL.Path,
			rec.statusCode,
			http.StatusText(rec.statusCode),
			time.Since(start).Milliseconds(),
		)
	}
}
