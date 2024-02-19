package types_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/kylelemons/godebug/pretty"
)

func TestEvent_Format_GivenEmptyEventType_ShouldSetCheck(t *testing.T) {
	event := getEvent()
	event.EventType = ""
	event.Format()
	if event.EventType != types.EventTypeCheck {
		t.Errorf("expected %q but got %q", types.EventTypeCheck, event.EventType)
	}
}

func TestEvent_Format_GivenEmptyTs_ShouldSetTs(t *testing.T) {
	event := getEvent()
	event.Timestamp = types.CpsTime{}
	event.Format()
	if event.Timestamp.Unix() <= 0 {
		t.Errorf("expected ts but nothing")
	}
}

func TestEvent_IsValid(t *testing.T) {
	dataSet := []struct {
		Event types.Event
		Err   error
	}{
		{
			Event: getEvent(),
		},
		{
			Event: types.Event{
				State:         types.AlarmStateMajor,
				Connector:     "centreon",
				ConnectorName: "centreon",
				Component:     "host",
				Resource:      "nginx",
				SourceType:    types.SourceTypeResource,
			},
			Err: fmt.Errorf("wrong event type: "),
		},
		{
			Event: types.Event{
				EventType:     types.EventTypeCheck,
				State:         types.AlarmStateMajor,
				ConnectorName: "centreon",
				Component:     "host",
				Resource:      "nginx",
				SourceType:    types.SourceTypeResource,
			},
			Err: errors.New("missing connector"),
		},
		{
			Event: types.Event{
				EventType:     types.EventTypeCheck,
				State:         types.AlarmStateMajor,
				Connector:     "centreon",
				ConnectorName: "centreon",
				Component:     "host",
				SourceType:    types.SourceTypeResource,
			},
			Err: errors.New("missing resource"),
		},
		{
			Event: types.Event{
				EventType:     types.EventTypeCheck,
				State:         types.AlarmStateMajor,
				Connector:     "centreon",
				ConnectorName: "centreon",
				Resource:      "nginx",
				SourceType:    types.SourceTypeResource,
			},
			Err: errors.New("missing component"),
		},
	}

	for i, data := range dataSet {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if err := data.Event.IsValid(); data.Err == nil && err != nil || data.Err != nil && err == nil ||
				data.Err != nil && err != nil && data.Err.Error() != err.Error() {
				t.Errorf("expected %q but got %q", data.Err, err)
			}
		})
	}
}

func TestEvent_IsMatched(t *testing.T) {
	event := getEvent()
	event.Component = "foo"
	fields := []string{"Fear", "Resource", "Component"}
	if !event.IsMatched(".*foo", fields) {
		t.Errorf("expected true but got false")
	}
}

func TestEvent_GetRequiredKeys(t *testing.T) {
	event := getEvent()
	expected := []string{
		"_id",
		"event_type",
		"connector",
		"connector_name",
		"component",
		"resource",
		"source_type",
		"status",
		"state",
	}
	result := event.GetRequiredKeys()
	for _, expectedField := range expected {
		matched := true
		for _, field := range result {
			if field == expectedField {
				matched = true
				break
			}
		}
		if !matched {
			t.Errorf("expected %q", expectedField)
		}
	}
}

func TestEvent_InjectExtraInfos(t *testing.T) {
	str := `{
		"event_type":"check",
		"component":"bla",
		"resource":"blurk",
		"state":3,
		"connector":"bla",
		"connector_name":"bla",
		"personnemeconnait":"ulyss31",
		"personnemeconnait2":"ulyss62"
	}`
	event := types.Event{}
	err := json.Unmarshal([]byte(str), &event)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	err = event.InjectExtraInfos([]byte(str))
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if event.ExtraInfos["personnemeconnait"] != "ulyss31" {
		t.Errorf("expected %q but got %q", "ulyss31", event.ExtraInfos["personnemeconnait"])
	}
	if event.ExtraInfos["personnemeconnait2"] != "ulyss62" {
		t.Errorf("expected %q but got %q", "ulyss62", event.ExtraInfos["personnemeconnait2"])
	}
}

