package middleware

import (
	"bytes"
	"net/http"
	"sync"
)

var bufPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

type loggingResponseWriter struct {
	w          http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.w.WriteHeader(statusCode)
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}

	return w.w.Write(b)
}

func (w *loggingResponseWriter) Header() http.Header {
	return w.w.Header()
}
