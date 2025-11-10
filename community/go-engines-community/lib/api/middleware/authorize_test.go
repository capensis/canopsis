package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	mock_httperror "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/httperror"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/gin-gonic/gin"
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

func okHandler(c *gin.Context) {
	c.Status(http.StatusOK)
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
	mockErrResponder := mock_httperror.NewMockResponder(ctrl)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		authctx.SetUserKey(c, subj)
	})
	router.GET(
		okURL,
		Authorize(obj, act, mockEnforcer, mockErrResponder),
		okHandler,
	)

	w := performRequest(router, "GET", okURL)

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
	mockErrResponder := mock_httperror.NewMockResponder(ctrl)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		authctx.SetUserKey(c, subj)
	})
	router.GET(
		okURL,
		Authorize(obj, act, mockEnforcer, mockErrResponder),
		okHandler,
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}
