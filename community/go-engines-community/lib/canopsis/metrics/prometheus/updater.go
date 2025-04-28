package prometheus

import (
	"context"
	"errors"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/websocket"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	redisStep = 1000
)

const (
	DynamicInfosLabel      = "dynamic-infos"
	EventFilterLabel       = "event-filter"
	FlappingRuleLabel      = "flapping-rule"
	IdleRuleLabel          = "idle-rule"
	LinkRuleLabel          = "link-rule"
	MetaAlarmRuleLabel     = "meta-alarm-rule"
	PbehaviorLabel         = "pbehavior"
	ResolveRuleLabel       = "resolve-rule"
	SnmpRuleLabel          = "snmp-rule"
	ScenarioLabel          = "scenario"
	DeclareTicketRuleLabel = "declare-ticket-rule"
)

type Updater interface {
	Update(ctx context.Context, m *Metrics)
}

type updater struct {
	runInfoManager     engine.RunInfoManager
	healthCheckAdapter config.HealthCheckAdapter
	pbhRedisClient     *redis.Client
	websocketStore     websocket.Store
	logger             zerolog.Logger

	entityMongoCollection        libmongo.DbCollection
	userCollection               libmongo.DbCollection
	alarmCollection              libmongo.DbCollection
	resolvedAlarmCollection      libmongo.DbCollection
	eventfilterCollection        libmongo.DbCollection
	metaAlarmRulesCollection     libmongo.DbCollection
	eventfilterFailureCollection libmongo.DbCollection
	dynamicInfosCollection       libmongo.DbCollection

	rulesLabelCollMap map[string]libmongo.DbCollection
}

func NewUpdater(
	client libmongo.DbClient,
	manager engine.RunInfoManager,
	adapter config.HealthCheckAdapter,
	pbhClient *redis.Client,
	store websocket.Store,
	logger zerolog.Logger,
) Updater {
	return &updater{
		runInfoManager:     manager,
		healthCheckAdapter: adapter,
		pbhRedisClient:     pbhClient,
		websocketStore:     store,
		logger:             logger,

		entityMongoCollection:        client.Collection(libmongo.EntityMongoCollection),
		userCollection:               client.Collection(libmongo.UserCollection),
		alarmCollection:              client.Collection(libmongo.AlarmMongoCollection),
		resolvedAlarmCollection:      client.Collection(libmongo.ResolvedAlarmMongoCollection),
		eventfilterCollection:        client.Collection(libmongo.EventFilterRuleCollection),
		metaAlarmRulesCollection:     client.Collection(libmongo.MetaAlarmRulesMongoCollection),
		eventfilterFailureCollection: client.Collection(libmongo.EventFilterFailureCollection),
		dynamicInfosCollection:       client.Collection(libmongo.DynamicInfosRulesMongoCollection),

		rulesLabelCollMap: map[string]libmongo.DbCollection{
			DynamicInfosLabel:      client.Collection(libmongo.DynamicInfosRulesMongoCollection),
			EventFilterLabel:       client.Collection(libmongo.EventFilterRuleCollection),
			FlappingRuleLabel:      client.Collection(libmongo.FlappingRuleMongoCollection),
			IdleRuleLabel:          client.Collection(libmongo.IdleRuleMongoCollection),
			LinkRuleLabel:          client.Collection(libmongo.LinkRuleMongoCollection),
			MetaAlarmRuleLabel:     client.Collection(libmongo.MetaAlarmRulesMongoCollection),
			PbehaviorLabel:         client.Collection(libmongo.PbehaviorMongoCollection),
			ResolveRuleLabel:       client.Collection(libmongo.ResolveRuleMongoCollection),
			SnmpRuleLabel:          client.Collection(libmongo.SnmpRulesCollection),
			ScenarioLabel:          client.Collection(libmongo.ScenarioMongoCollection),
			DeclareTicketRuleLabel: client.Collection(libmongo.DeclareTicketRuleMongoCollection),
		},
	}
}

