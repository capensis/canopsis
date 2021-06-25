package eventfilter_test

import (
	"reflect"
	"testing"
	"text/template"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAction_UnmarshalBSONValue_GivenInvalidData_ShouldReturnError(t *testing.T) {
	dataSets := map[string]bson.M{
		"Given an action without type Should return error":    {},
		"Given an action an invalid type Should return error": {"type": "invalid_type"},
		"Given a set_field action with an unexpected field Should return error": {
			"type":             "set_field",
			"name":             "Output",
			"value":            "output",
			"unexpected_field": "",
		},
		"Given a set_field_from_template action with an unexpected field Should return error": {
			"type":             "set_field_from_template",
			"name":             "Output",
			"value":            "output",
			"unexpected_field": "",
		},
		"Given a set_entity_info_from_template action with an unexpected field Should return error": {
			"type":             "set_entity_info_from_template",
			"name":             "info_name",
			"value":            "info_value",
			"unexpected_field": "",
		},
		"Given a copy action with an unexpected field Should return error": {
			"type":             "copy",
			"from":             "ExternalData.Entity",
			"to":               "Entity",
			"unexpected_field": "",
		},
		"Given a set_entity_info action with an unexpected field Should return error": {
			"type":             "set_entity_info",
			"name":             "Output",
			"value":            "output",
			"unexpected_field": "",
		},
		"Given a copy_to_entity_info action with an unexpected field Should return error": {
			"type":             "copy_to_entity_info",
			"name":             "Output",
			"from":             "Output",
			"unexpected_field": "",
		},
	}

	for testCase, data := range dataSets {
		t.Run(testCase, func(t *testing.T) {
			b, err := bson.Marshal(data)
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}
			var action eventfilter.Action
			err = bson.Unmarshal(b, &action)
			if err == nil {
				t.Errorf("expected no error but got nothing")
			}
		})
	}
}

func TestAction_UnmarshalBSONValue_GivenValidData_ShouldReturnProcessor(t *testing.T) {
	dataSets := map[string]struct {
		data       bson.M
		actionType eventfilter.ActionType
		processor  eventfilter.ActionProcessor
	}{
		"Given a valid set_field action Should return action": {
			data: bson.M{
				"type":  "set_field",
				"name":  "testname",
				"value": "testval",
			},
			actionType: eventfilter.SetField,
			processor: eventfilter.SetFieldProcessor{
				Name: types.OptionalString{
					Set:   true,
					Value: "testname",
				},
				Value: types.OptionalInterface{
					Set:   true,
					Value: "testval",
				},
				UnexpectedFields: nil,
			},
		},
		"Given a valid set_entity_info action Should return action": {
			data: bson.M{
				"type":  "set_entity_info",
				"name":  "testname",
				"value": "testval",
			},
			actionType: eventfilter.SetEntityInfo,
			processor: eventfilter.SetEntityInfoProcessor{
				Name: types.OptionalString{
					Set:   true,
					Value: "testname",
				},
				Value: types.OptionalInterface{
					Set:   true,
					Value: "testval",
				},
				UnexpectedFields: nil,
			},
		},
		"Given a valid copy_to_entity_info action Should return action": {
			data: bson.M{
				"type": "copy_to_entity_info",
				"name": "testname",
				"from": "testfrom",
			},
			actionType: eventfilter.CopyToEntityInfo,
			processor: eventfilter.CopyToEntityInfoProcessor{
				Name: types.OptionalString{
					Set:   true,
					Value: "testname",
				},
				From: types.OptionalString{
					Set:   true,
					Value: "testfrom",
				},
				UnexpectedFields: nil,
			},
		},
		"Given a valid copy action Should return action": {
			data: bson.M{
				"type": "copy",
				"from": "testfrom",
				"to":   "testto",
			},
			actionType: eventfilter.Copy,
			processor: eventfilter.CopyProcessor{
				From: types.OptionalString{
					Set:   true,
					Value: "testfrom",
				},
				To: types.OptionalString{
					Set:   true,
					Value: "testto",
				},
				UnexpectedFields: nil,
			},
		},
	}

	for testCase, data := range dataSets {
		t.Run(testCase, func(t *testing.T) {
			b, err := bson.Marshal(data.data)
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}
			var action eventfilter.Action
			err = bson.Unmarshal(b, &action)
			if err != nil {
				t.Errorf("expected no error but got %v", err)
			}
			if action.Type != data.actionType {
				t.Errorf("expected type %v but got %v", data.actionType, action.Type)
			}
			if !reflect.DeepEqual(action.ActionProcessor, data.processor) {
				t.Errorf("expected processor %+v but got %+v", data.processor, action.ActionProcessor)
			}
		})
	}
}

