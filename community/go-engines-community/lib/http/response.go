package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/valyala/fastjson"
)

var ErrResponseTooLong = errors.New("response too long")

const buffChunk = 512

type ResponseError struct {
	failReason string
	err        error
}

func NewResponseError(err error, failReason string) *ResponseError {
	return &ResponseError{failReason: failReason, err: err}
}

func (e *ResponseError) Error() string {
	return e.err.Error()
}

func (e *ResponseError) FailReason() string {
	return e.failReason
}

func (e *ResponseError) Unwrap() error {
	return e.err
}

func ReadResponse(response *http.Response, maxSize int64) ([]byte, error) {
	if maxSize <= 0 {
		return nil, &ResponseError{
			failReason: ErrResponseTooLong.Error(),
			err:        ErrResponseTooLong,
		}
	}

	b := make([]byte, 0, buffChunk)
	for {
		if len(b) == cap(b) {
			b = append(b, 0)[:len(b)]
		}
		n, err := response.Body.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		if int64(len(b)) > maxSize {
			return nil, &ResponseError{
				failReason: ErrResponseTooLong.Error(),
				err:        ErrResponseTooLong,
			}
		}
	}

	return b, nil
}

func IsValidStatusCode(code int) bool {
	switch code {
	case http.StatusOK,
		http.StatusCreated,
		http.StatusAccepted,
		http.StatusNoContent,
		http.StatusResetContent,
		http.StatusPartialContent:
		return true
	}

	return false
}

// ValidateStatusCode checks response status code and generates error if status code is not allowed.
func ValidateStatusCode(
	request *http.Request,
	response *http.Response,
	resErrMsgKey string,
	maxSize int64,
) error {
	if IsValidStatusCode(response.StatusCode) {
		return nil
	}

	reqUrl := request.URL.String()
	errMsg := ""
	if resErrMsgKey != "" {
		body, err := ReadResponse(response, maxSize)
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

// FlattenResponse reads response body to a new map[string]any with a flat hierarchy.
func FlattenResponse(request *http.Request, response *http.Response, maxSize int64) (map[string]any, error) {
	body, err := ReadResponse(response, maxSize)
	if err != nil {
		return nil, err
	}

	v, err := fastjson.ParseBytes(body)
	if err != nil {
		failReason := ""
		if request != nil {
			failReason = "response of " + request.Method + " " + request.URL.String() + " is not valid JSON"
		}

		return nil, &ResponseError{
			failReason: failReason,
			err:        err,
		}
	}

	return flatten(v, ""), nil
}

func flatten(in *fastjson.Value, prevKey string) map[string]any {
	out := make(map[string]any)

	switch in.Type() {
	case fastjson.TypeObject:
		in.GetObject().Visit(func(key []byte, v *fastjson.Value) {
			newPrevKey := prevKey + "." + string(key)
			if prevKey == "" {
				newPrevKey = string(key)
			}

			nm := flatten(v, newPrevKey)
			for nk, nv := range nm {
				out[nk] = nv
			}
		})
	case fastjson.TypeArray:
		for idx, v := range in.GetArray() {
			newPrevKey := fmt.Sprintf("%s.%d", prevKey, idx)
			if prevKey == "" {
				newPrevKey = newPrevKey[1:]
			}

			nm := flatten(v, newPrevKey)
			for nk, nv := range nm {
				out[nk] = nv
			}
		}
	case fastjson.TypeNull:
		out[prevKey] = nil
	case fastjson.TypeString:
		out[prevKey] = string(in.GetStringBytes())
	case fastjson.TypeNumber:
		var err error
		out[prevKey], err = in.Int()
		if err != nil {
			out[prevKey] = in.GetFloat64()
		}
	case fastjson.TypeTrue, fastjson.TypeFalse:
		out[prevKey] = in.GetBool()
	default:
		out[prevKey] = string(in.GetStringBytes())
	}

	return out
}
