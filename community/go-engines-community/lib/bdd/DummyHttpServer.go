package bdd

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

func RunDummyHttpServer(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	dummyRoutes := getDummyRoutes("http://" + addr)
	mux.HandleFunc("/", dummyHandler(dummyRoutes))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			panic(err)
		}
	}()

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func dummyHandler(dummyRoutes map[string]dummyResponse) func(w http.ResponseWriter, r *http.Request) {
	dummyRoutesMx := sync.Mutex{}

	return func(w http.ResponseWriter, r *http.Request) {
		dummyRoutesMx.Lock()
		response, ok := dummyRoutes[r.URL.Path]
		if !ok {
			http.Error(w, fmt.Sprintf("[%s][%+v]", r.URL.Path, dummyRoutes), http.StatusNotFound)
			dummyRoutesMx.Unlock()

			return
		}
		dummyRoutesMx.Unlock()

		if response.Method != r.Method {
			http.Error(w, r.Method, http.StatusNotFound)
			return
		}

		if response.Timeout > 0 {
			time.Sleep(response.Timeout)
		}

		if response.Code != http.StatusOK && response.Code != http.StatusNoContent {
			http.Error(w, response.Body, response.Code)
			return
		}

		if response.Username != "" && response.Password != "" {
			header := r.Header.Get(headerAuthorization)
			if len(header) <= len(basicPrefix) || header[0:len(basicPrefix)] != basicPrefix {
				http.Error(w, r.Method, http.StatusUnauthorized)
				return
			}
			header = strings.TrimSpace(header[len(basicPrefix):])
			base, err := base64.StdEncoding.DecodeString(header)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			username, password, _ := strings.Cut(string(base), ":")
			if username != response.Username || password != response.Password {
				http.Error(w, r.Method, http.StatusUnauthorized)
				return
			}
		}

		if len(response.Headers) > 0 {
			for k, v := range response.Headers {
				w.Header().Set(k, v)
			}
		} else if response.ForwardRequestHeader {
			for k, v := range r.Header {
				if len(v) > 0 {
					w.Header().Set(k, v[0])
				}
			}
		}

		w.WriteHeader(response.Code)

		if len(response.BodySequence) != 0 {
			_, err := fmt.Fprintf(w, response.BodySequence[response.BodySequenceIndex])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			response.BodySequenceIndex++
			if response.BodySequenceIndex == len(response.BodySequence) {
				response.BodySequenceIndex = 0
			}

			dummyRoutesMx.Lock()
			dummyRoutes[r.URL.Path] = response
			dummyRoutesMx.Unlock()
		} else if response.Body != "" {
			_, err := fmt.Fprintf(w, response.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if response.ForwardRequestBody {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = w.Write(body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

type dummyResponse struct {
	Code    int
	Method  string
	Body    string
	Headers map[string]string
	Timeout time.Duration

	ForwardRequestBody   bool
	ForwardRequestHeader bool

	Username string
	Password string

	BodySequence      map[int]string
	BodySequenceIndex int
}

func getDummyRoutes(addr string) map[string]dummyResponse {
	return map[string]dummyResponse{
		// Rundeck
		"/api/35/job/test-job-succeeded/run": {
			Code:   http.StatusOK,
			Method: http.MethodPost,
			Body:   fmt.Sprintf("{\"id\":\"test-job-execution-succeeded\",\"href\":\"/api/35/execution/test-job-execution-succeeded\",\"permalink\":\"%s/rundeck/execution/show/test-job-execution-succeeded\"}", addr),
		},
		"/api/35/execution/test-job-execution-succeeded": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "{\"id\":\"test-job-execution-succeeded\",\"status\":\"succeeded\"}",
		},
		"/api/35/execution/test-job-execution-succeeded/output": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "test-job-execution-succeeded-output",
		},
		"/api/35/job/test-job-failed/run": {
			Code:   http.StatusOK,
			Method: http.MethodPost,
			Body:   fmt.Sprintf("{\"id\":\"test-job-execution-failed\",\"href\":\"/api/35/execution/test-job-execution-failed\",\"permalink\":\"%s/rundeck/execution/show/test-job-execution-failed\"}", addr),
		},
		"/api/35/execution/test-job-execution-failed": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "{\"id\":\"test-job-execution-failed\",\"status\":\"failed\"}",
		},
		"/api/35/execution/test-job-execution-failed/output": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "test-job-execution-failed-output",
		},
		"/api/35/job/test-job-aborted/run": {
			Code:   http.StatusOK,
			Method: http.MethodPost,
			Body:   fmt.Sprintf("{\"id\":\"test-job-execution-aborted\",\"href\":\"/api/35/execution/test-job-execution-aborted\",\"permalink\":\"%s/rundeck/execution/show/test-job-execution-aborted\"}", addr),
		},
		"/api/35/execution/test-job-execution-aborted": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "{\"id\":\"test-job-execution-aborted\",\"status\":\"aborted\"}",
		},
		"/api/35/execution/test-job-execution-aborted/output": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "test-job-execution-aborted-output",
		},
		"/api/35/job/test-job-http-error/run": {
			Code:   http.StatusOK,
			Method: http.MethodPost,
			Body:   fmt.Sprintf("{\"id\":\"test-job-execution-http-error\",\"href\":\"/api/35/execution/test-job-execution-http-error\",\"permalink\":\"%s/rundeck/execution/show/test-job-execution-http-error\"}", addr),
		},
		"/api/35/execution/test-job-execution-http-error": {
			Code:   http.StatusBadRequest,
			Method: http.MethodGet,
			Body:   "{\"message\":\"http-error\"}",
		},
		"/api/35/execution/test-job-execution-http-error/output": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "test-job-execution-http-error-output",
		},
		"/api/35/job/test-job-long-succeeded/run": {
			Code:   http.StatusOK,
			Method: http.MethodPost,
			Body:   fmt.Sprintf("{\"id\":\"test-job-execution-long-succeeded\",\"href\":\"/api/35/execution/test-job-execution-long-succeeded\",\"permalink\":\"%s/rundeck/execution/show/test-job-execution-long-succeeded\"}", addr),
		},
		"/api/35/execution/test-job-execution-long-succeeded": {
			Code:    http.StatusOK,
			Method:  http.MethodGet,
			Body:    "{\"id\":\"test-job-execution-long-succeeded\",\"status\":\"succeeded\"}",
			Timeout: 3 * time.Second,
		},
		"/api/35/execution/test-job-execution-long-succeeded/output": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "test-job-execution-long-succeeded-output",
		},
		"/api/35/job/test-job-long-failed/run": {
			Code:   http.StatusOK,
			Method: http.MethodPost,
			Body:   fmt.Sprintf("{\"id\":\"test-job-execution-long-failed\",\"href\":\"/api/35/execution/test-job-execution-long-failed\",\"permalink\":\"%s/rundeck/execution/show/test-job-execution-long-failed\"}", addr),
		},
		"/api/35/execution/test-job-execution-long-failed": {
			Code:    http.StatusOK,
			Method:  http.MethodGet,
			Body:    "{\"id\":\"test-job-execution-long-failed\",\"status\":\"failed\"}",
			Timeout: 3 * time.Second,
		},
		"/api/35/execution/test-job-execution-long-failed/output": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "test-job-execution-long-failed-output",
		},
		"/api/35/job/test-job-running/run": {
			Code:   http.StatusOK,
			Method: http.MethodPost,
			Body:   fmt.Sprintf("{\"id\":\"test-job-execution-running\",\"href\":\"/api/35/execution/test-job-execution-running\",\"permalink\":\"%s/rundeck/execution/show/test-job-execution-running\"}", addr),
		},
		"/api/35/execution/test-job-execution-running": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "{\"id\":\"test-job-execution-running\",\"status\":\"running\"}",
		},
		// AWX
		"/api/v2/job_templates/test-job-succeeded/launch/": {
			Code:   http.StatusOK,
			Method: http.MethodPost,
			Body:   "{\"url\":\"/api/v2/jobs/test-job-execution-succeeded\"}",
		},
		"/api/v2/jobs/test-job-execution-succeeded": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "{\"id\":\"test-job-execution-succeeded\",\"status\":\"successful\"}",
		},
		"/api/v2/jobs/test-job-execution-succeeded/stdout": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "test-job-execution-succeeded-output",
		},
		// Jenkins
		"/job/test-job-succeeded/build": {
			Code:   http.StatusNoContent,
			Method: http.MethodPost,
			Headers: map[string]string{
				"Location": "/queue/item/test-job-queue-succeeded/",
			},
		},
		"/queue/item/test-job-queue-succeeded/api/json": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "{\"executable\":{\"url\":\"/job/test-job-succeeded/test-job-execution-succeeded\"}}",
		},
		"/job/test-job-succeeded/test-job-execution-succeeded/api/json": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "{\"id\":\"test-job-execution-succeeded\",\"result\":\"SUCCESS\"}",
		},
		"/job/test-job-succeeded/test-job-execution-succeeded/consoleText": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "test-job-execution-succeeded-output",
		},
		"/job/test-job-params-succeeded/buildWithParameters": {
			Code:   http.StatusNoContent,
			Method: http.MethodPost,
			Headers: map[string]string{
				"Location": "/queue/item/test-job-queue-params-succeeded/",
			},
		},
		"/queue/item/test-job-queue-params-succeeded/api/json": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "{\"executable\":{\"url\":\"/job/test-job-params-succeeded/test-job-execution-params-succeeded\"}}",
		},
		"/job/test-job-params-succeeded/test-job-execution-params-succeeded/api/json": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "{\"id\":\"test-job-execution-params-succeeded\",\"result\":\"SUCCESS\"}",
		},
		"/job/test-job-params-succeeded/test-job-execution-params-succeeded/consoleText": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "test-job-execution-params-succeeded-output",
		},
		// VTOM
		"/vtom/public/monitoring/1.0/environments/CANOPSIS/applications/CANOPSIS_1/jobs/test-job-succeeded/action": {
			Code:   http.StatusOK,
			Method: http.MethodPost,
		},
		"/vtom/public/monitoring/1.0/environments/CANOPSIS/applications/CANOPSIS_1/jobs/test-job-succeeded/status": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			BodySequence: map[int]string{
				0: "{\"status\":\"Waiting\"}",
				1: "{\"status\":\"Finished\"}",
			},
		},
		"/vtom/public/monitoring/1.0/environments/CANOPSIS/applications/CANOPSIS_1/jobs/test-job-succeeded/logs/last/stdout": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "test-job-execution-succeeded-output",
		},
		"/vtom/public/monitoring/1.0/environments/CANOPSIS/applications/CANOPSIS_1/jobs/test-job-failed/action": {
			Code:   http.StatusOK,
			Method: http.MethodPost,
		},
		"/vtom/public/monitoring/1.0/environments/CANOPSIS/applications/CANOPSIS_1/jobs/test-job-failed/status": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			BodySequence: map[int]string{
				0: "{\"status\":\"Waiting\"}",
				1: "{\"status\":\"Error\"}",
			},
		},
		"/vtom/public/monitoring/1.0/environments/CANOPSIS/applications/CANOPSIS_1/jobs/test-job-failed/logs/last/stderr": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "test-job-execution-failed-output",
		},
		"/vtom/public/monitoring/1.0/environments/CANOPSIS/applications/CANOPSIS_2/jobs/test-job-succeeded/action": {
			Code:   http.StatusOK,
			Method: http.MethodPost,
		},
		"/vtom/public/monitoring/1.0/environments/CANOPSIS/applications/CANOPSIS_2/jobs/test-job-succeeded/status": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			BodySequence: map[int]string{
				0: "{\"status\":\"Waiting\"}",
				1: "{\"status\":\"Finished\"}",
			},
		},
		"/vtom/public/monitoring/1.0/environments/CANOPSIS/applications/CANOPSIS_2/jobs/test-job-succeeded/logs/last/stdout": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "test-job-execution-succeeded-output",
		},
		"/vtom/public/monitoring/1.0/environments/CANOPSIS/applications/CANOPSIS_2/jobs/test-job-failed/action": {
			Code:   http.StatusOK,
			Method: http.MethodPost,
		},
		"/vtom/public/monitoring/1.0/environments/CANOPSIS/applications/CANOPSIS_2/jobs/test-job-failed/status": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			BodySequence: map[int]string{
				0: "{\"status\":\"Waiting\"}",
				1: "{\"status\":\"Error\"}",
			},
		},
		"/vtom/public/monitoring/1.0/environments/CANOPSIS/applications/CANOPSIS_2/jobs/test-job-failed/logs/last/stderr": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "test-job-execution-failed-output",
		},
		// Webhook
		"/webhook/request": {
			Code:                 http.StatusOK,
			Method:               http.MethodPost,
			ForwardRequestBody:   true,
			ForwardRequestHeader: true,
		},
		"/webhook/auth-request": {
			Code:                 http.StatusOK,
			Method:               http.MethodPost,
			ForwardRequestBody:   true,
			ForwardRequestHeader: true,
			Username:             "test",
			Password:             "test",
		},
		"/webhook/long-request": {
			Code:                 http.StatusOK,
			Method:               http.MethodPost,
			ForwardRequestBody:   true,
			ForwardRequestHeader: true,
			Timeout:              2 * time.Second,
		},
		"/webhook/long-auth-request": {
			Code:                 http.StatusOK,
			Method:               http.MethodPost,
			ForwardRequestBody:   true,
			ForwardRequestHeader: true,
			Username:             "test",
			Password:             "test",
			Timeout:              2 * time.Second,
		},
	}
}
