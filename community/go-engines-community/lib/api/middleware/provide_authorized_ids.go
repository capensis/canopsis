package middleware

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
	"net/http"
)

const AuthorizedIds = "authorized_ids"

// ProvideAuthorizedIds determines on which objects current subject has been authorized to take
// an action.
func ProvideAuthorizedIds(
	act string,
	enforcer security.Enforcer,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		subj, ok := c.Get(auth.UserKey)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
			return
		}

		roles, err := enforcer.GetRolesForUser(subj.(string))
		if err != nil {
			panic(err)
		}

		ids := make([]string, 0)
		for _, role := range roles {
			perms := enforcer.GetPermissionsForUser(role)
			for _, perm := range perms {
				if len(perm) != 3 {
					continue
				}

				if perm[2] == act {
					ids = append(ids, perm[1])
				}
			}
		}

		c.Set(AuthorizedIds, ids)

		c.Next()
	}
}