func TestEvent_SetField(t *testing.T) {
	str := `{
		"event_type":"check",
		"component":"bla",
		"resource":"blurk",
		"state":3,
		"connector":"bla",
		"connector_name":"bla",
		"extra_info":"ulyss31",
		"tags": {
			"tag3": "value3",
			"tag2": "value2a"
		}
	}`
	event := types.Event{}
	err := json.Unmarshal([]byte(str), &event)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	err = event.InjectExtraInfos([]byte(str))
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	dataSet := []struct {
		Field    string
		Value    any
		Err      error
		Expected any
	}{
		{
			Field: "extra_info",
			Value: 12,
		},
		{
			Field: "extra_info",
			Value: int32(12),
		},
		{
			Field: "extra_info",
			Value: int64(12),
		},
		{
			Field: "new_info",
			Value: "twelve",
		},
		{
			Field: "State",
			Value: 2,
		},
		{
			Field: "State",
			Value: int32(2),
		},
		{
			Field: "State",
			Value: int64(2),
		},
		{
			Field: "State",
			Value: types.CpsNumber(2),
		},
		{
			Field: "Status",
			Value: 2,
		},
		{
			Field: "Status",
			Value: int32(2),
		},
		{
			Field: "Status",
			Value: int64(2),
		},
		{
			Field: "Status",
			Value: types.CpsNumber(2),
		},
		{
			Field: "Timestamp",
			Value: 12,
		},
		{
			Field: "EventType",
			Value: "test",
		},
		{
			Field: "Debug",
			Value: true,
		},
		{
			Field: "Tags",
			Value: map[string]any{
				"tag1": "value1",
				"tag2": "value2",
			},
			Expected: map[string]string{
				"tag1": "value1",
				"tag2": "value2",
			},
		},
	}

	for i, data := range dataSet {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err = event.SetField(data.Field, data.Value)
			if !errors.Is(err, data.Err) {
				t.Fatalf("expected %v but got %v", data.Err, err)
			}
			if err == nil {
				switch data.Field {
				case "State":
					if diff := pretty.Compare(event.State, data.Value); diff != "" {
						t.Errorf("expected %v but got %v", data.Value, event.State)
					}
				case "Status":
					if diff := pretty.Compare(event.Status, data.Value); diff != "" {
						t.Errorf("expected %v but got %v", data.Value, event.Status)
					}
				case "Timestamp":
					if diff := pretty.Compare(event.Timestamp.Unix(), data.Value); diff != "" {
						t.Errorf("expected %v but got %v", data.Value, event.Timestamp.Unix())
					}
				case "EventType":
					if diff := pretty.Compare(event.EventType, data.Value); diff != "" {
						t.Errorf("expected %v but got %v", data.Value, event.EventType)
					}
				case "Debug":
					if diff := pretty.Compare(event.Debug, data.Value); diff != "" {
						t.Errorf("expected %v but got %v", data.Value, event.Debug)
					}
				case "Tags":
					expected := data.Expected
					if expected == nil {
						expected = data.Value
					}
					if diff := pretty.Compare(event.Tags, expected); diff != "" {
						t.Errorf("expected %v but got %v", expected, event.Tags)
					}
				default:
					if diff := pretty.Compare(event.ExtraInfos[data.Field], data.Value); diff != "" {
						t.Errorf("expected %v but got %v", data.Value, event.Status)
					}
				}
			}
		})
	}
}

func getEvent() types.Event {
	return types.Event{
		EventType:     types.EventTypeCheck,
		State:         types.AlarmStateMajor,
		Connector:     "centreon",
		ConnectorName: "centreon",
		Component:     "host",
		Resource:      "nginx",
		SourceType:    types.SourceTypeResource,
	}
}
