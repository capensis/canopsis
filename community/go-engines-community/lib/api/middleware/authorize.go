package middleware

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/auth"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/security"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Authorize determines if current subject has been authorized to take
// an action on an object. Use OnlyAuth middleware to only check if user is authenticated.
func Authorize(
	obj string,
	act string,
	enforcer security.Enforcer,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		subj, ok := c.Get(auth.UserKey)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
			return
		}

		ok, err := enforcer.Enforce(subj.(string), obj, act)

		if err != nil {
			panic(err)
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}

		c.Next()
	}
}
