package contextgraph

import (
	"context"
	"slices"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"github.com/rs/zerolog"
	"golang.org/x/sync/singleflight"
)

type testEntityServiceStorage struct {
	services []entityservice.EntityService
	err      error
}

func (s testEntityServiceStorage) GetAll(context.Context) ([]entityservice.EntityService, error) {
	return s.services, s.err
}

func (s testEntityServiceStorage) Get(context.Context, string) (entityservice.EntityService, error) {
	return entityservice.EntityService{}, nil
}

func TestManagerLoadServicesBuildsIndexes(t *testing.T) {
	services := []entityservice.EntityService{
		{
			Entity: entityservice.EntityService{}.Entity,
			EntityPatternFields: savedpattern.EntityPatternFields{
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
		},
		{
			EntityPatternFields: savedpattern.EntityPatternFields{
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
		},
		{
			EntityPatternFields: savedpattern.EntityPatternFields{
				EntityPattern: [][]pattern.FieldCondition{{
					{
						Field:     "infos.team",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "sre"),
					},
				}},
			},
		},
	}

	services[0].ID = "serv-env-team"
	services[1].ID = "serv-env-owner"
	services[2].ID = "serv-team-only"

	manager := &manager{
		storage: testEntityServiceStorage{services: services},
		logger:  zerolog.Nop(),
		sfGroup: &singleflight.Group{},
	}

	if err := manager.LoadServices(t.Context()); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got, want := len(manager.services), 3; got != want {
		t.Fatalf("expected %d services in cache, got %d", want, got)
	}

	if _, ok := manager.services["serv-env-team"]; !ok {
		t.Fatal("expected serv-env-team in services cache")
	}

	if _, ok := manager.services["serv-env-owner"]; !ok {
		t.Fatal("expected serv-env-owner in services cache")
	}

	if _, ok := manager.services["serv-team-only"]; !ok {
		t.Fatal("expected serv-team-only in services cache")
	}

	if got, want := len(manager.servicesByComponentInfos), 1; got != want {
		t.Fatalf("expected %d component infos keys, got %d", want, got)
	}

	envServices := append([]string(nil), manager.servicesByComponentInfos["env"]...)
	slices.Sort(envServices)
	if !slices.Equal(envServices, []string{"serv-env-owner", "serv-env-team"}) {
		t.Fatalf("unexpected services for component_infos.env: %v", envServices)
	}

	if got, want := len(manager.servicesByInfos), 2; got != want {
		t.Fatalf("expected %d infos keys, got %d", want, got)
	}

	ownerServices := append([]string(nil), manager.servicesByInfos["owner"]...)
	slices.Sort(ownerServices)
	if !slices.Equal(ownerServices, []string{"serv-env-owner"}) {
		t.Fatalf("unexpected services for infos.owner: %v", ownerServices)
	}

	teamServices := append([]string(nil), manager.servicesByInfos["team"]...)
	slices.Sort(teamServices)
	if !slices.Equal(teamServices, []string{"serv-env-team", "serv-team-only"}) {
		t.Fatalf("unexpected services for infos.team: %v", teamServices)
	}
}
