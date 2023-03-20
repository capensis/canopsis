// Package bdd contains feature context utils.
package bdd

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go/types"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	libhttp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/http"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin/binding"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	ApiEnvURL           = "API_URL"
	requestTimeout      = 10 * time.Second
	userPass            = "test"
	headerAuthorization = "Authorization"
	headerContentType   = "Content-Type"
	basicPrefix         = "Basic"
	bearerPrefix        = "Bearer"

	startRepeatRequestInterval = time.Millisecond * 10
	totalRepeatRequestInterval = time.Second * 10
)

// ApiClient represents utility struct which implements API steps to feature context.
type ApiClient struct {
	// url is base API url.
	url string
	// client is http client to make API requests.
	client *http.Client
	// db is db client.
	db            mongo.DbClient
	requestLogger zerolog.Logger
	templater     *Templater
}

// NewApiClient creates new API client.
func NewApiClient(db mongo.DbClient, url string, requestLogger zerolog.Logger, templater *Templater) *ApiClient {
	return &ApiClient{
		url: url,
		client: &http.Client{
			Timeout: requestTimeout,
		},
		db:            db,
		requestLogger: requestLogger,
		templater:     templater,
	}
}

// GetApiURL retrieves API url from env var.
func GetApiURL() (string, error) {
	legacy := os.Getenv(ApiEnvURL)
	if legacy == "" {
		return "", fmt.Errorf("environment variable %s empty", ApiEnvURL)
	}

	parsed, err := url.Parse(legacy)
	if err != nil {
		return "", fmt.Errorf("cannot parse api url: %w", err)
	}

	return parsed.String(), nil
}

/*
*
Step example:

	Then the response code should be 200
*/
func (a *ApiClient) TheResponseCodeShouldBe(ctx context.Context, code int) error {
	responseStatusCode, ok := getResponseStatusCode(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}
	responseBodyOutput, ok := getResponseBodyOutput(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}

	if code != responseStatusCode {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d\nresponse body: %v",
			code,
			responseStatusCode,
			responseBodyOutput,
		)
	}

	return nil
}

/*
*
Step example:

	Then the response body should be:
	"""
	{
		"_id": "441d896b-c0bd-40f4-9926-0568f4a94ec7",
		"name": "Test name",
		"created": 1603882800
	}
	"""
*/
func (a *ApiClient) TheResponseBodyShouldBe(ctx context.Context, doc string) error {
	responseBody, ok := getResponseBody(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}

	// Try to execute template on expected body
	b, err := a.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}

	content := b.Bytes()
	// Try to unmarshal expected body as json
	expectedBody, err := unmarshalJson(content)
	if err != nil {
		return fmt.Errorf("cannot decode expected response body: %w", err)
	}

	if err := checkResponse(responseBody, expectedBody); err != nil {
		return err
	}

	return nil
}

/*
*
Step example:

	Then the response raw body should be:
	"""
	Test
	"""
*/
func (a *ApiClient) TheResponseRawBodyShouldBe(ctx context.Context, doc string) error {
	responseBodyOutput, ok := getResponseBodyOutput(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}

	// Try to execute template on expected body
	b, err := a.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}

	expectedBody := b.String()
	if responseBodyOutput != expectedBody {
		return fmt.Errorf("expected response body to be:\n%v\n but actual is:\n%v",
			expectedBody, responseBodyOutput)
	}

	return nil
}

/*
*
If some fields are not defined in step content they are ignored.

Step example:

	Then the response body should contain:
	"""
	{
		"name": "Test name"
	}
	"""
*/
func (a *ApiClient) TheResponseBodyShouldContain(ctx context.Context, doc string) error {
	responseBody, ok := getResponseBody(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}

	// Try to execute template on expected body
	b, err := a.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}

	content := b.Bytes()
	// Try to umarshal expected body as json
	expectedBody, err := unmarshalJson(content)
	if err != nil {
		return fmt.Errorf("cannot unmarshal json %v: %s", err, content)
	}

	partialBody := getPartialResponse(responseBody, expectedBody)

	if err := checkResponse(partialBody, expectedBody); err != nil {
		return err
	}

	return nil
}

/*
*
Step example:

	Then the response key "data.0.created_at" should not be "0"
*/
func (a *ApiClient) TheResponseKeyShouldNotBe(ctx context.Context, path, value string) error {
	responseBody, ok := getResponseBody(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}

	if nestedVal, ok := getNestedJsonVal(responseBody, strings.Split(path, ".")); ok {
		switch v := nestedVal.(type) {
		case types.Nil:
			if value != "null" {
				return nil
			}
		case string:
			if v != value {
				return nil
			}
		case int:
			if i, err := strconv.ParseInt(value, 10, 0); err != nil || v != int(i) {
				return nil
			}
		case int32:
			if i, err := strconv.ParseInt(value, 10, 0); err != nil || v != int32(i) {
				return nil
			}
		case int64:
			if i, err := strconv.ParseInt(value, 10, 0); err != nil || v != i {
				return nil
			}
		case float32:
			if f, err := strconv.ParseFloat(value, 32); err != nil || v != float32(f) {
				return nil
			}
		case float64:
			if f, err := strconv.ParseFloat(value, 64); err != nil || v != f {
				return nil
			}
		case bool:
			if b, err := strconv.ParseBool(value); err != nil || v != b {
				return nil
			}
		}

		return fmt.Errorf("%v is equal to %v", value, nestedVal)
	}

	responseBodyOutput, ok := getResponseBodyOutput(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}

	return fmt.Errorf("%s not exists in response:\n%v", path, responseBodyOutput)
}

