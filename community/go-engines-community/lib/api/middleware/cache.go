package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/url"
	"sort"
	"time"

	libhttp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/http"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
)

func Cache(
	defaultExpire time.Duration,
	getExpire func() time.Duration,
) gin.HandlerFunc {
	redisStore := persist.NewMemoryStore(defaultExpire)
	return cache.Cache(redisStore, defaultExpire, cache.WithPrefixKey(libredis.ApiCacheRequestKey), cache.WithCacheStrategyByRequest(func(c *gin.Context) (bool, cache.Strategy) {
		buff := bytes.Buffer{}
		getRequestUriIgnoreQueryOrder(&buff, c.Request.URL)
		if c.Request.URL.Fragment != "" {
			buff.WriteRune('#')
			buff.WriteString(c.Request.URL.EscapedFragment())
		}

		var body io.ReadCloser
		var err error
		if c.Request.Body != nil {
			body, c.Request.Body, err = libhttp.DrainBody(c.Request.Body)
			if err != nil {
				panic(err)
			}

			b, err := io.ReadAll(body)
			if err != nil {
				panic(err)
			}

			buff.WriteRune('|')
			getRequestBodyIgnoreJsonOrder(&buff, b)
		}

		var cacheDuration time.Duration
		if getExpire != nil {
			cacheDuration = getExpire()
		}

		cacheKey := sha256.Sum256(buff.Bytes())
		return true, cache.Strategy{
			CacheKey:      hex.EncodeToString(cacheKey[:]),
			CacheDuration: cacheDuration,
		}
	}))
}

func getRequestUriIgnoreQueryOrder(buff *bytes.Buffer, u *url.URL) {
	buff.WriteString(u.Path)
	values := u.Query()
	if len(values) == 0 {
		return
	}

	keys := make([]string, len(values))
	i := 0
	for k := range values {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	added := false
	buff.WriteRune('?')
	for _, k := range keys {
		sort.Strings(values[k])
		escapedKey := url.QueryEscape(k)
		for _, v := range values[k] {
			if added {
				buff.WriteRune('&')
			}
			buff.WriteString(escapedKey)
			buff.WriteRune('=')
			buff.WriteString(url.QueryEscape(v))
			added = true
		}
	}
}

func getRequestBodyIgnoreJsonOrder(buff *bytes.Buffer, body []byte) {
	values := make(map[string]any)
	err := json.Unmarshal(body, &values)
	if err != nil || len(values) == 0 {
		buff.Write(body)
		return
	}

	keys := make([]string, len(values))
	i := 0
	for k := range values {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	sortedValues := make([]map[string]any, len(keys))
	for i, k := range keys {
		sortedValues[i] = map[string]any{
			k: values[k],
		}
	}

	err = json.NewEncoder(buff).Encode(sortedValues)
	if err != nil {
		buff.Write(body)
		return
	}
}
