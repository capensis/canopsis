package pbehavior_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	mock_pbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/pbehavior"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestServiceComputerAndResolver(t *testing.T) {
	ctx := t.Context()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockResolver := mock_pbehavior.NewMockComputedEntityTypeResolver(ctrl)
	mockResolver.EXPECT().Resolve(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, entity types.Entity, _ time.Time) (pbehavior.ResolveResult, error) {
		switch entity.ID {
		case "parent-1":
			return pbehavior.ResolveResult{Type: pbehavior.Type{Name: "test-type-parent-1", Priority: 10}, Inherited: true}, nil
		case "parent-2":
			return pbehavior.ResolveResult{Type: pbehavior.Type{Name: "test-type-parent-2", Priority: 20}, Inherited: true}, nil
		case "parent-3":
			return pbehavior.ResolveResult{Type: pbehavior.Type{Name: "test-type-parent-3", Priority: 30}, Inherited: true}, nil
		case "parent-4":
			return pbehavior.ResolveResult{Type: pbehavior.Type{Name: "test-type-parent-4", Priority: 40}, Inherited: true}, nil
		case "parent-5":
			return pbehavior.ResolveResult{Type: pbehavior.Type{Name: "test-type-parent-5", Priority: 50}, Inherited: false}, nil
		case "child-1":
			return pbehavior.ResolveResult{}, nil
		case "child-2":
			return pbehavior.ResolveResult{}, nil
		case "child-3":
			return pbehavior.ResolveResult{Type: pbehavior.Type{Name: "test-type-child-3", Priority: 5}}, nil
		case "child-4":
			return pbehavior.ResolveResult{Type: pbehavior.Type{Name: "test-type-child-4", Priority: 100}}, nil
		default:
			return pbehavior.ResolveResult{}, errors.New("unexpected error")
		}
	}).AnyTimes()

	dataSets := []struct {
		name           string
		relations      map[string][]string
		expectedResult map[string]string
		expectedSkip   map[string]bool
	}{
		{
			name: "simple parent child relation",
			relations: map[string][]string{
				"parent-1": {"child-1"},
			},
			expectedResult: map[string]string{
				"parent-1": "test-type-parent-1",
				"child-1":  "test-type-parent-1",
			},
		},
		{
			name: "simple several children relation",
			relations: map[string][]string{
				"parent-1": {"child-1", "child-2"},
			},
			expectedResult: map[string]string{
				"parent-1": "test-type-parent-1",
				"child-1":  "test-type-parent-1",
				"child-2":  "test-type-parent-1",
			},
		},
		{
			name: "several children relation, one child has pbh with lower priority than parent",
			relations: map[string][]string{
				"parent-1": {"child-1", "child-3"},
			},
			expectedResult: map[string]string{
				"parent-1": "test-type-parent-1",
				"child-1":  "test-type-parent-1",
				"child-3":  "test-type-parent-1",
			},
		},
		{
			name: "several children relation, one child has pbh with higher priority than parent",
			relations: map[string][]string{
				"parent-1": {"child-1", "child-4"},
			},
			expectedResult: map[string]string{
				"parent-1": "test-type-parent-1",
				"child-1":  "test-type-parent-1",
				"child-4":  "test-type-child-4",
			},
		},
		{
			name: "several children relation, parent has not inherited pbh",
			relations: map[string][]string{
				"parent-5": {"child-3", "child-4"},
			},
			expectedResult: map[string]string{},
			expectedSkip: map[string]bool{
				"parent-5": true,
			},
		},
		{
			name: "nested graph relation, root service has higher priority",
			relations: map[string][]string{
				"parent-1": {"child-1", "child-2"},
				// kinda graph lookup: child-1, child-2 are parent-1 children and parent-2 is parent-1 parent.
				"parent-2": {"parent-1", "child-1", "child-2"},
			},
			expectedResult: map[string]string{
				"parent-1": "test-type-parent-2",
				"parent-2": "test-type-parent-2",
				"child-1":  "test-type-parent-2",
				"child-2":  "test-type-parent-2",
			},
		},
		{
			name: "nested graph relation, root service has lower priority",
			relations: map[string][]string{
				"parent-2": {"child-1", "child-2"},
				"parent-1": {"parent-2", "child-1", "child-2"},
			},
			expectedResult: map[string]string{
				"parent-1": "test-type-parent-1",
				"parent-2": "test-type-parent-2",
				"child-1":  "test-type-parent-2",
				"child-2":  "test-type-parent-2",
			},
		},
		{
			name: "nested graph relation, mid parent has not inherited pbh",
			relations: map[string][]string{
				"parent-5": {"child-1", "child-2"},
				"parent-1": {"parent-5", "child-1", "child-2"},
			},
			expectedResult: map[string]string{
				"parent-1": "test-type-parent-1",
				"parent-5": "test-type-parent-5",
				"child-1":  "test-type-parent-1",
				"child-2":  "test-type-parent-1",
			},
			expectedSkip: map[string]bool{
				"parent-5": true,
			},
		},
		{
			name: "nested graph relation, mid parent has not inherited pbh, children has their own pbhs",
			relations: map[string][]string{
				"parent-5": {"child-3", "child-4"},
				"parent-1": {"parent-5", "child-3", "child-4"},
			},
			expectedResult: map[string]string{
				"parent-1": "test-type-parent-1",
				"parent-5": "test-type-parent-5",
				"child-3":  "test-type-parent-1",
				"child-4":  "test-type-child-4",
			},
			expectedSkip: map[string]bool{
				"parent-5": true,
			},
		},
		{
			name: "nested graph relation, two parents should set higher priority",
			relations: map[string][]string{
				"parent-3": {"parent-1", "child-1", "child-2"},
				"parent-1": {"child-1", "child-2"},
				"parent-2": {"parent-1", "child-1", "child-2"},
			},
			expectedResult: map[string]string{
				"parent-1": "test-type-parent-3",
				"parent-2": "test-type-parent-2",
				"parent-3": "test-type-parent-3",
				"child-1":  "test-type-parent-3",
				"child-2":  "test-type-parent-3",
			},
		},
		{
			name: "nested graph relation, two parents, should set higher priority",
			relations: map[string][]string{
				"parent-3": {"parent-1", "child-1", "child-2"},
				"parent-1": {"child-1", "child-2"},
				"parent-2": {"parent-1", "child-1", "child-2"},
			},
			expectedResult: map[string]string{
				"parent-1": "test-type-parent-3",
				"parent-2": "test-type-parent-2",
				"parent-3": "test-type-parent-3",
				"child-1":  "test-type-parent-3",
				"child-2":  "test-type-parent-3",
			},
		},
		{
			name: "nested graph relation, two parents, mid parent has not inh pbh, should set higher priority",
			relations: map[string][]string{
				"parent-3": {"parent-5", "child-1", "child-2"},
				"parent-5": {"child-1", "child-2"},
				"parent-2": {"parent-5", "child-1", "child-2"},
			},
			expectedResult: map[string]string{
				"parent-5": "test-type-parent-5",
				"parent-2": "test-type-parent-2",
				"parent-3": "test-type-parent-3",
				"child-1":  "test-type-parent-3",
				"child-2":  "test-type-parent-3",
			},
			expectedSkip: map[string]bool{
				"parent-5": true,
			},
		},
		{
			name: "nested graph relation, two parents, mid parent has not inh pbh, should set higher priority, one child has another parent",
			relations: map[string][]string{
				"parent-3": {"parent-5", "child-1", "child-2"},
				"parent-5": {"child-1", "child-2"},
				"parent-2": {"parent-5", "child-1", "child-2"},
				"parent-4": {"child-1"},
			},
			expectedResult: map[string]string{
				"parent-5": "test-type-parent-5",
				"parent-2": "test-type-parent-2",
				"parent-3": "test-type-parent-3",
				"child-1":  "test-type-parent-4",
				"child-2":  "test-type-parent-3",
			},
			expectedSkip: map[string]bool{
				"parent-5": true,
			},
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			inheritedResolveResult := pbehavior.InheritedServicesPbhResolveResult{
				IDs:          make([]string, 0),
				PersonalPbh:  map[string]pbehavior.ResolveResult{},
				InheritedPbh: map[string]pbehavior.ResolveResult{},
			}

			for parent, children := range dSet.relations {
				parentResolveResult, err := inheritedResolveResult.ResolveParentServicePbh(ctx, mockResolver, types.Entity{ID: parent})
				if err != nil {
					t.Error(fmt.Errorf("error is not expected: %w", err))
				}

				if !parentResolveResult.Inherited {
					if !dSet.expectedSkip[parent] {
						t.Errorf("skip is not expected")
					}

					continue
				}

				for _, child := range children {
					err = inheritedResolveResult.ResolveChildServicePbh(ctx, mockResolver, types.Entity{ID: child}, parentResolveResult)
					if err != nil {
						t.Error(fmt.Errorf("error is not expected: %w", err))
					}
				}
			}

			for id, expectedName := range dSet.expectedResult {
				if expectedName != inheritedResolveResult.PersonalPbh[id].Type.Name {
					t.Errorf("expected pbh type = %s, but got = %s\n", expectedName, inheritedResolveResult.PersonalPbh[id].Type.Name)
				}
			}
		})
	}
}

