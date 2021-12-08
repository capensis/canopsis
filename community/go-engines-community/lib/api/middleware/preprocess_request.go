package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"io"
	"io/ioutil"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"github.com/gin-gonic/gin"
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

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(encodedStr))

		c.Next()
	}
}

// PreProcessBulk middleware checks if bulk has valid size and sets authorized user id to author field to bulk request body. Use it for create and update model endpoints.
func PreProcessBulk(cfg config.CanopsisConf) func(c *gin.Context) {
	return func(c *gin.Context) {
		var body []map[string]interface{}

		encodedBody := json.NewDecoder(c.Request.Body)
		err := encodedBody.Decode(&body)
		if err != nil {
			if err == io.EOF {
				c.Next()
				return
			}
			panic(err)
		}

		if len(body) > cfg.API.BulkMaxSize {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(fmt.Errorf("number of elements shouldn't be greater than %d", cfg.API.BulkMaxSize)))
			return
		}

		userId := c.MustGet(auth.UserKey)
		for _, item := range body {
			item["author"] = userId
		}

		encodedStr, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(encodedStr))

		c.Next()
	}
}
