package middleware

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
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
	errorResponder httperror.Responder,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		subj, err := authctx.GetUserKey(c)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		roles, err := enforcer.GetRolesForUser(subj)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		ids := make([]string, 0)
		for _, role := range roles {
			perms, err := enforcer.GetPermissionsForUser(role)
			if err != nil {
				errorResponder.Respond(c, err)

				return
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
				errorResponder.Respond(c, err)

				return
			}

			c.Set(OwnedIds, ownedIds)
		}

		c.Set(AuthorizedIds, ids)

		c.Next()
	}
}
