package axe

import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/flappingrule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statistics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func BenchmarkMessageProcessor_Process_GivenNewAlarm(b *testing.B) {
	now := types.NewCpsTime()
	entity := &types.Entity{
		ID:        "test-resource/test-component",
		Name:      "test-resource",
		Enabled:   true,
		Type:      types.EntityTypeResource,
		Created:   now,
		Impacts:   []string{"test-component"},
		Depends:   []string{"test-connector/test-connector-name"},
		Connector: "test-connector/test-connector-name",
		Component: "test-component",
	}
	benchmarkMessageProcessor(b, "./testdata/fixtures/new_alarm.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     "test-connector",
			ConnectorName: fmt.Sprintf("test-connector-name-%d", i),
			Component:     "test-component",
			Resource:      "test-resource",
			SourceType:    types.SourceTypeResource,
			State:         types.AlarmStateCritical,
			Entity:        entity,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenOldAlarm(b *testing.B) {
	now := types.NewCpsTime()
	entity := &types.Entity{
		ID:        "test-resource/test-component",
		Name:      "test-resource",
		Enabled:   true,
		Type:      types.EntityTypeResource,
		Created:   now,
		Impacts:   []string{"test-component"},
		Depends:   []string{"test-connector/test-connector-name"},
		Connector: "test-connector/test-connector-name",
		Component: "test-component",
	}
	benchmarkMessageProcessor(b, "./testdata/fixtures/old_alarm.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     "test-connector",
			ConnectorName: "test-connector-name",
			Component:     "test-component",
			Resource:      "test-resource",
			SourceType:    types.SourceTypeResource,
			State:         types.AlarmStateCritical,
			Entity:        entity,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenNoAlarm(b *testing.B) {
	now := types.NewCpsTime()
	entity := &types.Entity{
		ID:        "test-resource/test-component",
		Name:      "test-resource",
		Enabled:   true,
		Type:      types.EntityTypeResource,
		Created:   now,
		Impacts:   []string{"test-component"},
		Depends:   []string{"test-connector/test-connector-name"},
		Connector: "test-connector/test-connector-name",
		Component: "test-component",
	}
	benchmarkMessageProcessor(b, "./testdata/fixtures/no_alarm.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     "test-connector",
			ConnectorName: "test-connector-name",
			Component:     "test-component",
			Resource:      "test-resource",
			SourceType:    types.SourceTypeResource,
			State:         types.AlarmStateOK,
			Entity:        entity,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenNewAlarmState(b *testing.B) {
	const alarmsCount = 1000
	now := types.NewCpsTime()
	benchmarkMessageProcessor(b, "./testdata/fixtures/new_alarm_state.yml", func(i int) types.Event {
		alarmIndex := (i % alarmsCount) + 1
		state := ((i/alarmsCount + 1) % 3) + 1
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     "test-connector",
			ConnectorName: "test-connector-name",
			Component:     "test-component",
			Resource:      fmt.Sprintf("test-resource-%d", alarmIndex),
			SourceType:    types.SourceTypeResource,
			State:         types.CpsNumber(state),
			Entity: &types.Entity{
				ID:        fmt.Sprintf("test-resource-%d/test-component", alarmIndex),
				Name:      fmt.Sprintf("test-resource-%d", alarmIndex),
				Enabled:   true,
				Type:      types.EntityTypeResource,
				Created:   now,
				Impacts:   []string{"test-component"},
				Depends:   []string{"test-connector/test-connector-name"},
				Connector: "test-connector/test-connector-name",
				Component: "test-component",
			},
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenNewComment(b *testing.B) {
	const alarmsCount = 1000
	now := types.NewCpsTime()
	benchmarkMessageProcessor(b, "./testdata/fixtures/new_comment.yml", func(i int) types.Event {
		alarmIndex := (i % alarmsCount) + 1
		return types.Event{
			EventType:     types.EventTypeComment,
			Connector:     "test-connector",
			ConnectorName: "test-connector-name",
			Component:     "test-component",
			Resource:      fmt.Sprintf("test-resource-%d", alarmIndex),
			SourceType:    types.SourceTypeResource,
			Output:        fmt.Sprintf("test-output-%d", i),
			Entity: &types.Entity{
				ID:        fmt.Sprintf("test-resource-%d/test-component", alarmIndex),
				Name:      fmt.Sprintf("test-resource-%d", alarmIndex),
				Enabled:   true,
				Type:      types.EntityTypeResource,
				Created:   now,
				Impacts:   []string{"test-component"},
				Depends:   []string{"test-connector/test-connector-name"},
				Connector: "test-connector/test-connector-name",
				Component: "test-component",
			},
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenNewMetaAlarm(b *testing.B) {
	const alarmsCount = 100
	now := types.NewCpsTime()
	children := make([]string, alarmsCount)
	for i := 0; i < alarmsCount; i++ {
		children[i] = fmt.Sprintf("test-resource-%d/test-component", i+1)
	}
	benchmarkMessageProcessor(b, "./testdata/fixtures/new_meta_alarm.yml", func(i int) types.Event {
		resource := "meta-alarm-entity-" + utils.NewID()
		return types.Event{
			EventType:         types.EventTypeMetaAlarm,
			Component:         "metaalarm",
			Connector:         "engine",
			ConnectorName:     "correlation",
			Resource:          resource,
			SourceType:        types.SourceTypeResource,
			MetaAlarmChildren: &children,
			MetaAlarmRuleID:   "test-metaalarm-rule",
			Entity: &types.Entity{
				ID:      fmt.Sprintf("%s/metaalarm", resource),
				Name:    resource,
				Enabled: true,
				Type:    types.EntityTypeResource,
				Created: now,
			},
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenManyAlarmSteps(b *testing.B) {
	const alarmsCount = 1000
	now := types.NewCpsTime()
	benchmarkMessageProcessor(b, "./testdata/fixtures/many_alarm_steps.yml", func(i int) types.Event {
		alarmIndex := (i % alarmsCount) + 1
		state := ((i/alarmsCount + 1) % 3) + 1
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     "test-connector",
			ConnectorName: "test-connector-name",
			Component:     "test-component",
			Resource:      fmt.Sprintf("test-resource-%d", alarmIndex),
			SourceType:    types.SourceTypeResource,
			State:         types.CpsNumber(state),
			Entity: &types.Entity{
				ID:        fmt.Sprintf("test-resource-%d/test-component", alarmIndex),
				Name:      fmt.Sprintf("test-resource-%d", alarmIndex),
				Enabled:   true,
				Type:      types.EntityTypeResource,
				Created:   now,
				Impacts:   []string{"test-component"},
				Depends:   []string{"test-connector/test-connector-name"},
				Connector: "test-connector/test-connector-name",
				Component: "test-component",
			},
		}
	}, func(ctx context.Context, dbClient mongo.DbClient) error {
		stepsCount := types.AlarmStepsHardLimit - 50
		steps := make([]types.AlarmStep, stepsCount)
		for i := 0; i < stepsCount; i++ {
			steps[i] = types.AlarmStep{
				Type:      types.AlarmStepComment,
				Timestamp: now,
			}
		}
		_, err := dbClient.Collection(mongo.AlarmMongoCollection).UpdateMany(ctx, bson.M{}, bson.M{
			"$push": bson.M{
				"v.steps": bson.M{"$each": steps},
			},
		})
		return err
	})
}

func benchmarkMessageProcessor(
	b *testing.B,
	fixturesPath string,
	genEvent func(i int) types.Event,
	adjustFixtures ...func(ctx context.Context, dbClient mongo.DbClient) error,
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dbClient, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}
	b.Cleanup(func() {
		err := dbClient.Disconnect(context.Background())
		if err != nil {
			b.Errorf("unexpected error %v", err)
		}
	})
	redisClient, err := redis.NewSession(ctx, redis.PBehaviorLockStorage, zerolog.Nop(), 0, 0)
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}
	b.Cleanup(func() {
		err := redisClient.Close()
		if err != nil {
			b.Errorf("unexpected error %v", err)
		}
	})

	loader := fixtures.NewLoader(dbClient, []string{fixturesPath},
		fixtures.NewParser(fixtures.NewFaker(password.NewSha1Encoder())), zerolog.Nop())
	err = loader.Load(ctx)
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}
	b.Cleanup(func() {
		err := loader.Clean(context.Background())
		if err != nil {
			b.Errorf("unexpected error %v", err)
		}
	})

	cfg := config.CanopsisConf{}
	logger := zerolog.Nop()
	metricsSender := metrics.NewNullSender()
	alarmConfigProvider := config.NewAlarmConfigProvider(cfg, logger)
	tzConfigProvider := config.NewTimezoneConfigProvider(cfg, logger)
	alarmStatusService := alarmstatus.NewService(flappingrule.NewAdapter(dbClient), alarmConfigProvider, logger)
	metaAlarmEventProcessor := alarm.NewMetaAlarmEventProcessor(dbClient, alarm.NewAdapter(dbClient), correlation.NewRuleAdapter(dbClient),
		alarmStatusService, alarmConfigProvider, json.NewEncoder(), nil, canopsis.FIFOExchangeName, canopsis.FIFOQueueName,
		metricsSender, logger)

	p := messageProcessor{
		FeaturePrintEventOnError: true,
		EventProcessor: alarm.NewEventProcessor(
			dbClient,
			alarm.NewAdapter(dbClient),
			libentity.NewAdapter(dbClient),
			correlation.NewRuleAdapter(dbClient),
			alarmConfigProvider,
			DependencyMaker{}.DepOperationExecutor(dbClient, alarmConfigProvider, alarmStatusService, metricsSender),
			alarmStatusService,
			metrics.NewNullSender(),
			metaAlarmEventProcessor,
			statistics.NewEventStatisticsSender(dbClient, logger, tzConfigProvider),
			pbehavior.NewEntityTypeResolver(pbehavior.NewStore(redisClient, json.NewEncoder(), json.NewDecoder()), pbehavior.NewEntityMatcher(dbClient), logger),
			logger,
		),
		RemediationRpcClient:   nil,
		TimezoneConfigProvider: tzConfigProvider,
		PbehaviorAdapter:       pbehavior.NewAdapter(dbClient),
		Encoder:                json.NewEncoder(),
		Decoder:                json.NewDecoder(),
		Logger:                 logger,
	}

	encoder := json.NewEncoder()

	for _, f := range adjustFixtures {
		err = f(ctx, dbClient)
	}

	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		event := genEvent(i)
		body, err := encoder.Encode(event)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
		_, err = p.Process(ctx, amqp.Delivery{
			Body: body,
		})
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
		count++
	}
}
