package request

//todo: copy from webhook package, webhook package should use this package instead of its own models

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
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
