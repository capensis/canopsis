package eventfilter_test

import (
	"context"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestChangeEntityApply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	applicator := eventfilter.NewChangeEntityApplicator(eventfilter.NewExternalDataGetterContainer(), tplExecutor)

	var dataSets = []struct {
		testName      string
		rule          eventfilter.Rule
		event         types.Event
		expectedEvent types.Event
		regexMatches  eventfilter.RegexMatch
	}{
		{
			testName: "given event and rule, resource should be changed",
			rule: eventfilter.Rule{
				Config: eventfilter.RuleConfig{
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
			rule: eventfilter.Rule{
				Config: eventfilter.RuleConfig{
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
			rule: eventfilter.Rule{
				Config: eventfilter.RuleConfig{
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
			rule: eventfilter.Rule{
				Config: eventfilter.RuleConfig{
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
			rule: eventfilter.Rule{
				Config: eventfilter.RuleConfig{
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
			rule: eventfilter.Rule{
				Config: eventfilter.RuleConfig{
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
			rule: eventfilter.Rule{
				Config: eventfilter.RuleConfig{
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
			rule: eventfilter.Rule{
				Config: eventfilter.RuleConfig{
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
			rule: eventfilter.Rule{
				Config: eventfilter.RuleConfig{
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
			regexMatches: eventfilter.RegexMatch{
				EventRegexMatches: pattern.EventRegexMatches{
					ExtraInfos: map[string]pattern.RegexMatches{
						"data": map[string]string{
							"match": "new value",
						},
					},
				},
			},
		},
		{
			testName: "given event and rule, component should be changed by regexMatches template",
			rule: eventfilter.Rule{
				Config: eventfilter.RuleConfig{
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
			regexMatches: eventfilter.RegexMatch{
				EventRegexMatches: pattern.EventRegexMatches{
					ExtraInfos: map[string]pattern.RegexMatches{
						"data": map[string]string{
							"match": "new value",
						},
					},
				},
			},
		},
		{
			testName: "given event and rule, connector should be changed by regexMatches template",
			rule: eventfilter.Rule{
				Config: eventfilter.RuleConfig{
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
			regexMatches: eventfilter.RegexMatch{
				EventRegexMatches: pattern.EventRegexMatches{
					ExtraInfos: map[string]pattern.RegexMatches{
						"data": map[string]string{
							"match": "new value",
						},
					},
				},
			},
		},
		{
			testName: "given event and rule, connector_name should be changed by regexMatches template",
			rule: eventfilter.Rule{
				Config: eventfilter.RuleConfig{
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
			regexMatches: eventfilter.RegexMatch{
				EventRegexMatches: pattern.EventRegexMatches{
					ExtraInfos: map[string]pattern.RegexMatches{
						"data": map[string]string{
							"match": "new value",
						},
					},
				},
			},
		},
	}

	for _, dataSet := range dataSets {
		t.Run(dataSet.testName, func(t *testing.T) {
			outcome, resultEvent, err := applicator.Apply(context.Background(), dataSet.rule, dataSet.event, eventfilter.RegexMatchWrapper{
				BackwardCompatibility: false,
				RegexMatch:            dataSet.regexMatches,
			})

			if err != nil {
				t.Errorf("expected not error but got %v", err)
			}

			if outcome != eventfilter.OutcomePass {
				t.Errorf("expected outcome %s, but got %s", eventfilter.OutcomePass, outcome)
			}

			if !reflect.DeepEqual(dataSet.expectedEvent, resultEvent) {
				t.Errorf("expected event %v, but got %v", dataSet.expectedEvent, resultEvent)
			}
		})
	}
}

func TestChangeEntityApplyWithExternalData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetter := mock_eventfilter.NewMockExternalDataGetter(ctrl)
	mockGetter.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(types.Entity{ID: "test_value"}, nil)

	externalDataContainer := eventfilter.NewExternalDataGetterContainer()
	externalDataContainer.Set("test", mockGetter)

	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))

	applicator := eventfilter.NewChangeEntityApplicator(externalDataContainer, tplExecutor)

	externalData := make(map[string]eventfilter.ExternalDataParameters)
	externalData["test"] = eventfilter.ExternalDataParameters{
		Type: "test",
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
		eventfilter.Rule{
			ExternalData: externalData,
			Config: eventfilter.RuleConfig{
				Resource: "{{.ExternalData.test.ID}}",
			},
		},
		event,
		eventfilter.RegexMatchWrapper{},
	)

	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if outcome != eventfilter.OutcomePass {
		t.Errorf("expected outcome %s, but got %s", eventfilter.OutcomePass, outcome)
	}

	if !reflect.DeepEqual(expectedEvent, resultEvent) {
		t.Errorf("expected event %v, but got %v", expectedEvent, resultEvent)
	}
}
