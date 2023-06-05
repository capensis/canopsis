package che_test

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"github.com/rs/zerolog"
)

func BenchmarkEventfilterIntervalsWorker_Work_Given10Exdates(b *testing.B) {
	benchmarkRulesChangesWatcher(b, "./testdata/fixtures/eventfilters_with_rrule_10_exdates.yml")
}

func BenchmarkEventfilterIntervalsWorker_Work_Given100Exdates(b *testing.B) {
	benchmarkRulesChangesWatcher(b, "./testdata/fixtures/eventfilters_with_rrule_100_exdates.yml")
}

func BenchmarkEventfilterIntervalsWorker_Work_Given1000Exdates(b *testing.B) {
	benchmarkRulesChangesWatcher(b, "./testdata/fixtures/eventfilters_with_rrule_1000_exdates.yml")
}

func BenchmarkEventfilterIntervalsWorker_Work_Given10Exdates_10ExceptionsWith_1Exdate(b *testing.B) {
	benchmarkRulesChangesWatcher(b, "./testdata/fixtures/eventfilters_with_rrule_10_exdates_10_exceptions_1_exdate.yml")
}

func BenchmarkEventfilterIntervalsWorker_Work_Given10Exdates_10ExceptionsWith_10Exdates(b *testing.B) {
	benchmarkRulesChangesWatcher(b, "./testdata/fixtures/eventfilters_with_rrule_10_exdates_10_exceptions_10_exdate.yml")
}

func BenchmarkEventfilterIntervalsWorker_Work_Given10Exdates_10ExceptionsWith_100Exdates(b *testing.B) {
	benchmarkRulesChangesWatcher(b, "./testdata/fixtures/eventfilters_with_rrule_10_exdates_10_exceptions_100_exdate.yml")
}

func BenchmarkEventfilterIntervalsWorker_Work_Given10Exdates_100ExceptionsWith_100Exdates(b *testing.B) {
	benchmarkRulesChangesWatcher(b, "./testdata/fixtures/eventfilters_with_rrule_10_exdates_100_exceptions_100_exdate.yml")
}

func BenchmarkEventfilterIntervalsWorker_Work_Given100Exdates_100ExceptionsWith_100Exdates(b *testing.B) {
	benchmarkRulesChangesWatcher(b, "./testdata/fixtures/eventfilters_with_rrule_100_exdates_100_exceptions_100_exdate.yml")
}

func benchmarkRulesChangesWatcher(b *testing.B, fixturesPath string) {
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

	worker := che.NewEventfilterIntervalsWorker(dbClient, config.NewTimezoneConfigProvider(cfg, zerolog.Nop()), time.Minute, zerolog.Nop())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		worker.Work(ctx)
	}
}
