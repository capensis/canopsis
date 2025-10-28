package eventfilter_test

import (
	"context"
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestProcessEventSuccess(t *testing.T) {
	ctx := t.Context()

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
	applicator1.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ eventfilter.ParsedRule, event *types.Event, _ map[string]eventfilter.UpdatedValue, _ eventfilter.RegexMatch) (eventfilter.RuleResult, error) {
			event.Resource = "apply 1"

			return eventfilter.RuleResult{Outcome: eventfilter.OutcomePass}, nil
		})
	applicator2 := mock_eventfilter.NewMockRuleApplicator(ctrl)
	applicator2.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ eventfilter.ParsedRule, event *types.Event, _ map[string]eventfilter.UpdatedValue, _ eventfilter.RegexMatch) (eventfilter.RuleResult, error) {
			event.Component = "apply 2"

			return eventfilter.RuleResult{Outcome: eventfilter.OutcomePass}, nil
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

	mockEventCounter := mock_eventfilter.NewMockEventCounter(ctrl)
	mockEventCounter.EXPECT().Add(gomock.Any(), gomock.Any()).AnyTimes()
	mockFailureService := mock_eventfilter.NewMockFailureService(ctrl)
	mockFailureService.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}, zerolog.Nop()), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	ruleService := eventfilter.NewRuleService(adapter, container, mockEventCounter, mockFailureService, tplExecutor, zerolog.Logger{})
	err := ruleService.LoadRules(ctx, []string{"rule-1", "rule-2"})
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	expectedEvent := types.Event{
		Resource:  "apply 1",
		Component: "apply 2",
	}

	event := types.Event{
		Resource:  "test resource",
		Component: "test component",
	}

	_, _, _, err = ruleService.ProcessEvent(ctx, &event)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if diff := pretty.Compare(expectedEvent, event); diff != "" {
		t.Fatalf("unexpected event (-want +got):\n%s", diff)
	}
}

func TestProcessEventBreakOutcome(t *testing.T) {
	ctx := t.Context()

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
	applicator1.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ eventfilter.ParsedRule, event *types.Event, _ map[string]eventfilter.UpdatedValue, _ eventfilter.RegexMatch) (string, bool, map[string]int64, error) {
			event.Resource = "apply 1"

			return eventfilter.OutcomeBreak, false, nil, nil
		})
	applicator2 := mock_eventfilter.NewMockRuleApplicator(ctrl)
	applicator2.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ eventfilter.ParsedRule, event *types.Event, _ map[string]eventfilter.UpdatedValue, _ eventfilter.RegexMatch) (string, bool, map[string]int64, error) {
			event.Component = "apply 2"

			return eventfilter.OutcomePass, false, nil, nil
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

	mockEventCounter := mock_eventfilter.NewMockEventCounter(ctrl)
	mockFailureService := mock_eventfilter.NewMockFailureService(ctrl)
	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}, zerolog.Nop()), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	ruleService := eventfilter.NewRuleService(adapter, container, mockEventCounter, mockFailureService, tplExecutor, zerolog.Logger{})
	err := ruleService.LoadRules(ctx, []string{"rule-1", "rule-2"})
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	//since first applicator returns break outcome, second applicator should be skipped
	expectedEvent := types.Event{
		Resource:  "apply 1",
		Component: "test component",
	}

	event := types.Event{
		Resource:  "test resource",
		Component: "test component",
	}

	_, _, _, err = ruleService.ProcessEvent(ctx, &event)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if diff := pretty.Compare(expectedEvent, event); diff != "" {
		t.Fatalf("unexpected event (-want +got):\n%s", diff)
	}
}

func TestProcessEventDropOutcome(t *testing.T) {
	ctx := t.Context()

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
	applicator1.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ eventfilter.ParsedRule, event *types.Event, _ map[string]eventfilter.UpdatedValue, _ eventfilter.RegexMatch) (eventfilter.RuleResult, error) {
			event.Resource = "apply 1"

			return eventfilter.RuleResult{Outcome: eventfilter.OutcomeDrop}, nil
		})
	applicator2 := mock_eventfilter.NewMockRuleApplicator(ctrl)
	applicator2.EXPECT().Apply(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		DoAndReturn(func(_ context.Context, _ eventfilter.ParsedRule, event *types.Event, _ map[string]eventfilter.UpdatedValue, _ eventfilter.RegexMatch) (eventfilter.RuleResult, error) {
			event.Component = "apply 2"

			return eventfilter.RuleResult{Outcome: eventfilter.OutcomePass}, nil
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

	mockEventCounter := mock_eventfilter.NewMockEventCounter(ctrl)
	mockFailureService := mock_eventfilter.NewMockFailureService(ctrl)
	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}, zerolog.Nop()), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	ruleService := eventfilter.NewRuleService(adapter, container, mockEventCounter, mockFailureService, tplExecutor, zerolog.Logger{})
	err := ruleService.LoadRules(ctx, []string{"rule-1", "rule-2"})
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	//since first applicator returns drop outcome, second applicator should be skipped
	expectedEvent := types.Event{
		Resource:  "apply 1",
		Component: "test component",
	}

	event := types.Event{
		Resource:  "test resource",
		Component: "test component",
	}

	_, _, _, err = ruleService.ProcessEvent(ctx, &event)
	if !errors.Is(err, eventfilter.ErrDropOutcome) {
		t.Errorf("expected error %v, but got %v", eventfilter.ErrDropOutcome, err)
	}

	if diff := pretty.Compare(expectedEvent, event); diff != "" {
		t.Fatalf("unexpected event (-want +got):\n%s", diff)
	}
}