func BenchmarkInheritedServicesPbhResolve10Services(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/10_services.yml")
}

func BenchmarkInheritedServicesPbhResolve100Services(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/100_services.yml")
}

func BenchmarkInheritedServicesPbhResolve1000Services(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/1000_services.yml")
}

func BenchmarkInheritedServicesPbhResolve10Parents1Child(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/10_parents_1_child.yml")
}

func BenchmarkInheritedServicesPbhResolve100Parents1Child(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/100_parents_1_child.yml")
}

func BenchmarkInheritedServicesPbhResolve1000Parents1Child(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/1000_parents_1_child.yml")
}

func BenchmarkInheritedServicesPbhResolve10Parents1MidChild1Child(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/10_parents_1_midchild_1_child.yml")
}

func BenchmarkInheritedServicesPbhResolve100Parents1MidChild1Child(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/100_parents_1_midchild_1_child.yml")
}

func BenchmarkInheritedServicesPbhResolve1000Parents1MidChild1Child(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/1000_parents_1_midchild_1_child.yml")
}

func BenchmarkInheritedServicesPbhResolve1Parent10Children(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/1_parent_10_children.yml")
}

func BenchmarkInheritedServicesPbhResolve1Parent100Children(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/1_parent_100_children.yml")
}

func BenchmarkInheritedServicesPbhResolve1Parent1000Children(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/1_parent_1000_children.yml")
}

