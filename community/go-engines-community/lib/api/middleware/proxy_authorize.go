package middleware

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/proxy"
	"github.com/gin-gonic/gin"
)

// ProxyAuthorize determines if current subject has been authorized to take
// an action on an object for proxy routes.
func ProxyAuthorize(
	enforcer security.Enforcer,
	accessConfig proxy.AccessConfig,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if obj, act := accessConfig.Get(c.Request.RequestURI, c.Request.Method); obj != "" && act != "" {
			subj := c.MustGet(auth.UserKey)

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
