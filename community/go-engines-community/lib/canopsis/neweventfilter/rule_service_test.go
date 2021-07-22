package neweventfilter_test

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/neweventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_neweventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/neweventfilter"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"reflect"
	"testing"
)

func TestProcessEventSuccess(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := mock_neweventfilter.NewMockRuleAdapter(ctrl)
	adapter.EXPECT().GetAll(gomock.Any()).Return([]neweventfilter.Rule{
		{Type: "rule-1"},
		{Type: "rule-2"},
	}, nil)

	applicator1 := mock_neweventfilter.NewMockRuleApplicator(ctrl)
	applicator1.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ neweventfilter.Rule, event types.Event, _ pattern.EventRegexMatches) (int, types.Event, error) {
			event.Resource = "apply 1"

			return neweventfilter.OutcomePass, event, nil
		})
	applicator2 := mock_neweventfilter.NewMockRuleApplicator(ctrl)
	applicator2.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ neweventfilter.Rule, event types.Event, _ pattern.EventRegexMatches) (int, types.Event, error) {
			event.Component = "apply 2"

			return neweventfilter.OutcomePass, event, nil
		})

	container := mock_neweventfilter.NewMockRuleApplicatorContainer(ctrl)
	container.EXPECT().Get(gomock.Any()).AnyTimes().DoAndReturn(func(ruleType string) (neweventfilter.RuleApplicator, bool) {
		switch ruleType {
		case "rule-1":
			return applicator1, true
		case "rule-2":
			return applicator2, true
		}

		return nil, false
	})

	ruleService := neweventfilter.NewRuleService(adapter, container, zerolog.Logger{})
	err := ruleService.LoadRules(ctx)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	expectedEvent := types.Event{
		Resource:  "apply 1",
		Component: "apply 2",
	}

	resultEvent, err := ruleService.ProcessEvent(ctx, types.Event{})
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if !reflect.DeepEqual(expectedEvent, resultEvent) {
		t.Errorf("expected event %v, but got %v", expectedEvent, resultEvent)
	}
}

func TestProcessEventBreakOutcome(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := mock_neweventfilter.NewMockRuleAdapter(ctrl)
	adapter.EXPECT().GetAll(gomock.Any()).Return([]neweventfilter.Rule{
		{Type: "rule-1"},
		{Type: "rule-2"},
	}, nil)

	applicator1 := mock_neweventfilter.NewMockRuleApplicator(ctrl)
	applicator1.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ neweventfilter.Rule, event types.Event, _ pattern.EventRegexMatches) (int, types.Event, error) {
			event.Resource = "apply 1"

			return neweventfilter.OutcomeBreak, event, nil
		})
	applicator2 := mock_neweventfilter.NewMockRuleApplicator(ctrl)
	applicator2.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ neweventfilter.Rule, event types.Event, _ pattern.EventRegexMatches) (int, types.Event, error) {
			event.Component = "apply 2"

			return neweventfilter.OutcomePass, event, nil
		})

	container := mock_neweventfilter.NewMockRuleApplicatorContainer(ctrl)
	container.EXPECT().Get(gomock.Any()).AnyTimes().DoAndReturn(func(ruleType string) (neweventfilter.RuleApplicator, bool) {
		switch ruleType {
		case "rule-1":
			return applicator1, true
		case "rule-2":
			return applicator2, true
		}

		return nil, false
	})

	ruleService := neweventfilter.NewRuleService(adapter, container, zerolog.Logger{})
	err := ruleService.LoadRules(ctx)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	//since first applicator returns break outcome, second applicator should be skipped
	expectedEvent := types.Event{
		Resource: "apply 1",
	}

	resultEvent, err := ruleService.ProcessEvent(ctx, types.Event{})
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if !reflect.DeepEqual(expectedEvent, resultEvent) {
		t.Errorf("expected event %v, but got %v", expectedEvent, resultEvent)
	}
}

func TestProcessEventDropOutcome(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := mock_neweventfilter.NewMockRuleAdapter(ctrl)
	adapter.EXPECT().GetAll(gomock.Any()).Return([]neweventfilter.Rule{
		{Type: "rule-1"},
		{Type: "rule-2"},
	}, nil)

	applicator1 := mock_neweventfilter.NewMockRuleApplicator(ctrl)
	applicator1.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ neweventfilter.Rule, event types.Event, _ pattern.EventRegexMatches) (int, types.Event, error) {
			event.Resource = "apply 1"

			return neweventfilter.OutcomeDrop, event, nil
		})
	applicator2 := mock_neweventfilter.NewMockRuleApplicator(ctrl)
	applicator2.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ neweventfilter.Rule, event types.Event, _ pattern.EventRegexMatches) (int, types.Event, error) {
			event.Component = "apply 2"

			return neweventfilter.OutcomePass, event, nil
		})

	container := mock_neweventfilter.NewMockRuleApplicatorContainer(ctrl)
	container.EXPECT().Get(gomock.Any()).AnyTimes().DoAndReturn(func(ruleType string) (neweventfilter.RuleApplicator, bool) {
		switch ruleType {
		case "rule-1":
			return applicator1, true
		case "rule-2":
			return applicator2, true
		}

		return nil, false
	})

	ruleService := neweventfilter.NewRuleService(adapter, container, zerolog.Logger{})
	err := ruleService.LoadRules(ctx)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	//since first applicator returns drop outcome, second applicator should be skipped
	expectedEvent := types.Event{
		Resource: "apply 1",
	}

	resultEvent, err := ruleService.ProcessEvent(ctx, types.Event{})
	if !errors.Is(err, neweventfilter.ErrDropOutcome) {
		t.Errorf("expected error %v, but got %v", neweventfilter.ErrDropOutcome, err)
	}

	if !reflect.DeepEqual(expectedEvent, resultEvent) {
		t.Errorf("expected event %v, but got %v", expectedEvent, resultEvent)
	}
}
