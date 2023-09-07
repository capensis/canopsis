package middleware

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

// Auth middleware uses http providers to authenticate user.
// It checks auth only if request contains credentials.
func Auth(providers []security.HttpProvider, maintenanceAdapter config.MaintenanceAdapter, enforcer security.Enforcer) gin.HandlerFunc {
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

				maintenanceConf, err := maintenanceAdapter.GetConfig(c)
				if err != nil {
					panic(err)
				}

				if maintenanceConf.Enabled {
					ok, err = enforcer.Enforce(user.ID, apisecurity.PermMaintenance, model.PermissionCan)
					if err != nil {
						panic(err)
					}

					if !ok {
						c.AbortWithStatusJSON(http.StatusServiceUnavailable, common.CanopsisUnderMaintenanceResponse)
						return
					}
				}

				// The user credentials was found, set user's id to key UserKey in this context,
				// the user's id can be read later using c.MustGet(auth.UserKey).
				c.Set(auth.Username, user.DisplayName)
				c.Set(auth.UserKey, user.ID)
				c.Set(auth.RolesKey, user.Roles)
				c.Set(auth.ApiKey, user.AuthApiKey)
				break
			}
		}

		c.Next()
	}
}
