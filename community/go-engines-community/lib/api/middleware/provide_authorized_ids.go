package middleware

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
)

const AuthorizedIds = "authorized_ids"
const OwnedIds = "owned_ids"

// ProvideAuthorizedIds determines on which objects current subject has been authorized to take
// an action.
func ProvideAuthorizedIds(
	act string,
	enforcer security.Enforcer,
	provider apisecurity.OwnedObjectsProvider,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawSubj, ok := c.Get(authctx.UserKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
			return
		}

		subj, ok := rawSubj.(string)
		if !ok {
			panic("user key is not a string")
		}

		roles, err := enforcer.GetRolesForUser(subj)
		if err != nil {
			panic(err)
		}

		ids := make([]string, 0)
		for _, role := range roles {
			perms, err := enforcer.GetPermissionsForUser(role)
			if err != nil {
				panic(err)
			}
			for _, perm := range perms {
				if len(perm) != 3 {
					continue
				}

				if perm[2] == act {
					ids = append(ids, perm[1])
				}
			}
		}

		if provider != nil {
			ownedIds, err := provider.GetOwnedIDs(c, subj)
			if err != nil {
				panic(err)
			}
			c.Set(OwnedIds, ownedIds)
		}

		c.Set(AuthorizedIds, ids)

		c.Next()
	}
}
