package scheduler

import (
	"context"
	"fmt"
	"strings"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	amqp "github.com/rabbitmq/amqp091-go"
	redismod "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var (
	redisSubscriptionPattern = fmt.Sprintf("__key*@%d__:*", redis.LockStorage)
)

// Scheduler ...
type Scheduler interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
	ProcessEvent(context.Context, types.Event) error
	AckEvent(context.Context, types.Event) error
}

type scheduler struct {
	redisConn      redismod.UniversalClient
	channelPub     libamqp.Channel
	publishToQueue string

	decoder encoding.Decoder
	encoder encoding.Encoder

	queueLock QueueLock

	logger zerolog.Logger

	ch     <-chan *redismod.Message
	pubsub *redismod.PubSub
}

// NewSchedulerService ...
func NewSchedulerService(
	redisLockStorage redismod.UniversalClient,
	redisQueueStorage redismod.UniversalClient,
	channelPub libamqp.Channel,
	publishToQueue string,
	logger zerolog.Logger,
	lockTtl int,
	decoder encoding.Decoder,
	encoder encoding.Encoder,
) Scheduler {
	s := scheduler{
		redisConn:      redisLockStorage,
		channelPub:     channelPub,
		publishToQueue: publishToQueue,
		logger:         logger,

		decoder: decoder,
		encoder: encoder,

		queueLock: NewQueueLock(
			redisLockStorage,
			time.Second*time.Duration(lockTtl),
			redisQueueStorage,
			logger,
		),
	}

	return &s
}

func (s *scheduler) subscribe(ctx context.Context) {
	s.redisConn.ConfigSet(ctx, "notify-keyspace-events", "KEx")
	s.pubsub = s.redisConn.PSubscribe(ctx, redisSubscriptionPattern)

	_, err := s.pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}
	s.ch = s.pubsub.Channel()
	go s.listen(ctx)
	s.logger.Debug().Msg("subscribed")
}

func (s *scheduler) Start(ctx context.Context) {
	s.subscribe(ctx)
}

func (s *scheduler) Stop(ctx context.Context) {
	if s.pubsub == nil {
		return
	}
	err := s.pubsub.PUnsubscribe(ctx, redisSubscriptionPattern)
	if err != nil {
		s.logger.Error().Err(err).Msg("unsubscribe pubsub")
	}
	err = s.pubsub.Close()
	if err != nil {
		s.logger.Error().Err(err).Msg("close pubsub")
	}
}

func (s *scheduler) ProcessEvent(ctx context.Context, event types.Event) error {
	lockID := event.GetLockID()
	bevent, err := s.encoder.Encode(event)
	if err != nil {
		return err
	}

	locked, err := s.queueLock.LockOrPush(ctx, lockID, bevent)
	if err != nil {
		return err
	}

	if event.Healthcheck {
		_, err := s.queueLock.PopOrUnlock(ctx, lockID, false)
		return err
	}

	if !locked {
		return nil
	}

	return s.publishToNext(ctx, bevent)
}

func (s *scheduler) AckEvent(ctx context.Context, event types.Event) error {
	lockID := event.GetLockID()
	s.logger.Debug().Str("lockID", lockID).Msg("AckEvent")

	nextEvent, err := s.queueLock.PopOrUnlock(ctx, lockID, true)

	if err != nil {
		return err
	}

	if nextEvent == nil {
		return nil
	}

	return s.publishToNext(ctx, nextEvent)
}

func (s *scheduler) publishToNext(ctx context.Context, eventByte []byte) error {
	return s.channelPub.PublishWithContext(
		ctx,
		"",
		s.publishToQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         eventByte,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (s *scheduler) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-s.ch:
			if !ok {
				return
			}

			if msg.Payload == "expired" {
				parsedStr := strings.SplitN(msg.Channel, ":", 2)
				lockID := parsedStr[1]

				s.logger.
					Info().
					Str("lockID", lockID).
					Msg("alarm lock has been expired, processing next event")

				s.processExpiredLock(ctx, lockID)
			}
		}
	}
}

func (s *scheduler) processExpiredLock(ctx context.Context, lockID string) {
	s.logger.Debug().Str("lockID", lockID).Msg("processExpireLock")
	nextEvent, err := s.queueLock.LockAndPop(ctx, lockID, true)
	if err != nil {
		s.logger.
			Err(err).
			Str("lockID", lockID).
			Msg("error on popping event from queue")
		return
	}
	if nextEvent == nil {
		return
	}

	err = s.publishToNext(ctx, nextEvent)
	if err != nil {
		s.logger.
			Err(err).
			Str("lockID", lockID).
			Msg("error on publishsing event to queue")
	}
}
