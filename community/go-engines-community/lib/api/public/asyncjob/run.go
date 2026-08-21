package asyncjob

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

const pingInterval = 10 * time.Second

func RunWithPing(
	ctx context.Context,
	ping func(ctx context.Context) error,
	f func(ctx context.Context) error,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		err := f(ctx)
		cancel()

		return err
	})

	g.Go(func() error {
		ticker := time.NewTicker(pingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				err := ping(ctx)
				if err != nil {
					return err
				}
			}
		}
	})

	return g.Wait()
}
