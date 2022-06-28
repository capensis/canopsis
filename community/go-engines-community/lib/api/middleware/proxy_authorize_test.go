package middleware

import (
	"net/http"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	mock_proxy "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security/proxy"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestProxyAuthorize_GivenUnprotectedRoute_ShouldReturnResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedCode := http.StatusOK
	mockAccessConfig := mock_proxy.NewMockAccessConfig(ctrl)
	mockAccessConfig.
		EXPECT().
		Get(okURL, "GET").
		Return("", "")
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		Enforce(gomock.Any()).
		Times(0)
	router := gin.New()
	router.GET(
		okURL,
		ProxyAuthorize(mockEnforcer, mockAccessConfig),
		okHandler,
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestProxyAuthorize_GivenAuthorizedUser_ShouldReturnResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	subj := "testsubj"
	obj := "testobj"
	act := "testact"
	expectedCode := http.StatusOK
	mockAccessConfig := mock_proxy.NewMockAccessConfig(ctrl)
	mockAccessConfig.
		EXPECT().
		Get(okURL, "GET").
		Return(obj, act)
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		Enforce(subj, obj, act).
		Return(true, nil)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(auth.UserKey, subj)
	})
	router.GET(
		okURL,
		ProxyAuthorize(mockEnforcer, mockAccessConfig),
		okHandler,
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestProxyAuthorize_GivenNotAuthorizedUser_ShouldForbiddenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	subj := "testsubj"
	obj := "testobj"
	act := "testact"
	expectedCode := http.StatusForbidden
	mockAccessConfig := mock_proxy.NewMockAccessConfig(ctrl)
	mockAccessConfig.
		EXPECT().
		Get(okURL, "GET").
		Return(obj, act)
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		Enforce(subj, obj, act).
		Return(false, nil)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(auth.UserKey, subj)
	})
	router.GET(
		okURL,
		ProxyAuthorize(mockEnforcer, mockAccessConfig),
		okHandler,
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}
