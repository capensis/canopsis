package healthcheck

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type Handler struct {
	inProgressMx sync.Mutex
	inProgress   bool
	checker      Checker
	waitInterval time.Duration
	logger       zerolog.Logger
}

func NewHandler(checker Checker, logger zerolog.Logger) http.Handler {
	return &Handler{
		checker:      checker,
		waitInterval: 10 * time.Second,
		logger:       logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.setInProgress() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer h.unsetInProgress()

	ctx, cancel := context.WithTimeout(r.Context(), h.waitInterval)
	defer cancel()
	err := h.checker.Check(ctx)
	if err != nil {
		h.logger.Err(err).Msg("cannot process healthcheck event")
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

func (h *Handler) unsetInProgress() {
	h.inProgressMx.Lock()
	defer h.inProgressMx.Unlock()
	h.inProgress = false
}
