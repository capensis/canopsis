package request

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
