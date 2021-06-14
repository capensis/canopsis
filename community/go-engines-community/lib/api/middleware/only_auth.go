package middleware

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/auth"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

// OnlyAuth determines if user is authenticated.
// Use Authorize middleware to check user permissions.
func OnlyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get(auth.UserKey)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
			return
		}

		c.Next()
	}
}
