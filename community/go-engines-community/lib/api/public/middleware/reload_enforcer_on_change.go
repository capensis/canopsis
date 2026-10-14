package middleware

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/rs/zerolog"
)

// ReloadEnforcerPolicyOnChange loads security policy if request changes policy.
func ReloadEnforcerPolicyOnChange(
	enforcer security.Enforcer,
	logger zerolog.Logger,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapperW := &loggingResponseWriter{
				w: w,
			}
			next.ServeHTTP(wrapperW, r)
			switch wrapperW.statusCode {
			case http.StatusOK, http.StatusNoContent, http.StatusCreated, http.StatusMultiStatus:
				err := enforcer.LoadPolicy()
				if err != nil {
					logger.Err(err).Msg("cannot reload enforcer")

					return
				}
			}
		})
	}
}
