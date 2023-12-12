package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestAuth_GivenCredentials_ShouldReturnResponseAndSetUserDataToContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMaintenanceAdapter := mock_config.NewMockMaintenanceAdapter(ctrl)
	mockMaintenanceAdapter.EXPECT().GetConfig(gomock.Any()).Return(config.MaintenanceConf{}, nil).AnyTimes()

	enforcer := mock_security.NewMockEnforcer(ctrl)

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
		Auth(gomock.Eq(req)).
		Return(user, nil, true)
	router := gin.New()
	router.GET(
		okURL,
		Auth([]security.HttpProvider{mockProvider}, mockMaintenanceAdapter, enforcer),
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

	mockMaintenanceAdapter := mock_config.NewMockMaintenanceAdapter(ctrl)
	mockMaintenanceAdapter.EXPECT().GetConfig(gomock.Any()).Return(config.MaintenanceConf{}, nil).AnyTimes()

	enforcer := mock_security.NewMockEnforcer(ctrl)

	expectedCode := http.StatusOK
	req := httptest.NewRequest(http.MethodGet, okURL, nil)
	mockProvider := mock_security.NewMockHttpProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Eq(req)).
		Return(nil, nil, false)
	router := gin.New()
	router.GET(
		okURL,
		Auth([]security.HttpProvider{mockProvider}, mockMaintenanceAdapter, enforcer),
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

	mockMaintenanceAdapter := mock_config.NewMockMaintenanceAdapter(ctrl)
	mockMaintenanceAdapter.EXPECT().GetConfig(gomock.Any()).Return(config.MaintenanceConf{}, nil).AnyTimes()

	enforcer := mock_security.NewMockEnforcer(ctrl)

	expectedCode := http.StatusUnauthorized
	req := httptest.NewRequest(http.MethodGet, okURL, nil)
	mockProvider := mock_security.NewMockHttpProvider(ctrl)
	mockProvider.
		EXPECT().
		Auth(gomock.Eq(req)).
		Return(nil, nil, true)
	router := gin.New()
	router.GET(
		okURL,
		Auth([]security.HttpProvider{mockProvider}, mockMaintenanceAdapter, enforcer),
		okHandler,
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}