/*
*
Step example:

	Then the response key "data.0.created_at" should not exist
*/
func (a *ApiClient) TheResponseKeyShouldNotExist(ctx context.Context, path string) error {
	responseBody, ok := getResponseBody(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}
	responseBodyOutput, ok := getResponseBodyOutput(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}

	if _, ok := getNestedJsonVal(responseBody, strings.Split(path, ".")); ok {
		return fmt.Errorf("%s exists in response:\n%v", path, responseBodyOutput)
	}

	return nil
}

/*
*
Step example:

	Then the response key "data.0.created_at" should exist
*/
func (a *ApiClient) TheResponseKeyShouldExist(ctx context.Context, path string) error {
	responseBody, ok := getResponseBody(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}
	responseBodyOutput, ok := getResponseBodyOutput(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}
	if _, ok := getNestedJsonVal(responseBody, strings.Split(path, ".")); !ok {
		return fmt.Errorf("%s not exists in response:\n%v", path, responseBodyOutput)
	}

	return nil
}

/*
Step example:

	Then the difference between metaalarmLastEventDate createTimestamp is in range -2,2
*/
func (a *ApiClient) TheDifferenceBetweenValues(ctx context.Context, var1, var2 string, left, right float64) error {
	val1, err := parseFloatVar(ctx, var1)
	if err != nil {
		return fmt.Errorf("first variable %s", err)
	}
	val2, err := parseFloatVar(ctx, var2)
	if err != nil {
		return fmt.Errorf("second variable %s", err)
	}
	d := val1 - val2
	if d < left || right < d {
		return fmt.Errorf("difference is %f and out of range %f, %f", d, left, right)
	}

	return nil
}

/*
*
Step example:

	Then the response key "data.0.duration" should be greater or equal than 3
*/
func (a *ApiClient) TheResponseKeyShouldBeGreaterOrEqualThan(ctx context.Context, path string, value float64) error {
	responseBody, ok := getResponseBody(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}
	responseBodyOutput, ok := getResponseBodyOutput(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}
	if nestedVal, ok := getNestedJsonVal(responseBody, strings.Split(path, ".")); ok {
		var fieldVal float64
		switch v := nestedVal.(type) {
		case int:
			fieldVal = float64(v)
		case int32:
			fieldVal = float64(v)
		case int64:
			fieldVal = float64(v)
		case float32:
			fieldVal = float64(v)
		case float64:
			fieldVal = v
		default:
			return fmt.Errorf("%v is not number", nestedVal)
		}

		if fieldVal >= value {
			return nil
		}

		return fmt.Errorf("%v is lesser then %v", fieldVal, value)
	}

	return fmt.Errorf("%s not exists in response:\n%v", path, responseBodyOutput)
}

/*
*
Step example:

	Then the response key "data.0.duration" should be greater than 3
*/
func (a *ApiClient) TheResponseKeyShouldBeGreaterThan(ctx context.Context, path string, value float64) error {
	responseBody, ok := getResponseBody(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}
	responseBodyOutput, ok := getResponseBodyOutput(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}

	if nestedVal, ok := getNestedJsonVal(responseBody, strings.Split(path, ".")); ok {
		var fieldVal float64
		switch v := nestedVal.(type) {
		case int:
			fieldVal = float64(v)
		case int32:
			fieldVal = float64(v)
		case int64:
			fieldVal = float64(v)
		case float32:
			fieldVal = float64(v)
		case float64:
			fieldVal = v
		default:
			return fmt.Errorf("%v is not number", nestedVal)
		}

		if fieldVal > value {
			return nil
		}

		return fmt.Errorf("%v is lesser or equal then %v", fieldVal, value)
	}

	return fmt.Errorf("%s not exists in response:\n%v", path, responseBodyOutput)
}

// TheResponseArrayKeyShouldContain
// Step example:
//
//	Then the response array key "data.0.v.steps" should contain:
//	"""
//	[
//	  {
//	    "_t": "stateinc"
//	  }
//	]
//	"""
func (a *ApiClient) TheResponseArrayKeyShouldContain(ctx context.Context, path string, doc string) error {
	responseBody, ok := getResponseBody(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}
	responseBodyOutput, ok := getResponseBodyOutput(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}

	b, err := a.templater.Execute(ctx, path)
	if err != nil {
		return err
	}

	path = b.String()

	b, err = a.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}

	if nestedVal, ok := getNestedJsonVal(responseBody, strings.Split(path, ".")); ok {
		receivedStr, _ := json.MarshalIndent(nestedVal, "", "  ")

		switch received := nestedVal.(type) {
		case []interface{}:
			expected := make([]map[string]interface{}, 0)
			err := json.Unmarshal(b.Bytes(), &expected)
			if err != nil {
				expected := make([]interface{}, 0)
				err := json.Unmarshal(b.Bytes(), &expected)
				if err != nil {
					return err
				}
				for _, ev := range expected {
					found := false
					for _, v := range received {
						if err := checkResponse(v, ev); err == nil {
							found = true
							break
						}
					}

					if !found {
						return fmt.Errorf("%s\nis not in:\n%s", ev, receivedStr)
					}
				}

				return nil
			}

			if len(expected) == 0 && len(received) != 0 {
				return fmt.Errorf("%s is not empty", path)
			}

			for _, ev := range expected {
				if len(ev) == 0 {
					return fmt.Errorf("%s contains empty element", doc)
				}

				found := false
				for _, v := range received {
					if err := checkResponse(getPartialResponse(v, ev), ev); err == nil {
						found = true
						break
					}
				}

				if !found {
					expectedStr, _ := json.MarshalIndent(ev, "", "  ")
					return fmt.Errorf("%s\nis not in:\n%s", expectedStr, receivedStr)
				}
			}

			return nil
		}

		return fmt.Errorf("%s is not array but %T:\n%s", path, nestedVal, receivedStr)
	}

	return fmt.Errorf("%s not exists in response:\n%v", path, responseBodyOutput)
}

