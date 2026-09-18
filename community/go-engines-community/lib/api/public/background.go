package api

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/asyncjob"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
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
		config.RunReload(ctx, s.configAdapter, canopsis.PeriodicalWaitTime, s.Logger, s.apiConfigProvider, s.timezoneConfigProvider)

		return nil
	})
}
