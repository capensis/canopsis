package middleware

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/auth"
	mock_security "git.canopsis.net/canopsis/go-engines/mocks/lib/security"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
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

func okHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}

func TestAuthorize_GivenAuthorizedUser_ShouldReturnResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	subj := "testsubj"
	obj := "testobj"
	act := "testact"
	expectedCode := http.StatusOK
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
		Authorize(obj, act, mockEnforcer),
		okHandler,
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestAuthorize_GivenNoUser_ShouldReturnUnauthorizedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	obj := "testobj"
	act := "testact"
	expectedCode := http.StatusUnauthorized
	mockEnforcer := mock_security.NewMockEnforcer(ctrl)
	mockEnforcer.
		EXPECT().
		Enforce(gomock.Any()).
		Times(0)
	router := gin.New()
	router.GET(
		okURL,
		Authorize(obj, act, mockEnforcer),
		okHandler,
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestAuthorize_GivenNotAuthorizedUser_ShouldForbiddenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	subj := "testsubj"
	obj := "testobj"
	act := "testact"
	expectedCode := http.StatusForbidden
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
		Authorize(obj, act, mockEnforcer),
		okHandler,
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}
