package middleware

import (
	"encoding/json"
	"io"
	"net/url"
	"sort"
	"strings"
	"time"

	libhttp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/http"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func Cache(
	redisClient *redis.Client,
	defaultExpire time.Duration,
	getExpire func() time.Duration,
) gin.HandlerFunc {
	redisStore := persist.NewRedisStore(redisClient)
	return cache.Cache(redisStore, defaultExpire, cache.WithPrefixKey(libredis.ApiCacheRequestKey), cache.WithCacheStrategyByRequest(func(c *gin.Context) (bool, cache.Strategy) {
		cacheKey := getRequestUriIgnoreQueryOrder(c.Request.URL)
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

			cacheKey += getRequestBodyIgnoreJsonOrder(b)
		}

		var cacheDuration time.Duration
		if getExpire != nil {
			cacheDuration = getExpire()
		}

		return true, cache.Strategy{
			CacheKey:      cacheKey,
			CacheDuration: cacheDuration,
		}
	}))
}

func getRequestUriIgnoreQueryOrder(url *url.URL) string {
	values := url.Query()
	if len(values) == 0 {
		return url.Path
	}

	keys := make([]string, len(values))
	i := 0
	for k := range values {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	sortedValues := make([]string, 0, len(keys))
	for _, k := range keys {
		sort.Strings(values[k])
		for _, v := range values[k] {
			sortedValues = append(sortedValues, k+"="+v)
		}
	}

	return url.Path + "?" + strings.Join(sortedValues, "&")
}

func getRequestBodyIgnoreJsonOrder(body []byte) string {
	values := make(map[string]any)
	err := json.Unmarshal(body, &values)
	if err != nil || len(values) == 0 {
		return string(body)
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

	b, err := json.Marshal(sortedValues)
	if err != nil {
		return string(body)
	}

	return string(b)
}
