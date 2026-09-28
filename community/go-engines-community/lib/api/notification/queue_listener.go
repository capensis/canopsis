package notification

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type QueueListener interface {
	Listen(ctx context.Context) error
}

func NewQueueListener(
	dbClient mongo.DbClient,
	amqpChannel amqp.ChannelPool,
	prefetchCount, prefetchSize int,
	websocketHub websocket.Hub,
	store Store,
	decoder encoding.Decoder,
	configProvider config.ApiConfigProvider,
	logger zerolog.Logger,
) QueueListener {
	return &queueListener{
		amqpChannelPool: amqpChannel,
		prefetchCount:   prefetchCount,
		prefetchSize:    prefetchSize,
		websocketHub:    websocketHub,
		store:           store,
		userCollection:  dbClient.Collection(mongo.UserCollection),
		workers:         10,
		decoder:         decoder,
		configProvider:  configProvider,
		logger:          logger,
	}
}

type queueListener struct {
	amqpChannelPool amqp.ChannelPool
	prefetchCount   int
	prefetchSize    int
	websocketHub    websocket.Hub
	store           Store
	userCollection  mongo.DbCollection
	workers         int
	decoder         encoding.Decoder
	configProvider  config.ApiConfigProvider
	logger          zerolog.Logger
}

func (s *queueListener) Listen(ctx context.Context) error {
	ch, err := s.amqpChannelPool.Get(ctx)
	if err != nil {
		return err
	}

	defer s.amqpChannelPool.Put(ch)

	opts := amqp.SubscribeOptions{
		Exchange:       canopsis.ApiNotificationExchangeName,
		QueueExclusive: true,
		PrefetchCount:  s.prefetchCount,
		PrefetchSize:   s.prefetchSize,
	}

	return amqp.SubscribeWithReconnect(ctx, ch, opts, s.workers, s.processMsg, s.logger)
}

func (s *queueListener) processMsg(ctx context.Context, msg amqp091.Delivery) (amqp.AckAction, error) {
	event := rpc.ApiNotificationEvent{}
	err := s.decoder.Decode(msg.Body, &event)
	if err != nil {
		s.logger.Err(err).Msg("failed to decode remediation result event")

		return amqp.Ack, nil
	}

	err = s.processEvent(ctx, event)
	if err != nil {
		s.logger.Err(err).Msg("failed to process notification event")
		if mongo.IsConnectionError(err) {
			return amqp.Nack, nil
		}
	}

	return amqp.Ack, nil
}

func (s *queueListener) processEvent(ctx context.Context, event rpc.ApiNotificationEvent) error {
	connectedUserIDs := s.websocketHub.ConnectedUserIDs()
	if len(connectedUserIDs) == 0 || len(event.Users) == 0 && len(event.Roles) == 0 {
		return nil
	}

	orCond := make([]bson.M, 0)
	if len(event.Users) > 0 {
		orCond = append(orCond, bson.M{"_id": bson.M{"$in": event.Users}})
	}

	if len(event.Roles) > 0 {
		orCond = append(orCond, bson.M{"roles": bson.M{"$in": event.Roles}})
	}

	cursor, err := s.userCollection.Find(ctx, bson.M{
		"_id": bson.M{"$in": connectedUserIDs},
		"$or": orCond,
	}, options.Find().SetProjection(bson.M{"roles": 1}))
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	r := pagination.GetDefaultQuery()
	r.Limit = int64(s.configProvider.Get().NotificationDisplayCount)
	for cursor.Next(ctx) {
		u := struct {
			ID    string   `bson:"_id"`
			Roles []string `bson:"roles"`
		}{}
		err := cursor.Decode(&u)
		if err != nil {
			return err
		}

		notifs, err := s.store.Find(ctx, r, u.ID, u.Roles)
		if err != nil {
			return err
		}

		s.websocketHub.SendMessage(ctx, map[string]any{
			"data":        notifs.Data,
			"total_count": notifs.TotalCount,
		}, websocket.ToUser(websocket.RoomNotifications, u.ID))
	}

	if err = cursor.Err(); err != nil {
		return err
	}

	return nil
}
