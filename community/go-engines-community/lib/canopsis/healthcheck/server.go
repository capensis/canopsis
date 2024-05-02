package healthcheck

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

const (
	envPORT                = "CPS_HEALTHCHECK_PORT"
	serverShutdownInterval = 5 * time.Second
	readHeaderTimeout      = 5 * time.Second
)

func Start(
	ctx context.Context,
	checker Checker,
	logger zerolog.Logger,
) {
	port := getHttpPort()
	if port <= 0 {
		return
	}

	go func() {
		err := runHttpServer(ctx, checker, port, logger)
		if err != nil {
			logger.Err(err).Msg("http server failed")
		}
	}()
}

func getHttpPort() int {
	v := os.Getenv(envPORT)
	if v == "" {
		return 0
	}

	port, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}

	return port
}

func runHttpServer(
	ctx context.Context,
	checker Checker,
	port int,
	logger zerolog.Logger,
) error {
	mux := http.NewServeMux()
	mux.Handle("GET /", NewHandler(checker, logger))

	server := &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), serverShutdownInterval)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Err(err).Msg("cannot shutdown http server")
		}
	}()

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
