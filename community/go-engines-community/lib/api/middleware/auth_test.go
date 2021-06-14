package middleware

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth_GivenCredentials_ShouldReturnResponseAndSetUserDataToContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedCode := http.StatusOK
	user := &security.User{
		ID:             "testid",
		AuthApiKey:     "testkey",
		HashedPassword: "testhash",
	}
	req := httptest.NewRequest("GET", okURL, nil)
	mockProvider := mock_security.NewMockHttpProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Eq(req)).
		Return(user, nil, true)
	router := gin.New()
	router.GET(
		okURL,
		Auth([]security.HttpProvider{mockProvider}),
		func(c *gin.Context) {
			c.String(
				expectedCode,
				"test %v %v",
				c.MustGet(auth.UserKey).(string),
				c.MustGet(auth.ApiKey).(string),
			)
		},
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

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
	expectedCode := http.StatusOK
	req := httptest.NewRequest("GET", okURL, nil)
	mockProvider := mock_security.NewMockHttpProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Eq(req)).
		Return(nil, nil, false)
	router := gin.New()
	router.GET(
		okURL,
		Auth([]security.HttpProvider{mockProvider}),
		okHandler,
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestAuth_GivenInvalidCredentials_ShouldReturnUnauthorizedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedCode := http.StatusUnauthorized
	req := httptest.NewRequest("GET", okURL, nil)
	mockProvider := mock_security.NewMockHttpProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Eq(req)).
		Return(nil, nil, true)
	router := gin.New()
	router.GET(
		okURL,
		Auth([]security.HttpProvider{mockProvider}),
		okHandler,
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}
