package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"github.com/gin-gonic/gin"
)

// AuthorizeOwnership determines if current subject is the owner of an object.
func AuthorizeOwnership(strategy security.OwnershipStrategy) gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := c.Param("id")
		if obj == "" {
			panic(errors.New("missing id parameter"))
		}

		if strategy == nil {
			panic(errors.New("missing ownership strategy"))
		}

		rawSubj, ok := c.Get(authctx.UserKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
			return
		}

		subj, ok := rawSubj.(string)
		if !ok {
			panic("user key is not a string")
		}

		ownership, err := strategy.IsOwner(c, obj, subj)
		if err != nil {
			panic(err)
		}

		switch ownership {
		case security.OwnershipPublic, security.OwnershipOwner:
			break
		case security.OwnershipNotOwner:
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		case security.OwnershipNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
			return
		default:
			panic(fmt.Errorf("unexpected ownership: %d", ownership))
		}

		c.Next()
	}
}
