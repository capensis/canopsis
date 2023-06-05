package eventfilter_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestEnrichmentApplyOnSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedOutcome := eventfilter.OutcomePass
	expectedEvent := types.Event{Resource: "updated"}

	actionProcessor := mock_eventfilter.NewMockActionProcessor(ctrl)
	actionProcessor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedEvent, nil)

	applicator := eventfilter.NewEnrichmentApplicator(eventfilter.NewExternalDataGetterContainer(), actionProcessor)
	resOutcome, resEvent, resError := applicator.Apply(context.Background(), eventfilter.Rule{Config: eventfilter.RuleConfig{Actions: []eventfilter.Action{{}}, OnSuccess: expectedOutcome}}, types.Event{}, eventfilter.RegexMatchWrapper{})
	if resError != nil {
		t.Errorf("expected not error but got %v", resError)
	}

	if resOutcome != expectedOutcome {
		t.Errorf("expected outcome %s, but got %s", expectedOutcome, resOutcome)
	}

	if !reflect.DeepEqual(expectedEvent, resEvent) {
		t.Errorf("expected event %v, but got %v", expectedEvent, resEvent)
	}
}

func TestEnrichmentApplyOnFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedOutcome := eventfilter.OutcomeBreak
	expectedEvent := types.Event{}

	actionProcessor := mock_eventfilter.NewMockActionProcessor(ctrl)
	actionProcessor.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedEvent, fmt.Errorf("error"))

	applicator := eventfilter.NewEnrichmentApplicator(eventfilter.NewExternalDataGetterContainer(), actionProcessor)
	resOutcome, resEvent, resError := applicator.Apply(context.Background(), eventfilter.Rule{Config: eventfilter.RuleConfig{Actions: []eventfilter.Action{{}}, OnFailure: expectedOutcome}}, types.Event{}, eventfilter.RegexMatchWrapper{})
	if resError == nil {
		t.Errorf("expected error but nothing")
	}

	if resOutcome != expectedOutcome {
		t.Errorf("expected outcome %s, but got %s", expectedOutcome, resOutcome)
	}

	if !reflect.DeepEqual(expectedEvent, resEvent) {
		t.Errorf("expected event %v, but got %v", expectedEvent, resEvent)
	}
}

func TestApplyWithExternalData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	getter := mock_eventfilter.NewMockExternalDataGetter(ctrl)
	getter.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(types.Entity{ID: "test_value"}, nil)

	externalDataContainer := eventfilter.NewExternalDataGetterContainer()
	externalDataContainer.Set("test", getter)

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
