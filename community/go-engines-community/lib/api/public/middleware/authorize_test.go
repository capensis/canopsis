package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	mock_response "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/public/response"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
)

type header struct {
	Key   string
	Value string
}

func performRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

const okURL = "/ok"

func okHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestAuthorize_GivenAuthorizedUser_ShouldReturnResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	subj := "testsubj"
	obj := "testobj"
	act := "testact"
	expectedCode := http.StatusOK
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		Enforce(subj, obj, act).
		Return(true, nil)
	mockResponder := mock_response.NewMockResponder(ctrl)
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(authctx.SetUserKey(r.Context(), subj))
			next.ServeHTTP(w, r)
		})
	})
	r.Use(Authorize(obj, act, mockEnforcer, mockResponder))
	r.Get(okURL, okHandler)
	w := performRequest(r, "GET", okURL)
	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestAuthorize_GivenNotAuthorizedUser_ShouldForbiddenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	subj := "testsubj"
	obj := "testobj"
	act := "testact"
	expectedCode := http.StatusForbidden
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		Enforce(subj, obj, act).
		Return(false, nil)
	mockResponder := mock_response.NewMockResponder(ctrl)
	mockResponder.EXPECT().RespondError(gomock.Any(), gomock.Any(), gomock.Eq(httperror.NewForbiddenError(""))).Do(func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(expectedCode)
	})
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(authctx.SetUserKey(r.Context(), subj))
			next.ServeHTTP(w, r)
		})
	})
	r.Use(Authorize(obj, act, mockEnforcer, mockResponder))
	r.Get(okURL, okHandler)
	w := performRequest(r, "GET", okURL)
	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}
