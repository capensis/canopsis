package alarm_test

import (
	"context"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"github.com/rs/zerolog"
)

func BenchmarkStore_Find_GivenRequestWithBookmarksFilterWithoutUser(b *testing.B) {
	benchmarkStoreFind(b, "./testdata/fixtures/bookmarks_filter.yml", alarm.ListRequestWithPagination{
		Query: pagination.Query{
			Page:     1,
			Limit:    100,
			Paginate: true,
		},
		ListRequest: alarm.ListRequest{
			FilterRequest: alarm.FilterRequest{
				BaseFilterRequest: alarm.BaseFilterRequest{},
			},
		},
	}, "")
}

func BenchmarkStore_Find_GivenRequestWithBookmarksFilterWithUser(b *testing.B) {
	benchmarkStoreFind(b, "./testdata/fixtures/bookmarks_filter.yml", alarm.ListRequestWithPagination{
		Query: pagination.Query{
			Page:     1,
			Limit:    100,
			Paginate: true,
		},
		ListRequest: alarm.ListRequest{
			FilterRequest: alarm.FilterRequest{
				BaseFilterRequest: alarm.BaseFilterRequest{},
			},
		},
	}, "user_500")
}

func BenchmarkStore_Find_GivenRequestWithIncludeInstructionsFilter(b *testing.B) {
	benchmarkStoreFind(b, "./testdata/fixtures/include_instructions_filter.yml", alarm.ListRequestWithPagination{
		Query: pagination.Query{
			Page:     1,
			Limit:    100,
			Paginate: true,
		},
		ListRequest: alarm.ListRequest{
			FilterRequest: alarm.FilterRequest{
				BaseFilterRequest: alarm.BaseFilterRequest{
					Instructions: []alarm.InstructionFilterRequest{
						{
							IncludeTypes: []int{alarm.InstructionTypeManual},
						},
					},
				},
			},
		},
	}, "test")
}

func BenchmarkStore_Find_GivenRequestWithExcludeInstructionsFilter(b *testing.B) {
	benchmarkStoreFind(b, "./testdata/fixtures/exclude_instructions_filter.yml", alarm.ListRequestWithPagination{
		Query: pagination.Query{
			Page:     1,
			Limit:    100,
			Paginate: true,
		},
		ListRequest: alarm.ListRequest{
			FilterRequest: alarm.FilterRequest{
				BaseFilterRequest: alarm.BaseFilterRequest{
					Instructions: []alarm.InstructionFilterRequest{
						{
							ExcludeTypes: []int{alarm.InstructionTypeManual},
						},
					},
				},
			},
		},
	}, "test")
}

func benchmarkStoreFind(b *testing.B, fixturesPath string, request alarm.ListRequestWithPagination, userID string) {
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

	authorProvider := author.NewProvider(config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	s := alarm.NewStore(dbClient, dbClient, nil, config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()),
		authorProvider, nil, json.NewDecoder(), zerolog.Nop())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := s.Find(ctx, request, userID)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}
