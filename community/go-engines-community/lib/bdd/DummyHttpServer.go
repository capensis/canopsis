package bdd

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func RunDummyHttpServer(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	dummyRoutes := getDummyRoutes(addr)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response, ok := dummyRoutes[r.URL.Path]
		if !ok {
			http.Error(w, fmt.Sprintf("[%s][%+v]", r.URL.Path, dummyRoutes), http.StatusNotFound)
			return
		}

		if response.Timeout > 0 {
			time.Sleep(response.Timeout)
		}

		if response.Method != r.Method {
			http.Error(w, r.Method, http.StatusNotFound)
			return
		}

		if response.Code != http.StatusOK && response.Code != http.StatusNoContent {
			http.Error(w, response.Body, response.Code)
			return
		}

		for k, v := range response.Headers {
			w.Header().Set(k, v)
		}

		w.WriteHeader(response.Code)

		if response.Body != "" {
			_, err := fmt.Fprintf(w, response.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})

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

type dummyResponse struct {
	Code    int
	Method  string
	Body    string
	Headers map[string]string
	Timeout time.Duration
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
			Timeout: 2 * time.Second,
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
			Timeout: 2 * time.Second,
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
		// External data
		"/api/external_data": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "{\"id\": 1,\"title\": \"test title\"}",
		},
		"/api/external_data_document_with_array": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "{\"array\":[{\"id\":\"1\",\"title\":\"test title 1\"},{\"id\":\"2\",\"title\":\"test title 2\"}]}",
		},
		"/api/external_data_response_is_array": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "[{\"id\":\"1\",\"title\":\"test title 1\"},{\"id\":\"2\",\"title\":\"test title 2\"}]",
		},
		"/api/external_data_response_is_nested_documents": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "{\"objects\":{\"server\":{\"code\":200,\"message\":\"test message\",\"fields\":{\"name\":\"test name\"}}},\"code\":400,\"message\":\"test message\"}",
		},
		// Webhook
		"/webhook/ticket": {
			Code:   http.StatusOK,
			Method: http.MethodPost,
			Body:   "{\"ticket_id\":\"testticket\",\"ticket_data\":\"testdata\"}",
		},
		"/webhook/document_with_array": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "{\"array\":[{\"elem1\":\"test1\",\"elem2\":\"test2\"},{\"elem1\":\"test3\",\"elem2\":\"test4\"}]}",
		},
		"/webhook/response_is_array": {
			Code:   http.StatusOK,
			Method: http.MethodGet,
			Body:   "[{\"elem1\":\"test1\",\"elem2\":\"test2\"},{\"elem1\":\"test3\",\"elem2\":\"test4\"}]",
		},
	}
}
