package healthcheck

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Handler struct {
	inProgressMx sync.Mutex
	inProgress   bool
	checker      Checker
	waitInterval time.Duration
}

func NewHandler(checker Checker) http.Handler {
	return &Handler{
		checker:      checker,
		waitInterval: 10 * time.Second,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.setInProgress() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.waitInterval)
	defer cancel()
	err := h.checker.Check(ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) setInProgress() bool {
	h.inProgressMx.Lock()
	defer h.inProgressMx.Unlock()
	if h.inProgress {
		return false
	}

	h.inProgress = true
	return true
}