// TheResponseArrayKeyShouldContainOnly
// Step example:
//
//	Then the response array key "data.0.v.steps" should contain only:
//	[
//	  {
//	    "_t": "stateinc"
//	  }
//	]
//	"""
func (a *ApiClient) TheResponseArrayKeyShouldContainOnly(ctx context.Context, path string, doc string) error {
	responseBody, ok := getResponseBody(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}
	responseBodyOutput, ok := getResponseBodyOutput(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}

	b, err := a.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}

	if nestedVal, ok := getNestedJsonVal(responseBody, strings.Split(path, ".")); ok {
		receivedStr, _ := json.MarshalIndent(nestedVal, "", "  ")

		switch received := nestedVal.(type) {
		case []interface{}:
			expected := make([]map[string]interface{}, 0)
			err := json.Unmarshal(b.Bytes(), &expected)
			if err != nil {
				expected := make([]interface{}, 0)
				err := json.Unmarshal(b.Bytes(), &expected)
				if err != nil {
					return err
				}

				if len(expected) != len(received) {
					return fmt.Errorf("expected %d items but receieved:\n%s", len(expected), receivedStr)
				}

				for _, ev := range expected {
					found := false
					for _, v := range received {
						if err := checkResponse(v, ev); err == nil {
							found = true
							break
						}
					}

					if !found {
						return fmt.Errorf("%s\nis not in:\n%s", ev, receivedStr)
					}
				}

				return nil
			}

			if len(expected) == 0 {
				return fmt.Errorf("%s is empty", doc)
			}

			if len(expected) != len(received) {
				return fmt.Errorf("expected %d items but receieved:\n%s", len(expected), receivedStr)
			}

			for _, ev := range expected {
				if len(ev) == 0 {
					return fmt.Errorf("%s contains empty element", doc)
				}

				found := false
				for _, v := range received {
					if err := checkResponse(getPartialResponse(v, ev), ev); err == nil {
						found = true
						break
					}
				}

				if !found {
					expectedStr, _ := json.MarshalIndent(ev, "", "  ")
					return fmt.Errorf("%s\nis not in:\n%s", expectedStr, receivedStr)
				}
			}

			return nil
		}

		return fmt.Errorf("%s is not array but %T:\n%s", path, nestedVal, receivedStr)
	}

	return fmt.Errorf("%s not exists in response:\n%v", path, responseBodyOutput)
}

// TheResponseArrayKeyShouldContainInOrder
// Step example:
//
//	Then the response array key "data.0.v.steps" should contain in order:
//	"""
//	[
//	  {
//	    "_t": "stateinc"
//	  }
//	]
//	"""
func (a *ApiClient) TheResponseArrayKeyShouldContainInOrder(ctx context.Context, path string, doc string) error {
	responseBody, ok := getResponseBody(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}
	responseBodyOutput, ok := getResponseBodyOutput(ctx)
	if !ok {
		return fmt.Errorf("response is nil")
	}

	b, err := a.templater.Execute(ctx, doc)
	if err != nil {
		return err
	}

	if nestedVal, ok := getNestedJsonVal(responseBody, strings.Split(path, ".")); ok {
		receivedStr, _ := json.MarshalIndent(nestedVal, "", "  ")

		switch received := nestedVal.(type) {
		case []interface{}:
			expected := make([]map[string]interface{}, 0)
			err := json.Unmarshal(b.Bytes(), &expected)
			if err != nil {
				expected := make([]interface{}, 0)
				err := json.Unmarshal(b.Bytes(), &expected)
				if err != nil {
					return err
				}
				prevIndex := -1
				for _, ev := range expected {
					foundIndex := -1
					for j, v := range received {
						if err := checkResponse(v, ev); err == nil {
							foundIndex = j
							break
						}
					}

					if foundIndex < 0 {
						return fmt.Errorf("%s\nis not in:\n%s", ev, receivedStr)
					}

					if prevIndex > foundIndex {
						return fmt.Errorf("invalid order:\n%s", pretty.Compare(received, expected))
					}

					prevIndex = foundIndex
				}

				return nil
			}

			if len(expected) == 0 {
				return fmt.Errorf("%s is empty", doc)
			}

			prevIndex := -1
			for _, ev := range expected {
				if len(ev) == 0 {
					return fmt.Errorf("%s contains empty element", doc)
				}

				foundIndex := -1
				for j, v := range received {
					if err := checkResponse(getPartialResponse(v, ev), ev); err == nil {
						foundIndex = j
						break
					}
				}

				if foundIndex < 0 {
					expectedStr, _ := json.MarshalIndent(ev, "", "  ")
					return fmt.Errorf("%s\nis not in:\n%s", expectedStr, receivedStr)
				}

				if prevIndex > foundIndex {
					return fmt.Errorf("invalid order:\n%s", pretty.Compare(received, expected))
				}

				prevIndex = foundIndex
			}

			return nil
		}

		return fmt.Errorf("%s is not array but %T:\n%s", path, nestedVal, receivedStr)
	}

	return fmt.Errorf("%s not exists in response:\n%v", path, responseBodyOutput)
}

