package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	mock_response "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/public/response"
	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
)

func TestSetAuthor_ShouldUpdateAuthor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedCode := http.StatusOK
	expectedAuthorValue := "test-author"
	noAuthorBody := map[string]interface{}{
		"test_key": "test_value",
	}
	noAuthorEncodedBody, _ := json.Marshal(noAuthorBody)
	req := httptest.NewRequest(http.MethodPost, okURL, bytes.NewReader(noAuthorEncodedBody))
	mockResponder := mock_response.NewMockResponder(ctrl)
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(authctx.SetUserKey(r.Context(), expectedAuthorValue))
			next.ServeHTTP(w, r)
		})
	})
	r.Use(SetAuthor(mockResponder))
	r.Post(okURL, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		encodedBody := json.NewDecoder(r.Body)
		err := encodedBody.Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))

			return
		}

		w.WriteHeader(expectedCode)
		_, _ = fmt.Fprintf(w, "author %v", body["author"])
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
