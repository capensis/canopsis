package engine

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"github.com/rs/zerolog"
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

func (w *runInfoPeriodicalWorker) Work(ctx context.Context) {
	updateInstanceRunInfo(ctx, w.GetInterval(), w.manager, w.info, w.channel, w.logger)
}

func updateInstanceRunInfo(
	ctx context.Context,
	interval time.Duration,
	manager RunInfoManager,
	info InstanceRunInfo,
	channel amqp.Channel,
	logger zerolog.Logger,
) InstanceRunInfo {
	if info.ConsumeQueue != "" {
		queue, err := channel.QueueInspect(info.ConsumeQueue)
		if err != nil {
			logger.Err(err).Msg("cannot get consume queue length")
			return InstanceRunInfo{}
		}

		info.QueueLength = queue.Messages
	}

	info.Time = datetime.NewCpsTime()
	err := manager.SaveInstance(ctx, info, 2*interval)
	if err != nil {
		logger.Err(err).Msg("cannot save run info")
	}

	return info
}