func (u *updater) Update(ctx context.Context, m *Metrics) {
	u.logger.Debug().Msg("fetching metrics from the db")

	metricsValues := newMetricsValues()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		count, err := u.entityMongoCollection.CountDocuments(ctx, bson.M{"enabled": true})
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to count active entities from db")
		}

		metricsValues.SetGauge(ActiveEntitiesGauge, float64(count))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		count, err := u.entityMongoCollection.CountDocuments(ctx, bson.M{"enabled": false})
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to count number of disabled entities from db")
		}

		metricsValues.SetGauge(DisabledEntitiesGauge, float64(count))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		count, err := u.userCollection.CountDocuments(ctx, bson.M{"enable": true})
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to count number of active users from db")
		}

		metricsValues.SetGauge(EnabledUsersGauge, float64(count))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		count, err := u.eventfilterCollection.CountDocuments(ctx, bson.M{"enabled": true})
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to count number of event filter rules from db")
		}

		metricsValues.SetGauge(EventFiltersGauge, float64(count))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		count, err := u.metaAlarmRulesCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to count number of meta alarm rules from db")
		}

		metricsValues.SetGauge(MetaAlarmsRulesGauge, float64(count))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		count, err := u.dynamicInfosCollection.CountDocuments(ctx, bson.M{"enabled": true})
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to count number of dynamic infos rules from db")
		}

		metricsValues.SetGauge(DynamicInfosRulesGauge, float64(count))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		count, err := u.alarmCollection.CountDocuments(ctx, bson.M{"v.resolved": nil})
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to count number of active alarms from db")
		}

		metricsValues.SetGauge(OpenedAlarmsGauge, float64(count))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		count, err := u.resolvedAlarmCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to count number of closed alarms from db")
		}

		metricsValues.SetGauge(ResolvedAlarmsGauge, float64(count))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		count, err := u.eventfilterFailureCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to get event filter failures from db")
		}

		metricsValues.SetGauge(EventfilterErrorsGauge, float64(count))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		count, err := u.websocketStore.GetActiveConnections(ctx)
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to get active connections")
		}

		metricsValues.SetGauge(UserConnectionsGauge, float64(count))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		cfg, err := u.healthCheckAdapter.GetConfig(ctx)
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to get healthcheck config")
		}

		runInfo, err := u.runInfoManager.GetEngines(ctx)
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to get engines run info")
		}

		runInfoMap := make(map[string]bool, len(runInfo))
		for _, info := range runInfo {
			runInfoMap[info.Name] = true
		}

		processed := make(map[string]bool)
		for _, pair := range cfg.EngineOrder {
			if !processed[pair.From] {
				processed[pair.From] = true

				if runInfoMap[pair.From] {
					metricsValues.SetGaugeVector(EngineStatusGaugeVector, pair.From, 1)
				} else {
					metricsValues.SetGaugeVector(EngineStatusGaugeVector, pair.From, 0)
				}
			}

			if !processed[pair.To] {
				processed[pair.To] = true

				if runInfoMap[pair.To] {
					metricsValues.SetGaugeVector(EngineStatusGaugeVector, pair.To, 1)
				} else {
					metricsValues.SetGaugeVector(EngineStatusGaugeVector, pair.To, 0)
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		var cursor uint64
		var activePbh int

		for {
			res := u.pbhRedisClient.Scan(ctx, cursor, libredis.PbehaviorComputedKey+"*", redisStep)
			if err := res.Err(); err != nil {
				u.logger.Error().Err(err).Msg("cannot scan active pbehavior keys")
				break
			}

			var keys []string
			keys, cursor = res.Val()

			activePbh += len(keys)

			if cursor == 0 {
				break
			}
		}

		metricsValues.SetGauge(ActivePBehaviorsGauge, float64(activePbh))
	}()

	opts := options.FindOne().SetSort(bson.M{"updated": -1}).SetProjection(bson.M{"updated": 1})

	wg.Add(len(u.rulesLabelCollMap))
	for label, coll := range u.rulesLabelCollMap {
		go func() {
			defer wg.Done()

			var doc struct {
				Updated float64 `bson:"updated"`
			}

			val := 0.0

			err := coll.FindOne(ctx, bson.M{}, opts).Decode(&doc)
			if err != nil {
				if !errors.Is(err, mongo.ErrNoDocuments) {
					u.logger.Error().Err(err).Msgf("failed to get latest update timestamp for %s", label)
				}
			} else {
				val = doc.Updated
			}

			metricsValues.SetGaugeVector(LastExploitationModTimeGaugeVector, label, val)
		}()
	}

	wg.Wait()

	m.Set(metricsValues)
}
