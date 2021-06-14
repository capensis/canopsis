package middleware

import (
	"net/http"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"github.com/gin-gonic/gin"
)

func TestOnlyAuth_GivenAuthUser_ShouldReturnResponse(t *testing.T) {
	subj := "testsubj"
	expectedCode := http.StatusOK
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(auth.UserKey, subj)
	})
	router.GET(
		okURL,
		OnlyAuth(),
		okHandler,
	)

	w := performRequest(router, "GET", okURL)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}
