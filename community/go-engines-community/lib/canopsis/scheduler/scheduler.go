package scheduler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	redismod "github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
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

	if !locked || err != nil {
		return err
	}

	return s.publishToNext(bevent)
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

	return s.publishToNext(nextEvent)
}

func (s *scheduler) publishToNext(eventByte []byte) error {
	return s.channelPub.Publish(
		"",
		s.publishToQueue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json", // this type is mandatory to avoid bad conversions into Python.
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

	err = s.publishToNext(nextEvent)
	if err != nil {
		s.logger.
			Err(err).
			Str("lockID", lockID).
			Msg("error on publishsing event to queue")
	}
}

func getChildren(b []byte) ([]string, error) {
	jsonEvent, err := fastjson.ParseBytes(b)
	if err != nil {
		return nil, err
	}

	jsonChildren := jsonEvent.GetArray("ma_children")
	children := make([]string, len(jsonChildren))
	for idx, child := range jsonChildren {
		if child == nil {
			continue
		}

		children[idx], err = strconv.Unquote(child.String())
		if err != nil {
			return nil, err
		}
	}

	return children, nil
}
