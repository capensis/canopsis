package api

import (
	"context"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/sharetoken"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/token"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func updateConfig(
	timezoneConfigProvider *config.BaseTimezoneConfigProvider,
	dataStorageConfigProvider *config.BaseDataStorageConfigProvider,
	apiConfigProvider *config.BaseApiConfigProvider,
	templateConfigProvider *config.BaseTemplateConfigProvider,
	techMetricsConfigProvider *config.BaseTechMetricsConfigProvider,
	configAdapter config.Adapter,
	userInterfaceConfigProvider *config.BaseUserInterfaceConfigProvider,
	userInterfaceAdapter config.UserInterfaceAdapter,
	interval time.Duration,
	logger zerolog.Logger,
) func(context.Context) {
	return func(ctx context.Context) {
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

				userInterfaceConfig, err := userInterfaceAdapter.GetConfig(ctx)
				if err != nil {
					logger.Err(err).Msg("fail to load user interface config")
					continue
				}
				userInterfaceConfigProvider.Update(userInterfaceConfig)
			case <-ctx.Done():
				return
			}
		}
	}
}

func updateTokenActivity(
	interval time.Duration,
	tokenStore *token.MongoStore,
	shareTokenStore *sharetoken.MongoStore,
	websocketHub websocket.Hub,
	logger zerolog.Logger,
) func(context.Context) {
	return func(ctx context.Context) {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, t := range websocketHub.GetUserTokens() {
					err := tokenStore.Access(ctx, t)
					if err != nil {
						logger.Err(err).Msg("cannot update token access")
					}
					err = shareTokenStore.Access(ctx, t)
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
	tokenStore *token.MongoStore,
	shareTokenStore *sharetoken.MongoStore,
	logger zerolog.Logger,
) func(context.Context) {
	return func(ctx context.Context) {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := tokenStore.DeleteExpired(ctx)
				if err != nil {
					logger.Err(err).Msg("cannot delete expired tokens")
				}
				err = shareTokenStore.DeleteExpired(ctx)
				if err != nil {
					logger.Err(err).Msg("cannot delete expired share tokens")
				}
			}
		}
	}
}

func updateWebsocketConns(
	interval time.Duration,
	websocketHub websocket.Hub,
	websocketStore websocket.Store,
	logger zerolog.Logger,
) func(context.Context) {
	return func(ctx context.Context) {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := websocketStore.UpdateConnections(ctx, websocketHub.GetConnections())
				if err != nil {
					logger.Err(err).Msg("cannot update websocket connections")
					continue
				}

				c, err := websocketStore.GetActiveConnections(ctx)
				if err != nil {
					logger.Err(err).Msg("cannot get active websocket connections")
					continue
				}

				websocketHub.Send(ctx, websocket.RoomLoggedUserCount, c)
			}
		}
	}
}

func sendPbhRecomputeEvents(
	pbhComputeChan <-chan rpc.PbehaviorRecomputeEvent,
	encoder encoding.Encoder,
	publisher libamqp.Publisher,
	logger zerolog.Logger,
) func(context.Context) {
	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-pbhComputeChan:
				if !ok {
					return
				}

				body, err := encoder.Encode(event)
				if err != nil {
					logger.Err(err).Msg("cannot encode event")
					continue
				}
				err = publisher.PublishWithContext(
					ctx,
					"",
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
