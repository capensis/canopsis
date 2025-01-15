package scheduler

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	amqp "github.com/rabbitmq/amqp091-go"
	redismod "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
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
	redisConn            redismod.UniversalClient
	channelPub           libamqp.Channel
	publishToQueuePrefix string

	decoder encoding.Decoder
	encoder encoding.Encoder

	queueLock QueueLock

	logger zerolog.Logger

	pubsubMx sync.Mutex
	pubsub   *redismod.PubSub
}

// NewSchedulerService ...
func NewSchedulerService(
	redisLockStorage redismod.UniversalClient,
	redisQueueStorage redismod.UniversalClient,
	channelPub libamqp.Channel,
	publishToQueuePrefix string,
	logger zerolog.Logger,
	lockTtl int,
	decoder encoding.Decoder,
	encoder encoding.Encoder,
) Scheduler {
	return &scheduler{
		redisConn:            redisLockStorage,
		channelPub:           channelPub,
		publishToQueuePrefix: publishToQueuePrefix,
		logger:               logger,

		decoder: decoder,
		encoder: encoder,

		queueLock: NewQueueLock(
			redisLockStorage,
			time.Second*time.Duration(lockTtl),
			redisQueueStorage,
			logger,
		),
	}
}

func (s *scheduler) Start(ctx context.Context) {
	s.pubsubMx.Lock()
	defer s.pubsubMx.Unlock()
	if s.pubsub != nil {
		panic("scheduler already started")
	}

	s.redisConn.ConfigSet(ctx, "notify-keyspace-events", "KEx")
	s.pubsub = s.redisConn.PSubscribe(ctx, redisSubscriptionPattern)
	_, err := s.pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	go s.listen(ctx, s.pubsub.Channel())
	s.logger.Debug().Msg("subscribed")
}

func (s *scheduler) Stop(ctx context.Context) {
	s.pubsubMx.Lock()
	defer s.pubsubMx.Unlock()
	if s.pubsub == nil {
		return
	}

	err := s.pubsub.PUnsubscribe(ctx, redisSubscriptionPattern)
	if err != nil {
		s.logger.Error().Err(err).Msg("unsubscribe pubsub")
	}

	err = s.pubsub.Close()
	if err != nil && !errors.Is(err, redismod.ErrClosed) {
		s.logger.Error().Err(err).Msg("close pubsub")
	}

	s.pubsub = nil
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

	return s.publishToNext(ctx, bevent, event.Initiator)
}

func (s *scheduler) AckEvent(ctx context.Context, event types.Event) error {
	lockID := event.GetLockID()
	s.logger.Debug().Str("lockID", lockID).Msg("AckEvent")

	nextEvent, err := s.queueLock.PopOrUnlock(ctx, lockID, true)
	if err != nil {
		return fmt.Errorf("failed to ack event: %w", err)
	}

	if nextEvent == nil {
		return nil
	}

	initiator, err := s.getInitiator(nextEvent)
	if err != nil {
		return fmt.Errorf("failed to ack event: %w", err)
	}

	return s.publishToNext(ctx, nextEvent, initiator)
}

func (s *scheduler) publishToNext(ctx context.Context, eventByte []byte, initiator string) error {
	return s.channelPub.PublishWithContext(
		ctx,
		canopsis.EngineExchangeName,
		libamqp.BuildRoutingKey(s.publishToQueuePrefix, initiator),
		false,
		false,
		amqp.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         eventByte,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (s *scheduler) listen(ctx context.Context, ch <-chan *redismod.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
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

	initiator, err := s.getInitiator(nextEvent)
	if err != nil {
		s.logger.
			Err(err).
			Str("lockID", lockID).
			Msg("error on getting initiator from the event")
	}

	err = s.publishToNext(ctx, nextEvent, initiator)
	if err != nil {
		s.logger.
			Err(err).
			Str("lockID", lockID).
			Msg("error on publishing event to queue")
	}
}

func (s *scheduler) getInitiator(event []byte) (string, error) {
	msg, err := fastjson.ParseBytes(event)
	if err != nil {
		return "", fmt.Errorf("failed to get initiator: %w", err)
	}

	initiator := string(msg.GetStringBytes(types.InitiatorJSONTag))
	if !types.IsValidInitiator(initiator) {
		initiator = types.InitiatorExternal
	}

	return initiator, nil
}