func BenchmarkInheritedServicesPbhResolve1Parent10000Children(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/1_parent_10000_children.yml")
}

func BenchmarkInheritedServicesPbhResolve1Parent100000Children(b *testing.B) {
	benchmarkInheritedRecompute(b, "./testdata/fixtures/1_parent_100000_children.yml")
}

func benchmarkInheritedRecompute(b *testing.B, fixturesPath string) {
	ctx := b.Context()

	dbClient, err := mongo.NewClient(ctx)
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}

	pbhRedisSession, err := redis.NewSession(ctx, redis.PBehaviorLockStorage, zerolog.Nop(), 0, 0)
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}

	b.Cleanup(func() {
		err = pbhRedisSession.FlushDB(context.WithoutCancel(ctx)).Err()
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	})

	b.Cleanup(func() {
		err := dbClient.Disconnect(context.WithoutCancel(ctx))
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
		err := loader.Clean(context.WithoutCancel(ctx))
		if err != nil {
			b.Errorf("unexpected error %v", err)
		}
	})

	now := time.Now()
	newSpan := timespan.New(now, now.Add(time.Hour*2))

	pbhLockerClient := redis.NewLockClient(pbhRedisSession)

	authorProvider := author.NewProvider(config.NewApiConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	pbhService := pbehavior.NewService(
		dbClient,
		pbehavior.NewTypeComputer(pbehavior.NewModelProvider(dbClient, authorProvider), json.NewDecoder()),
		pbehavior.NewStore(pbhRedisSession, json.NewEncoder(), json.NewDecoder()),
		pbhLockerClient,
		zerolog.Nop(),
	)

	resolver, _, err := pbhService.Compute(ctx, newSpan, time.UTC)
	if err != nil {
		b.Errorf("unexpected error %v", err)
	}

	inhResolver := pbehavior.NewInheritedServicePbhResolver(
		dbClient,
		pbehavior.NewEventManager(libevent.NewGenerator(canopsis.PBehaviorConnector, canopsis.PBehaviorConnector)),
		pbehavior.NewStore(pbhRedisSession, json.NewEncoder(), json.NewDecoder()),
		pbhLockerClient,
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err = inhResolver.ComputeAndResolveInheritedServicePbh(ctx, resolver)
		if err != nil {
			b.Errorf("unexpected error %v", err)
		}
	}
}
