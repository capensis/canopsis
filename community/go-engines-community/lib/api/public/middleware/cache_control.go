package middleware

import (
	"net/http"
)

const (
	CacheControlHeaderKey          = "Cache-Control"
	DefaultCacheControlHeaderValue = "public, no-cache"
)

// CacheControl middleware adds default cache headers to response.
func CacheControl() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(CacheControlHeaderKey, DefaultCacheControlHeaderValue)
			next.ServeHTTP(w, r)
		})
	}
}
