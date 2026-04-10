package prometheus

import (
	"context"
	"errors"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/wsconn"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	libredis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

const (
	instrTypeLabelManual           = "manual"
	instrTypeLabelAuto             = "auto"
	instrTypeLabelSimplifiedManual = "simplified_manual"
)

type Updater interface {
	Update(ctx context.Context, m *DbCollectionsMetrics)
}

type updater struct {
	runInfoManager     engine.RunInfoManager
	healthCheckAdapter config.HealthCheckAdapter
	pbhRedisClient     *redis.Client
	websocketStore     wsconn.Store
	logger             zerolog.Logger

	entityMongoCollection        libmongo.DbCollection
	userCollection               libmongo.DbCollection
	alarmCollection              libmongo.DbCollection
	resolvedAlarmCollection      libmongo.DbCollection
	eventfilterCollection        libmongo.DbCollection
	metaAlarmRulesCollection     libmongo.DbCollection
	eventfilterFailureCollection libmongo.DbCollection
	dynamicInfosCollection       libmongo.DbCollection
	instructionCollection        libmongo.DbCollection

	rulesLabelCollMap map[string]libmongo.DbCollection

	entityTypeLabels      []string
	instructionTypeLabels []string
}

func NewUpdater(
	client libmongo.DbClient,
	manager engine.RunInfoManager,
	adapter config.HealthCheckAdapter,
	pbhClient *redis.Client,
	store wsconn.Store,
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
		instructionCollection:        client.Collection(libmongo.InstructionMongoCollection),

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
			ScenarioLabel:          client.Collection(libmongo.ScenarioCollection),
			DeclareTicketRuleLabel: client.Collection(libmongo.DeclareTicketRuleCollection),
		},

		entityTypeLabels: []string{
			types.EntityTypeResource,
			types.EntityTypeComponent,
			types.EntityTypeConnector,
			types.EntityTypeService,
		},
		instructionTypeLabels: []string{
			instrTypeLabelManual,
			instrTypeLabelAuto,
			instrTypeLabelSimplifiedManual,
		},
	}
}

func (u *updater) Update(ctx context.Context, m *DbCollectionsMetrics) {
	u.logger.Debug().Msg("fetching metrics from the db")

	metricsValues := newDbCollectionsMetricsValues()

	var wg sync.WaitGroup

	wg.Add(len(u.entityTypeLabels))
	for _, label := range u.entityTypeLabels {
		go func() {
			defer wg.Done()

			count, err := u.entityMongoCollection.CountDocuments(ctx, bson.M{"enabled": true, "type": label})
			if err != nil {
				u.logger.Err(err).Str("label", label).Msg("failed to count entities from db")
			}

			metricsValues.SetGaugeVector(ActiveEntitiesGaugeVector, label, float64(count))
		}()
	}

	wg.Add(len(u.instructionTypeLabels))
	for t, label := range u.instructionTypeLabels {
		go func() {
			defer wg.Done()

			count, err := u.instructionCollection.CountDocuments(ctx, bson.M{"enabled": true, "type": t})
			if err != nil {
				u.logger.Err(err).Str("label", label).Msg("failed to count instructions from db")
			}

			metricsValues.SetGaugeVector(InstructionsGaugeVector, label, float64(count))
		}()
	}

	wg.Go(func() {
		count, err := u.entityMongoCollection.CountDocuments(ctx, bson.M{"enabled": false})
		if err != nil {
			u.logger.Err(err).Msg("failed to count number of disabled entities from db")
		}

		metricsValues.SetGauge(DisabledEntitiesGauge, float64(count))
	})

	wg.Go(func() {
		count, err := u.userCollection.CountDocuments(ctx, bson.M{"enable": true})
		if err != nil {
			u.logger.Err(err).Msg("failed to count number of active users from db")
		}

		metricsValues.SetGauge(EnabledUsersGauge, float64(count))
	})

	wg.Go(func() {
		count, err := u.eventfilterCollection.CountDocuments(ctx, bson.M{"enabled": true})
		if err != nil {
			u.logger.Err(err).Msg("failed to count number of event filter rules from db")
		}

		metricsValues.SetGauge(EventFiltersGauge, float64(count))
	})

	wg.Go(func() {
		count, err := u.metaAlarmRulesCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			u.logger.Err(err).Msg("failed to count number of meta alarm rules from db")
		}

		metricsValues.SetGauge(MetaAlarmsRulesGauge, float64(count))
	})

	wg.Go(func() {
		count, err := u.dynamicInfosCollection.CountDocuments(ctx, bson.M{"enabled": true})
		if err != nil {
			u.logger.Err(err).Msg("failed to count number of dynamic infos rules from db")
		}

		metricsValues.SetGauge(DynamicInfosRulesGauge, float64(count))
	})

	wg.Go(func() {
		countActive, err := u.alarmCollection.CountDocuments(ctx, bson.M{"v.resolved": nil, "v.activation_date": bson.M{"$exists": true}})
		if err != nil {
			u.logger.Err(err).Msg("failed to count number of active alarms from db")
		}

		countInactive, err := u.alarmCollection.CountDocuments(ctx, bson.M{"v.resolved": nil, "v.activation_date": bson.M{"$exists": false}})
		if err != nil {
			u.logger.Err(err).Msg("failed to count number of active alarms from db")
		}

		metricsValues.SetGaugeVector(OpenedAlarmsGaugeVector, "true", float64(countActive))
		metricsValues.SetGaugeVector(OpenedAlarmsGaugeVector, "false", float64(countInactive))
	})

	wg.Go(func() {
		count, err := u.resolvedAlarmCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			u.logger.Err(err).Msg("failed to count number of closed alarms from db")
		}

		metricsValues.SetGauge(ResolvedAlarmsGauge, float64(count))
	})

	wg.Go(func() {
		count, err := u.eventfilterFailureCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			u.logger.Err(err).Msg("failed to get event filter failures from db")
		}

		metricsValues.SetGauge(EventfilterErrorsGauge, float64(count))
	})

	wg.Go(func() {
		count, err := u.websocketStore.CountActiveConnections(ctx)
		if err != nil {
			u.logger.Err(err).Msg("failed to get active connections")
		}

		metricsValues.SetGauge(UserConnectionsGauge, float64(count))
	})

	wg.Go(func() {
		cfg, err := u.healthCheckAdapter.GetConfig(ctx)
		if err != nil {
			u.logger.Err(err).Msg("failed to get healthcheck config")
		}

		runInfo, err := u.runInfoManager.GetEngines(ctx)
		if err != nil {
			u.logger.Err(err).Msg("failed to get engines run info")
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
	})

	wg.Go(func() {
		var cursor uint64
		var activePbh int

		for {
			res := u.pbhRedisClient.Scan(ctx, cursor, libredis.PbehaviorComputedKey+"*", redisStep)
			if err := res.Err(); err != nil {
				u.logger.Err(err).Msg("cannot scan active pbehavior keys")
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
	})

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
					u.logger.Err(err).Msgf("failed to get latest update timestamp for %s", label)
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
