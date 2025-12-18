package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fastjson"
)

// SetAuthor middleware sets authorized user id to author field to request body. Use it for create and update model endpoints.
func SetAuthor(errorResponder httperror.Responder) func(c *gin.Context) {
	return func(c *gin.Context) {
		var body map[string]interface{}

		encodedBody := json.NewDecoder(c.Request.Body)
		err := encodedBody.Decode(&body)
		if err != nil {
			var syntaxError *json.SyntaxError
			var unmarshalTypeError *json.UnmarshalTypeError
			if errors.Is(err, io.EOF) || errors.As(err, &syntaxError) || errors.As(err, &unmarshalTypeError) {
				c.Next()
				return
			}

			errorResponder.Respond(c, err)

			return
		}

		body["author"], err = authctx.GetUserKey(c)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		encodedStr, err := json.Marshal(body)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(encodedStr))

		c.Next()
	}
}

// PreProcessBulk middleware checks if bulk has valid size and sets authorized user id to author field to bulk request body. Use it for create and update model endpoints.
func PreProcessBulk(configProvider config.ApiConfigProvider, errorResponder httperror.Responder, addAuthor bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		var ar fastjson.Arena
		raw, err := c.GetRawData()
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		if len(raw) == 0 {
			c.Next()
			return
		}

		jsonValue, err := fastjson.ParseBytes(raw)
		if err != nil {
			errorResponder.Respond(c, validation.NewInvalidRequestBodyError(err))

			return
		}

		rawObjects, err := jsonValue.Array()
		if err != nil {
			errorResponder.Respond(c, validation.NewInvalidRequestBodyError(err))

			return
		}

		bulkMaxSize := configProvider.Get().BulkMaxSize
		if len(rawObjects) > bulkMaxSize {
			errorResponder.Respond(c, httperror.ErrRequestEntityTooLarge)

			return
		}

		if addAuthor {
			userID, err := authctx.GetUserKey(c)
			if err != nil {
				errorResponder.Respond(c, err)

				return
			}

			for _, object := range rawObjects {
				object.Set("author", ar.NewString(userID))
			}
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonValue.MarshalTo(nil)))

		c.Next()
	}
}
