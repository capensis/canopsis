package link_test

import (
	"context"
	"strconv"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"github.com/rs/zerolog"
)

func BenchmarkGenerator_GenerateForAlarms_GivenRulesWithLinks(b *testing.B) {
	benchmarkGeneratorGenerateForAlarms(b, "./testdata/fixtures/rules_with_links.yml")
}

func BenchmarkGenerator_GenerateForAlarms_GivenRulesWithSourceCode(b *testing.B) {
	benchmarkGeneratorGenerateForAlarms(b, "./testdata/fixtures/rules_with_source_code.yml")
}

func benchmarkGeneratorGenerateForAlarms(
	b *testing.B,
	fixturesPath string,
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
	generator := link.NewGenerator(dbClient, template.NewExecutor(config.NewTemplateConfigProvider(cfg, zerolog.Nop()),
		config.NewTimezoneConfigProvider(cfg, zerolog.Nop())), zerolog.Nop())
	user := link.User{}
	ids := make([]string, 100)
	for i := 0; i < len(ids); i++ {
		ids[i] = "test-alarm-" + strconv.Itoa(i+1)
	}

	err = generator.Load(ctx)
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := generator.GenerateForAlarms(ctx, ids, user)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}