// getNestedJsonVal returns val by path.
func getNestedJsonVal(v interface{}, path []string) (interface{}, bool) {
	field := path[0]

	if i, err := strconv.Atoi(field); err == nil {
		if ar, ok := v.([]interface{}); ok {
			if i >= 0 && i < len(ar) {
				fv := ar[i]
				if len(path) == 1 {
					return fv, true
				}

				return getNestedJsonVal(fv, path[1:])
			}
		}

		return nil, false
	}

	if m, ok := v.(map[string]interface{}); ok {
		if fv, ok := m[field]; ok {
			if len(path) == 1 {
				return fv, true
			}

			return getNestedJsonVal(fv, path[1:])
		}
	}

	return nil, false
}

/*
*
Step example:

	Given I am admin
*/
func (a *ApiClient) IAm(ctx context.Context, role string) (context.Context, error) {
	var line model.Rbac
	res := a.db.Collection(mongo.RightsMongoCollection).FindOne(ctx, bson.M{
		"crecord_type": model.LineTypeRole,
		"crecord_name": role,
	})
	if err := res.Err(); err != nil {
		return ctx, fmt.Errorf("cannot fetch role: %w", err)
	}

	err := res.Decode(&line)
	if err != nil {
		return ctx, fmt.Errorf("cannot decode role: %w", err)
	}

	res = a.db.Collection(mongo.RightsMongoCollection).FindOne(ctx, bson.M{
		"crecord_type": model.LineTypeSubject,
		"role":         line.ID,
	})
	if err := res.Err(); err != nil {
		return ctx, fmt.Errorf("cannot fetch user: %w", err)
	}

	err = res.Decode(&line)
	if err != nil {
		return ctx, fmt.Errorf("cannot decode user: %w", err)
	}

	uri := fmt.Sprintf("%s/api/v4/login", a.url)
	body, err := json.Marshal(map[string]string{
		"username": line.Name,
		"password": userPass,
	})
	if err != nil {
		return ctx, err
	}
	request, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(body))
	if err != nil {
		return ctx, fmt.Errorf("cannot create login request: %w", err)
	}

	request.Header.Set(headerContentType, binding.MIMEJSON)

	response, err := a.client.Do(request)
	if err != nil {
		return ctx, fmt.Errorf("cannot do login request: %w", err)
	}

	buf, err := io.ReadAll(response.Body)
	if err != nil {
		return ctx, fmt.Errorf("cannot fetch login response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return ctx, fmt.Errorf("unexpected response status %v %s", response.StatusCode, buf)
	}

	responseBody := make(map[string]string)
	err = json.Unmarshal(buf, &responseBody)
	if err != nil {
		return ctx, fmt.Errorf("cannot decode login response: %w", err)
	}

	token := responseBody["access_token"]
	if token == "" {
		return ctx, fmt.Errorf("unexpected login response %v", buf)
	}

	headers, ok := getHeaders(ctx)
	if !ok {
		headers = make(map[string]string, 1)
	}
	ctx = setApiAuthToken(ctx, token)
	headers[headerAuthorization] = bearerPrefix + " " + token
	ctx = setHeaders(ctx, headers)

	return ctx, nil
}

/*
*
Step example:

	When I am authenticated with username "user" password "pass"
*/
func (a *ApiClient) IAmAuthenticatedByBasicAuth(ctx context.Context, username, password string) (context.Context, error) {
	headers, ok := getHeaders(ctx)
	if !ok {
		headers = make(map[string]string, 1)
	}
	headers[headerAuthorization] = basicPrefix + " " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
	ctx = setHeaders(ctx, headers)

	return ctx, nil
}

/*
*
Step example:

	When I send an event:
	"""
	  {
		"connector" : "test_post_connector",
		"connector_name" : "test_post_connector_name",
		"source_type" : "resource",
		"event_type" : "check",
		"component" : "test_post_component",
		"resource" : "test_post_resource",
		"state" : 1,
		"output" : "noveo alarm"
	  }
	"""
*/
func (a *ApiClient) ISendAnEvent(ctx context.Context, doc string) (context.Context, error) {
	uri := fmt.Sprintf("%s/api/v4/event", a.url)
	b, err := a.templater.Execute(ctx, doc)
	if err != nil {
		return ctx, err
	}

	events := make([]map[string]interface{}, 0)
	err = json.Unmarshal(b.Bytes(), &events)
	if err != nil {
		event := make(map[string]interface{})
		err = json.Unmarshal(b.Bytes(), &event)
		if err != nil {
			return ctx, err
		}
		events = append(events, event)
	}

	for i := range events {
		events[i]["debug"] = true
	}

	body, err := json.Marshal(events)
	if err != nil {
		return ctx, err
	}
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return ctx, fmt.Errorf("cannot create event request: %w", err)
	}

	req.Header.Set(headerContentType, binding.MIMEJSON)
	ctx, err = a.doRequest(ctx, req)
	if err != nil {
		return ctx, err
	}

	err = a.TheResponseCodeShouldBe(ctx, http.StatusOK)
	if err != nil {
		return ctx, err
	}

	err = a.TheResponseBodyShouldContain(ctx, fmt.Sprintf("{\"sent_events\":%s}", body))
	return ctx, err
}

