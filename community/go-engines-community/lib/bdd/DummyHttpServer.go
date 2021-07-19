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

		if response.Code != http.StatusOK {
			http.Error(w, response.Body, response.Code)
			return
		}

		_, err := fmt.Fprintf(w, response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
	}
}
