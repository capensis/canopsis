package api

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"golang.org/x/sync/errgroup"
)

const jobWorkers = 5

// AddBackground starts background works.
// cancel stops the whole api when the job worker fails.
func AddBackground(
	ctx context.Context,
	g *errgroup.Group,
	s Services,
	opts Options,
	jobHandlers map[string]asyncjob.Handler,
	jobMaintainers []asyncjob.Maintainer,
	cancel context.CancelFunc,
) {
	jobWorker := asyncjob.NewWorker(
		s.AMQPConsumeChPool,
		json.NewDecoder(),
		jobHandlers,
		jobWorkers,
		s.Logger,
	)

	g.Go(func() error {
		err := jobWorker.Run(ctx)
		if err != nil {
			cancel()

			return fmt.Errorf("async job worker stopped with error: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		asyncjob.RunMaintenance(ctx, jobMaintainers, s.Logger)

		return nil
	})

	g.Go(func() error {
		config.RunReload(ctx, s.configAdapter, opts.PeriodicalWaitTime, s.Logger,
			s.apiConfigProvider, s.timezoneConfigProvider, s.templateConfigProvider)

		return nil
	})

	g.Go(func() error {
		ticker := time.NewTicker(opts.PeriodicalWaitTime)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				err := s.LinkGenerator.Load(ctx)
				if err != nil {
					s.Logger.Err(err).Msg("error loading link rules")
				}
			}
		}
	})
}
