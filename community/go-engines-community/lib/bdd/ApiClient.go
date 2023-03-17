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
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"

	libtypes "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin/binding"
	"github.com/kylelemons/godebug/pretty"
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
	url *url.URL
	// client is http client to make API requests.
	client *http.Client
	// db is db client.
	db mongo.DbClient
	// response is http response of made API request.
	response           *http.Response
	responseBody       interface{}
	responseBodyOutput string
	// cookies is http cookies which are retrieved from API response and used in following steps.
	// todo remove after session remove
	cookies []*http.Cookie
	// vars is used to save data between steps.
	vars map[string]string
	// request header
	headers map[string]string
}

// NewApiClient creates new API client.
func NewApiClient(db mongo.DbClient) (*ApiClient, error) {
	apiUrl, err := GetApiURL()
	if err != nil {
		return nil, err
	}

	var apiClient ApiClient
	apiClient.client = &http.Client{
		Timeout: requestTimeout,
	}
	apiClient.url = apiUrl
	apiClient.db = db
	apiClient.headers = make(map[string]string)

	return &apiClient, nil
}

// GetApiURL retrieves API url from env var.
func GetApiURL() (*url.URL, error) {
	legacy := os.Getenv(ApiEnvURL)
	if legacy == "" {
		return nil, fmt.Errorf("environment variable %s empty", ApiEnvURL)
	}

	parsed, err := url.Parse(legacy)
	if err != nil {
		return nil, fmt.Errorf("cannot parse api url: %w", err)
	}

	return parsed, nil
}

// ResetResponse clears all saved response data.
func (a *ApiClient) ResetResponse(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	a.response = nil
	a.responseBody = nil
	a.responseBodyOutput = ""
	a.cookies = nil
	a.vars = nil
	a.headers = make(map[string]string)

	return ctx, nil
}

