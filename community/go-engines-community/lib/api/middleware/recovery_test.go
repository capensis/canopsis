package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestRecovery_GivenNoUserKey_ShouldReturnUnauthorizedResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedCode := http.StatusUnauthorized

	req := httptest.NewRequest(http.MethodGet, okURL, nil)

	router := gin.New()
	router.GET(
		okURL,
		Recovery(zerolog.Nop()),
		func(c *gin.Context) {
			c.MustGet(auth.UserKey)
		},
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestRecovery_GivenPanicErr_ShouldReturnInternalErrorResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedCode := http.StatusInternalServerError

	req := httptest.NewRequest(http.MethodGet, okURL, nil)

	router := gin.New()
	router.GET(
		okURL,
		Recovery(zerolog.Nop()),
		func(c *gin.Context) {
			panic(errors.New("test error"))
		},
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestRecovery_GivenPanicStr_ShouldReturnInternalErrorResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedCode := http.StatusInternalServerError

	req := httptest.NewRequest(http.MethodGet, okURL, nil)

	router := gin.New()
	router.GET(
		okURL,
		Recovery(zerolog.Nop()),
		func(c *gin.Context) {
			panic("test string panic")
		},
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestRecovery_GivenPanicSysCallErr_ShouldReturnInternalErrorResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedContextAbort := true
	expectedCode := http.StatusOK

	req := httptest.NewRequest(http.MethodGet, okURL, nil)

	router := gin.New()
	router.GET(
		okURL,
		func(c *gin.Context) {
			defer func() {
				if c.IsAborted() != expectedContextAbort {
					t.Errorf("expected context abort: %v but got %v\n", expectedContextAbort, c.IsAborted())
				}
			}()
			c.Next()
		},
		Recovery(zerolog.Nop()),
		func(c *gin.Context) {
			err := os.NewSyscallError("broken pipe", errors.New("test error"))
			panic(err)
		},
		func(c *gin.Context) {
			c.JSON(http.StatusCreated, "test")
		},
	)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check if no status header is written
	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}

	// Check if no response body is written
	if w.Body.String() != "" {
		t.Errorf("expected response to be empty but got \"%v\"", w.Body.String())
	}
}
