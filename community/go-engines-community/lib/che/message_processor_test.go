package che

import (
	"context"
	"fmt"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func BenchmarkMessageProcessor_Process_GivenOldEntity(b *testing.B) {
	benchmarkMessageProcessor(b, "./testdata/fixtures/old_entity.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     "test-connector",
			ConnectorName: "test-connector-name",
			Component:     "test-component",
			Resource:      "test-resource",
			SourceType:    types.SourceTypeResource,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenNewEntity(b *testing.B) {
	benchmarkMessageProcessor(b, "./testdata/fixtures/new_entity.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     fmt.Sprintf("test-connector-%d", i),
			ConnectorName: fmt.Sprintf("test-connector-name-%d", i),
			Component:     fmt.Sprintf("test-component-%d", i),
			Resource:      fmt.Sprintf("test-resource-%d", i),
			SourceType:    types.SourceTypeResource,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenConnectorAndNewResource(b *testing.B) {
	benchmarkMessageProcessor(b, "./testdata/fixtures/connector_and_new_resource.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     "test-connector",
			ConnectorName: "test-connector-name",
			Component:     fmt.Sprintf("test-new-component-%d", i),
			Resource:      fmt.Sprintf("test-new-resource-%d", i),
			SourceType:    types.SourceTypeResource,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenServiceAndNewResource(b *testing.B) {
	benchmarkMessageProcessor(b, "./testdata/fixtures/service_and_new_resource.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     "test-connector",
			ConnectorName: "test-connector-name",
			Component:     fmt.Sprintf("test-new-component-%d", i),
			Resource:      fmt.Sprintf("test-new-resource-%d", i),
			SourceType:    types.SourceTypeResource,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenOldEntityAndUnmatchedEventFilters(b *testing.B) {
	benchmarkMessageProcessor(b, "./testdata/fixtures/old_entity_and_unmatched_event_filters.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     "test-connector",
			ConnectorName: "test-connector-name",
			Component:     "test-component",
			Resource:      "test-resource",
			SourceType:    types.SourceTypeResource,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenOldEntityAndMatchedEnrichmentEventFilters(b *testing.B) {
	benchmarkMessageProcessor(b, "./testdata/fixtures/old_entity_and_matched_enrichment_event_filters.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     "test-connector",
			ConnectorName: "test-connector-name",
			Component:     "test-component",
			Resource:      "test-resource",
			SourceType:    types.SourceTypeResource,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenOldEntityAndMatchedEnrichmentEntityEventFilters(b *testing.B) {
	benchmarkMessageProcessor(b, "./testdata/fixtures/old_entity_and_matched_enrichment_entity_event_filters.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     "test-connector",
			ConnectorName: "test-connector-name",
			Component:     "test-component",
			Resource:      "test-resource",
			SourceType:    types.SourceTypeResource,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenNewEntityAndMatchedEnrichmentEntityEventFilters(b *testing.B) {
	benchmarkMessageProcessor(b, "./testdata/fixtures/new_entity_and_matched_enrichment_entity_event_filters.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     fmt.Sprintf("test-connector-%d", i),
			ConnectorName: fmt.Sprintf("test-connector-name-%d", i),
			Component:     fmt.Sprintf("test-component-%d", i),
			Resource:      fmt.Sprintf("test-resource-%d", i),
			SourceType:    types.SourceTypeResource,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenNewEntityAndMatchedDropEventfilterWith10ResolvedExdates(b *testing.B) {
	benchmarkMessageProcessor(b, "./testdata/fixtures/eventfilters_with_rrule_10_resolved_exdates.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     fmt.Sprintf("test-connector-%d", i),
			ConnectorName: fmt.Sprintf("test-connector-name-%d", i),
			Component:     fmt.Sprintf("test-component-%d", i),
			Resource:      fmt.Sprintf("test-resource-%d", i),
			SourceType:    types.SourceTypeResource,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenNewEntityAndMatchedDropEventfilterWith100ResolvedExdates(b *testing.B) {
	benchmarkMessageProcessor(b, "./testdata/fixtures/eventfilters_with_rrule_100_resolved_exdates.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     fmt.Sprintf("test-connector-%d", i),
			ConnectorName: fmt.Sprintf("test-connector-name-%d", i),
			Component:     fmt.Sprintf("test-component-%d", i),
			Resource:      fmt.Sprintf("test-resource-%d", i),
			SourceType:    types.SourceTypeResource,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenNewEntityAndMatchedEnrichmentEntityEventFiltersWithTpl(b *testing.B) {
	benchmarkMessageProcessor(b, "./testdata/fixtures/new_entity_and_matched_enrichment_entity_event_filters_with_tpl.yml", func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     fmt.Sprintf("test-connector-%d", i),
			ConnectorName: fmt.Sprintf("test-connector-name-%d", i),
			Component:     fmt.Sprintf("test-component-%d", i),
			Resource:      fmt.Sprintf("test-resource-%d", i),
			SourceType:    types.SourceTypeResource,
		}
	})
}

func BenchmarkMessageProcessor_Process_GivenNewEntityAndMatchedEnrichmentEntityEventFiltersWithTplWithEnvVars(b *testing.B) {
	cfg := config.CanopsisConf{
		Template: config.SectionTemplate{
			Vars: map[string]any{
				"Location": "FR",
			},
		},
	}
	benchmarkMessageProcessorWithConfig(b, "./testdata/fixtures/new_entity_and_matched_enrichment_entity_event_filters_with_tpl_with_env_vars.yml", cfg, func(i int) types.Event {
		return types.Event{
			EventType:     types.EventTypeCheck,
			Connector:     fmt.Sprintf("test-connector-%d", i),
			ConnectorName: fmt.Sprintf("test-connector-name-%d", i),
			Component:     fmt.Sprintf("test-component-%d", i),
			Resource:      fmt.Sprintf("test-resource-%d", i),
			SourceType:    types.SourceTypeResource,
		}
	})
}

func benchmarkMessageProcessor(
	b *testing.B,
	fixturesPath string,
	genEvent func(i int) types.Event,
) {
	benchmarkMessageProcessorWithConfig(b, fixturesPath, config.CanopsisConf{}, genEvent)
}

func benchmarkMessageProcessorWithConfig(
	b *testing.B,
	fixturesPath string,
	cfg config.CanopsisConf,
	genEvent func(i int) types.Event,
) {
	defer func() {
		if r := recover(); r != nil {
			b.Fatal("benchmark failed due to panic:", r)
		}
	}()
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
	redisClient, err := redis.NewSession(ctx, redis.EntityServiceStorage, zerolog.Nop(), 0, 0)
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

	alarmConfigProvider := config.NewAlarmConfigProvider(cfg, zerolog.Nop())
	failureService := eventfilter.NewFailureService(dbClient, time.Hour, zerolog.Nop())
	eventCounter := eventfilter.NewEventCounter(dbClient, time.Hour, zerolog.Nop())
	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(cfg, zerolog.Nop()), config.NewTimezoneConfigProvider(cfg, zerolog.Nop()))
	techMetricsConfigProvider := config.NewTechMetricsConfigProvider(cfg, zerolog.Nop())
	techMetricsSender := techmetrics.NewSender(canopsis.CheEngineName+"/"+utils.NewID(), techMetricsConfigProvider, canopsis.TechMetricsFlushInterval,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), zerolog.Nop())
	ruleApplicatorContainer := eventfilter.NewRuleApplicatorContainer()
	ruleApplicatorContainer.Set(eventfilter.RuleTypeChangeEntity, eventfilter.NewChangeEntityApplicator(eventfilter.NewExternalDataGetterContainer(), failureService, tplExecutor))
	ruleApplicatorContainer.Set(eventfilter.RuleTypeEnrichment, eventfilter.NewEnrichmentApplicator(eventfilter.NewExternalDataGetterContainer(), eventfilter.NewActionProcessor(alarmConfigProvider, failureService, tplExecutor, techMetricsSender), failureService))
	ruleApplicatorContainer.Set(eventfilter.RuleTypeDrop, eventfilter.NewDropApplicator())
	ruleApplicatorContainer.Set(eventfilter.RuleTypeBreak, eventfilter.NewBreakApplicator())
	ruleService := eventfilter.NewRuleService(eventfilter.NewRuleAdapter(dbClient), ruleApplicatorContainer, eventCounter, failureService, tplExecutor, zerolog.Nop())
	err = ruleService.LoadRules(ctx, []string{eventfilter.RuleTypeDrop, eventfilter.RuleTypeEnrichment, eventfilter.RuleTypeBreak})
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}

	pgPoolProvider := postgres.NewPoolProvider(cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout())
	metricsConfigProvider := config.NewMetricsConfigProvider(cfg, zerolog.Nop())
	p := messageProcessor{
		FeaturePrintEventOnError: true,
		AlarmConfigProvider:      alarmConfigProvider,
		MetricsSender:            metrics.NewTimescaleDBSender(pgPoolProvider, metricsConfigProvider, zerolog.Nop()),
		MetaUpdater:              metrics.NewNullMetaUpdater(),
		TechMetricsSender:        techMetricsSender,
		EntityCollection:         dbClient.Collection(mongo.EntityMongoCollection),
		Encoder:                  json.NewEncoder(),
		Decoder:                  json.NewDecoder(),
		Logger:                   zerolog.Nop(),
		// AmqpPublisher field has not accessed by test paths, otherwise it has to be initialized with mock value
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		event := genEvent(i)
		body, err := p.Encoder.Encode(event)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
		b.StartTimer()
		_, err = p.Process(ctx, amqp.Delivery{
			Body: body,
		})
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}
