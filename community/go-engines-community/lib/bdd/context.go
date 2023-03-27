package bdd

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type contextKey int

const (
	contextKeyScenarioName contextKey = iota
	contextKeyScenarioUri
	contextKeyApiAuthToken
	contextKeyRequestBody
	contextKeyResponseStatusCode
	contextKeyResponseBody
	contextKeyResponseBodyOutput
	contextKeyHeaders
	contextKeyCookies
	contextKeyVars
	contextConsumer
	contextWebsocketConn
)

func GetScenarioName(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(contextKeyScenarioName).(string)
	return v, ok
}

func SetScenarioName(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, contextKeyScenarioName, v)
}

func GetScenarioUri(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(contextKeyScenarioUri).(string)
	return v, ok
}

func SetScenarioUri(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, contextKeyScenarioUri, v)
}

func getApiAuthToken(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(contextKeyApiAuthToken).(string)
	return v, ok
}

func setApiAuthToken(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, contextKeyApiAuthToken, v)
}

func getRequestBody(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(contextKeyRequestBody).(string)
	return v, ok
}

func setRequestBody(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, contextKeyRequestBody, v)
}

func getResponseStatusCode(ctx context.Context) (int, bool) {
	v, ok := ctx.Value(contextKeyResponseStatusCode).(int)
	return v, ok
}

func setResponseStatusCode(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, contextKeyResponseStatusCode, v)
}

func getResponseBody(ctx context.Context) (interface{}, bool) {
	v := ctx.Value(contextKeyResponseBody)
	if v == nil {
		return nil, false
	}
	return v, true
}

func setResponseBody(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, contextKeyResponseBody, v)
}

func getResponseBodyOutput(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(contextKeyResponseBodyOutput).(string)
	return v, ok
}

func setResponseBodyOutput(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, contextKeyResponseBodyOutput, v)
}

func getHeaders(ctx context.Context) (map[string]string, bool) {
	v, ok := ctx.Value(contextKeyHeaders).(map[string]string)
	return v, ok
}

func setHeaders(ctx context.Context, v map[string]string) context.Context {
	return context.WithValue(ctx, contextKeyHeaders, v)
}

func getCookies(ctx context.Context) ([]*http.Cookie, bool) {
	v, ok := ctx.Value(contextKeyCookies).([]*http.Cookie)
	return v, ok
}

func setCookies(ctx context.Context, v []*http.Cookie) context.Context {
	return context.WithValue(ctx, contextKeyCookies, v)
}

func getVars(ctx context.Context) (map[string]string, bool) {
	v, ok := ctx.Value(contextKeyVars).(map[string]string)
	return v, ok
}

func getVar(ctx context.Context, k string) (string, error) {
	vars, ok := getVars(ctx)
	if !ok {
		return "", fmt.Errorf("%q doesn't exist", k)
	}

	s, ok := vars[k]
	if !ok {
		return "", fmt.Errorf("%q doesn't exist", k)
	}

	return s, nil
}

func parseFloatVar(ctx context.Context, k string) (float64, error) {
	v, err := getVar(ctx, k)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(v, 64)
}

func setVar(ctx context.Context, k, v string) context.Context {
	vars, ok := getVars(ctx)
	if ok {
		vars[k] = v
	} else {
		vars = map[string]string{k: v}
	}

	ctx = context.WithValue(ctx, contextKeyVars, vars)
	return ctx
}

func getConsumer(ctx context.Context) (int, bool) {
	v, ok := ctx.Value(contextConsumer).(int)
	return v, ok
}

func setConsumer(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, contextConsumer, v)
}

func setWebsocketConn(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, contextWebsocketConn, v)
}

func getWebsocketConn(ctx context.Context) (int, bool) {
	v, ok := ctx.Value(contextWebsocketConn).(int)
	return v, ok
}
