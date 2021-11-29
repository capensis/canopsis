package middleware

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Auth middleware uses http providers to authenticate user.
// It checks auth only if request contains credentials.
func Auth(providers []security.HttpProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, p := range providers {
			user, err, ok := p.Auth(c.Request)
			if err != nil {
				panic(err)
			}

			if ok {
				if user == nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
					return
				}

				// The user credentials was found, set user's id to key UserKey in this context,
				// the user's id can be read later using c.MustGet(auth.UserKey).
				c.Set(auth.Username, user.Name)
				c.Set(auth.UserKey, user.ID)
				c.Set(auth.RoleKey, user.Role)
				c.Set(auth.ApiKey, user.AuthApiKey)
				break
			}
		}

		c.Next()
	}
}
