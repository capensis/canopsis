package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libprometheus "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics/prometheus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

const (
	shutdownTimeout   = 5 * time.Second
	readHeaderTimeout = 5 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var flags Flags
	flags.ParseArgs()

	if flags.Version {
		canopsis.PrintVersionInfo()
		return
	}

	logger := log.NewLogger(ctx, flags.Options)

	mongoClient, err := mongo.NewClient(ctx, mongo.ClientOptions{
		ReadPreference: mongo.SecondaryPreferred(),
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create mongo client")
	}

	defer func() {
		err = mongoClient.Disconnect(context.WithoutCancel(ctx))
		if err != nil {
			logger.Err(err).Msg("failed to close mongo")
		}
	}()

	runInfoClient, err := redis.NewSession(ctx, redis.EngineRunInfo, logger, 0, 0)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create redis client")
	}

	defer func() {
		err = runInfoClient.Close()
		if err != nil {
			logger.Err(err).Msg("failed to close run info redis client")
		}
	}()

	pbhClient, err := redis.NewSession(ctx, redis.PBehaviorLockStorage, logger, 0, 0)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create redis client")
	}

	defer func() {
		err = pbhClient.Close()
		if err != nil {
			logger.Err(err).Msg("failed to close pbehavior redis client")
		}
	}()

	m := libprometheus.NewMetrics()

	reg := prometheus.NewRegistry()
	err = reg.Register(m)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to register metrics")
	}

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", flags.Port),
		ReadHeaderTimeout: readHeaderTimeout,
	}
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		logger.Debug().Msg("GET /metrics request")
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	})

	updater := libprometheus.NewUpdater(
		mongoClient,
		engine.NewRunInfoManager(runInfoClient),
		config.NewHealthCheckAdapter(mongoClient),
		pbhClient,
		websocket.NewStore(mongoClient, flags.UpdateMetricsInterval),
		logger,
	)

	// init update to fill metrics.
	updater.Update(ctx, m)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(flags.UpdateMetricsInterval):
				updater.Update(ctx, m)
			}
		}
	})
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

	logger.Info().Msg("prometheus exporter started")

	err = g.Wait()
	if err != nil {
		logger.Fatal().Err(err).Msg("prometheus exporter exited with error")
	}

	logger.Info().Msg("prometheus exporter stopped")
}
