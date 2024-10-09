package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fastjson"
)

// SetAuthor middleware sets authorized user id to author field to request body. Use it for create and update model endpoints.
func SetAuthor() func(c *gin.Context) {
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
			panic(err)
		}

		userID := c.MustGet(auth.UserKey)
		body["author"] = userID

		encodedStr, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(encodedStr))

		c.Next()
	}
}

// PreProcessBulk middleware checks if bulk has valid size and sets authorized user id to author field to bulk request body. Use it for create and update model endpoints.
func PreProcessBulk(configProvider config.ApiConfigProvider, addAuthor bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		var ar fastjson.Arena

		raw, err := c.GetRawData()
		if err != nil {
			panic(err)
		}

		if len(raw) == 0 {
			c.Next()
			return
		}

		jsonValue, err := fastjson.ParseBytes(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		rawObjects, err := jsonValue.Array()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		bulkMaxSize := configProvider.Get().BulkMaxSize
		if len(rawObjects) > bulkMaxSize {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(fmt.Errorf("number of elements shouldn't be greater than %d", bulkMaxSize)))
			return
		}

		if addAuthor {
			userID, ok := c.MustGet(auth.UserKey).(string)
			if !ok {
				panic(fmt.Errorf("unknown type of %s", auth.UserKey))
			}

			for _, object := range rawObjects {
				object.Set("author", ar.NewString(userID))
			}
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonValue.MarshalTo(nil)))

		c.Next()
	}
}
