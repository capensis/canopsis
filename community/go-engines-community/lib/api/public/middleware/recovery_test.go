package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestRecovery_GivenPanicErr_ShouldReturnInternalErrorResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedCode := http.StatusInternalServerError
	req := httptest.NewRequest(http.MethodGet, okURL, nil)
	r := chi.NewRouter()
	r.Use(Recovery(zerolog.Nop()))
	r.Get(okURL, func(w http.ResponseWriter, r *http.Request) {
		panic(errors.New("test error"))
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestRecovery_GivenPanicStr_ShouldReturnInternalErrorResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedCode := http.StatusInternalServerError
	req := httptest.NewRequest(http.MethodGet, okURL, nil)
	r := chi.NewRouter()
	r.Use(Recovery(zerolog.Nop()))
	r.Get(okURL, func(w http.ResponseWriter, r *http.Request) {
		panic("test string panic")
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}
}

func TestRecovery_GivenPanicSysCallErr_ShouldReturnInternalErrorResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	expectedCode := http.StatusOK
	req := httptest.NewRequest(http.MethodGet, okURL, nil)
	r := chi.NewRouter()
	r.Use(Recovery(zerolog.Nop()))
	r.
		With(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err := os.NewSyscallError("broken pipe", errors.New("test error"))
				panic(err)
			})
		}).
		Get(okURL, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte("test"))
		})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// Check if no status header is written
	if w.Code != expectedCode {
		t.Errorf("expected code: %v but got %v", expectedCode, w.Code)
	}

	// Check if no response body is written
	if w.Body.String() != "" {
		t.Errorf("expected response to be empty but got \"%v\"", w.Body.String())
	}
}
