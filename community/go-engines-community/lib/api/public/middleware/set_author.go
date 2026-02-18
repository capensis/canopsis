package middleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"github.com/valyala/fastjson"
)

// SetAuthor middleware sets authorized user id to author field to request body.
// Use it for create and update model endpoints.
func SetAuthor(responder response.Responder) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body == nil {
				next.ServeHTTP(w, r)

				return
			}

			b, err := io.ReadAll(r.Body)
			if err != nil {
				if errors.Is(err, io.EOF) {
					next.ServeHTTP(w, r)

					return
				}

				responder.RespondError(w, r, err)

				return
			}

			p, err := fastjson.ParseBytes(b)
			if err != nil {
				next.ServeHTTP(w, r)

				return
			}

			a, err := authctx.GetUserKey(r.Context())
			if err != nil {
				responder.RespondError(w, r, err)

				return
			}

			ar := fastjson.Arena{}
			p.Set("author", ar.NewString(a))
			r.Body = io.NopCloser(bytes.NewBuffer(p.MarshalTo(nil)))
			next.ServeHTTP(w, r)
		})
	}
}
