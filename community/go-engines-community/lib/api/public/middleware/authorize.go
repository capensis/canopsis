package middleware

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
)

// Authorize determines if current subject has been authorized to take
// an action on an object. Use OnlyAuth middleware to only check if user is authenticated.
//
// Note: if new user is created, then enforcer.LoadPolicy() should be called to reload security policies,
// it throws http.StatusForbidden otherwise!
func Authorize(
	obj, act string,
	enforcer security.Enforcer,
	responder response.Responder,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			subj, err := authctx.GetUserKey(r.Context())
			if err != nil {
				responder.RespondError(w, r, err)

				return
			}

			ok, err := enforcer.Enforce(subj, obj, act)
			if err != nil {
				responder.RespondError(w, r, err)

				return
			}

			if !ok {
				responder.RespondError(w, r, httperror.NewForbiddenError(""))

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
