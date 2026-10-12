package eventfilter_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
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
		gomock.Any(),
	).DoAndReturn(func(
		_ context.Context,
		_ eventfilter.ParsedRule,
		_ eventfilter.ParsedAction,
		event *types.Event,
		_ map[string]eventfilter.UpdatedValue,
		_ eventfilter.RegexMatch,
		_ map[string]any,
	) (bool, error) {
		event.Resource = "updated"
		return false, nil
	})

	applicator := eventfilter.NewEnrichmentApplicator(mockActionProcessor, mockFailureService)

	event := types.Event{}
	res, resError := applicator.Apply(
		t.Context(),
		eventfilter.ParsedRule{Config: eventfilter.ParsedRuleConfig{Actions: []eventfilter.ParsedAction{{}}, OnSuccess: expectedOutcome}},
		&event,
		nil,
		eventfilter.RegexMatch{},
		nil,
	)
	if resError != nil {
		t.Errorf("expected not error but got %v", resError)
	}

	if res.Outcome != expectedOutcome {
		t.Errorf("expected outcome %s, but got %s", expectedOutcome, res.Outcome)
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
		gomock.Any(),
	).DoAndReturn(func(
		_ context.Context,
		_ eventfilter.ParsedRule,
		_ eventfilter.ParsedAction,
		_ *types.Event,
		_ map[string]eventfilter.UpdatedValue,
		_ eventfilter.RegexMatch,
		_ map[string]any,
	) (map[string]eventfilter.UpdatedValue, error) {
		return nil, errors.New("error")
	})

	event := types.Event{}
	applicator := eventfilter.NewEnrichmentApplicator(mockActionProcessor, mockFailureService)
	res, resError := applicator.Apply(t.Context(), eventfilter.ParsedRule{Config: eventfilter.ParsedRuleConfig{Actions: []eventfilter.ParsedAction{{}}, OnFailure: expectedOutcome}}, &event, nil, eventfilter.RegexMatch{}, nil)
	if resError == nil {
		t.Errorf("expected error but nothing")
	}

	if res.Outcome != expectedOutcome {
		t.Errorf("expected outcome %s, but got %s", expectedOutcome, res.Outcome)
	}

	if !reflect.DeepEqual(expectedEvent, event) {
		t.Errorf("expected event %v, but got %v", expectedEvent, event)
	}
}
