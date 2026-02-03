package middleware

import (
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
)

// AuthorizeByID determines if current subject has been authorized to take
// an action on an object by id.
func AuthorizeByID(
	act string,
	enforcer security.Enforcer,
	errorResponder httperror.Responder,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := c.Param("id")
		if obj == "" {
			errorResponder.Respond(c, errors.New("missing id parameter"))

			return
		}

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
			errorResponder.Respond(c, httperror.NewForbiddenError(""))

			return
		}

		c.Next()
	}
}
