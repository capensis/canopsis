package che

import (
	"context"
	"fmt"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libcontext "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
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

	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(cfg), config.NewTimezoneConfigProvider(cfg, zerolog.Nop()))
	techMetricsConfigProvider := config.NewTechMetricsConfigProvider(cfg, zerolog.Nop())
	techMetricsSender := techmetrics.NewSender(techMetricsConfigProvider, canopsis.TechMetricsFlushInterval,
		cfg.Global.ReconnectRetries, cfg.Global.GetReconnectTimeout(), zerolog.Nop())
	ruleApplicatorContainer := eventfilter.NewRuleApplicatorContainer()
	ruleApplicatorContainer.Set(eventfilter.RuleTypeChangeEntity, eventfilter.NewChangeEntityApplicator(eventfilter.NewExternalDataGetterContainer(), tplExecutor))
	ruleApplicatorContainer.Set(eventfilter.RuleTypeEnrichment, eventfilter.NewEnrichmentApplicator(eventfilter.NewExternalDataGetterContainer(), eventfilter.NewActionProcessor(tplExecutor, techMetricsSender)))
	ruleApplicatorContainer.Set(eventfilter.RuleTypeDrop, eventfilter.NewDropApplicator())
	ruleApplicatorContainer.Set(eventfilter.RuleTypeBreak, eventfilter.NewBreakApplicator())
	ruleService := eventfilter.NewRuleService(eventfilter.NewRuleAdapter(dbClient), ruleApplicatorContainer, zerolog.Nop())
	err = ruleService.LoadRules(ctx, []string{eventfilter.RuleTypeDrop, eventfilter.RuleTypeEnrichment, eventfilter.RuleTypeBreak})
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}

	p := messageProcessor{
		FeatureContextCreation:   true,
		FeatureEventProcessing:   true,
		FeaturePrintEventOnError: true,

		AlarmConfigProvider: config.NewAlarmConfigProvider(cfg, zerolog.Nop()),
		EnrichmentCenter: libcontext.NewEnrichmentCenter(
			entity.NewAdapter(dbClient),
			dbClient,
			entityservice.NewManager(
				entityservice.NewAdapter(dbClient),
				entityservice.NewStorage(
					entityservice.NewAdapter(dbClient),
					redisClient,
					json.NewEncoder(),
					json.NewDecoder(),
					zerolog.Nop(),
				),
				zerolog.Nop(),
			),
			metrics.NewNullMetaUpdater(),
		),
		EventFilterService: ruleService,
		TechMetricsSender:  techMetricsSender,
		Encoder:            json.NewEncoder(),
		Decoder:            json.NewDecoder(),
		Logger:             zerolog.Nop(),
	}

	encoder := json.NewEncoder()

	b.ResetTimer()

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
	}
}
