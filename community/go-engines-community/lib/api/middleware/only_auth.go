package middleware

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"github.com/gin-gonic/gin"
)

// OnlyAuth determines if user is authenticated.
// Use Authorize middleware to check user permissions.
func OnlyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.MustGet(authctx.UserKey)

		c.Next()
	}
}
