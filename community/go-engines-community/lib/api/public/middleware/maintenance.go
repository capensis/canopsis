package middleware

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
)

func Maintenance(
	maintenanceAdapter config.MaintenanceAdapter,
	enforcer security.Enforcer,
	responder response.Responder,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			maintenanceConf, err := maintenanceAdapter.GetConfig(ctx)
			if err != nil {
				responder.RespondError(w, r, err)

				return
			}

			if maintenanceConf.Enabled {
				userID, err := authctx.GetUserKey(ctx)
				if err != nil {
					responder.RespondError(w, r, err)

					return
				}

				ok, err := enforcer.Enforce(userID, apisecurity.PermMaintenance, model.PermissionCan)
				if err != nil {
					responder.RespondError(w, r, err)

					return
				}

				if !ok {
					responder.RespondError(w, r, httperror.ErrMaintenance)

					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
