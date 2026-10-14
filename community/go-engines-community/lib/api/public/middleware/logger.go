package middleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	libhttp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/http"
	"github.com/rs/zerolog"
)

func Logger(
	responder response.Responder,
	logger zerolog.Logger,
	logBody bool,
	logBodyOnError bool,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.URL.Path
			raw := r.URL.RawQuery
			if raw != "" {
				path = path + "?" + raw
			}

			wrapperW := &loggingResponseWriter{
				w: w,
			}
			var requestBody io.ReadCloser
			if logBody || logBodyOnError {
				buf, ok := bufPool.Get().(*bytes.Buffer)
				if !ok {
					responder.RespondError(w, r, errors.New("unknown buffer type"))

					return
				}

				defer bufPool.Put(buf)
				buf.Reset()
				wrapperW.body = buf
				requestBody, r.Body, _ = libhttp.DrainBody(r.Body)
			}

			next.ServeHTTP(wrapperW, r)
			duration := time.Since(start)
			statusCode := wrapperW.statusCode
			var logEvent *zerolog.Event
			isResponseOk := false
			switch {
			case statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices:
				logEvent = logger.Info() // nolint:zerologlint
				isResponseOk = true
			case statusCode >= http.StatusMultipleChoices && statusCode < http.StatusInternalServerError:
				logEvent = logger.Warn() // nolint:zerologlint
			default:
				logEvent = logger.Error() // nolint:zerologlint
			}

			if logBody || logBodyOnError && !isResponseOk {
				if requestBody != nil {
					b, err := io.ReadAll(requestBody)
					if err == nil {
						logEvent = logEvent.Str("request_body", string(b))
					}
				}

				if wrapperW.body != nil {
					logEvent = logEvent.Str("response_body", wrapperW.body.String())
				}
			}

			logEvent.
				Str("duration", duration.String()).
				Str("remote_ip", r.RemoteAddr).
				Msg(strconv.Itoa(statusCode) + " " + r.Method + " " + path)
		})
	}
}
