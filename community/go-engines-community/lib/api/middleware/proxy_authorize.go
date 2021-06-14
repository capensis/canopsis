package middleware

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/auth"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/security"
	"git.canopsis.net/canopsis/go-engines/lib/security/proxy"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ProxyAuthorize determines if current subject has been authorized to take
// an action on an object for proxy routes.
func ProxyAuthorize(
	enforcer security.Enforcer,
	accessConfig proxy.AccessConfig,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if obj, act := accessConfig.Get(c.Request.RequestURI, c.Request.Method); obj != "" && act != "" {
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
		}

		c.Next()
	}
}