/*
*
Step example:

	Then the response code should be 200
*/
func (a *ApiClient) TheResponseCodeShouldBe(code int) error {
	if a.response == nil {
		return fmt.Errorf("response is nil")
	}

	if code != a.response.StatusCode {
		return fmt.Errorf("expected response code to be: %d, but actual is: %d\nresponse body: %v",
			code,
			a.response.StatusCode,
			a.responseBodyOutput,
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
func (a *ApiClient) TheResponseBodyShouldBe(doc string) error {
	if a.responseBody == nil {
		return fmt.Errorf("response is nil")
	}

	// Try to execute template on expected body
	b, err := a.executeTemplate(doc)
	if err != nil {
		return err
	}

	content := b.Bytes()
	// Try to unmarshal expected body as json
	expectedBody, err := unmarshalJson(content)
	if err != nil {
		return fmt.Errorf("cannot decode expected response body: %w", err)
	}

	if err := checkResponse(a.responseBody, expectedBody); err != nil {
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
func (a *ApiClient) TheResponseRawBodyShouldBe(doc string) error {
	// Try execute template on expected body
	b, err := a.executeTemplate(doc)
	if err != nil {
		return err
	}

	expectedBody := b.String()
	if a.responseBodyOutput != expectedBody {
		return fmt.Errorf("expected response body to be:\n%v\n but actual is:\n%v",
			expectedBody, a.responseBodyOutput)
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
func (a *ApiClient) TheResponseBodyShouldContain(doc string) error {
	if a.responseBody == nil {
		return fmt.Errorf("response is nil")
	}

	// Try execute template on expected body
	b, err := a.executeTemplate(doc)
	if err != nil {
		return err
	}

	content := b.Bytes()
	// Try to umarshal expected body as json
	expectedBody, err := unmarshalJson(content)
	if err != nil {
		return fmt.Errorf("cannot unmarshal json %v: %s", err, content)
	}

	partialBody := getPartialResponse(a.responseBody, expectedBody)

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
func (a *ApiClient) TheResponseKeyShouldNotBe(path, value string) error {
	if nestedVal, ok := getNestedJsonVal(a.responseBody, strings.Split(path, ".")); ok {
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

	return fmt.Errorf("%s not exists in response:\n%v", path, a.responseBodyOutput)
}

/*
*
Step example:

	Then the response key "data.0.created_at" should not exist
*/
func (a *ApiClient) TheResponseKeyShouldNotExist(path string) error {
	if _, ok := getNestedJsonVal(a.responseBody, strings.Split(path, ".")); ok {
		return fmt.Errorf("%s exists in response:\n%v", path, a.responseBodyOutput)
	}

	return nil
}

/*
*
Step example:

	Then the response key "data.0.created_at" should exist
*/
func (a *ApiClient) TheResponseKeyShouldExist(path string) error {
	if _, ok := getNestedJsonVal(a.responseBody, strings.Split(path, ".")); !ok {
		return fmt.Errorf("%s not exists in response:\n%v", path, a.responseBodyOutput)
	}

	return nil
}

/*
Step example:

	Then the difference between metaalarmLastEventDate createTimestamp is in range -2,2
*/
func (a *ApiClient) TheDifferenceBetweenValues(var1, var2 string, left, right float64) error {
	val1, err := a.getFloatVar(var1)
	if err != nil {
		return fmt.Errorf("first variable %s", err)
	}
	val2, err := a.getFloatVar(var2)
	if err != nil {
		return fmt.Errorf("second variable %s", err)
	}
	d := val1 - val2
	if d < left || right < d {
		return fmt.Errorf("difference is %f and out of range %f, %f", d, left, right)
	}

	return nil
}

func (a *ApiClient) getFloatVar(name string) (float64, error) {
	val, ok := a.vars[name]
	if !ok {
		return 0, fmt.Errorf("doesn't exist")
	}
	return strconv.ParseFloat(val, 64)
}

/*
*
Step example:

	Then the response key "data.0.duration" should be greater or equal than 3
*/
func (a *ApiClient) TheResponseKeyShouldBeGreaterOrEqualThan(path string, value float64) error {
	if nestedVal, ok := getNestedJsonVal(a.responseBody, strings.Split(path, ".")); ok {
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

	return fmt.Errorf("%s not exists in response:\n%v", path, a.responseBodyOutput)
}

/*
*
Step example:

	Then the response key "data.0.duration" should be greater than 3
*/
func (a *ApiClient) TheResponseKeyShouldBeGreaterThan(path string, value float64) error {
	if nestedVal, ok := getNestedJsonVal(a.responseBody, strings.Split(path, ".")); ok {
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

	return fmt.Errorf("%s not exists in response:\n%v", path, a.responseBodyOutput)
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
func (a *ApiClient) TheResponseArrayKeyShouldContain(path string, doc string) error {
	b, err := a.executeTemplate(doc)
	if err != nil {
		return err
	}

	if nestedVal, ok := getNestedJsonVal(a.responseBody, strings.Split(path, ".")); ok {
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

	return fmt.Errorf("%s not exists in response:\n%v", path, a.responseBodyOutput)
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
func (a *ApiClient) TheResponseArrayKeyShouldContainOnly(path string, doc string) error {
	b, err := a.executeTemplate(doc)
	if err != nil {
		return err
	}

	if nestedVal, ok := getNestedJsonVal(a.responseBody, strings.Split(path, ".")); ok {
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

	return fmt.Errorf("%s not exists in response:\n%v", path, a.responseBodyOutput)
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
func (a *ApiClient) IAm(ctx context.Context, role string) error {
	var line model.Rbac
	res := a.db.Collection(mongo.RightsMongoCollection).FindOne(ctx, bson.M{
		"crecord_type": model.LineTypeRole,
		"crecord_name": role,
	})
	if err := res.Err(); err != nil {
		return fmt.Errorf("cannot fetch role: %w", err)
	}

	err := res.Decode(&line)
	if err != nil {
		return fmt.Errorf("cannot decode role: %w", err)
	}

	res = a.db.Collection(mongo.RightsMongoCollection).FindOne(ctx, bson.M{
		"crecord_type": model.LineTypeSubject,
		"role":         line.ID,
	})
	if err := res.Err(); err != nil {
		return fmt.Errorf("cannot fetch user: %w", err)
	}

	err = res.Decode(&line)
	if err != nil {
		return fmt.Errorf("cannot decode user: %w", err)
	}

	uri := fmt.Sprintf("%s/api/v4/login", a.url)
	body, err := json.Marshal(map[string]string{
		"username": line.Name,
		"password": userPass,
	})
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cannot create login request: %w", err)
	}

	request.Header.Set(headerContentType, binding.MIMEJSON)

	response, err := a.client.Do(request)
	if err != nil {
		return fmt.Errorf("cannot do login request: %w", err)
	}

	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("cannot fetch login response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status %v %s", response.StatusCode, buf)
	}

	responseBody := make(map[string]string)
	err = json.Unmarshal(buf, &responseBody)
	if err != nil {
		return fmt.Errorf("cannot decode login response: %w", err)
	}

	token := responseBody["access_token"]
	if token == "" {
		return fmt.Errorf("unexpected login response %v", buf)
	}

	a.headers[headerAuthorization] = bearerPrefix + " " + token

	return nil
}

/*
*
Step example:

	When I am authenticated with username "user" password "pass"
*/
func (a *ApiClient) IAmAuthenticatedByBasicAuth(username, password string) error {
	a.headers[headerAuthorization] = basicPrefix + " " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

	return nil
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
func (a *ApiClient) ISendAnEvent(doc string) (err error) {
	uri := fmt.Sprintf("%s/api/v4/event", a.url)
	body, err := a.executeTemplate(doc)
	if err != nil {
		return err
	}

	responseStr := strings.TrimSpace(body.String())
	if responseStr == "" || responseStr[0] != '[' {
		responseStr = "[" + responseStr + "]"
	}

	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return fmt.Errorf("cannot create event request: %w", err)
	}

	req.Header.Set(headerContentType, binding.MIMEJSON)
	err = a.doRequest(req)
	if err != nil {
		return err
	}

	err = a.TheResponseCodeShouldBe(http.StatusOK)
	if err != nil {
		return err
	}

	return a.TheResponseBodyShouldContain(fmt.Sprintf("{\"sent_events\":%s}", responseStr))
}

/*
*
Step example:

	When I do GET /api/v4/alarms
	When I do GET /api/v4/entitybasic/{{ .lastResponse._id}}
*/
func (a *ApiClient) IDoRequest(method, uri string) error {
	if strings.Contains(uri, "until") {
		return fmt.Errorf("step is wrongly matched to IDoRequest")
	}

	uri, err := a.getRequestURL(uri)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	return a.doRequest(req)
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
func (a *ApiClient) IDoRequestWithBody(method, uri string, doc string) error {
	if doc == "" {
		return fmt.Errorf("body is empty")
	}
	if strings.Contains(uri, "until") {
		return fmt.Errorf("step is wrongly matched to IDoRequestWithBody")
	}

	uri, err := a.getRequestURL(uri)
	if err != nil {
		return err
	}

	body, err := a.getRequestBody(doc)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		method,
		uri,
		body,
	)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	if _, ok := a.headers[headerContentType]; !ok {
		req.Header.Set(headerContentType, binding.MIMEJSON)
	}

	return a.doRequest(req)
}

/*
*
Step example:

	When I do GET /api/v4/entitybasic/{{ .lastResponse._id}} until response code is 200
*/
func (a *ApiClient) IDoRequestUntilResponseCode(method, uri string, code int) error {
	uri, err := a.getRequestURL(uri)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	timeout := startRepeatRequestInterval
	start := time.Now()
	for {
		err := a.doRequest(req)
		if err != nil {
			return err
		}

		if code == a.response.StatusCode {
			return nil
		}

		if time.Since(start) > totalRepeatRequestInterval {
			break
		}

		time.Sleep(timeout)
		timeout *= 2
	}

	return fmt.Errorf("max retries exceeded, expected response code to be: %d, but actual is: %d\nresponse body: %v",
		code,
		a.response.StatusCode,
		a.responseBodyOutput,
	)
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
func (a *ApiClient) IDoRequestUntilResponse(method, uri string, code int, doc string) error {
	if doc == "" {
		return fmt.Errorf("body is empty")
	}
	uri, err := a.getRequestURL(uri)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	b, err := a.executeTemplate(doc)
	if err != nil {
		return err
	}
	content := b.Bytes()
	expectedBody, err := unmarshalJson(content)
	if err != nil {
		return fmt.Errorf("cannot decode expected response body: %w", err)
	}

	var resDiffErr error
	timeout := startRepeatRequestInterval
	start := time.Now()
	for {
		err := a.doRequest(req)
		if err != nil {
			return err
		}

		if code == a.response.StatusCode {
			resDiffErr = checkResponse(a.responseBody, expectedBody)
			if resDiffErr == nil {
				return nil
			}
		}

		if time.Since(start) > totalRepeatRequestInterval {
			break
		}

		time.Sleep(timeout)
		timeout *= 2
	}

	if code != a.response.StatusCode {
		return fmt.Errorf("max retries exceeded: expected response code to be: %d, but actual is: %d\nresponse body: %v",
			code,
			a.response.StatusCode,
			a.responseBodyOutput,
		)
	}

	return fmt.Errorf("max retries exceeded: %w", resDiffErr)
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
func (a *ApiClient) IDoRequestUntilResponseContains(method, uri string, code int, doc string) error {
	if doc == "" {
		return fmt.Errorf("body is empty")
	}
	uri, err := a.getRequestURL(uri)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	b, err := a.executeTemplate(doc)
	if err != nil {
		return err
	}
	content := b.Bytes()
	expectedBody, err := unmarshalJson(content)
	if err != nil {
		return fmt.Errorf("cannot decode expected response body: %w", err)
	}

	var resDiffErr error
	timeout := startRepeatRequestInterval
	start := time.Now()
	for {
		err := a.doRequest(req)
		if err != nil {
			return err
		}

		if code == a.response.StatusCode {
			partialBody := getPartialResponse(a.responseBody, expectedBody)
			resDiffErr = checkResponse(partialBody, expectedBody)

			if resDiffErr == nil {
				return nil
			}
		}

		if time.Since(start) > totalRepeatRequestInterval {
			break
		}

		time.Sleep(timeout)
		timeout *= 2
	}

	if code != a.response.StatusCode {
		return fmt.Errorf("max retries exceeded: expected response code to be: %d, but actual is: %d\nresponse body: %v",
			code,
			a.response.StatusCode,
			a.responseBodyOutput,
		)
	}

	return fmt.Errorf("max retries exceeded: %w", resDiffErr)
}

/*
*
Step example:

	When I do GET /api/v4/contextgraph/import/status/{{ .lastResponse._id}} until response code is 200 and response key "data.0.duration" is greater or equal than 3
	"""
*/
func (a *ApiClient) IDoRequestUntilResponseKeyIsGreaterOrEqualThan(method, uri string, code int, path string, value float64) error {
	uri, err := a.getRequestURL(uri)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	var resDiffErr error
	timeout := startRepeatRequestInterval
	start := time.Now()
	for {
		err := a.doRequest(req)
		if err != nil {
			return err
		}

		if code == a.response.StatusCode {
			resDiffErr = a.TheResponseKeyShouldBeGreaterOrEqualThan(path, value)

			if resDiffErr == nil {
				return nil
			}
		}

		if time.Since(start) > totalRepeatRequestInterval {
			break
		}

		time.Sleep(timeout)
		timeout *= 2
	}

	if code != a.response.StatusCode {
		return fmt.Errorf("max retries exceeded: expected response code to be: %d, but actual is: %d\nresponse body: %v",
			code,
			a.response.StatusCode,
			a.responseBodyOutput,
		)
	}

	return fmt.Errorf("max retries exceeded: %w", resDiffErr)
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
func (a *ApiClient) IDoRequestUntilResponseArrayKeyContains(method, uri string, code int, path string, doc string) error {
	uri, err := a.getRequestURL(uri)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	var resDiffErr error
	timeout := startRepeatRequestInterval
	start := time.Now()
	for {
		err := a.doRequest(req)
		if err != nil {
			return err
		}

		if code == a.response.StatusCode {
			resDiffErr = a.TheResponseArrayKeyShouldContain(path, doc)

			if resDiffErr == nil {
				return nil
			}
		}

		if time.Since(start) > totalRepeatRequestInterval {
			break
		}

		time.Sleep(timeout)
		timeout *= 2
	}

	if code != a.response.StatusCode {
		return fmt.Errorf("max retries exceeded: expected response code to be: %d, but actual is: %d\nresponse body: %v",
			code,
			a.response.StatusCode,
			a.responseBodyOutput,
		)
	}

	return fmt.Errorf("max retries exceeded: %w", resDiffErr)
}

/*
*
Step example:

	When I set header Content-Type=application/json
*/
func (a *ApiClient) ISetRequestHeader(key, value string) error {
	b, err := a.executeTemplate(value)
	if err != nil {
		return err
	}

	a.headers[key] = b.String()

	return nil
}

/*
*
Step example:

	When I save response id={{ .lastResponse._id }}
*/
func (a *ApiClient) ISaveResponse(key, value string) error {
	b, err := a.executeTemplate(value)
	if err != nil {
		return err
	}

	if a.vars == nil {
		a.vars = make(map[string]string)
	}

	a.vars[key] = b.String()

	return nil
}

// ValueShouldBeGteLteThan
// Step example:
//
//	Then "value1" > "value2"
//	Then "value1" <= "value2"
func (a *ApiClient) ValueShouldBeGteLteThan(left, op, right string) error {
	leftV, err := a.getFloatVar(left)
	if err != nil {
		return err
	}
	rightV, err := a.getFloatVar(right)
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

// doRequest adds auth credentials and makes request.
func (a *ApiClient) doRequest(req *http.Request) error {
	for k, v := range a.headers {
		req.Header.Set(k, v)
	}

	// Add session's cookies
	if a.cookies != nil {
		for _, c := range a.cookies {
			req.AddCookie(c)
		}
	}

	var err error
	a.responseBody = nil
	a.responseBodyOutput = ""
	a.response, err = a.client.Do(req)
	// Read response
	if err != nil {
		return fmt.Errorf("cannot do request: %w", err)
	}
	buf, err := ioutil.ReadAll(a.response.Body)
	if err != nil {
		return fmt.Errorf("cannot fetch response: %w", err)
	}

	// Parse response
	if len(buf) > 0 {
		a.responseBody, err = unmarshalJson(buf)
		if err == nil {
			ibuf, _ := json.MarshalIndent(a.responseBody, "", "  ")
			a.responseBodyOutput = string(ibuf)
		} else {
			a.responseBodyOutput = string(buf)
		}
	}

	// Save session
	cookies := a.response.Cookies()
	if len(cookies) > 0 {
		a.cookies = make([]*http.Cookie, 0, len(cookies))
		for _, cookie := range cookies {
			if cookie.MaxAge > 0 {
				a.cookies = append(a.cookies, cookie)
			}
		}
	}

	return nil
}

// getRequestURL applies template uri to last response data.
func (a *ApiClient) getRequestURL(uri string) (string, error) {
	b, err := a.executeTemplate(uri)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", a.url, b.String()), nil
}

// getRequestBody executes template body.
func (a *ApiClient) getRequestBody(body string) (io.Reader, error) {
	return a.executeTemplate(body)
}

// executeTemplate executes provided template with last response data and time functions.
func (a *ApiClient) executeTemplate(tpl string) (*bytes.Buffer, error) {
	t, err := template.New("tpl").
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"now": func() int64 {
				return time.Now().Unix()
			},
			"nowAdd": func(s string) (int64, error) {
				d, err := libtypes.ParseDurationWithUnit(s)
				if err != nil {
					return 0, err
				}

				return d.AddTo(libtypes.NewCpsTime()).Unix(), nil
			},
			"nowDate": func() int64 {
				y, m, d := time.Now().UTC().Date()

				return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix()
			},
			"nowDateAdd": func(s string) (int64, error) {
				d, err := libtypes.ParseDurationWithUnit(s)
				if err != nil {
					return 0, err
				}

				year, month, day := time.Now().UTC().Date()
				now := libtypes.CpsTime{Time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}

				return d.AddTo(now).Unix(), nil
			},
			"parseTime": func(s string) (int64, error) {
				t, err := time.ParseInLocation("02-01-2006 15:04", s, time.UTC)
				if err != nil {
					return 0, err
				}

				return t.Unix(), nil
			},
			"sumTime": func(args ...interface{}) (int64, error) {
				var sum int64
				for _, arg := range args {
					switch v := arg.(type) {
					case string:
						i, err := strconv.Atoi(v)
						if err != nil {
							return 0, err
						}

						sum += int64(i)
					case int:
						sum += int64(v)
					case int64:
						sum += v
					default:
						return 0, fmt.Errorf("unexpected type %T of argument %v", arg, arg)
					}
				}

				return sum, nil
			},
		}).
		Parse(tpl)
	if err != nil {
		return nil, fmt.Errorf("cannot parse template: %w", err)
	}

	data := map[string]interface{}{
		"lastResponse": a.responseBody,
		"apiURL":       a.url,
	}

	for k, v := range a.vars {
		data[k] = v
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return nil, fmt.Errorf("cannot execute template: %w", err)
	}

	return buf, nil
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
