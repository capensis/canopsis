package engine

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"time"
)

func NewRunInfoPeriodicalWorker(
	periodicalInterval time.Duration,
	manager RunInfoManager,
	info InstanceRunInfo,
	channel amqp.Channel,
	logger zerolog.Logger,
) PeriodicalWorker {
	return &runInfoPeriodicalWorker{
		periodicalInterval: periodicalInterval,
		manager:            manager,
		info:               info,
		channel:            channel,
		logger:             logger,
	}
}

type runInfoPeriodicalWorker struct {
	periodicalInterval time.Duration
	manager            RunInfoManager
	info               InstanceRunInfo
	channel            amqp.Channel
	logger             zerolog.Logger
}

func (w *runInfoPeriodicalWorker) GetInterval() time.Duration {
	return w.periodicalInterval
}

func (w *runInfoPeriodicalWorker) Work(ctx context.Context) error {
	if w.info.ConsumeQueue != "" {
		queue, err := w.channel.QueueInspect(w.info.ConsumeQueue)
		if err != nil {
			return fmt.Errorf("cannot get consume queue length: %w", err)
		}

		w.info.QueueLength = queue.Messages
	}

	w.info.Time = types.CpsTime{Time: time.Now()}
	err := w.manager.SaveInstance(ctx, w.info, w.GetInterval())
	if err != nil {
		w.logger.Error().Err(err).Msg("cannot save run info")
	}

	return nil
}
