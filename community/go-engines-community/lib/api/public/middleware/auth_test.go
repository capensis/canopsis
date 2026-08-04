package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	mock_response "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/public/response"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
)

func TestAuth_GivenCredentials_ShouldReturnResponseAndSetUserDataToContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockResponder := mock_response.NewMockResponder(ctrl)
	expectedCode := http.StatusOK
	user := &security.User{
		ID:             "testid",
		AuthApiKey:     "testkey",
		HashedPassword: "testhash",
	}
	req := httptest.NewRequest(http.MethodGet, okURL, nil)
	mockProvider := mock_security.NewMockHttpProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Any()).
		Return(user, nil, true)
	r := chi.NewRouter()
	r.Use(Auth([]security.HttpProvider{mockProvider}, mockResponder))
	r.Get(
		okURL,
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := authctx.GetUserKey(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))

				return
			}

			apiKey, err := authctx.GetAPIKey(r.Context())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))

				return
			}

			w.WriteHeader(expectedCode)
			_, _ = w.Write([]byte("test " + userID + " " + apiKey))
		},
	)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}

	expectedResponse := fmt.Sprintf("test %v %v", user.ID, user.AuthApiKey)
	if w.Body.String() != expectedResponse {
		t.Errorf("expected response: \"%v\" but got \"%v\"", expectedResponse, w.Body.String())
	}
}

func TestAuth_GivenNoCredentials_ShouldReturnResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockResponder := mock_response.NewMockResponder(ctrl)
	expectedCode := http.StatusOK
	req := httptest.NewRequest(http.MethodGet, okURL, nil)
	mockProvider := mock_security.NewMockHttpProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Any()).
		Return(nil, nil, false)
	r := chi.NewRouter()
	r.Use(Auth([]security.HttpProvider{mockProvider}, mockResponder))
	r.Get(okURL, okHandler)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestAuth_GivenInvalidCredentials_ShouldReturnUnauthorizedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedCode := http.StatusUnauthorized
	req := httptest.NewRequest(http.MethodGet, okURL, nil)
	mockProvider := mock_security.NewMockHttpProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Any()).
		Return(nil, nil, true)
	mockResponder := mock_response.NewMockResponder(ctrl)
	mockResponder.EXPECT().RespondError(gomock.Any(), gomock.Any(), gomock.Eq(httperror.ErrUnauthorized)).Do(func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(expectedCode)
	})
	r := chi.NewRouter()
	r.Use(Auth([]security.HttpProvider{mockProvider}, mockResponder))
	r.Get(okURL, okHandler)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}
