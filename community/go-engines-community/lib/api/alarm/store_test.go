package alarm

import (
	"context"
	"testing"

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
	benchmarkStoreFind(b, "./testdata/fixtures/bookmarks_filter.yml", ListRequestWithPagination{
		Query: pagination.Query{
			Page:     1,
			Limit:    100,
			Paginate: true,
		},
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{},
			},
		},
	}, "")
}

func BenchmarkStore_Find_GivenRequestWithBookmarksFilterWithUser(b *testing.B) {
	benchmarkStoreFind(b, "./testdata/fixtures/bookmarks_filter.yml", ListRequestWithPagination{
		Query: pagination.Query{
			Page:     1,
			Limit:    100,
			Paginate: true,
		},
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{},
			},
		},
	}, "user_500")
}

func BenchmarkStore_Find_GivenRequestWithIncludeInstructionsFilter(b *testing.B) {
	instructionFilterType := instrFilterHasInstructions
	instructionType := InstructionTypeManual

	benchmarkStoreFind(b, "./testdata/fixtures/include_instructions_filter.yml", ListRequestWithPagination{
		Query: pagination.Query{
			Page:     1,
			Limit:    100,
			Paginate: true,
		},
		ListRequest: ListRequest{
			FilterRequest: FilterRequest{
				BaseFilterRequest: BaseFilterRequest{
					InstructionFilterType: &instructionFilterType,
					InstructionType:       &instructionType,
				},
			},
		},
	}, "test")
}

func benchmarkStoreFind(b *testing.B, fixturesPath string, request ListRequestWithPagination, userID string) {
	ctx := b.Context()

	dbClient, err := mongo.NewClient(ctx)
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}
	b.Cleanup(func() {
		err := dbClient.Disconnect(context.WithoutCancel(ctx))
		if err != nil {
			b.Errorf("unexpected error %v", err)
		}
	})

	loader := fixtures.NewLoader(dbClient, []string{fixturesPath},
		fixtures.NewParser(fixtures.NewFaker(password.NewBcryptEncoder())), zerolog.Nop())
	err = loader.Load(ctx)
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}
	b.Cleanup(func() {
		err := loader.Clean(context.WithoutCancel(ctx))
		if err != nil {
			b.Errorf("unexpected error %v", err)
		}
	})

	authorProvider := author.NewProvider(config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	s := NewStore(dbClient, dbClient, nil, config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()),
		authorProvider, nil, json.NewDecoder(), zerolog.Nop())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := s.Find(ctx, request, userID)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}
