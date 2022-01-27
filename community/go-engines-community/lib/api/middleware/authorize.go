package middleware

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
)

// Authorize determines if current subject has been authorized to take
// an action on an object. Use OnlyAuth middleware to only check if user is authenticated.
//
// Note: if new user is created, then enforcer.LoadPolicy() should be called to reload security policies,
// it throws http.StatusForbidden otherwise!
func Authorize(
	obj string,
	act string,
	enforcer security.Enforcer,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		subj := c.MustGet(auth.UserKey)

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
