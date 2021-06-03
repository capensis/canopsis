// bdd contains feature context utils.
package bdd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

	"github.com/gin-gonic/gin/binding"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/cucumber/messages-go/v10"
	"go.mongodb.org/mongo-driver/bson"
)

const authHeader = "x-canopsis-authkey"
const ApiEnvURL = "API_URL"
const requestTimeout = 10 * time.Second

// ApiClient represents utility struct which implements API steps to feature context.
type ApiClient struct {
	// url is base API url.
	url *url.URL
	// client is http client to make API requests.
	client *http.Client
	// db is db client.
	db mongo.DbClient
	// authApiKey is api key which can be set using corresponding step.
	authApiKey string
	// basicAuth is username and password which can be set using corresponding step.
	basicAuth *basicAuth
	// response is http response of made API request.
	response           *http.Response
	responseBody       interface{}
	responseBodyOutput string
	// cookies is http cookies which are retrieved from API response and used in following steps.
	cookies []*http.Cookie
	// vars is used to save data between steps.
	vars map[string]string
	// request header
	contentType string
}

// basicAuth represents Basic Auth credentials.
type basicAuth struct {
	username, password string
}

// NewApiClient creates new API client.
func NewApiClient() (*ApiClient, error) {
	apiUrl, err := GetApiURL()
	if err != nil {
		return nil, err
	}

	db, err := mongo.NewClient(0, 0)
	if err != nil {
		return nil, err
	}

	var apiClient ApiClient
	apiClient.client = &http.Client{
		Timeout: requestTimeout,
	}
	apiClient.url = apiUrl
	apiClient.db = db
	apiClient.contentType = binding.MIMEJSON

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
		return nil, err
	}

	return parsed, nil
}

// ResetResponse clears all saved response data.
func (a *ApiClient) ResetResponse(_ *messages.Pickle) {
	a.response = nil
	a.responseBody = nil
	a.responseBodyOutput = ""
	a.authApiKey = ""
	a.basicAuth = nil
	a.cookies = nil
	a.vars = nil
	a.contentType = binding.MIMEJSON
}

/**
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

/**
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
func (a *ApiClient) TheResponseBodyShouldBe(doc *messages.PickleStepArgument_PickleDocString) error {
	if a.responseBody == nil {
		return fmt.Errorf("response is nil")
	}

	// Try execute template on expected body
	b, err := a.executeTemplate(doc.Content)
	if err != nil {
		return err
	}

	content := b.Bytes()
	// Try to umarshal expected body as json
	var expectedBody interface{}
	err = json.Unmarshal(content, &expectedBody)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(a.responseBody, expectedBody) {
		expectedBodyOutput, _ := json.MarshalIndent(expectedBody, "", "  ")
		return fmt.Errorf("expected response body to be:\n%v\n, but actual is:\n%v",
			string(expectedBodyOutput), a.responseBodyOutput)
	}

	return nil
}

/**
Step example:
	Then the response raw body should be:
	"""
	Test
	"""
*/
func (a *ApiClient) TheResponseRawBodyShouldBe(doc *messages.PickleStepArgument_PickleDocString) error {
	// Try execute template on expected body
	b, err := a.executeTemplate(doc.Content)
	if err != nil {
		return err
	}

	expectedBody := b.String()
	if a.responseBodyOutput != expectedBody {
		return fmt.Errorf("expected response body to be:\n%v\n, but actual is:\n%v",
			expectedBody, a.responseBodyOutput)
	}

	return nil
}

