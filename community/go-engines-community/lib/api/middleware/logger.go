package middleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	libhttp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/http"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var bufPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

type responseBodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	return w.ResponseWriter.Write(b)
}

func Logger(logger zerolog.Logger, logBody bool, logBodyOnError bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		var responseWriter *responseBodyLogWriter
		var requestBody io.ReadCloser
		if logBody || logBodyOnError {
			buf, ok := bufPool.Get().(*bytes.Buffer)
			if !ok {
				panic(errors.New("unknown buffer type"))
			}

			defer bufPool.Put(buf)
			responseWriter = &responseBodyLogWriter{
				body:           buf,
				ResponseWriter: c.Writer,
			}
			c.Writer = responseWriter

			requestBody, c.Request.Body, _ = libhttp.DrainBody(c.Request.Body)
		}

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
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

			if responseWriter != nil {
				logEvent = logEvent.Str("response_body", responseWriter.body.String())
			}
		}

		logEvent.
			Str("duration", duration.String()).
			Str("client_ip", c.ClientIP()).
			Msg(strconv.Itoa(statusCode) + " " + c.Request.Method + " " + path)
	}
}
