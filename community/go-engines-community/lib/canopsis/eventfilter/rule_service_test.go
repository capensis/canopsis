package eventfilter_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestProcessEventSuccess(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adapter := mock_eventfilter.NewMockRuleAdapter(ctrl)
	adapter.EXPECT().GetByTypes(gomock.Any(), gomock.Any()).Return([]eventfilter.Rule{
		{Type: "rule-1", EventPattern: [][]pattern.FieldCondition{{pattern.FieldCondition{
			Field:     "resource",
			Condition: pattern.NewStringCondition("eq", "test resource"),
		}}}},
		{Type: "rule-2", EventPattern: [][]pattern.FieldCondition{{pattern.FieldCondition{
			Field:     "component",
			Condition: pattern.NewStringCondition("eq", "test component"),
		}}}},
	}, nil)

	applicator1 := mock_eventfilter.NewMockRuleApplicator(ctrl)
	applicator1.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ eventfilter.Rule, event types.Event, _ eventfilter.RegexMatchWrapper) (string, types.Event, error) {
			event.Resource = "apply 1"

			return eventfilter.OutcomePass, event, nil
		})
	applicator2 := mock_eventfilter.NewMockRuleApplicator(ctrl)
	applicator2.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ eventfilter.Rule, event types.Event, _ eventfilter.RegexMatchWrapper) (string, types.Event, error) {
			event.Component = "apply 2"

			return eventfilter.OutcomePass, event, nil
		})

	container := mock_eventfilter.NewMockRuleApplicatorContainer(ctrl)
	container.EXPECT().Get(gomock.Any()).AnyTimes().DoAndReturn(func(ruleType string) (eventfilter.RuleApplicator, bool) {
		switch ruleType {
		case "rule-1":
			return applicator1, true
		case "rule-2":
			return applicator2, true
		}

		return nil, false
	})

	ruleService := eventfilter.NewRuleService(adapter, container, zerolog.Logger{})
	err := ruleService.LoadRules(ctx, []string{"rule-1", "rule-2"})
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	expectedEvent := types.Event{
		Resource:  "apply 1",
		Component: "apply 2",
	}

	resultEvent, err := ruleService.ProcessEvent(ctx, types.Event{Resource: "test resource", Component: "test component"})
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

	adapter := mock_eventfilter.NewMockRuleAdapter(ctrl)
	adapter.EXPECT().GetByTypes(gomock.Any(), gomock.Any()).Return([]eventfilter.Rule{
		{Type: "rule-1", EventPattern: [][]pattern.FieldCondition{{pattern.FieldCondition{
			Field:     "resource",
			Condition: pattern.NewStringCondition("eq", "test resource"),
		}}}},
		{Type: "rule-2", EventPattern: [][]pattern.FieldCondition{{pattern.FieldCondition{
			Field:     "component",
			Condition: pattern.NewStringCondition("eq", "test component"),
		}}}},
	}, nil)

	applicator1 := mock_eventfilter.NewMockRuleApplicator(ctrl)
	applicator1.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ eventfilter.Rule, event types.Event, _ eventfilter.RegexMatchWrapper) (string, types.Event, error) {
			event.Resource = "apply 1"

			return eventfilter.OutcomeBreak, event, nil
		})
	applicator2 := mock_eventfilter.NewMockRuleApplicator(ctrl)
	applicator2.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ eventfilter.Rule, event types.Event, _ eventfilter.RegexMatchWrapper) (string, types.Event, error) {
			event.Component = "apply 2"

			return eventfilter.OutcomePass, event, nil
		})

	container := mock_eventfilter.NewMockRuleApplicatorContainer(ctrl)
	container.EXPECT().Get(gomock.Any()).AnyTimes().DoAndReturn(func(ruleType string) (eventfilter.RuleApplicator, bool) {
		switch ruleType {
		case "rule-1":
			return applicator1, true
		case "rule-2":
			return applicator2, true
		}

		return nil, false
	})

	ruleService := eventfilter.NewRuleService(adapter, container, zerolog.Logger{})
	err := ruleService.LoadRules(ctx, []string{"rule-1", "rule-2"})
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	//since first applicator returns break outcome, second applicator should be skipped
	expectedEvent := types.Event{
		Resource:  "apply 1",
		Component: "test component",
	}

	resultEvent, err := ruleService.ProcessEvent(ctx, types.Event{Resource: "test resource", Component: "test component"})
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

	adapter := mock_eventfilter.NewMockRuleAdapter(ctrl)
	adapter.EXPECT().GetByTypes(gomock.Any(), gomock.Any()).Return([]eventfilter.Rule{
		{Type: "rule-1", EventPattern: [][]pattern.FieldCondition{{pattern.FieldCondition{
			Field:     "resource",
			Condition: pattern.NewStringCondition("eq", "test resource"),
		}}}},
		{Type: "rule-2", EventPattern: [][]pattern.FieldCondition{{pattern.FieldCondition{
			Field:     "component",
			Condition: pattern.NewStringCondition("eq", "test component"),
		}}}},
	}, nil)

	applicator1 := mock_eventfilter.NewMockRuleApplicator(ctrl)
	applicator1.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ eventfilter.Rule, event types.Event, _ eventfilter.RegexMatchWrapper) (string, types.Event, error) {
			event.Resource = "apply 1"

			return eventfilter.OutcomeDrop, event, nil
		})
	applicator2 := mock_eventfilter.NewMockRuleApplicator(ctrl)
	applicator2.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ eventfilter.Rule, event types.Event, _ eventfilter.RegexMatchWrapper) (string, types.Event, error) {
			event.Component = "apply 2"

			return eventfilter.OutcomePass, event, nil
		})

	container := mock_eventfilter.NewMockRuleApplicatorContainer(ctrl)
	container.EXPECT().Get(gomock.Any()).AnyTimes().DoAndReturn(func(ruleType string) (eventfilter.RuleApplicator, bool) {
		switch ruleType {
		case "rule-1":
			return applicator1, true
		case "rule-2":
			return applicator2, true
		}

		return nil, false
	})

	ruleService := eventfilter.NewRuleService(adapter, container, zerolog.Logger{})
	err := ruleService.LoadRules(ctx, []string{"rule-1", "rule-2"})
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	//since first applicator returns drop outcome, second applicator should be skipped
	expectedEvent := types.Event{
		Resource:  "apply 1",
		Component: "test component",
	}

	resultEvent, err := ruleService.ProcessEvent(ctx, types.Event{Resource: "test resource", Component: "test component"})
	if !errors.Is(err, eventfilter.ErrDropOutcome) {
		t.Errorf("expected error %v, but got %v", eventfilter.ErrDropOutcome, err)
	}

	if !reflect.DeepEqual(expectedEvent, resultEvent) {
		t.Errorf("expected event %v, but got %v", expectedEvent, resultEvent)
	}
}