/*
*
Step example:

	When I do GET /api/v4/alarms
	When I do GET /api/v4/entitybasic/{{ .lastResponse._id}}
*/
func (a *ApiClient) IDoRequest(ctx context.Context, method, uri string) (context.Context, error) {
	if strings.Contains(uri, "until") {
		return ctx, fmt.Errorf("step is wrongly matched to IDoRequest")
	}

	req, err := a.createRequest(ctx, method, uri, "")
	if err != nil {
		return ctx, err
	}

	return a.doRequest(ctx, req)
}

/*
*
Step example:

	When I do POST /api/v4/event:
	"""
	  {
		"connector" : "test_post_connector",
		"connector_name" : "test_post_connector_name",
		"source_type" : "resource",
		"event_type" : "check",
		"component" : "test_post_component",
		"resource" : "test_post_resource",
		"state" : 1,
		"output" : "noveo alarm"
	  }
	"""
	When I do PUT /api/v4/entitybasics/{{ .lastResponse._id}}:
	"""
	  {
		"state": 1
	  }
	"""
*/
func (a *ApiClient) IDoRequestWithBody(ctx context.Context, method, uri string, doc string) (context.Context, error) {
	if doc == "" {
		return ctx, fmt.Errorf("body is empty")
	}
	if strings.Contains(uri, "until") {
		return ctx, fmt.Errorf("step is wrongly matched to IDoRequestWithBody")
	}

	req, err := a.createRequest(ctx, method, uri, doc)
	if err != nil {
		return ctx, err
	}

	if headers, ok := getHeaders(ctx); ok {
		if _, ok := headers[headerContentType]; !ok {
			req.Header.Set(headerContentType, binding.MIMEJSON)
		}
	} else {
		req.Header.Set(headerContentType, binding.MIMEJSON)
	}

	return a.doRequest(ctx, req)
}

/*
*
Step example:

	When I do GET /api/v4/entitybasic/{{ .lastResponse._id}} until response code is 200
*/
func (a *ApiClient) IDoRequestUntilResponseCode(ctx context.Context, method, uri string, code int) (context.Context, error) {
	req, ctx, err := a.createRequestWithSavedRequest(ctx, method, uri)
	if err != nil {
		return ctx, err
	}

	ok, ctx, err := a.doRequestUntil(ctx, req, func(ctx context.Context) bool {
		responseStatusCode, _ := getResponseStatusCode(ctx)
		return code == responseStatusCode
	})

	if err != nil || ok {
		return ctx, err
	}

	responseStatusCode, _ := getResponseStatusCode(ctx)
	responseBodyOutput, _ := getResponseBodyOutput(ctx)

	return ctx, fmt.Errorf("max retries exceeded, expected response code to be: %d, but actual is: %d\nresponse body: %v",
		code,
		responseStatusCode,
		responseBodyOutput,
	)
}

func (a *ApiClient) ISaveRequest(ctx context.Context, doc string) (context.Context, error) {
	ctx = setRequestBody(ctx, doc)
	return ctx, nil
}

/*
*
Step example:

	When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body is:
	"""
	{
	  "status": "done"
	}
	"""
*/
func (a *ApiClient) IDoRequestUntilResponse(ctx context.Context, method, uri string, code int, doc string) (context.Context, error) {
	if doc == "" {
		return ctx, fmt.Errorf("body is empty")
	}

	req, ctx, err := a.createRequestWithSavedRequest(ctx, method, uri)
	if err != nil {
		return ctx, err
	}

	b, err := a.templater.Execute(ctx, doc)
	if err != nil {
		return ctx, err
	}
	content := b.Bytes()
	expectedBody, err := unmarshalJson(content)
	if err != nil {
		return ctx, fmt.Errorf("cannot decode expected response body: %w", err)
	}

	ok, ctx, err := a.doRequestUntil(ctx, req, func(ctx context.Context) bool {
		responseStatusCode, _ := getResponseStatusCode(ctx)
		responseBody, _ := getResponseBody(ctx)

		return code == responseStatusCode && checkResponse(responseBody, expectedBody) == nil
	})

	if err != nil || ok {
		return ctx, err
	}

	responseStatusCode, _ := getResponseStatusCode(ctx)
	if code != responseStatusCode {
		responseBodyOutput, _ := getResponseBodyOutput(ctx)

		return ctx, fmt.Errorf("max retries exceeded: expected response code to be: %d, but actual is: %d\nresponse body: %v",
			code,
			responseStatusCode,
			responseBodyOutput,
		)
	}

	responseBody, _ := getResponseBody(ctx)

	return ctx, fmt.Errorf("max retries exceeded: %w", checkResponse(responseBody, expectedBody))
}

