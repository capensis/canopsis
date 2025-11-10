package middleware

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"github.com/gin-gonic/gin"
)

// OnlyAuth determines if user is authenticated.
// Use Authorize middleware to check user permissions.
func OnlyAuth(errorResponder httperror.Responder) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := authctx.GetUserKey(c)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		c.Next()
	}
}
