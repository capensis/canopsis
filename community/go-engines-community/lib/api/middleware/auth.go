package middleware

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

// Auth middleware uses http providers to authenticate user.
// It checks auth only if request contains credentials.
func Auth(
	providers []security.HttpProvider,
	maintenanceAdapter config.MaintenanceAdapter,
	enforcer security.Enforcer,
	errorResponder httperror.Responder,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, p := range providers {
			user, err, ok := p.Auth(c.Request)
			if err != nil {
				errorResponder.Respond(c, err)

				return
			}

			if ok {
				if user == nil {
					errorResponder.Respond(c, httperror.ErrUnauthorized)

					return
				}

				maintenanceConf, err := maintenanceAdapter.GetConfig(c)
				if err != nil {
					errorResponder.Respond(c, err)

					return
				}

				if maintenanceConf.Enabled {
					ok, err = enforcer.Enforce(user.ID, apisecurity.PermMaintenance, model.PermissionCan)
					if err != nil {
						errorResponder.Respond(c, err)

						return
					}

					if !ok {
						errorResponder.Respond(c, httperror.ErrMaintenance)

						return
					}
				}

				// The user credentials was found, set user's id to key UserKey in this context,
				// the user's id can be read later using c.MustGet(auth.UserKey).
				authctx.SetUsername(c, user.DisplayName)
				authctx.SetUserKey(c, user.ID)
				authctx.SetRoles(c, user.Roles)
				authctx.SetAPIKey(c, user.AuthApiKey)
				authctx.SetLocale(c, user.Language)
				break
			}
		}

		c.Next()
	}
}
