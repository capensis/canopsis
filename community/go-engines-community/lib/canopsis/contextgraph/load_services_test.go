package contextgraph

import (
	"context"
	"slices"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"github.com/rs/zerolog"
	"golang.org/x/sync/singleflight"
)

type testEntityServiceStorage struct {
	services []EntityService
	err      error
}

func (s testEntityServiceStorage) GetAll(_ context.Context) ([]EntityService, error) {
	return s.services, s.err
}

func (s testEntityServiceStorage) Get(_ context.Context, _ string) (entityservice.EntityService, error) {
	return entityservice.EntityService{}, nil
}

func TestManagerLoadServices(t *testing.T) {
	services := []EntityService{
		{
			ID: "serv-env-team",
			EntityPattern: [][]pattern.FieldCondition{{
				{
					Field:     "component_infos.env",
					FieldType: pattern.FieldTypeString,
					Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
				},
				{
					Field:     "infos.team",
					FieldType: pattern.FieldTypeString,
					Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
				},
			}},
		},
		{
			ID: "serv-env-owner",
			EntityPattern: [][]pattern.FieldCondition{{
				{
					Field:     "component_infos.env",
					FieldType: pattern.FieldTypeString,
					Condition: pattern.NewStringCondition(pattern.ConditionEqual, "dev"),
				},
				{
					Field:     "infos.owner",
					FieldType: pattern.FieldTypeString,
					Condition: pattern.NewStringCondition(pattern.ConditionEqual, "ops"),
				},
			}},
		},
		{
			ID: "serv-team-only",
			EntityPattern: [][]pattern.FieldCondition{{
				{
					Field:     "infos.team",
					FieldType: pattern.FieldTypeString,
					Condition: pattern.NewStringCondition(pattern.ConditionEqual, "sre"),
				},
			}},
		},
	}

	m := &manager{
		storage: testEntityServiceStorage{services: services},
		logger:  zerolog.Nop(),
		sfGroup: &singleflight.Group{},
	}

	if err := m.LoadServices(t.Context()); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got, want := len(m.services), 3; got != want {
		t.Fatalf("expected %d services in cache, got %d", want, got)
	}

	if _, ok := m.services["serv-env-team"]; !ok {
		t.Fatal("expected serv-env-team in services cache")
	}

	if _, ok := m.services["serv-env-owner"]; !ok {
		t.Fatal("expected serv-env-owner in services cache")
	}

	if _, ok := m.services["serv-team-only"]; !ok {
		t.Fatal("expected serv-team-only in services cache")
	}

	if got, want := len(m.servicesByComponentInfos), 1; got != want {
		t.Fatalf("expected %d component infos keys, got %d", want, got)
	}

	envServices := serviceIDs(m.servicesByComponentInfos["env"])
	if !slices.Equal(envServices, []string{"serv-env-owner", "serv-env-team"}) {
		t.Fatalf("unexpected services for component_infos.env: %v", envServices)
	}

	if got, want := len(m.servicesByInfos), 2; got != want {
		t.Fatalf("expected %d infos keys, got %d", want, got)
	}

	ownerServices := serviceIDs(m.servicesByInfos["owner"])
	if !slices.Equal(ownerServices, []string{"serv-env-owner"}) {
		t.Fatalf("unexpected services for infos.owner: %v", ownerServices)
	}

	teamServices := serviceIDs(m.servicesByInfos["team"])
	if !slices.Equal(teamServices, []string{"serv-env-team", "serv-team-only"}) {
		t.Fatalf("unexpected services for infos.team: %v", teamServices)
	}
}

func serviceIDs(services map[string]EntityService) []string {
	ids := make([]string, 0, len(services))
	for _, s := range services {
		ids = append(ids, s.ID)
	}

	slices.Sort(ids)

	return ids
}
