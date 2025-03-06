package eventfilter_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
	mock_externaldata "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/externaldata"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestEnrichmentApplyOnSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedOutcome := eventfilter.OutcomePass
	expectedEvent := types.Event{Resource: "updated"}

	mockFailureService := mock_eventfilter.NewMockFailureService(ctrl)
	mockActionProcessor := mock_eventfilter.NewMockActionProcessor(ctrl)
	mockActionProcessor.EXPECT().Process(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).DoAndReturn(func(
		ctx context.Context,
		ruleID string,
		action eventfilter.ParsedAction,
		event *types.Event,
		regexMatch eventfilter.RegexMatch,
		externalData map[string]interface{},
	) (bool, error) {
		event.Resource = "updated"
		return false, nil
	})

	applicator := eventfilter.NewEnrichmentApplicator(externaldata.NewGetterContainer(), mockActionProcessor, mockFailureService)

	event := types.Event{}
	resOutcome, _, _, resError := applicator.Apply(
		t.Context(),
		eventfilter.ParsedRule{Config: eventfilter.ParsedRuleConfig{Actions: []eventfilter.ParsedAction{{}}, OnSuccess: expectedOutcome}},
		&event,
		eventfilter.RegexMatch{})
	if resError != nil {
		t.Errorf("expected not error but got %v", resError)
	}

	if resOutcome != expectedOutcome {
		t.Errorf("expected outcome %s, but got %s", expectedOutcome, resOutcome)
	}

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected event %v, but got %v", expectedEvent, event)
	}
}

func TestEnrichmentApplyOnFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedOutcome := eventfilter.OutcomeBreak
	expectedEvent := types.Event{}

	mockFailureService := mock_eventfilter.NewMockFailureService(ctrl)
	mockActionProcessor := mock_eventfilter.NewMockActionProcessor(ctrl)
	mockActionProcessor.EXPECT().Process(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).DoAndReturn(func(
		ctx context.Context,
		ruleID string,
		action eventfilter.ParsedAction,
		event *types.Event,
		regexMatch eventfilter.RegexMatch,
		externalData map[string]interface{},
	) (bool, error) {
		return false, errors.New("error")
	})

	event := types.Event{}
	applicator := eventfilter.NewEnrichmentApplicator(externaldata.NewGetterContainer(), mockActionProcessor, mockFailureService)
	resOutcome, _, _, resError := applicator.Apply(t.Context(), eventfilter.ParsedRule{Config: eventfilter.ParsedRuleConfig{Actions: []eventfilter.ParsedAction{{}}, OnFailure: expectedOutcome}}, &event, eventfilter.RegexMatch{})
	if resError == nil {
		t.Errorf("expected error but nothing")
	}

	if resOutcome != expectedOutcome {
		t.Errorf("expected outcome %s, but got %s", expectedOutcome, resOutcome)
	}

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected event %v, but got %v", expectedEvent, event)
	}
}

func TestApplyWithExternalData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetter := mock_externaldata.NewMockGetter(ctrl)
	mockGetter.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(map[string]any{"ID": "test_value"}, nil)

	externalDataContainer := externaldata.NewGetterContainer()
	externalDataContainer.Set("test", mockGetter)
	mockFailureService := mock_eventfilter.NewMockFailureService(ctrl)
	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}, zerolog.Nop()), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))

	applicator := eventfilter.NewChangeEntityApplicator(externalDataContainer, mockFailureService, tplExecutor)

	externalData := []externaldata.ParsedRefParameters{
		{
			Reference: "test",
			Type:      "test",
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

	outcome, _, _, err := applicator.Apply(
		t.Context(),
		eventfilter.ParsedRule{
			ExternalData: externalData,
			Config: eventfilter.ParsedRuleConfig{
				Resource: tplExecutor.Parse("{{.ExternalData.test.ID}}"),
			},
		},
		&event,
		eventfilter.RegexMatch{},
	)

	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if outcome != eventfilter.OutcomePass {
		t.Errorf("expected outcome %s, but got %s", eventfilter.OutcomePass, outcome)
	}

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected event %v, but got %v", expectedEvent, event)
	}
}
