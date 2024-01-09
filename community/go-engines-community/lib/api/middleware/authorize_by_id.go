package middleware

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
)

// AuthorizeByID determines if current subject has been authorized to take
// an action on a object by id.
func AuthorizeByID(
	act string,
	enforcer security.Enforcer,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := c.Param("id")
		if obj == "" {
			panic(errors.New("missing id parameter"))
		}

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