/*
*
Step example:

	When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and body contains:
	"""
	{
	  "status": "done"
	}
	"""
*/
func (a *ApiClient) IDoRequestUntilResponseContains(ctx context.Context, method, uri string, code int, doc string) (context.Context, error) {
	if doc == "" {
		return ctx, fmt.Errorf("body is empty")
	}

	req, ctx, err := a.createRequestWithSavedRequest(ctx, method, uri)
	if err != nil {
		return ctx, err
	}

	b, err := a.templater.Execute(ctx, doc)
	if err != nil {
		return ctx, err
	}
	content := b.Bytes()
	expectedBody, err := unmarshalJson(content)
	if err != nil {
		return ctx, fmt.Errorf("cannot decode expected response body: %w", err)
	}

	ok, ctx, err := a.doRequestUntil(ctx, req, func(ctx context.Context) bool {
		responseStatusCode, _ := getResponseStatusCode(ctx)
		if code == responseStatusCode {
			responseBody, _ := getResponseBody(ctx)
			partialBody := getPartialResponse(responseBody, expectedBody)

			return checkResponse(partialBody, expectedBody) == nil
		}

		return false
	})

	if err != nil || ok {
		return ctx, err
	}

	responseStatusCode, _ := getResponseStatusCode(ctx)
	if code != responseStatusCode {
		responseBodyOutput, _ := getResponseBodyOutput(ctx)

		return ctx, fmt.Errorf("max retries exceeded: expected response code to be: %d, but actual is: %d\nresponse body: %v",
			code,
			responseStatusCode,
			responseBodyOutput,
		)
	}

	responseBody, _ := getResponseBody(ctx)
	partialBody := getPartialResponse(responseBody, expectedBody)
	return ctx, fmt.Errorf("max retries exceeded: %w", checkResponse(partialBody, expectedBody))
}

/*
*
Step example:

	When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and response key "data.0.duration" is greater or equal than 3
	"""
*/
func (a *ApiClient) IDoRequestUntilResponseKeyIsGreaterOrEqualThan(ctx context.Context, method, uri string, code int, path string, value float64) (context.Context, error) {
	req, ctx, err := a.createRequestWithSavedRequest(ctx, method, uri)
	if err != nil {
		return ctx, err
	}

	ok, ctx, err := a.doRequestUntil(ctx, req, func(ctx context.Context) bool {
		responseStatusCode, _ := getResponseStatusCode(ctx)
		if code == responseStatusCode {
			err := a.TheResponseKeyShouldBeGreaterOrEqualThan(ctx, path, value)
			return err == nil
		}

		return false
	})

	if err != nil || ok {
		return ctx, err
	}

	responseStatusCode, _ := getResponseStatusCode(ctx)
	if code != responseStatusCode {
		responseBodyOutput, _ := getResponseBodyOutput(ctx)

		return ctx, fmt.Errorf("max retries exceeded: expected response code to be: %d, but actual is: %d\nresponse body: %v",
			code,
			responseStatusCode,
			responseBodyOutput,
		)
	}

	return ctx, fmt.Errorf("max retries exceeded: %w", a.TheResponseKeyShouldBeGreaterOrEqualThan(ctx, path, value))
}

// IDoRequestUntilResponseArrayKeyContains
// Step example:
//
//	When I do GET /api/v4/alarms until response code is 200 and response array key "data.0.v.steps" contains:
//	"""
//	[
//	  {
//	    "_t": "stateinc"
//	  }
//	]
//	"""
func (a *ApiClient) IDoRequestUntilResponseArrayKeyContains(ctx context.Context, method, uri string, code int, path string, doc string) (context.Context, error) {
	req, ctx, err := a.createRequestWithSavedRequest(ctx, method, uri)
	if err != nil {
		return ctx, err
	}

	ok, ctx, err := a.doRequestUntil(ctx, req, func(ctx context.Context) bool {
		responseStatusCode, _ := getResponseStatusCode(ctx)
		if code == responseStatusCode {
			err := a.TheResponseArrayKeyShouldContain(ctx, path, doc)
			return err == nil
		}

		return false
	})

	if err != nil || ok {
		return ctx, err
	}

	responseStatusCode, _ := getResponseStatusCode(ctx)
	if code != responseStatusCode {
		responseBodyOutput, _ := getResponseBodyOutput(ctx)

		return ctx, fmt.Errorf("max retries exceeded: expected response code to be: %d, but actual is: %d\nresponse body: %v",
			code,
			responseStatusCode,
			responseBodyOutput,
		)
	}

	return ctx, fmt.Errorf("max retries exceeded: %w", a.TheResponseArrayKeyShouldContain(ctx, path, doc))
}

// IDoRequestUntilResponseArrayKeyContainsOnly
// Step example:
//
//	When I do GET /api/v4/alarms until response code is 200 and response array key "data.0.v.steps" contains only:
//	"""
//	[
//	  {
//	    "_t": "stateinc"
//	  }
//	]
//	"""
func (a *ApiClient) IDoRequestUntilResponseArrayKeyContainsOnly(ctx context.Context, method, uri string, code int, path string, doc string) (context.Context, error) {
	req, ctx, err := a.createRequestWithSavedRequest(ctx, method, uri)
	if err != nil {
		return ctx, err
	}

	ok, ctx, err := a.doRequestUntil(ctx, req, func(ctx context.Context) bool {
		responseStatusCode, _ := getResponseStatusCode(ctx)
		if code == responseStatusCode {
			err := a.TheResponseArrayKeyShouldContainOnly(ctx, path, doc)
			return err == nil
		}

		return false
	})

	if err != nil || ok {
		return ctx, err
	}

	responseStatusCode, _ := getResponseStatusCode(ctx)
	if code != responseStatusCode {
		responseBodyOutput, _ := getResponseBodyOutput(ctx)

		return ctx, fmt.Errorf("max retries exceeded: expected response code to be: %d, but actual is: %d\nresponse body: %v",
			code,
			responseStatusCode,
			responseBodyOutput,
		)
	}

	return ctx, fmt.Errorf("max retries exceeded: %w", a.TheResponseArrayKeyShouldContainOnly(ctx, path, doc))
}

