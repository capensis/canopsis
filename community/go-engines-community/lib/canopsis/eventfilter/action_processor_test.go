package eventfilter_test

import (
	"context"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	mock_eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
	mock_techmetrics "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/techmetrics"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestActionProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}, zerolog.Nop()), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	dataSets := []struct {
		testName      string
		action        eventfilter.ParsedAction
		event         types.Event
		regexMatches  eventfilter.RegexMatch
		externalData  map[string]interface{}
		expectedEvent types.Event
		expectedError bool
	}{
		{
			testName: "given set_field action should return success",
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetFieldFromTemplate,
				Name:        "Output",
				ParsedValue: tplExecutor.Parse("{{.ExternalData.data_1}}"),
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
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetFieldFromTemplate,
				Name:        "Output",
				ParsedValue: tplExecutor.Parse("{{.Some.data_1}}"),
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with int type",
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with bool type",
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with string slice type",
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with slice of interfaces but all items are strings",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       []interface{}{"test2", "test"},
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with primitive.A",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       primitive.A{"test2", "test"},
			},
			event: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       []string{"test"},
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
							Value:       []string{"test", "test2"},
						},
					},
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with float64 as a whole number",
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return success with float32 as a whole number",
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should return error with float value",
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info action should not return updated entity true, if info is not changed",
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfoFromTemplate,
				Name:        "Info 1",
				Description: "Test description",
				ParsedValue: tplExecutor.Parse("{{.ExternalData.data_1}}"),
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info_from_template action should return error, because of wrong template",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfoFromTemplate,
				Name:        "Info 1",
				Description: "Test description",
				ParsedValue: tplExecutor.Parse("{{.Some.data}}"),
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfoFromTemplate,
				Name:        "Info 1",
				Description: "Test description",
				ParsedValue: tplExecutor.Parse("new info"),
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info_from_template action should not return updated entity true, if info is not changed",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfoFromTemplate,
				Name:        "Info 1",
				Description: "Test description",
				ParsedValue: tplExecutor.Parse("new info"),
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with string type",
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with int type",
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with bool type",
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with string type",
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with string slice type",
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
				ExtraInfos: map[string]interface{}{
					"Test": []string{"test", "test2"},
				},
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with slice of interfaces but all items are strings",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.Test",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": []interface{}{"test2", "test"},
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
					IsUpdated: true,
				},
				ExtraInfos: map[string]interface{}{
					"Test": []interface{}{"test2", "test"},
				},
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success with primitive.A",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionCopyToEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.Test",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]interface{}{
					"Test": primitive.A{"test2", "test"},
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
					IsUpdated: true,
				},
				ExtraInfos: map[string]interface{}{
					"Test": primitive.A{"test2", "test"},
				},
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success  with float64 as a whole number",
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
				ExtraInfos: map[string]interface{}{
					"Test": float64(2),
				},
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return success  with float32 as a whole number",
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
				ExtraInfos: map[string]interface{}{
					"Test": float32(2),
				},
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should return error with float value",
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
					IsUpdated: true,
				},
			},
			expectedError: false,
		},
		{
			testName: "given copy_to_entity_info action should not return entity updated true, if info is not changed",
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
			action: eventfilter.ParsedAction{
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
		{
			testName: "given set_entity_info action with string slice with another order of items then in entity should return is updated false",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfo,
				Name:        "Info 1",
				Description: "Test description",
				Value:       []interface{}{"test2", "test3", "test1"},
			},
			event: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"Info 1": {
							Name:        "Info 1",
							Description: "Test description",
							Value:       primitive.A{"test3", "test1", "test2"},
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
							Value:       primitive.A{"test3", "test1", "test2"},
						},
					},
					IsUpdated: false,
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info_from_dictionary action should return success",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfoFromDictionary,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.dict",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]any{
					"dict": map[string]any{
						"key1": "value1",
						"key2": []string{"value2", "value3"},
						"key3": float64(3),
					},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]any{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"key1": {
							Name:        "key1",
							Description: "Test description",
							Value:       "value1",
						},
						"key2": {
							Name:        "key2",
							Description: "Test description",
							Value:       []string{"value2", "value3"},
						},
						"key3": {
							Name:        "key3",
							Description: "Test description",
							Value:       float64(3),
						},
					},
					IsUpdated: true,
				},
				ExtraInfos: map[string]any{
					"dict": map[string]any{
						"key1": "value1",
						"key2": []string{"value2", "value3"},
						"key3": float64(3),
					},
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info_from_dictionary action should update ExtraInfos and return success",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfoFromDictionary,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.dict",
			},
			event: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"key1": {
							Name:  "key1",
							Value: float32(1),
						},
					},
				},
				ExtraInfos: map[string]any{
					"dict": map[string]any{
						"key1": "value1",
						"key2": []string{"value2", "value3"},
						"key3": float64(3),
					},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]any{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"key1": {
							Name:        "key1",
							Description: "Test description",
							Value:       "value1",
						},
						"key2": {
							Name:        "key2",
							Description: "Test description",
							Value:       []string{"value2", "value3"},
						},
						"key3": {
							Name:        "key3",
							Description: "Test description",
							Value:       float64(3),
						},
					},
					IsUpdated: true,
				},
				ExtraInfos: map[string]any{
					"dict": map[string]any{
						"key1": "value1",
						"key2": []string{"value2", "value3"},
						"key3": float64(3),
					},
				},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info_from_dictionary action and without a dictionary should return success without EntityUpdated",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfoFromDictionary,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.dict",
			},
			event: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"key1": {
							Name:  "key1",
							Value: float32(1),
						},
					},
				},
				ExtraInfos: map[string]any{},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]any{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"key1": {
							Name:        "key1",
							Description: "",
							Value:       float32(1),
						},
					},
					IsUpdated: false,
				},
				ExtraInfos: map[string]any{},
			},
			expectedError: false,
		},
		{
			testName: "given set_entity_info_from_dictionary action and dictionary is not a map should return error",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfoFromDictionary,
				Description: "Test description",
				Value:       "Event.ExtraInfos.dict",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]any{
					"dict": "not a map",
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]any{},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]any{
					"dict": "not a map",
				},
			},
			expectedError: true,
		},
		{
			testName: "given set_entity_info_from_dictionary action and dictionary has invalid field should return error",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfoFromDictionary,
				Description: "Test description",
				Value:       "Event.ExtraInfos.dict",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]any{
					"dict": map[string]any{
						"key1": map[string]any{
							"subkey": "value",
						},
					},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]any{},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]any{
					"dict": map[string]any{
						"key1": map[string]any{
							"subkey": "value",
						},
					},
				},
			},
			expectedError: true,
		},
		{
			testName: "given set_entity_info_from_dictionary action and dictionary is nil should return error",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfoFromDictionary,
				Description: "Test description",
				Value:       "Event.ExtraInfos.dict",
			},
			event: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]any{
					"dict": nil,
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]any{},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
				ExtraInfos: map[string]any{
					"dict": nil,
				},
			},
			expectedError: true,
		},
		{
			testName: "given set_entity_info_from_dictionary action without infos change should return entity not updated",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetEntityInfoFromDictionary,
				Name:        "Info 1",
				Description: "Test description",
				Value:       "Event.ExtraInfos.dict",
			},
			event: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"key1": {
							Name:        "key1",
							Description: "Test description",
							Value:       "value1",
						},
						"key2": {
							Name:        "key2",
							Description: "Test description",
							Value:       []string{"value2", "value3"},
						},
						"key3": {
							Name:        "key3",
							Description: "Test description",
							Value:       float64(3),
						},
					},
				},
				ExtraInfos: map[string]any{
					"dict": map[string]any{
						"key1": "value1",
						"key2": []string{"value2", "value3"},
						"key3": float64(3),
					},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]any{},
			expectedEvent: types.Event{
				Entity: &types.Entity{
					Infos: map[string]types.Info{
						"key1": {
							Name:        "key1",
							Description: "Test description",
							Value:       "value1",
						},
						"key2": {
							Name:        "key2",
							Description: "Test description",
							Value:       []string{"value2", "value3"},
						},
						"key3": {
							Name:        "key3",
							Description: "Test description",
							Value:       float64(3),
						},
					},
				},
				ExtraInfos: map[string]any{
					"dict": map[string]any{
						"key1": "value1",
						"key2": []string{"value2", "value3"},
						"key3": float64(3),
					},
				},
			},
			expectedError: false,
		},
	}

	mockAlarmConfigProvider := mock_config.NewMockAlarmConfigProvider(ctrl)
	mockAlarmConfigProvider.EXPECT().Get().Return(config.AlarmConfig{
		EnableArraySortingInEntityInfos: true,
	}).AnyTimes()
	mockFailureService := mock_eventfilter.NewMockFailureService(ctrl)
	mockFailureService.EXPECT().Add(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockTechMetricsSender := mock_techmetrics.NewMockSender(ctrl)
	mockTechMetricsSender.EXPECT().SendCheEntityInfo(gomock.Any(), gomock.Any()).AnyTimes()
	processor := eventfilter.NewActionProcessor(mockAlarmConfigProvider, mockFailureService, tplExecutor, mockTechMetricsSender)
	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			resultEvent, resultErr := processor.Process(context.Background(), "test", dataset.action, dataset.event, eventfilter.RegexMatchWrapper{
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
				t.Errorf("expected no error but got %+v", resultErr)
			}
		})
	}
}
