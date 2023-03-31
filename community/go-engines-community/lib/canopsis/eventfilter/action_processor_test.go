package eventfilter_test

import (
	"context"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_techmetrics "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/techmetrics"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestActionProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dataSets := []struct {
		testName      string
		action        eventfilter.Action
		event         types.Event
		regexMatches  eventfilter.RegexMatch
		externalData  map[string]interface{}
		expectedEvent types.Event
		expectedError bool
	}{
		{
			testName: "given set_field action should return success",
			action: eventfilter.Action{
				Type:  eventfilter.ActionSetField,
				Name:  "Output",
				Value: "test output",
			},
			event:        types.Event{},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Output: "test output",
			},
			expectedError: false,
		},
		{
			testName: "given set_field action should return error, because of wrong value field type",
			action: eventfilter.Action{
				Type:  eventfilter.ActionSetField,
				Name:  "Output",
				Value: 5,
			},
			event:         types.Event{},
			regexMatches:  eventfilter.RegexMatch{},
			externalData:  map[string]interface{}{},
			expectedEvent: types.Event{},
			expectedError: true,
		},
		{
			testName: "given set_field_from_template action should return success",
			action: eventfilter.Action{
				Type:  eventfilter.ActionSetFieldFromTemplate,
				Name:  "Output",
				Value: "{{.ExternalData.data_1}}",
			},
			event:        types.Event{},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{
				"data_1": "test output",
			},
			expectedEvent: types.Event{
				Output: "test output",
			},
			expectedError: false,
		},
		{
			testName: "given set_field_from_template action should return error, because of wrong template",
			action: eventfilter.Action{
				Type:  eventfilter.ActionSetFieldFromTemplate,
				Name:  "Output",
				Value: "{{.Some.data_1}}",
			},
			event:        types.Event{},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{
				"data_1": "test output",
			},
			expectedEvent: types.Event{},
			expectedError: true,
		},
		{
			testName: "given set_field_from_template action should return error, because value should be a string",
			action: eventfilter.Action{
				Type:  eventfilter.ActionSetFieldFromTemplate,
				Name:  "Output",
				Value: 123,
			},
			event:        types.Event{},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{
				"data_1": "test output",
			},
			expectedEvent: types.Event{},
			expectedError: true,
		},
		{
			testName: "given set_entity_info action should return success with string type",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "test output",
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "test output",
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with int type",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       123,
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       123,
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with bool type",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       true,
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       true,
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with string slice type",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       []string{"test", "test2"},
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       []string{"test", "test2"},
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with slice of interfaces but all items are strings",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       []interface{}{"test", "test2"},
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       []interface{}{"test", "test2"},
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with primitive.A",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       primitive.A{"test", "test2"},
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       primitive.A{"test", "test2"},
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with float64 as a whole number",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       float64(2),
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       float64(2),
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with float32 as a whole number",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       float32(2),
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       float32(2),
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return error with float value",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       1.2,
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
			},
			expectedError: true,
		},
		{
			testName: "given set_entity_info action should return error with slice of interfaces, where some are not strings",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       []interface{}{"test", 1, "test2"},
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
			},
			expectedError: true,
		},
		{
			testName: "given set_entity_info action should return error with primitive.A, where some are not strings",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       primitive.A{"test", 1, "test2"},
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
			},
			expectedError: true,
		},
		{
			testName: "given set_entity_info action should return error with structs",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value: struct {
					test string
				}{test: "test"},
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
			},
			expectedError: true,
		},
		{
			testName: "given set_entity_info action should return updated entity true, if infos is changed",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "new info",
			},
			event: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "old info",
						},
					},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "new info",
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should not return updated entity true, if info is not changed",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "new info",
			},
			event: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "new info",
						},
					},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "new info",
						},
					},
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info_from_template action should return success",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfoFromTemplate,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "{{.ExternalData.data_1}}",
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{
				"data_1": "test output",
			},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "test output",
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info_from_template action should return error, because of wrong template",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfoFromTemplate,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "{{.Some.data}}",
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
			},
			expectedError: true,
		},
		{
			testName: "given set_entity_info_from_template action should return error, because value should be a string",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfoFromTemplate,
				Name:        "Info 1",
				Description: "Test description",
				Value:       123,
			},
			event: types.Event{
				Entity: &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
			},
			expectedError: true,
		},
		{
			testName: "given set_entity_info_from_template action should return updated entity true, if infos is changed",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfoFromTemplate,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "new info",
			},
			event: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "old info",
						},
					},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "new info",
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info_from_template action should not return updated entity true, if info is not changed",
			action: eventfilter.Action{
				Type:        eventfilter.ActionSetEntityInfoFromTemplate,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "new info",
			},
			event: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "new info",
						},
					},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "new info",
						},
					},
				},
			},
			expectedError: false,
		},
		{
			testName: "given copy action should return success",
			action: eventfilter.Action{
				Type:  eventfilter.ActionCopy,
				Name:  "Output",
				Value: "Event.Resource",
			},
			event: types.Event{
				Resource: "test resource",
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Resource: "test resource",
				Output:   "test resource",
			},
			expectedError: false,
		},
		{
			testName: "given copy action should return error, because value should be a string",
			action: eventfilter.Action{
				Type:  eventfilter.ActionCopy,
				Name:  "Output",
				Value: 123,
			},
			event: types.Event{
				Resource: "test resource",
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Resource: "test resource",
			},
			expectedError: true,
		},
		{
			testName: "given copy action should return error, because get field doesn't exist",
			action: eventfilter.Action{
				Type:  eventfilter.ActionCopy,
				Name:  "Output",
				Value: "Some",
			},
			event: types.Event{
				Resource: "test resource",
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Resource: "test resource",
			},
			expectedError: true,
		},
		{
			testName: "given copy action should return error, because set field doesn't exist",
			action: eventfilter.Action{
				Type:  eventfilter.ActionCopy,
				Name:  "Some",
				Value: "Event.Resource",
			},
			event: types.Event{
				Resource: "test resource",
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Resource: "test resource",
			},
			expectedError: true,
		},
		{
			testName: "given copy_to_entity_info action should return success with string type",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.Resource",
			},
			event: types.Event{
				Resource: "test resource",
				Entity:   &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Resource: "test resource",
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "test resource",
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with string type",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.Resource",
			},
			event: types.Event{
				Resource: "test resource",
				Entity:   &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Resource: "test resource",
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "test resource",
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with int type",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "ExternalData.Info",
			},
			event: types.Event{
				Resource: "test resource",
				Entity:   &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{
				"Info": 123,
			},
			expectedEvent: types.Event{
				Resource: "test resource",
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       123,
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with bool type",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "ExternalData.Info",
			},
			event: types.Event{
				Resource: "test resource",
				Entity:   &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{
				"Info": true,
			},
			expectedEvent: types.Event{
				Resource: "test resource",
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       true,
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with string type",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.Resource",
			},
			event: types.Event{
				Resource: "test resource",
				Entity:   &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Resource: "test resource",
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "test resource",
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with string slice type",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.Test",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": []string{"test", "test2"},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       []string{"test", "test2"},
						},
					},
				},
				ExtraInfos: map[string]interface{}{
					"Test": []string{"test", "test2"},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with slice of interfaces but all items are strings",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.Test",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": []interface{}{"test", "test2"},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       []interface{}{"test", "test2"},
						},
					},
				},
				ExtraInfos: map[string]interface{}{
					"Test": []interface{}{"test", "test2"},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with primitive.A",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.Test",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": primitive.A{"test", "test2"},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       primitive.A{"test", "test2"},
						},
					},
				},
				ExtraInfos: map[string]interface{}{
					"Test": primitive.A{"test", "test2"},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success  with float64 as a whole number",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.Test",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": float64(2),
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       float64(2),
						},
					},
				},
				ExtraInfos: map[string]interface{}{
					"Test": float64(2),
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success  with float32 as a whole number",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.Test",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": float32(2),
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       float32(2),
						},
					},
				},
				ExtraInfos: map[string]interface{}{
					"Test": float32(2),
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return error with float value",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.Test",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": 1.2,
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": 1.2,
				},
			},
			expectedError: true,
		},
		{
			testName: "given copy_to_entity_info action should return error with slice of interfaces, where some are not strings",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.Test",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": []interface{}{"test1", 1, "test2"},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": []interface{}{"test1", 1, "test2"},
				},
			},
			expectedError: true,
		},
		{
			testName: "given copy_to_entity_info action should return error with primitive.A, where some are not strings",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.Test",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": primitive.A{"test1", 1, "test2"},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": primitive.A{"test1", 1, "test2"},
				},
			},
			expectedError: true,
		},
		{
			testName: "given copy_to_entity_info action should return error with structs",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.Test",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": struct {
						Test string
					}{Test: "test"},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": struct {
						Test string
					}{Test: "test"},
				},
			},
			expectedError: true,
		},
		{
			testName: "given copy_to_entity_info action should change existing info",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.Resource",
			},
			event: types.Event{
				Resource: "test resource",
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "old resource",
						},
					},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Resource: "test resource",
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "test resource",
						},
					},
				},
				IsEntityUpdated: true,
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should not return entity updated true, if info is not changed",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.Resource",
			},
			event: types.Event{
				Resource: "test resource",
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "test resource",
						},
					},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Resource: "test resource",
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       "test resource",
						},
					},
				},
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return error, because value is not a string",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       123,
			},
			event: types.Event{
				Resource: "test resource",
				Entity:   &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Resource: "test resource",
				Entity:   &types.Entity{},
			},
			expectedError: true,
		},
		{
			testName: "given copy_to_entity_info action should return error, because value field is not exist",
			action: eventfilter.Action{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Some",
			},
			event: types.Event{
				Resource: "test resource",
				Entity:   &types.Entity{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Resource: "test resource",
				Entity:   &types.Entity{},
			},
			expectedError: true,
		},
	}

	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	mockTechMetricsSender := mock_techmetrics.NewMockSender(ctrl)
	mockTechMetricsSender.EXPECT().SendCheEntityInfo(gomock.Any(), gomock.Any()).AnyTimes()
	processor := eventfilter.NewActionProcessor(tplExecutor, mockTechMetricsSender)
	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			resultEvent, resultErr := processor.Process(context.Background(), dataset.action, dataset.event, eventfilter.RegexMatchWrapper{
				BackwardCompatibility: false,
				RegexMatch:            dataset.regexMatches,
			}, dataset.externalData)
			if !reflect.DeepEqual(resultEvent, dataset.expectedEvent) {
				t.Errorf("expected an event = %v, but got %v", dataset.expectedEvent, resultEvent)
			}

			if dataset.expectedError && resultErr == nil {
				t.Error("expected an error")
			}

			if !dataset.expectedError && resultErr != nil {
				t.Error("expected no error")
			}
		})
	}
}