func TestSetFieldProcessor_Apply_GivenValidData_ShouldApplyIt(t *testing.T) {
	expectedState := 3
	processor := eventfilter.SetFieldProcessor{
		Name:  types.OptionalString{Set: true, Value: "State"},
		Value: types.OptionalInterface{Set: true, Value: expectedState},
	}
	report := eventfilter.Report{}
	event := types.Event{
		State: 1,
	}
	event, err := processor.Apply(event, eventfilter.ActionParameters{}, &report)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}
	if event.State != types.CpsNumber(expectedState) {
		t.Errorf("expected state %v but got %v", expectedState, event.State)
	}
	if report.EntityUpdated {
		t.Errorf("expected unchanged enitity")
	}
}

func TestSetFieldProcessor_Apply_GivenInvalidData_ShouldReturnError(t *testing.T) {
	processor := eventfilter.SetFieldProcessor{
		Name:  types.OptionalString{Set: true, Value: "State"},
		Value: types.OptionalInterface{Set: true, Value: "invalid state"},
	}
	report := eventfilter.Report{}
	event := types.Event{
		State: 1,
	}
	event, err := processor.Apply(event, eventfilter.ActionParameters{}, &report)
	if err == nil {
		t.Errorf("expected error but nothing")
	}
}

func TestSetFieldFromTemplateProcessor_Apply_GivenValidData_ShouldApplyIt(t *testing.T) {
	expectedOutput := "test-output (by test-author)"
	tpl, err := template.New("test").Parse("{{.Event.Output}} (by {{.Event.Author}})")
	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}

	processor := eventfilter.SetFieldFromTemplateProcessor{
		Name: types.OptionalString{Set: true, Value: "Output"},
		Value: types.OptionalTemplate{
			Set:   true,
			Value: tpl,
		},
	}
	report := eventfilter.Report{}
	event := types.Event{}
	event, err = processor.Apply(event, eventfilter.ActionParameters{
		Event: types.Event{
			Output: "test-output",
			Author: "test-author",
		},
	}, &report)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}
	if event.Output != expectedOutput {
		t.Errorf("expected output %v but got %v", expectedOutput, event.Output)
	}
	if report.EntityUpdated {
		t.Errorf("expected unchanged enitity")
	}
}

func TestSetFieldFromTemplateProcessor_Apply_GivenInvalidData_ShouldReturnError(t *testing.T) {
	report := eventfilter.Report{}
	validTpl, err := template.New("test").Parse("{{.Event.Output}} (by {{.Event.Author}})")
	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}
	invalidTpl, err := template.New("test").Parse("{{.Event.UnknownField}}")
	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}

	dataSets := []struct {
		processor eventfilter.SetFieldFromTemplateProcessor
		event     types.Event
	}{
		{
			processor: eventfilter.SetFieldFromTemplateProcessor{
				Name: types.OptionalString{Set: true, Value: "State"},
				Value: types.OptionalTemplate{
					Set:   true,
					Value: validTpl,
				},
			},
			event: types.Event{
				Output: "test-output",
				Author: "test-author",
			},
		},
		{
			processor: eventfilter.SetFieldFromTemplateProcessor{
				Name: types.OptionalString{Set: true, Value: "State"},
				Value: types.OptionalTemplate{
					Set:   true,
					Value: invalidTpl,
				},
			},
			event: types.Event{},
		},
	}

	for _, data := range dataSets {
		_, err := data.processor.Apply(data.event, eventfilter.ActionParameters{}, &report)
		if err == nil {
			t.Errorf("expected error but got nothing")
		}
	}
}

