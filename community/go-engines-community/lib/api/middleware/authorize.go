package middleware

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
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
	errorResponder httperror.Responder,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		subj, err := authctx.GetUserKey(c)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		ok, err := enforcer.Enforce(subj, obj, act)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}

		c.Next()
	}
}

// AuthorizeAtLeastOnePerm allows access if at least one PermCheck pair is permitted for the user
func AuthorizeAtLeastOnePerm(
	permChecks []apisecurity.PermCheck,
	enforcer security.Enforcer,
	errorResponder httperror.Responder,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		subj, err := authctx.GetUserKey(c)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		for _, permCheck := range permChecks {
			ok, err := enforcer.Enforce(subj, permCheck.Obj, permCheck.Act)
			if err != nil {
				errorResponder.Respond(c, err)

				return
			}

			if ok {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
	}
}
