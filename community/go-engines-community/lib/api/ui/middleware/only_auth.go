package middleware

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/response"
	"github.com/gin-gonic/gin"
)

// OnlyAuth determines if user is authenticated.
// Use Authorize middleware to check user permissions.
func OnlyAuth(errorResponder response.ErrorResponder) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := authctx.GetUserKey(c)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		c.Next()
	}
}
