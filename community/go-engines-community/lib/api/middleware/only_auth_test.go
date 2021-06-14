package middleware

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
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

func TestOnlyAuth_GivenNoUserUser_ShouldReturnUnauthorizedError(t *testing.T) {
	expectedCode := http.StatusUnauthorized
	router := gin.New()
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