/*
*
Step example:

	When I set header Content-Type=application/json
*/
func (a *ApiClient) ISetRequestHeader(ctx context.Context, key, value string) (context.Context, error) {
	b, err := a.templater.Execute(ctx, value)
	if err != nil {
		return ctx, err
	}

	headers, ok := getHeaders(ctx)
	if ok {
		headers[key] = b.String()
	} else {
		headers = map[string]string{key: b.String()}
	}
	ctx = setHeaders(ctx, headers)

	return ctx, nil
}

/*
*
Step example:

	When I save response id={{ .lastResponse._id }}
*/
func (a *ApiClient) ISaveResponse(ctx context.Context, key, value string) (context.Context, error) {
	b, err := a.templater.Execute(ctx, value)
	if err != nil {
		return ctx, err
	}

	return setVar(ctx, key, b.String()), nil
}

// ValueShouldBeGteLteThan
// Step example:
//
//	Then "value1" > "value2"
//	Then "value1" <= "value2"
func (a *ApiClient) ValueShouldBeGteLteThan(ctx context.Context, left, op, right string) error {
	leftV, err := parseFloatVar(ctx, left)
	if err != nil {
		return err
	}
	rightV, err := parseFloatVar(ctx, right)
	if err != nil {
		return err
	}
	switch op {
	case "<":
		if !(leftV < rightV) {
			return fmt.Errorf("%q is not lesser than %q (%v < %v)", left, right, leftV, rightV)
		}
	case "<=":
		if !(leftV <= rightV) {
			return fmt.Errorf("%q is not lesser or equal than %q (%v <= %v)", left, right, leftV, rightV)
		}
	case ">":
		if !(leftV > rightV) {
			return fmt.Errorf("%q is not greater than %q (%v > %v)", left, right, leftV, rightV)
		}
	case ">=":
		if !(leftV >= rightV) {
			return fmt.Errorf("%q is not greater or equal than %q (%v >= %v)", left, right, leftV, rightV)
		}
	default:
		return fmt.Errorf("unknown operator %q", op)
	}

	return nil
}

