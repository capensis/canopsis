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
	b.Helper()
	ctx := b.Context()

	s, err := getAlarmStoreWithFixtures(ctx, b, fixturesPath)
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}

	for b.Loop() {
		_, err := s.Find(ctx, request, userID)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}

const benchAlarmID = "bench-alarm-fixture"
const detailsFixture = "./testdata/fixtures/get_details_steps.yml"

func BenchmarkStore_GetDetails_GivenGroupedStepsReversed(b *testing.B) {
	benchmarkStoreGetDetails(b, detailsFixture, detailsRequest(true, true, 1))
}

func BenchmarkStore_GetDetails_GivenGroupedStepsForward(b *testing.B) {
	benchmarkStoreGetDetails(b, detailsFixture, detailsRequest(false, true, 1))
}

func BenchmarkStore_GetDetails_GivenUngroupedStepsReversed(b *testing.B) {
	benchmarkStoreGetDetails(b, detailsFixture, detailsRequest(true, false, 1))
}

func BenchmarkStore_GetDetails_GivenGroupedStepsReversedDeepPage(b *testing.B) {
	benchmarkStoreGetDetails(b, detailsFixture, detailsRequest(true, true, 200))
}

// detailsRequest builds a DetailsRequest matching the production payload shape.
// reversed and group control the steps query; page sets the steps page number.
func detailsRequest(reversed, group bool, page int64) DetailsRequest {
	opened := true
	r := DetailsRequest{
		ID:                 benchAlarmID,
		Opened:             &opened,
		WithInstructions:   true,
		WithDeclareTickets: true,
		WithDependencies:   true,
		Steps: &StepsRequest{
			Query: pagination.Query{
				Page:  page,
				Limit: 10,
			},
			Reversed: reversed,
			Group:    group,
		},
		Children: &ChildDetailsRequest{
			Query: pagination.Query{
				Page:  1,
				Limit: 10,
			},
			SortRequest: SortRequest{
				SortBy: "v.last_event_date",
				Sort:   "desc",
			},
		},
	}
	r.Format()
	return r
}

func benchmarkStoreGetDetails(b *testing.B, fixturesPath string, request DetailsRequest) {
	b.Helper()
	ctx := b.Context()

	s, err := getAlarmStoreWithFixtures(ctx, b, fixturesPath)
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}

	for b.Loop() {
		_, err := s.GetDetails(ctx, request, "")
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}

func getAlarmStoreWithFixtures(ctx context.Context, b *testing.B, fixturesPath string) (Store, error) {
	dbClient, err := mongo.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	loader := fixtures.NewLoader(dbClient, []string{fixturesPath},
		fixtures.NewParser(fixtures.NewFaker(password.NewBcryptEncoder())), zerolog.Nop())
	err = loader.Load(ctx)
	if err != nil {
		return nil, err
	}

	authorProvider := author.NewProvider(config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	s := NewStore(dbClient, dbClient, nil, common.NewPatternFieldsTransformer(dbClient), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()),
		authorProvider, nil, json.NewDecoder(), zerolog.Nop())

	b.Cleanup(func() {
		err := loader.Clean(context.WithoutCancel(ctx))
		if err != nil {
			b.Errorf("unexpected error %v", err)
		}
		err = dbClient.Disconnect(context.WithoutCancel(ctx))
		if err != nil {
			b.Errorf("unexpected error %v", err)
		}
	})

	return s, nil
}
