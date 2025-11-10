package middleware

import (
	"errors"
	"net/http"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	mock_httperror "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/httperror"
	mock_security "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/security"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
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
	mockErrResponder := mock_httperror.NewMockResponder(ctrl)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		authctx.SetUserKey(c, subj)
	})
	router.GET(
		"/obj/:id",
		AuthorizeByID(act, mockEnforcer, mockErrResponder),
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
	mockErrResponder := mock_httperror.NewMockResponder(ctrl)
	mockErrResponder.EXPECT().Respond(gomock.Any(), gomock.Any()).Do(func(c *gin.Context, err error) {
		if errors.Is(err, authctx.ErrNotFound) {
			c.AbortWithStatus(expectedCode)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	})
	router := gin.New()
	router.GET(
		"/obj/:id",
		AuthorizeByID(act, mockEnforcer, mockErrResponder),
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
	mockErrResponder := mock_httperror.NewMockResponder(ctrl)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		authctx.SetUserKey(c, subj)
	})
	router.GET(
		"/obj/:id",
		AuthorizeByID(act, mockEnforcer, mockErrResponder),
		okHandler,
	)

	w := performRequest(router, "GET", "/obj/"+obj)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}