func TestSetEntityInfoFromTemplateProcessor_Apply_GivenValidData_ShouldApplyIt(t *testing.T) {
	expectedCustomer := "test-output (by test-author)"
	tpl, err := template.New("test").Parse("{{.Event.Output}} (by {{.Event.Author}})")
	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}
	processor := eventfilter.SetEntityInfoFromTemplateProcessor{
		Name: types.OptionalString{Set: true, Value: "Customer"},
		Value: types.OptionalTemplate{
			Set:   true,
			Value: tpl,
		},
	}
	report := eventfilter.Report{}
	event := types.Event{
		Entity: &types.Entity{},
	}
	event, err = processor.Apply(event, eventfilter.ActionParameters{
		Event: types.Event{
			Output: "test-output",
			Author: "test-author",
		},
	}, &report)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}
	if event.Entity.Infos["Customer"].Value != expectedCustomer {
		t.Errorf("expected output %v but got %v", expectedCustomer, event.Entity.Infos["Customer"].Value)
	}
	if !report.EntityUpdated {
		t.Errorf("expected changed enitity")
	}
}

func TestSetEntityInfoFromTemplateProcessor_Apply_GivenInvalidData_ShouldReturnError(t *testing.T) {
	report := eventfilter.Report{}
	validTpl, err := template.New("test").Parse("{{.Event.Output}} (by {{.Event.Author}})")
	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}
	invalidTpl, err := template.New("test").Parse("{{.Event.UnknownField}}")
	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}

	dataSets := []struct {
		processor eventfilter.SetEntityInfoFromTemplateProcessor
		event     types.Event
	}{
		{
			processor: eventfilter.SetEntityInfoFromTemplateProcessor{
				Name: types.OptionalString{Set: true, Value: "Customer"},
				Value: types.OptionalTemplate{
					Set:   true,
					Value: invalidTpl,
				},
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
		},
		{
			processor: eventfilter.SetEntityInfoFromTemplateProcessor{
				Name: types.OptionalString{Set: true, Value: "Customer"},
				Value: types.OptionalTemplate{
					Set:   true,
					Value: validTpl,
				},
			},
			event: types.Event{
				Entity: nil,
			},
		},
	}

	for _, data := range dataSets {
		_, err := data.processor.Apply(data.event, eventfilter.ActionParameters{}, &report)
		if err == nil {
			t.Errorf("expected error but got nothing")
		}
	}
}

func TestSetEntityInfoProcessor_Apply_GivenValidData_ShouldApplyIt(t *testing.T) {
	expectedCustomer := "test-output (by test-author)"
	processor := eventfilter.SetEntityInfoProcessor{
		Name:  types.OptionalString{Set: true, Value: "Customer"},
		Value: types.OptionalInterface{Set: true, Value: expectedCustomer},
	}
	report := eventfilter.Report{}
	event := types.Event{
		Entity: &types.Entity{},
	}
	event, err := processor.Apply(event, eventfilter.ActionParameters{}, &report)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}
	if event.Entity.Infos["Customer"].Value != expectedCustomer {
		t.Errorf("expected output %v but got %v", expectedCustomer, event.Output)
	}
	if !report.EntityUpdated {
		t.Errorf("expected changed enitity")
	}
}

func TestSetEntityInfoProcessor_Apply_GivenInvalidData_ShouldReturnError(t *testing.T) {
	report := eventfilter.Report{}

	dataSets := []struct {
		processor eventfilter.SetEntityInfoProcessor
		event     types.Event
	}{
		{
			processor: eventfilter.SetEntityInfoProcessor{
				Name:  types.OptionalString{Set: true, Value: "Customer"},
				Value: types.OptionalInterface{Set: true, Value: "testcustomer"},
			},
			event: types.Event{
				Entity: nil,
			},
		},
	}

	for _, data := range dataSets {
		_, err := data.processor.Apply(data.event, eventfilter.ActionParameters{}, &report)
		if err == nil {
			t.Errorf("expected error but got nothing")
		}
	}
}

