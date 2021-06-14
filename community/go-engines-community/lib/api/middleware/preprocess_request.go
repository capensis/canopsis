package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"git.canopsis.net/canopsis/go-engines/lib/api/auth"
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

// SetAuthorToBulk middleware sets authorized user id to author field to bulk request body. Use it for create and update model endpoints.
func SetAuthorToBulk() func(c *gin.Context) {
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
