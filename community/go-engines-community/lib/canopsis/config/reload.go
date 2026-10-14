package config

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

// RunReload loads the config every interval and gives it to every updater, until ctx is done.
// A config which cannot be loaded is logged and skipped: the updaters keep the last one they got.
func RunReload(ctx context.Context, adapter Adapter, interval time.Duration, logger zerolog.Logger, updaters ...Updater) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cfg, err := adapter.GetConfig(ctx)
			if err != nil {
				logger.Warn().Err(err).Msg("cannot load config")

				continue
			}

			for _, updater := range updaters {
				updater.Update(cfg)
			}
		}
	}
}
