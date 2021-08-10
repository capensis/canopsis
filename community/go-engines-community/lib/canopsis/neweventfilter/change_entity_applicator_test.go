package neweventfilter_test

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/neweventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestApply(t *testing.T) {
	factories := make(map[string]eventfilter.DataSourceFactory)
	applicator := neweventfilter.NewChangeEntityApplicator(factories)

	var dataSets = []struct {
		testName      string
		rule          neweventfilter.Rule
		event         types.Event
		expectedEvent types.Event
		regexMatches  pattern.EventRegexMatches
	}{
		{
			testName: "given event and rule, resource should be changed",
			rule: neweventfilter.Rule{
				Config: neweventfilter.RuleConfig{
					Resource: "new value",
				},
			},
			event: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
			},
			expectedEvent: types.Event{
				Resource:      "new value",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
			},
		},
		{
			testName: "given event and rule, component should be changed",
			rule: neweventfilter.Rule{
				Config: neweventfilter.RuleConfig{
					Component: "new value",
				},
			},
			event: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
			},
			expectedEvent: types.Event{
				Resource:      "resource",
				Component:     "new value",
				Connector:     "connector",
				ConnectorName: "connector name",
			},
		},
		{
			testName: "given event and rule, connector should be changed",
			rule: neweventfilter.Rule{
				Config: neweventfilter.RuleConfig{
					Connector: "new value",
				},
			},
			event: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
			},
			expectedEvent: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "new value",
				ConnectorName: "connector name",
			},
		},
		{
			testName: "given event and rule, connector_name should be changed",
			rule: neweventfilter.Rule{
				Config: neweventfilter.RuleConfig{
					ConnectorName: "new value",
				},
			},
			event: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
			},
			expectedEvent: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "new value",
			},
		},
		{
			testName: "given event and rule, resource should be changed by template",
			rule: neweventfilter.Rule{
				Config: neweventfilter.RuleConfig{
					Resource: "{{.Event.Output}}",
				},
			},
			event: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
				Output:        "new value",
			},
			expectedEvent: types.Event{
				Resource:      "new value",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
				Output:        "new value",
			},
		},
		{
			testName: "given event and rule, component should be changed by template",
			rule: neweventfilter.Rule{
				Config: neweventfilter.RuleConfig{
					Component: "{{.Event.Output}}",
				},
			},
			event: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
				Output:        "new value",
			},
			expectedEvent: types.Event{
				Resource:      "resource",
				Component:     "new value",
				Connector:     "connector",
				ConnectorName: "connector name",
				Output:        "new value",
			},
		},
		{
			testName: "given event and rule, connector should be changed by template",
			rule: neweventfilter.Rule{
				Config: neweventfilter.RuleConfig{
					Connector: "{{.Event.Output}}",
				},
			},
			event: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
				Output:        "new value",
			},
			expectedEvent: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "new value",
				ConnectorName: "connector name",
				Output:        "new value",
			},
		},
		{
			testName: "given event and rule, connector_name should be changed by template",
			rule: neweventfilter.Rule{
				Config: neweventfilter.RuleConfig{
					ConnectorName: "{{.Event.Output}}",
				},
			},
			event: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
				Output:        "new value",
			},
			expectedEvent: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "new value",
				Output:        "new value",
			},
		},
		{
			testName: "given event and rule, resource should be changed by regexMatches template",
			rule: neweventfilter.Rule{
				Config: neweventfilter.RuleConfig{
					Resource: "{{.RegexMatch.ExtraInfos.data.match}}",
				},
			},
			event: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
			},
			expectedEvent: types.Event{
				Resource:      "new value",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
			},
			regexMatches: pattern.EventRegexMatches{
				ExtraInfos: map[string]pattern.RegexMatches{
					"data": map[string]string{
						"match": "new value",
					},
				},
			},
		},
		{
			testName: "given event and rule, component should be changed by regexMatches template",
			rule: neweventfilter.Rule{
				Config: neweventfilter.RuleConfig{
					Component: "{{.RegexMatch.ExtraInfos.data.match}}",
				},
			},
			event: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
			},
			expectedEvent: types.Event{
				Resource:      "resource",
				Component:     "new value",
				Connector:     "connector",
				ConnectorName: "connector name",
			},
			regexMatches: pattern.EventRegexMatches{
				ExtraInfos: map[string]pattern.RegexMatches{
					"data": map[string]string{
						"match": "new value",
					},
				},
			},
		},
		{
			testName: "given event and rule, connector should be changed by regexMatches template",
			rule: neweventfilter.Rule{
				Config: neweventfilter.RuleConfig{
					Connector: "{{.RegexMatch.ExtraInfos.data.match}}",
				},
			},
			event: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
			},
			expectedEvent: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "new value",
				ConnectorName: "connector name",
			},
			regexMatches: pattern.EventRegexMatches{
				ExtraInfos: map[string]pattern.RegexMatches{
					"data": map[string]string{
						"match": "new value",
					},
				},
			},
		},
		{
			testName: "given event and rule, connector_name should be changed by regexMatches template",
			rule: neweventfilter.Rule{
				Config: neweventfilter.RuleConfig{
					ConnectorName: "{{.RegexMatch.ExtraInfos.data.match}}",
				},
			},
			event: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "connector name",
			},
			expectedEvent: types.Event{
				Resource:      "resource",
				Component:     "component",
				Connector:     "connector",
				ConnectorName: "new value",
			},
			regexMatches: pattern.EventRegexMatches{
				ExtraInfos: map[string]pattern.RegexMatches{
					"data": map[string]string{
						"match": "new value",
					},
				},
			},
		},
	}

	for _, dataSet := range dataSets {
		t.Run(dataSet.testName, func(t *testing.T) {
			outcome, resultEvent, err := applicator.Apply(context.Background(), dataSet.rule, dataSet.event, dataSet.regexMatches)

			if err != nil {
				t.Errorf("expected not error but got %v", err)
			}

			if outcome != neweventfilter.OutcomePass {
				t.Errorf("expected outcome %s, but got %s", neweventfilter.OutcomePass, outcome)
			}

			if !reflect.DeepEqual(dataSet.expectedEvent, resultEvent) {
				t.Errorf("expected event %v, but got %v", dataSet.expectedEvent, resultEvent)
			}
		})
	}
}

