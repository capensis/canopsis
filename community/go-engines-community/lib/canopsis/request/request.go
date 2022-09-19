package request

//todo: copy from webhook package, webhook package should use this package instead of its own models

import (
	"bytes"
	"context"
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

func Flatten(m map[string]interface{}) map[string]interface{} {
	o := make(map[string]interface{})
	for k, v := range m {
		switch child := v.(type) {
		case map[string]interface{}:
			nm := Flatten(child)
			for nk, nv := range nm {
				o[k+"."+nk] = nv
			}
		default:
			o[k] = v
		}
	}

	return o
}
