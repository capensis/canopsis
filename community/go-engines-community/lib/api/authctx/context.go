package authctx

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

var ErrNotFound = errors.New("key not found in context")

func SetUsername(c *gin.Context, v string) {
	c.Set(Username, v)
}

func GetUsername(c *gin.Context) (string, error) {
	return get[string](c, Username)
}

func SetUserKey(c *gin.Context, v string) {
	c.Set(UserKey, v)
}

func GetUserKey(c *gin.Context) (string, error) {
	return get[string](c, UserKey)
}

func SetRoles(c *gin.Context, v []string) {
	c.Set(Roles, v)
}

func GetRoles(c *gin.Context) ([]string, error) {
	return get[[]string](c, Roles)
}

func SetAPIKey(c *gin.Context, v string) {
	c.Set(APIKey, v)
}

func GetAPIKey(c *gin.Context) (string, error) {
	return get[string](c, APIKey)
}

func SetLocale(c *gin.Context, v string) {
	c.Set(Locale, v)
}

func GetLocale(c *gin.Context) (string, error) {
	return get[string](c, Locale)
}

func get[T any](c *gin.Context, k string) (T, error) {
	var res T
	v, ok := c.Get(k)
	if !ok {
		return res, fmt.Errorf("key=%s: %w", k, ErrNotFound)
	}

	res, ok = v.(T)
	if !ok {
		return res, fmt.Errorf("unexpected type %[1]T %[1]v for %s", v, k)
	}

	return res, nil
}