func TestCopyToEntityInfoProcessor_Apply_GivenValidData_ShouldApplyIt(t *testing.T) {
	expectedCustomer := "test-output (by test-author)"
	processor := eventfilter.CopyToEntityInfoProcessor{
		Name: types.OptionalString{Set: true, Value: "Customer"},
		From: types.OptionalString{Set: true, Value: "Event.Output"},
	}
	report := eventfilter.Report{}
	event := types.Event{
		Entity: &types.Entity{},
	}
	event, err := processor.Apply(event, eventfilter.ActionParameters{
		Event: types.Event{
			Output: expectedCustomer,
		},
	}, &report)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}
	if event.Entity.Infos["Customer"].Value != expectedCustomer {
		t.Errorf("expected output %v but got %v", expectedCustomer, event.Entity.Infos["Customer"].Value)
	}
	if !report.EntityUpdated {
		t.Errorf("expected changed enitity")
	}
}

func TestCopyToEntityInfoProcessor_Apply_GivenInvalidData_ShouldReturnError(t *testing.T) {
	report := eventfilter.Report{}

	dataSets := []struct {
		processor eventfilter.CopyToEntityInfoProcessor
		event     types.Event
	}{
		{
			processor: eventfilter.CopyToEntityInfoProcessor{
				Name: types.OptionalString{Set: true, Value: "Customer"},
				From: types.OptionalString{Set: true, Value: "Event.UnknownField"},
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
		},
		{
			processor: eventfilter.CopyToEntityInfoProcessor{
				Name: types.OptionalString{Set: true, Value: "Customer"},
				From: types.OptionalString{Set: true, Value: "Event.Output"},
			},
			event: types.Event{},
		},
	}

	for _, data := range dataSets {
		_, err := data.processor.Apply(data.event, eventfilter.ActionParameters{}, &report)
		if err == nil {
			t.Errorf("expected error but got nothing")
		}
	}
}

func TestCopyProcessor_Apply_GivenValidData_ShouldApplyIt(t *testing.T) {
	processor := eventfilter.CopyProcessor{
		From: types.OptionalString{Set: true, Value: "ExternalData.Entity"},
		To:   types.OptionalString{Set: true, Value: "Entity"},
	}
	report := eventfilter.Report{}
	expectedEntityID := "test-entity-id"
	event := types.Event{}
	event, err := processor.Apply(event, eventfilter.ActionParameters{
		ExternalData: map[string]interface{}{
			"Entity": types.Entity{ID: expectedEntityID},
		},
	}, &report)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
		return
	}
	if event.Entity.ID != expectedEntityID {
		t.Errorf("expected entity %v but got %v", expectedEntityID, event.Entity)
	}
	if report.EntityUpdated {
		t.Errorf("expected unchanged enitity")
	}
}

func TestCopyProcessor_Apply_GivenInvalidData_ShouldReturnError(t *testing.T) {
	report := eventfilter.Report{}

	dataSets := []struct {
		processor eventfilter.CopyProcessor
		event     types.Event
	}{
		{
			processor: eventfilter.CopyProcessor{
				From: types.OptionalString{
					Set:   true,
					Value: "Event.UnknownField",
				},
				To: types.OptionalString{
					Set:   true,
					Value: "Event.Output",
				},
			},
			event: types.Event{},
		},
		{
			processor: eventfilter.CopyProcessor{
				To: types.OptionalString{
					Set:   true,
					Value: "Event.UnknownField",
				},
				From: types.OptionalString{
					Set:   true,
					Value: "Event.Output",
				},
			},
			event: types.Event{},
		},
	}

	for _, data := range dataSets {
		_, err := data.processor.Apply(data.event, eventfilter.ActionParameters{}, &report)
		if err == nil {
			t.Errorf("expected error but got nothing")
		}
	}
}
