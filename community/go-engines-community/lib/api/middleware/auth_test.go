package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	mock_httperror "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/httperror"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func TestAuth_GivenCredentials_ShouldReturnResponseAndSetUserDataToContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMaintenanceAdapter := mock_config.NewMockMaintenanceAdapter(ctrl)
	mockMaintenanceAdapter.EXPECT().GetConfig(gomock.Any()).Return(config.MaintenanceConf{}, nil).AnyTimes()

	enforcer := mock_security.NewMockEnforcer(ctrl)
	mockErrResponder := mock_httperror.NewMockResponder(ctrl)

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
		Auth([]security.HttpProvider{mockProvider}, mockMaintenanceAdapter, enforcer, mockErrResponder),
		func(c *gin.Context) {
			userID, err := authctx.GetUserKey(c)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, err)

				return
			}

			apiKey, err := authctx.GetAPIKey(c)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, err)

				return
			}

			c.String(
				expectedCode,
				"test %v %v",
				userID,
				apiKey,
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
	mockErrResponder := mock_httperror.NewMockResponder(ctrl)

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
		Auth([]security.HttpProvider{mockProvider}, mockMaintenanceAdapter, enforcer, mockErrResponder),
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
	mockErrResponder := mock_httperror.NewMockResponder(ctrl)
	mockErrResponder.EXPECT().Respond(gomock.Any(), gomock.Eq(httperror.ErrUnauthorized)).Do(func(c *gin.Context, err error) {
		c.AbortWithStatus(expectedCode)
	})
	router := gin.New()
	router.GET(
		okURL,
		Auth([]security.HttpProvider{mockProvider}, mockMaintenanceAdapter, enforcer, mockErrResponder),
		okHandler,
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}