func TestApplyWithExternalData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	getter := mock_eventfilter.NewMockDataSourceGetter(ctrl)
	getter.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(types.Entity{ID: "test_value"}, nil)

	factory := mock_eventfilter.NewMockDataSourceFactory(ctrl)
	factory.EXPECT().Create(gomock.Any()).Return(getter, nil)
	factories := make(map[string]eventfilter.DataSourceFactory)
	factories["test"] = factory

	applicator := neweventfilter.NewChangeEntityApplicator(factories)

	externalData := make(map[string]eventfilter.DataSource)
	externalData["test"] = eventfilter.DataSource{
		DataSourceBase: eventfilter.DataSourceBase{
			Type: "test",
		},
	}

	event := types.Event{
		Resource:      "resource",
		Component:     "component",
		Connector:     "connector",
		ConnectorName: "connector name",
	}

	expectedEvent := types.Event{
		Resource:      "test_value",
		Component:     "component",
		Connector:     "connector",
		ConnectorName: "connector name",
	}

	outcome, resultEvent, err := applicator.Apply(
		context.Background(),
		neweventfilter.Rule{
			ExternalData: externalData,
			Config: neweventfilter.RuleConfig{
				Resource: "{{.ExternalData.test.ID}}",
			},
		},
		event,
		pattern.EventRegexMatches{},
	)

	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if outcome != neweventfilter.OutcomePass {
		t.Errorf("expected outcome %s, but got %s", neweventfilter.OutcomePass, outcome)
	}

	if !reflect.DeepEqual(expectedEvent, resultEvent) {
		t.Errorf("expected event %v, but got %v", expectedEvent, resultEvent)
	}
}