/**
If some fields are not defined in step content they are ignored.

Step example:
	Then the response body should contain:
	"""
	{
		"name": "Test name"
	}
	"""
*/
func (a *ApiClient) TheResponseBodyShouldContain(doc *messages.PickleStepArgument_PickleDocString) error {
	if a.responseBody == nil {
		return fmt.Errorf("response is nil")
	}

	// Try execute template on expected body
	b, err := a.executeTemplate(doc.Content)
	if err != nil {
		return err
	}

	content := b.Bytes()
	// Try to umarshal expected body as json
	var expectedBody interface{}
	err = json.Unmarshal(content, &expectedBody)
	if err != nil {
		return fmt.Errorf("cannot unmarshal json %v: %s", err, content)
	}

	if !partialEqual(expectedBody, a.responseBody) {
		expectedBodyOutput, _ := json.MarshalIndent(expectedBody, "", "  ")
		return fmt.Errorf("expected response body to be:\n%v\n, but actual is:\n%v",
			string(expectedBodyOutput), a.responseBodyOutput)
	}

	return nil
}

/**
Step example:
	Then the response key "data.0.created_at" should not be "0"
*/
func (a *ApiClient) TheResponseKeyShouldNotBe(path, value string) error {
	if v, ok := getNestedJsonVal(a.responseBody, strings.Split(path, ".")); ok {
		if fmt.Sprintf("%v", v) == value {
			return fmt.Errorf("%v is equal to %v", value, v)
		} else {
			return nil
		}
	}

	return fmt.Errorf("%s not exists in response:\n%v", path, a.responseBodyOutput)
}

/**
Step example:
	Then the response key "data.0.created_at" should not exist
*/
func (a *ApiClient) TheResponseKeyShouldNotExist(path string) error {
	if _, ok := getNestedJsonVal(a.responseBody, strings.Split(path, ".")); ok {
		return fmt.Errorf("%s exists in response:\n%v", path, a.responseBodyOutput)
	}

	return nil
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

/**
Step example:
	Given I am admin
*/
func (a *ApiClient) IAm(role string) error {
	var line model.Rbac
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	res := a.db.Collection(mongo.RightsMongoCollection).FindOne(ctx, bson.M{
		"crecord_type": model.LineTypeRole,
		"crecord_name": role,
	})
	if err := res.Err(); err != nil {
		return err
	}

	err := res.Decode(&line)
	if err != nil {
		return err
	}

	res = a.db.Collection(mongo.RightsMongoCollection).FindOne(ctx, bson.M{
		"crecord_type": model.LineTypeSubject,
		"role":         line.ID,
	})
	if err := res.Err(); err != nil {
		return err
	}

	err = res.Decode(&line)
	if err != nil {
		return err
	}

	a.authApiKey = line.AuthApiKey
	return nil
}

/**
Step example:
	When I am authenticated with username "user" password "pass"
*/
func (a *ApiClient) IAmAuthenticatedByBasicAuth(username, password string) error {
	a.basicAuth = &basicAuth{
		username: username,
		password: password,
	}

	return nil
}

/**
Step example:
	When I am authenticated with api key "key"
*/
func (a *ApiClient) IAmAuthenticatedByApiKey(apiKey string) error {
	b, err := a.executeTemplate(apiKey)
	if err != nil {
		return err
	}

	a.authApiKey = b.String()

	return nil
}

/**
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
func (a *ApiClient) ISendAnEvent(doc *messages.PickleStepArgument_PickleDocString) (err error) {
	uri := fmt.Sprintf("%s/api/v4/event", a.url)
	body, err := a.executeTemplate(doc.Content)
	if err != nil {
		return err
	}

	responseStr := strings.TrimSpace(body.String())
	if responseStr == "" || responseStr[0] != '[' {
		responseStr = "[" + responseStr + "]"
	}

	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", a.contentType)
	err = a.doRequest(req)
	if err != nil {
		return err
	}

	err = a.TheResponseCodeShouldBe(http.StatusOK)
	if err != nil {
		return err
	}

	return a.TheResponseBodyShouldContain(&messages.PickleStepArgument_PickleDocString{
		Content: fmt.Sprintf("{\"sent_events\":%s}", responseStr),
	})
}

/**
Step example:
	When I do GET /api/v4/alarms
	When I do GET /api/v4/entitybasic/{{ .lastResponse._id}}
*/
func (a *ApiClient) IDoRequest(method, uri string) error {
	uri, err := a.getRequestURL(uri)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return err
	}

	return a.doRequest(req)
}

