package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	mock_response "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/public/response"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestPreProcessBulk_ShouldUpdateAuthorToAllItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedCode := http.StatusOK
	author := "test-author"
	expectedAuthorValue := "test-author test-author test-author"
	noAuthorBody := []map[string]interface{}{
		{
			"test_key-1": "test_value-1",
		},
		{
			"test_key-2": "test_value-2",
		},
		{
			"test_key-3": "test_value-3",
		},
	}
	noAuthorEncodedBody, _ := json.Marshal(noAuthorBody)
	req := httptest.NewRequest(http.MethodPost, okURL, bytes.NewReader(noAuthorEncodedBody))
	mockResponder := mock_response.NewMockResponder(ctrl)
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(authctx.SetUserKey(r.Context(), author))
			next.ServeHTTP(w, r)
		})
	})
	r.Use(PreProcessBulk(config.NewApiConfigProvider(config.CanopsisConf{API: config.SectionApi{BulkMaxSize: 100}}, zerolog.Nop()), mockResponder, true))
	r.Post(okURL, func(w http.ResponseWriter, r *http.Request) {
		var body []map[string]interface{}
		encodedBody := json.NewDecoder(r.Body)
		err := encodedBody.Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))

			return
		}

		var authorValues []string
		for _, item := range body {
			authorValues = append(authorValues, item["author"].(string))
		}

		w.WriteHeader(expectedCode)
		_, _ = w.Write([]byte("author " + strings.Join(authorValues, " ")))
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}

	expectedResponse := fmt.Sprintf("author %v", expectedAuthorValue)
	if w.Body.String() != expectedResponse {
		t.Errorf("expected response: \"%v\" but got \"%v\"", expectedResponse, w.Body.String())
	}
}

func TestPreProcessBulk_ShouldCheckBulkSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockResponder := mock_response.NewMockResponder(ctrl)
	mockResponder.EXPECT().RespondError(gomock.Any(), gomock.Any(), gomock.Eq(httperror.ErrRequestEntityTooLarge)).Do(func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
	})
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(authctx.SetUserKey(r.Context(), "test-author"))
			next.ServeHTTP(w, r)
		})
	})
	r.Use(PreProcessBulk(config.NewApiConfigProvider(config.CanopsisConf{API: config.SectionApi{BulkMaxSize: 3}}, zerolog.Nop()), mockResponder, true))
	r.Post(okURL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	valid := []map[string]interface{}{
		{
			"test_key-1": "test_value-1",
		},
		{
			"test_key-2": "test_value-2",
		},
		{
			"test_key-3": "test_value-3",
		},
	}
	body, _ := json.Marshal(valid)
	req := httptest.NewRequest(http.MethodPost, okURL, bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected code: %v but got %v", http.StatusOK, w.Code)
	}

	invalid := []map[string]interface{}{
		{
			"test_key-1": "test_value-1",
		},
		{
			"test_key-2": "test_value-2",
		},
		{
			"test_key-3": "test_value-3",
		},
		{
			"test_key-4": "test_value-4",
		},
	}
	body, _ = json.Marshal(invalid)
	req = httptest.NewRequest(http.MethodPost, okURL, bytes.NewReader(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected code: %v but got %v", http.StatusRequestEntityTooLarge, w.Code)
	}
}
