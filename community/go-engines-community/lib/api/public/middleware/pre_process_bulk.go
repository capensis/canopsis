package middleware

import (
	"bytes"
	"io"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/valyala/fastjson"
)

// PreProcessBulk middleware checks if bulk has valid size and sets authorized user id to author field to bulk request body. Use it for create and update model endpoints.
func PreProcessBulk(
	configProvider config.ApiConfigProvider,
	responder response.Responder,
	addAuthor bool,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body == nil {
				next.ServeHTTP(w, r)

				return
			}

			raw, err := io.ReadAll(r.Body)
			if err != nil {
				responder.RespondError(w, r, err)

				return
			}

			if len(raw) == 0 {
				next.ServeHTTP(w, r)

				return
			}

			jsonValue, err := fastjson.ParseBytes(raw)
			if err != nil {
				responder.RespondError(w, r, validation.NewInvalidRequestBodyError(err))

				return
			}

			rawObjects, err := jsonValue.Array()
			if err != nil {
				responder.RespondError(w, r, validation.NewInvalidRequestBodyError(err))

				return
			}

			bulkMaxSize := configProvider.Get().BulkMaxSize
			if len(rawObjects) > bulkMaxSize {
				responder.RespondError(w, r, httperror.ErrRequestEntityTooLarge)

				return
			}

			if addAuthor {
				var ar fastjson.Arena
				userID, err := authctx.GetUserKey(r.Context())
				if err != nil {
					responder.RespondError(w, r, err)

					return
				}

				for _, object := range rawObjects {
					object.Set("author", ar.NewString(userID))
				}
			}

			r.Body = io.NopCloser(bytes.NewBuffer(jsonValue.MarshalTo(nil)))
			next.ServeHTTP(w, r)
		})
	}
}