/**
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
func (a *ApiClient) IDoRequestWithBody(method, uri string, doc *messages.PickleStepArgument_PickleDocString) error {
	uri, err := a.getRequestURL(uri)
	if err != nil {
		return err
	}

	body, err := a.getRequestBody(doc.Content)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		method,
		uri,
		body,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", a.contentType)
	return a.doRequest(req)
}

/**
Step example:
*/
func (a *ApiClient) SetRequestContentType(contentType string) error {
	a.contentType = contentType
	return nil
}

/**
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

// doRequest adds auth credentials and makes request.
func (a *ApiClient) doRequest(req *http.Request) error {
	// Add auth credentials
	if a.authApiKey != "" {
		req.Header.Add(authHeader, a.authApiKey)
	}
	if a.basicAuth != nil {
		req.SetBasicAuth(a.basicAuth.username, a.basicAuth.password)
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
		return err
	}
	buf, err := ioutil.ReadAll(a.response.Body)
	if err != nil {
		return err
	}

	// Parse response
	if len(buf) > 0 {
		err = json.Unmarshal(buf, &a.responseBody)
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
		Funcs(template.FuncMap{
			"now": time.Now,
			"parseTime": func(s string) (time.Time, error) {
				v, err := time.Parse("02-01-2006 15:04", s)
				if err != nil {
					return time.Time{}, err
				}

				return v, nil
			},
			"parseDuration": func(s string) (time.Duration, error) {
				v, err := time.ParseDuration(s)
				if err != nil {
					return 0, err
				}

				return v, nil
			},
			"json": func(v interface{}) (string, error) {
				b, err := json.Marshal(v)
				if err != nil {
					return "", err
				}

				return string(b), nil
			},
			"sum": func(args ...interface{}) (float64, error) {
				sum := float64(0)
				for _, arg := range args {
					switch v := arg.(type) {
					case int:
						sum += float64(v)
					case int32:
						sum += float64(v)
					case int64:
						sum += float64(v)
					case float32:
						sum += float64(v)
					case float64:
						sum += v
					case string:
						i, err := strconv.Atoi(v)
						if err != nil {
							f, err := strconv.ParseFloat(v, 64)
							if err != nil {
								return 0, err
							}

							sum += f
						}

						sum += float64(i)
					default:
						return 0, fmt.Errorf("cannot process %v", arg)
					}
				}

				return sum, nil
			},
		}).
		Parse(tpl)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return buf, nil
}

// partialEqual compares two JSON unmarshal to interface{} results.
// If there is map in their structure only values by left map keys are compared.
// Extra keys from right map are ignored.
func partialEqual(left, right interface{}) bool {
	lval := reflect.ValueOf(left)
	rval := reflect.ValueOf(right)

	if !lval.IsValid() || !rval.IsValid() {
		return lval.IsValid() == rval.IsValid()
	}
	if lval.Type() != rval.Type() {
		return false
	}

	switch lval.Kind() {
	case reflect.Array:
		if lval.Len() != rval.Len() {
			return false
		}

		for i := 0; i < lval.Len(); i++ {
			if !partialEqual(lval.Index(i).Interface(), rval.Index(i).Interface()) {
				return false
			}
		}

		return true
	case reflect.Slice:
		if lval.IsNil() != rval.IsNil() {
			return false
		}
		if lval.Len() != rval.Len() {
			return false
		}
		if lval.Pointer() == rval.Pointer() {
			return true
		}

		for i := 0; i < lval.Len(); i++ {
			if !partialEqual(lval.Index(i).Interface(), rval.Index(i).Interface()) {
				return false
			}
		}

		return true
	case reflect.Map:
		if lval.IsNil() != rval.IsNil() {
			return false
		}
		if lval.Pointer() == rval.Pointer() {
			return true
		}
		// Compare only values by left map keys.
		for _, k := range lval.MapKeys() {
			l := lval.MapIndex(k)
			r := rval.MapIndex(k)

			if !r.IsValid() || !partialEqual(l.Interface(), r.Interface()) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(left, right)
	}
}