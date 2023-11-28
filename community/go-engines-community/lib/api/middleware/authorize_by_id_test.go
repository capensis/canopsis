package middleware

import (
	"net/http"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestAuthorizeByID_GivenAuthorizedUser_ShouldReturnResponse(t *testing.T) {
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
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(auth.UserKey, subj)
	})
	router.GET(
		"/obj/:id",
		AuthorizeByID(act, mockEnforcer),
		okHandler,
	)

	w := performRequest(router, "GET", "/obj/"+obj)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestAuthorizeByID_GivenNoUser_ShouldReturnUnauthorizedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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
		"/obj/:id",
		AuthorizeByID(act, mockEnforcer),
		okHandler,
	)

	w := performRequest(router, "GET", "/obj/"+obj)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestAuthorizeByID_GivenNotAuthorizedUser_ShouldForbiddenError(t *testing.T) {
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
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(auth.UserKey, subj)
	})
	router.GET(
		"/obj/:id",
		AuthorizeByID(act, mockEnforcer),
		okHandler,
	)

	w := performRequest(router, "GET", "/obj/"+obj)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}
