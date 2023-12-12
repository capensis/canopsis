package types_test

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAlarmSteps_Crop_GivenLessThenLimitSteps_ShouldDoNothing(t *testing.T) {
	statusStep := types.AlarmStep{
		Type:    types.AlarmStepStatusIncrease,
		Value:   1,
		Author:  "test",
		Message: "test",
	}
	stateIncStep := types.AlarmStep{
		Type:  types.AlarmStepStateIncrease,
		Value: 2,
	}
	stateDecStep := types.AlarmStep{
		Type:  types.AlarmStepStateDecrease,
		Value: 1,
	}
	steps := make(types.AlarmSteps, 0)
	err := steps.Add(statusStep)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	for i := 0; i < 10; i++ {
		err = steps.Add(stateIncStep)
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}
		err = steps.Add(stateDecStep)
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}
	}

	_, update := steps.Crop(&statusStep, 30)
	if update {
		t.Fatalf("expected no update")
	}
}

func TestAlarmSteps_Crop_GivenGreaterThenLimitSteps_ShouldCropSteps(t *testing.T) {
	statusStep := types.AlarmStep{
		Type:    types.AlarmStepStatusIncrease,
		Value:   1,
		Author:  "test",
		Message: "test",
	}
	stateIncStep := types.AlarmStep{
		Type:  types.AlarmStepStateIncrease,
		Value: 2,
	}
	stateDecStep := types.AlarmStep{
		Type:  types.AlarmStepStateDecrease,
		Value: 1,
	}
	steps := make(types.AlarmSteps, 0)
	err := steps.Add(statusStep)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	for i := 0; i < 20; i++ {
		err = steps.Add(stateIncStep)
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}
		err = steps.Add(stateDecStep)
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}
	}

	cropNum := 30
	steps, update := steps.Crop(&statusStep, cropNum)
	if !update {
		t.Fatalf("expected update")
	}
	if len(steps) != cropNum+2 {
		t.Fatalf("expected %d but got %d", cropNum+2, len(steps))
	}
	if steps[0].Type != statusStep.Type {
		t.Fatalf("expected %q but got %s", statusStep.Type, steps[0].Type)
	}
	if steps[1].Type != types.AlarmStepStateCounter {
		t.Fatalf("expected %q but got %s", types.AlarmStepStateCounter, steps[1].Type)
	}
	if steps[1].Author != "test" {
		t.Fatalf("expected %q but got %s", "test", steps[1].Author)
	}
	if steps[1].Message != "test" {
		t.Fatalf("expected %q but got %s", "test", steps[1].Message)
	}
	if steps[1].Value != statusStep.Value {
		t.Fatalf("expected %d but got %d", statusStep.Value, steps[1].Value)
	}
	expectedCounter := types.CropCounter{
		StateChanges:  10,
		Stateinc:      5,
		Statedec:      5,
		StateInfo:     0,
		StateMinor:    5,
		StateMajor:    5,
		StateCritical: 0,
	}
	if diff := pretty.Compare(steps[1].StateCounter, expectedCounter); diff != "" {
		t.Fatalf("unexpected counter %s", diff)
	}
}

func TestAlarmSteps_Crop_GivenCounterStep_ShouldSaveInDB(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbClient, err := mongo.NewClient(ctx, 0, 0, zerolog.Nop())
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	dbCollection := dbClient.Collection(mongo.AlarmMongoCollection)
	statusStep := types.AlarmStep{
		Type:      types.AlarmStepStatusIncrease,
		Value:     1,
		Message:   "coucou",
		Author:    "coucou",
		Timestamp: datetime.NewCpsTime(time.Now().Unix()),
	}
	stateCounter := types.CropCounter{}
	stateCounter.Stateinc = 1
	stateCounter.Statedec = 2
	counterStep := types.AlarmStep{
		Type:         types.AlarmStepStateCounter,
		StateCounter: stateCounter,
		Timestamp:    datetime.NewCpsTime(time.Now().Unix()),
	}

	steps := make(types.AlarmSteps, 0)
	err = steps.Add(statusStep)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}
	err = steps.Add(counterStep)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	alarm := types.Alarm{
		ID: utils.NewID(),
		Value: types.AlarmValue{
			Steps: steps,
		},
	}
	_, err = dbCollection.InsertOne(ctx, alarm)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	defer func() {
		_, err = dbCollection.DeleteOne(ctx, bson.M{"_id": alarm.ID})
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}
	}()

	fetchedAlarm := types.Alarm{}
	err = dbCollection.FindOne(ctx, bson.M{"_id": alarm.ID}).Decode(&fetchedAlarm)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if diff := pretty.Compare(alarm.Value.Steps, fetchedAlarm.Value.Steps); diff != "" {
		t.Errorf("unexpected steps %s", diff)
	}
}
