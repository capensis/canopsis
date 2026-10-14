package middleware

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
)

func Auth(
	providers []security.HttpProvider,
	responder response.Responder,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			for _, p := range providers {
				user, err, ok := p.Auth(r)
				if err != nil {
					responder.RespondError(w, r, err)

					return
				}

				if ok {
					if user == nil {
						responder.RespondError(w, r, httperror.ErrUnauthorized)

						return
					}

					ctx = authctx.SetUsername(ctx, user.DisplayName)
					ctx = authctx.SetUserKey(ctx, user.ID)
					ctx = authctx.SetRoles(ctx, user.Roles)
					ctx = authctx.SetAPIKey(ctx, user.AuthApiKey)
					ctx = authctx.SetLocale(ctx, user.Language)
					break
				}
			}

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
