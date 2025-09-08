package main

import (
	"context"
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
	"golang.org/x/sync/errgroup"
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

	logger := log.NewLogger(ctx, flags.Debug)

	mongoClient, err := mongo.NewClient(ctx, 0, 0, logger)
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

	m := libprometheus.NewDbCollectionsMetrics()

	reg := prometheus.NewRegistry()
	err = reg.Register(m)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to register metrics")
	}

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
		return libprometheus.RunPrometheusExporter(ctx, flags.Port, logger, m)
	})

	err = g.Wait()
	if err != nil {
		logger.Fatal().Err(err).Msg("prometheus exporter exited with error")
	}
}
