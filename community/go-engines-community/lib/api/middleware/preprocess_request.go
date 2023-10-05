package middleware

import (
	"bytes"
	"encoding/json"
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
			if err == io.EOF {
				c.Next()
				return
			}
			panic(err)
		}

		userId := c.MustGet(auth.UserKey)
		body["author"] = userId

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
			userId := c.MustGet(auth.UserKey)
			for _, object := range rawObjects {
				object.Set("author", ar.NewString(userId.(string)))
			}
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonValue.MarshalTo(nil)))

		c.Next()
	}
}
