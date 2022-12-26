package request

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/valyala/fastjson"
)

func CreateRequest(ctx context.Context, params Parameters) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, params.Method, params.URL, bytes.NewBufferString(params.Payload))
	if err != nil {
		return nil, err
	}

	for key, value := range params.Headers {
		request.Header.Set(key, value)
	}

	if params.Auth != nil {
		request.SetBasicAuth(params.Auth.Username, params.Auth.Password)
	}

	return request, nil
}

// Flatten creates a new map[string]interface{} with a flat hierarchy
func Flatten(in interface{}, prevKey string) map[string]interface{} {
	out := make(map[string]interface{})

	switch inVal := in.(type) {
	case map[string]interface{}:
		for k, v := range inVal {
			newPrevKey := prevKey + "." + k
			if prevKey == "" {
				newPrevKey = k
			}

			nm := Flatten(v, newPrevKey)
			for nk, nv := range nm {
				out[nk] = nv
			}
		}
	case []interface{}:
		for idx, v := range inVal {
			newPrevKey := fmt.Sprintf("%s.%d", prevKey, idx)
			if prevKey == "" {
				newPrevKey = newPrevKey[1:]
			}

			nm := Flatten(v, newPrevKey)
			for nk, nv := range nm {
				out[nk] = nv
			}
		}
	default:
		out[prevKey] = inVal
	}

	return out
}

// ValidateStatusCode checks response status code and generates error if status code is not allowed.
func ValidateStatusCode(
	request *http.Request,
	response *http.Response,
	resErrMsgKey string,
) error {
	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return nil
	}

	reqUrl := request.URL.String()
	errMsg := ""
	if resErrMsgKey != "" {
		body, err := io.ReadAll(response.Body)
		if err == nil {
			parsed, err := fastjson.ParseBytes(body)
			if err == nil {
				errFieldVal := parsed.GetStringBytes(strings.Split(resErrMsgKey, ".")...)
				if len(errFieldVal) > 0 {
					errMsg = string(errFieldVal)
				}
			}
		}
	}

	if errMsg == "" {
		switch response.StatusCode {
		case http.StatusNotFound:
			errMsg = fmt.Sprintf("url %s not found", reqUrl)
		case http.StatusMethodNotAllowed:
			errMsg = fmt.Sprintf("method %s not allowed for url %s", request.Method, reqUrl)
		case http.StatusUnauthorized:
			errMsg = fmt.Sprintf("url %s is unauthorized", reqUrl)
		case http.StatusForbidden:
			errMsg = fmt.Sprintf("url %s is forbidden", reqUrl)
		default:
			errMsg = fmt.Sprintf("request url %s failed with status code %d", reqUrl, response.StatusCode)
		}
	}

	return errors.New(errMsg)
}
