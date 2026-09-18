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

// RunBackground starts background works and returns the group to wait on once the http server stopped.
// cancel stops the whole api when the job worker fails.
func RunBackground(
	ctx context.Context,
	s Services,
	jobHandlers map[string]asyncjob.Handler,
	jobMaintainers []asyncjob.Maintainer,
	cancel context.CancelFunc,
) *errgroup.Group {
	jobWorker := asyncjob.NewWorker(
		s.AMQPConsumeChPool,
		json.NewDecoder(),
		jobHandlers,
		jobWorkers,
		s.Logger,
	)

	g, wctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		err := jobWorker.Run(wctx)
		if err != nil {
			cancel()

			return fmt.Errorf("async job worker stopped with error: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		asyncjob.RunMaintenance(wctx, jobMaintainers, s.Logger)

		return nil
	})

	g.Go(func() error {
		config.RunReload(wctx, s.configAdapter, canopsis.PeriodicalWaitTime, s.Logger, s.apiConfigProvider, s.timezoneConfigProvider)

		return nil
	})

	return g
}
