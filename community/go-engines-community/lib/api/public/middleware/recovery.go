package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

// Recovery recovers from any panics and if there was one it logs error and writes InternalServerError response.
func Recovery(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					var err, errToLog error
					var ok bool
					if err, ok = rec.(error); ok {
						errToLog = err
					} else {
						errToLog = fmt.Errorf("%v", rec)
					}

					// Check for a broken connection, as it is not really a
					// condition that warrants a panic stack trace.
					var se *os.SyscallError
					if err != nil && errors.As(err, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							logger.Err(err).Msgf("panic recovered")

							return
						}
					}

					var code int
					if err != nil && errors.Is(err, context.Canceled) {
						logger.Warn().Err(err).Msg("panic recovered")
						code = http.StatusRequestTimeout
					} else if err != nil && (errors.Is(err, context.DeadlineExceeded) || mongodriver.IsTimeout(err)) {
						logger.Err(err).Msg("panic recovered")
						code = http.StatusRequestTimeout
					} else {
						logger.Err(errToLog).Msgf("panic recovered\n%s\n", debug.Stack())
						code = http.StatusInternalServerError
					}

					w.Header().Set("Content-Type", response.MIMEJSON)
					w.WriteHeader(code)
					a := fastjson.Arena{}
					res := a.NewObject()
					res.Set("error", a.NewString(http.StatusText(code)))
					_, err = w.Write(res.MarshalTo(nil))
					if err != nil {
						logger.Err(err).Msg("failed to write response")
					}
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
