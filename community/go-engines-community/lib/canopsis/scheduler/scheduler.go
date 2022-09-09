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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	redismod "github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
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
	redisConn      redismod.UniversalClient
	channelPub     libamqp.Channel
	publishToQueue string

	decoder                   encoding.Decoder
	encoder                   encoding.Encoder
	enableMetaAlarmProcessing bool

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
	enableMetaAlarmProcessing bool,
) Scheduler {

	s := scheduler{
		redisConn:      redisLockStorage,
		channelPub:     channelPub,
		publishToQueue: publishToQueue,
		logger:         logger,

		decoder: decoder,
		encoder: encoder,

		enableMetaAlarmProcessing: enableMetaAlarmProcessing,

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

	if s.IsMetaAlarm(event) && s.enableMetaAlarmProcessing {
		return s.ProcessMetaAlarm(ctx, lockID, bevent)
	}

	locked, err := s.queueLock.LockOrPush(ctx, lockID, bevent)

	if !locked || err != nil {
		return err
	}

	return s.publishToNext(bevent)
}

func (s *scheduler) IsMetaAlarm(event types.Event) bool {
	return event.EventType != types.EventTypeMetaAlarm &&
		event.EventType != types.EventTypeMetaAlarmUpdated &&
		strings.HasPrefix(event.Resource, "meta-alarm-entity-")
}

func (s *scheduler) ProcessMetaAlarm(ctx context.Context, lockID string, bevent []byte) error {
	s.logger.Debug().Msg("Processing meta-alarm event")

	var event types.Event
	err := s.decoder.Decode(bevent, &event)
	if err != nil {
		return err
	}

	var lockIDList []string
	if event.MetaAlarmChildren != nil {
		lockIDList = append(lockIDList, *event.MetaAlarmChildren...)
	}

	lockIDList = utils.Unique(lockIDList)
	locked, err := s.queueLock.LockMultipleOrPush(ctx, lockIDList, lockID, bevent)

	if err != nil {
		s.logger.Err(err).Msg("error on setting meta-alarm lock")
	}

	if !locked {
		s.logger.Debug().Msg("meta-alarm is locked, pushing to queue")

		return nil
	}

	s.logger.Debug().Msg("sending meta-alarm event")

	return s.publishToNext(bevent)
}

func (s *scheduler) AckEvent(ctx context.Context, event types.Event) error {
	lockID := event.GetLockID()
	s.logger.Debug().Str("lockID", lockID).Msg("AckEvent")

	if s.IsMetaAlarm(event) && s.enableMetaAlarmProcessing {
		return s.processMetaAlarmUnlock(ctx, event)
	}

	asyncUnlock := !s.enableMetaAlarmProcessing || event.MetaAlarmParents == nil ||
		len(*event.MetaAlarmParents) == 0
	nextEvent, err := s.queueLock.PopOrUnlock(ctx, lockID, asyncUnlock)

	if err != nil {
		return err
	}

	if nextEvent == nil {
		go func() {
			if s.enableMetaAlarmProcessing {
				s.processMetaAlarmChildUnlock(ctx, event)
			}
		}()

		return nil
	}

	return s.publishToNext(nextEvent)
}

func (s *scheduler) processMetaAlarmUnlock(ctx context.Context, event types.Event) error {
	lockID := event.GetEID()

	if event.MetaAlarmChildren != nil && len(*event.MetaAlarmChildren) > 0 {
		lockIDList := *event.MetaAlarmChildren
		lockIDList = utils.Unique(lockIDList)
		nextEvents, err := s.queueLock.ExtendAndPopRelatedOrMultiple(ctx, lockIDList, lockID)
		if err != nil {
			s.logger.Err(err).
				Str("lockID", lockID).
				Msg("unable to process meta alarm")
		}

		for _, nextEvent := range nextEvents {
			err = s.publishToNext(nextEvent)
			if err != nil {
				return err
			}
		}

		if len(nextEvents) > 0 {
			return nil
		}
	}

	for _, relatedParentId := range event.MetaAlarmRelatedParents {
		nextMetaAlarmEvent, err := s.queueLock.ExtendAndPopMultiple(ctx, relatedParentId, getChildren)
		if err != nil {
			s.logger.Err(err).
				Str(relatedParentId, "lockID").
				Msg("unable to unlock alarm")

			continue
		}

		if nextMetaAlarmEvent != nil {
			return s.publishToNext(nextMetaAlarmEvent)
		}
	}

	return nil
}

func (s *scheduler) processMetaAlarmChildUnlock(ctx context.Context, event types.Event) {
	if metaAlarmParents := event.MetaAlarmParents; metaAlarmParents != nil && len(*metaAlarmParents) > 0 {
		for _, metaAlarmLock := range *metaAlarmParents {
			nextEvent, err := s.queueLock.ExtendAndPopMultiple(ctx, metaAlarmLock, getChildren)
			if err != nil {
				s.logger.Err(err).
					Str(metaAlarmLock, "meta-alarm-lockID").
					Msg("unable to get meta-alarms")
			}

			if nextEvent == nil {
				continue
			}

			err = s.publishToNext(nextEvent)
			if err != nil {
				s.logger.Err(err).
					Str(metaAlarmLock, "meta-alarm-lockID").
					Msg("unable to process meta-alarm")
			}

			break
		}
	}
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

	children = utils.Unique(children)
	return children, nil
}
