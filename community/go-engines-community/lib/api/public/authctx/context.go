package authctx

import (
	"context"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("key not found in context")

func SetUsername(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, Username, v)
}

func GetUsername(ctx context.Context) (string, error) {
	return get[string](ctx, Username)
}

func SetUserKey(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, UserKey, v)
}

func GetUserKey(ctx context.Context) (string, error) {
	return get[string](ctx, UserKey)
}

func SetRoles(ctx context.Context, v []string) context.Context {
	return context.WithValue(ctx, Roles, v)
}

func GetRoles(ctx context.Context) ([]string, error) {
	return get[[]string](ctx, Roles)
}

func SetAPIKey(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, APIKey, v)
}

func GetAPIKey(ctx context.Context) (string, error) {
	return get[string](ctx, APIKey)
}

func SetLocale(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, Locale, v)
}

func GetLocale(ctx context.Context) (string, error) {
	return get[string](ctx, Locale)
}

func get[T any](ctx context.Context, k contextKey) (T, error) {
	var res T
	v := ctx.Value(k)
	if v == nil {
		return res, fmt.Errorf("key=%v: %w", k, ErrNotFound)
	}

	res, ok := v.(T)
	if !ok {
		return res, fmt.Errorf("unexpected type %[1]T %[1]v for %d", v, k)
	}

	return res, nil
}
