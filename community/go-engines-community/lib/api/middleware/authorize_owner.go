package middleware

import (
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"github.com/gin-gonic/gin"
)

// AuthorizeOwnership determines if current subject is the owner of an object.
func AuthorizeOwnership(strategy security.OwnershipStrategy, errorResponder httperror.Responder) gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := c.Param("id")
		if obj == "" {
			errorResponder.Respond(c, errors.New("missing id parameter"))

			return
		}

		if strategy == nil {
			errorResponder.Respond(c, errors.New("missing ownership strategy"))

			return
		}

		subj, err := authctx.GetUserKey(c)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		ownership, err := strategy.IsOwner(c, obj, subj)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		switch ownership {
		case security.OwnershipPublic, security.OwnershipOwner:
			break
		case security.OwnershipNotOwner:
			errorResponder.Respond(c, httperror.NewForbiddenError(""))

			return
		case security.OwnershipNotFound:
			errorResponder.Respond(c, httperror.ErrNotFound)

			return
		default:
			errorResponder.Respond(c, fmt.Errorf("unexpected ownership: %d", ownership))

			return
		}

		c.Next()
	}
}
