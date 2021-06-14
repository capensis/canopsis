package middleware

import (
	"github.com/gin-gonic/gin"
)

const (
	CacheControlHeaderKey          = "Cache-Control"
	DefaultCacheControlHeaderValue = "public, no-cache"
)

// Cache middleware adds default cache headers to response.
func Cache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header(CacheControlHeaderKey, DefaultCacheControlHeaderValue)
		c.Next()
	}
}
