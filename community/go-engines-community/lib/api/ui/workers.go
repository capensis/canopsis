package ui

import (
	"context"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/security/token"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/wsconn"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func updateConfig(
	timezoneConfigProvider *config.BaseTimezoneConfigProvider,
	dataStorageConfigProvider *config.BaseDataStorageConfigProvider,
	apiConfigProvider *config.BaseApiConfigProvider,
	templateConfigProvider *config.BaseTemplateConfigProvider,
	techMetricsConfigProvider *config.BaseTechMetricsConfigProvider,
	alarmConfigProvider *config.BaseAlarmConfigProvider,
	configAdapter config.Adapter,
	userInterfaceConfigProvider *config.BaseUserInterfaceConfigProvider,
	userInterfaceAdapter config.UserInterfaceAdapter,
	interval time.Duration,
	logger zerolog.Logger,
) func(context.Context) error {
	return func(ctx context.Context) error {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cfg, err := configAdapter.GetConfig(ctx)
				if err != nil {
					logger.Err(err).Msg("fail to load config")
					continue
				}

				timezoneConfigProvider.Update(cfg)
				apiConfigProvider.Update(cfg)
				techMetricsConfigProvider.Update(cfg)
				dataStorageConfigProvider.Update(cfg)
				templateConfigProvider.Update(cfg)
				alarmConfigProvider.Update(cfg)

				userInterfaceConfig, err := userInterfaceAdapter.GetConfig(ctx)
				if err != nil {
					logger.Err(err).Msg("fail to load user interface config")
					continue
				}
				userInterfaceConfigProvider.Update(userInterfaceConfig)
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func updateTokenActivity(
	interval time.Duration,
	tokenTracker token.UsageTracker,
	shareTokenTracker token.UsageTracker,
	websocketHub websocket.Hub,
	logger zerolog.Logger,
) func(context.Context) error {
	return func(ctx context.Context) error {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				conns := websocketHub.Connections()
				seen := make(map[string]bool, len(conns))
				for _, c := range conns {
					if seen[c.Token] {
						continue
					}

					seen[c.Token] = true
					err := tokenTracker.Access(ctx, c.Token)
					if err != nil {
						logger.Err(err).Msg("cannot update token access")
					}

					err = shareTokenTracker.Access(ctx, c.Token)
					if err != nil {
						logger.Err(err).Msg("cannot update share token access")
					}
				}
			}
		}
	}
}

func removeExpiredTokens(
	interval time.Duration,
	client mongo.DbClient,
	logger zerolog.Logger,
) func(context.Context) error {
	return func(ctx context.Context) error {
		tokens := map[string]struct {
			collection mongo.DbCollection
			filter     func(now datetime.CpsTime) bson.M
		}{
			"expired tokens": {
				collection: client.Collection(mongo.TokenMongoCollection),
				filter: func(now datetime.CpsTime) bson.M {
					return bson.M{"expired": bson.M{"$ne": nil, "$lt": now}}
				},
			},
			"inactive tokens": {
				collection: client.Collection(mongo.TokenMongoCollection),
				filter: func(now datetime.CpsTime) bson.M {
					return bson.M{"expired_by_inactivity": bson.M{"$ne": nil, "$lt": now}}
				},
			},
			"expired share tokens": {
				collection: client.Collection(mongo.ShareTokenMongoCollection),
				filter: func(now datetime.CpsTime) bson.M {
					return bson.M{"expired": bson.M{"$ne": nil, "$lt": now}}
				},
			},
			"expired public tokens": {
				collection: client.Collection(mongo.PublicAPITokenCollection),
				filter: func(now datetime.CpsTime) bson.M {
					return bson.M{"expired": bson.M{"$ne": nil, "$lt": now}}
				},
			},
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				now := datetime.NewCpsTime()
				for k, t := range tokens {
					count, err := removeTokens(ctx, t.filter(now), t.collection)
					if err == nil {
						logger.Debug().Int64("deleted", count).Msg("deleted " + k)
					} else {
						logger.Err(err).Msg("cannot delete " + k)
					}
				}
			}
		}
	}
}

func removeTokens(ctx context.Context, filter bson.M, collection mongo.DbCollection) (int64, error) {
	cursor, err := collection.Find(ctx, filter, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)
	ids := make([]string, 0, canopsis.DefaultBulkSize)
	var count int64
	for cursor.Next(ctx) {
		doc := struct {
			ID string `json:"_id"`
		}{}
		err = cursor.Decode(&doc)
		if err != nil {
			return 0, err
		}

		ids = append(ids, doc.ID)
		if len(ids) == canopsis.DefaultBulkSize {
			d, err := collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
			if err != nil {
				return 0, err
			}

			count += d
			clear(ids)
		}
	}

	if len(ids) > 0 {
		d, err := collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
		if err != nil {
			return 0, err
		}

		count += d
	}

	return count, nil
}

func updateWebsocketConns(
	interval time.Duration,
	websocketHub websocket.Hub,
	websocketStore wsconn.Store,
	logger zerolog.Logger,
) func(context.Context) error {
	return func(ctx context.Context) error {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				err := websocketStore.SyncConnections(ctx, websocketHub.Connections())
				if err != nil {
					logger.Err(err).Msg("cannot update websocket connections")
					continue
				}

				c, err := websocketStore.CountActiveConnections(ctx)
				if err != nil {
					logger.Err(err).Msg("cannot get active websocket connections")
					continue
				}

				websocketHub.SendMessage(ctx, c, websocket.ToRoom(websocket.RoomLoggedUserCount))
			}
		}
	}
}

func sendPbhRecomputeEvents(
	pbhComputeChan <-chan rpc.PbehaviorRecomputeEvent,
	encoder encoding.Encoder,
	publisher libamqp.Publisher,
	logger zerolog.Logger,
) func(context.Context) error {
	return func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case event, ok := <-pbhComputeChan:
				if !ok {
					return nil
				}

				body, err := encoder.Encode(event)
				if err != nil {
					logger.Err(err).Msg("cannot encode event")
					continue
				}

				err = publisher.PublishWithContext(
					ctx,
					canopsis.DefaultExchangeName,
					canopsis.PBehaviorQueueRecomputeName,
					false,
					false,
					amqp.Publishing{
						ContentType:  canopsis.JsonContentType,
						Body:         body,
						DeliveryMode: amqp.Persistent,
					},
				)
				if err != nil {
					logger.Err(err).Msg("cannot send event")
				}
			}
		}
	}
}
