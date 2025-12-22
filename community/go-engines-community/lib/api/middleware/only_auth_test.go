package middleware

import (
	"net/http"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	mock_httperror "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/api/httperror"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func TestOnlyAuth_GivenAuthUser_ShouldReturnResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	subj := "testsubj"
	expectedCode := http.StatusOK
	mockErrResponder := mock_httperror.NewMockResponder(ctrl)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		authctx.SetUserKey(c, subj)
	})
	router.GET(
		okURL,
		OnlyAuth(mockErrResponder),
		okHandler,
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}
