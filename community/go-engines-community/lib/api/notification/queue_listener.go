package notification

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/sync/errgroup"
)

type QueueListener interface {
	Listen(ctx context.Context) error
}

func NewQueueListener(
	dbClient mongo.DbClient,
	amqpConn amqp.Connection,
	websocketHub websocket.Hub,
	store Store,
	decoder encoding.Decoder,
	configProvider config.ApiConfigProvider,
	logger zerolog.Logger,
) QueueListener {
	return &queueListener{
		amqpConn:       amqpConn,
		websocketHub:   websocketHub,
		store:          store,
		userCollection: dbClient.Collection(mongo.UserCollection),
		workers:        10,
		decoder:        decoder,
		configProvider: configProvider,
		logger:         logger,
	}
}

type queueListener struct {
	amqpConn       amqp.Connection
	websocketHub   websocket.Hub
	store          Store
	userCollection mongo.DbCollection
	workers        int
	decoder        encoding.Decoder
	configProvider config.ApiConfigProvider
	logger         zerolog.Logger
}

func (s *queueListener) Listen(ctx context.Context) error {
	channel, err := s.amqpConn.Channel()
	if err != nil {
		return fmt.Errorf("cannot create rmq channel: %w", err)
	}

	defer channel.Close()
	q, err := channel.QueueDeclare(
		"",    // name
		true,  // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("cannot declare queue: %w", err)
	}

	err = channel.QueueBind(
		q.Name,                               // name
		"",                                   // key
		canopsis.ApiNotificationExchangeName, // exchange
		false,                                // no-wait
		nil,                                  // arguments
	)
	if err != nil {
		return fmt.Errorf("cannot bind queue: %w", err)
	}

	ch, err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to consume events: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < s.workers; i++ {
		g.Go(func() error {
			var err error
			for {
				select {
				case <-ctx.Done():
					return nil
				case msg, ok := <-ch:
					if !ok {
						return errors.New("channel closed")
					}

					event := rpc.ApiNotificationEvent{}
					err = s.decoder.Decode(msg.Body, &event)
					if err != nil {
						s.logger.Err(err).Msg("failed to decode remediation result event")
						continue
					}

					err = s.processEvent(ctx, event)
					if err != nil {
						s.logger.Err(err).Msg("failed to process notification event")
						if mongo.IsConnectionError(err) {
							err = channel.Nack(msg.DeliveryTag, false, true)
							if err != nil {
								s.logger.Err(err).Msg("failed to negatively acknowledge message")
							}

							continue
						}
					}

					err = channel.Ack(msg.DeliveryTag, false)
					if err != nil {
						s.logger.Err(err).Msg("failed to acknowledge message")
					}
				}
			}
		})
	}

	return g.Wait()
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
