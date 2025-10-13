package prometheus

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

const (
	shutdownTimeout   = 5 * time.Second
	readHeaderTimeout = 5 * time.Second

	DefaultExporterPort = 9180
)

func RunPrometheusExporter(ctx context.Context, port int, logger zerolog.Logger, metrics prometheus.Collector) error {
	reg := prometheus.NewRegistry()

	err := reg.Register(metrics)
	if err != nil {
		return fmt.Errorf("failed to register metrics: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug().Msg("GET /metrics request")
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           mux,
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server listen error: %w", err)
		}

		return nil
	})
	g.Go(func() error {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("prometheus exporter forced to shutdown: %w", err)
		}

		return nil
	})

	logger.Debug().Msg("prometheus exporter started")

	err = g.Wait()
	if err != nil {
		return fmt.Errorf("prometheus exporter exited with error: %w", err)
	}

	logger.Debug().Msg("prometheus exporter stopped")

	return nil
}