func (a *ApiClient) createRequest(ctx context.Context, method, uri, body string) (*http.Request, error) {
	uri, err := a.getRequestURL(ctx, uri)
	if err != nil {
		return nil, err
	}

	var r io.Reader
	if body != "" {
		r, err = a.getRequestBody(ctx, body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, uri, r)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	return req, nil
}

func (a *ApiClient) createRequestWithSavedRequest(ctx context.Context, method, uri string) (*http.Request, context.Context, error) {
	uri, err := a.getRequestURL(ctx, uri)
	if err != nil {
		return nil, ctx, err
	}

	var r io.Reader
	body, _ := getRequestBody(ctx)
	if body != "" {
		r, err = a.getRequestBody(ctx, body)
		if err != nil {
			return nil, ctx, err
		}
		ctx = setRequestBody(ctx, "")
	}

	req, err := http.NewRequest(method, uri, r)
	if err != nil {
		return nil, ctx, fmt.Errorf("cannot create request: %w", err)
	}

	return req, ctx, nil
}

// doRequest adds auth credentials and makes request.
func (a *ApiClient) doRequest(ctx context.Context, req *http.Request) (context.Context, error) {
	scName, _ := GetScenarioName(ctx)
	scUri, _ := GetScenarioUri(ctx)

	if headers, ok := getHeaders(ctx); ok {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	// Add session's cookies
	if cookies, ok := getCookies(ctx); ok {
		for _, c := range cookies {
			req.AddCookie(c)
		}
	}

	var err error
	var responseBody interface{}
	var responseBodyOutput string
	dumpReq, _ := httputil.DumpRequest(req, true)
	response, err := a.client.Do(req)
	// Read response
	if err != nil {
		a.requestLogger.Err(err).
			Str("file", scUri).
			Str("scenario", scName).
			Str("request", string(dumpReq)).
			Msg("invalid called request")
		return ctx, fmt.Errorf("cannot do request: %w", err)
	}

	dumpRes, _ := httputil.DumpResponse(response, true)
	a.requestLogger.Info().
		Str("file", scUri).
		Str("scenario", scName).
		Str("request", string(dumpReq)).
		Str("response", string(dumpRes)).
		Msg("called request")

	buf, err := io.ReadAll(response.Body)
	if err != nil {
		return ctx, fmt.Errorf("cannot fetch response: %w", err)
	}

	// Parse response
	if len(buf) > 0 {
		responseBody, err = unmarshalJson(buf)
		if err == nil {
			ibuf, _ := json.MarshalIndent(responseBody, "", "  ")
			responseBodyOutput = string(ibuf)
		} else {
			responseBodyOutput = string(buf)
		}
	}

	// Save session
	resCookies := response.Cookies()
	cookies := make([]*http.Cookie, 0, len(resCookies))
	for _, cookie := range resCookies {
		if cookie.MaxAge > 0 {
			cookies = append(cookies, cookie)
		}
	}

	ctx = setResponseBody(ctx, responseBody)
	ctx = setResponseBodyOutput(ctx, responseBodyOutput)
	ctx = setResponseStatusCode(ctx, response.StatusCode)
	ctx = setCookies(ctx, cookies)

	return ctx, nil
}

func (a *ApiClient) doRequestUntil(
	ctx context.Context,
	req *http.Request,
	check func(context.Context) bool,
) (bool, context.Context, error) {
	body := req.Body
	var err error
	if body != nil {
		body, req.Body, err = libhttp.DrainBody(body)
		if err != nil {
			return false, ctx, err
		}
	}

	timeout := startRepeatRequestInterval
	start := time.Now()

	for {
		ctx, err = a.doRequest(ctx, req)
		if err != nil {
			return false, ctx, err
		}

		if check(ctx) {
			return true, ctx, nil
		}

		if time.Since(start) > totalRepeatRequestInterval {
			break
		}

		select {
		case <-ctx.Done():
			return false, ctx, ctx.Err()
		case <-time.After(timeout):
		}

		timeout *= 2

		if body != nil {
			body, req.Body, err = libhttp.DrainBody(body)
			if err != nil {
				return false, ctx, err
			}
		}
	}

	return false, ctx, nil
}

// getRequestURL applies template uri to last response data.
func (a *ApiClient) getRequestURL(ctx context.Context, uri string) (string, error) {
	b, err := a.templater.Execute(ctx, uri)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", a.url, b.String()), nil
}

// getRequestBody executes template body.
func (a *ApiClient) getRequestBody(ctx context.Context, body string) (io.Reader, error) {
	return a.templater.Execute(ctx, body)
}

// getPartialResponse removes fields from received which are not presented in expected.
func getPartialResponse(received, expected interface{}) interface{} {
	receviedRef := reflect.ValueOf(received)
	expectedRef := reflect.ValueOf(expected)

	if !receviedRef.IsValid() || !expectedRef.IsValid() || receviedRef.Type() != expectedRef.Type() {
		return received
	}

	switch receviedRef.Kind() {
	case reflect.Array:
		receivedLen := receviedRef.Len()
		expectedLen := expectedRef.Len()
		minLen := receivedLen
		if minLen > expectedLen {
			minLen = expectedLen
		}

		res := make([]interface{}, receivedLen)
		for i := 0; i < minLen; i++ {
			res[i] = getPartialResponse(receviedRef.Index(i).Interface(), expectedRef.Index(i).Interface())
		}
		for i := minLen; i < receivedLen; i++ {
			res[i] = receviedRef.Index(i).Interface()
		}

		return res
	case reflect.Slice:
		if receviedRef.IsNil() || expectedRef.IsNil() {
			return received
		}
		receivedLen := receviedRef.Len()
		expectedLen := expectedRef.Len()
		minLen := receivedLen
		if minLen > expectedLen {
			minLen = expectedLen
		}

		res := make([]interface{}, receivedLen)
		for i := 0; i < minLen; i++ {
			res[i] = getPartialResponse(receviedRef.Index(i).Interface(), expectedRef.Index(i).Interface())
		}
		for i := minLen; i < receivedLen; i++ {
			res[i] = receviedRef.Index(i).Interface()
		}

		return res
	case reflect.Map:
		if receviedRef.IsNil() || expectedRef.IsNil() {
			return received
		}
		res := make(map[string]interface{})
		for _, k := range expectedRef.MapKeys() {
			receivedVal := receviedRef.MapIndex(k)
			expectedVal := expectedRef.MapIndex(k)

			if receivedVal.IsValid() {
				res[k.String()] = getPartialResponse(receivedVal.Interface(), expectedVal.Interface())
			}
		}

		return res
	default:
		return received
	}
}

// checkResponse returns error which contains differences between received and expected.
func checkResponse(received, expected interface{}) error {
	if diff := pretty.Compare(received, expected); diff != "" {
		return fmt.Errorf("response doesn't match expected response body:\n%s\n", diff)
	}

	return nil
}

// unmarshalJson decodes JSON to structure where all numbers are decoded to int64 or float64.
func unmarshalJson(data []byte) (interface{}, error) {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	var v interface{}
	err := d.Decode(&v)
	if err != nil {
		return nil, err
	}

	v = convertJsonNumber(v)

	return v, nil
}

// convertJsonNumber transforms json.Number in all nested fields to int64 or float64 which value is fit.
func convertJsonNumber(v interface{}) interface{} {
	if n, ok := v.(json.Number); ok {
		i, err := n.Int64()
		if err == nil {
			return i
		}

		f, err := n.Float64()
		if err == nil {
			return f
		}

		return n.String()
	}

	refVal := reflect.ValueOf(v)
	if !refVal.IsValid() {
		return v
	}

	switch refVal.Kind() {
	case reflect.Array:
		l := refVal.Len()
		converted := make([]interface{}, l)
		for i := 0; i < l; i++ {
			converted[i] = convertJsonNumber(refVal.Index(i).Interface())
		}

		return converted
	case reflect.Slice:
		if refVal.IsNil() {
			return v
		}
		l := refVal.Len()
		converted := make([]interface{}, l)
		for i := 0; i < l; i++ {
			converted[i] = convertJsonNumber(refVal.Index(i).Interface())
		}

		return converted
	case reflect.Map:
		if refVal.IsNil() {
			return v
		}
		converted := make(map[string]interface{})
		for _, k := range refVal.MapKeys() {
			item := refVal.MapIndex(k)

			if item.IsValid() {
				converted[k.String()] = convertJsonNumber(item.Interface())
			}
		}

		return converted
	default:
		return v
	}
}
