package eventfilter_test

import (
	"context"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	mock_eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
	mock_techmetrics "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/techmetrics"
	"github.com/golang/mock/gomock"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestActionProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}, zerolog.Nop()), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	dataSets := []struct {
		testName              string
		action                eventfilter.ParsedAction
		event                 types.Event
		regexMatches          eventfilter.RegexMatch
		externalData          map[string]interface{}
		expectedEvent         types.Event
		expectedError         bool
		expectedEntityUpdated bool
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
			expectedError:         false,
			expectedEntityUpdated: false,
		},
		{
			testName: "given set_field action should return error, because of wrong value field type",
			action: eventfilter.ParsedAction{
				Type:  eventfilter.ActionSetField,
				Name:  "Output",
				Value: 5,
			},
			event:                 types.Event{},
			regexMatches:          eventfilter.RegexMatch{},
			externalData:          map[string]interface{}{},
			expectedEvent:         types.Event{},
			expectedError:         true,
			expectedEntityUpdated: false,
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
			expectedError:         false,
			expectedEntityUpdated: false,
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
			expectedEvent:         types.Event{},
			expectedError:         true,
			expectedEntityUpdated: false,
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
			expectedEvent:         types.Event{},
			expectedError:         true,
			expectedEntityUpdated: false,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
			expectedError:         false,
			expectedEntityUpdated: false,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
			expectedError:         false,
			expectedEntityUpdated: false,
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
			expectedError:         false,
			expectedEntityUpdated: false,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
			expectedError:         true,
			expectedEntityUpdated: false,
		},
		{
			testName: "given copy action should return success with Tags copied from ExtraInfos",
			action: eventfilter.ParsedAction{
				Type:  eventfilter.ActionCopy,
				Name:  "Tags",
				Value: "Event.ExtraInfos.newtags",
			},
			event: types.Event{
				Tags: map[string]string{
					"tag0": "",
					"tag1": "value1",
					"tag2": "value2",
					"tag3": "value3",
				},
				ExtraInfos: map[string]interface{}{
					"newtags": map[string]interface{}{
						"tag1": "value1a",
						"tag2": "value2a",
						"tag4": "",
					},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Tags: map[string]string{
					"tag0": "",
					"tag1": "value1a",
					"tag2": "value2a",
					"tag3": "value3",
					"tag4": "",
				},
				ExtraInfos: map[string]interface{}{
					"newtags": map[string]interface{}{
						"tag1": "value1a",
						"tag2": "value2a",
						"tag4": "",
					},
				},
			},
			expectedError:         false,
			expectedEntityUpdated: false,
		},
		{
			testName: "given copy action should return error because ExtraInfos value is incompatible type with Tags",
			action: eventfilter.ParsedAction{
				Type:  eventfilter.ActionCopy,
				Name:  "Tags",
				Value: "Event.ExtraInfos.newtags",
			},
			event: types.Event{
				Tags: map[string]string{
					"tag0": "",
					"tag1": "value1",
					"tag2": "value2",
					"tag3": "value3",
				},
				ExtraInfos: map[string]interface{}{
					"newtags": []string{"tag1", "tag2", "tag4"},
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{},
			expectedEvent: types.Event{
				Tags: map[string]string{
					"tag0": "",
					"tag1": "value1",
					"tag2": "value2",
					"tag3": "value3",
				},
				ExtraInfos: map[string]interface{}{
					"newtags": []string{"tag1", "tag2", "tag4"},
				},
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
				ExtraInfos: map[string]interface{}{
					"Test": []string{"test", "test2"},
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
				ExtraInfos: map[string]interface{}{
					"Test": []interface{}{"test2", "test"},
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
				ExtraInfos: map[string]interface{}{
					"Test": primitive.A{"test2", "test"},
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
				ExtraInfos: map[string]interface{}{
					"Test": float64(2),
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
				},
				ExtraInfos: map[string]interface{}{
					"Test": float32(2),
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
				},
			},
			expectedError:         false,
			expectedEntityUpdated: true,
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
			expectedError:         false,
			expectedEntityUpdated: false,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
			expectedError:         true,
			expectedEntityUpdated: false,
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
			testName: "given set_tags action should return success with tags assigned from regex",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetTags,
				Name:        "Tags",
				Description: "Test description",
				Value:       "Event.Output",
			},
			event: types.Event{
				Output: "Some text preceding tags. Prod ENV; Critical Severity;",
				Entity: &types.Entity{},
				Tags: map[string]string{
					"tag0": "",
					"tag1": "value1",
				},
			},
			regexMatches: eventfilter.RegexMatch{
				EventRegexMatches: match.EventRegexMatches{
					MatchedRegexp: func(s string) utils.RegexExpression {
						v, err := utils.NewRegexExpression(s)
						if err != nil {
							panic(err)
						}
						return v
					}(`(?P<value>[a-zA-Z]+)\s+(?P<name>[a-zA-Z]+);`),
					Output: pattern.RegexMatches{
						"value": "Prod",
					},
				},
			},
			expectedEvent: types.Event{
				Output: "Some text preceding tags. Prod ENV; Critical Severity;",
				Entity: &types.Entity{},
				Tags: map[string]string{
					"tag0":     "",
					"tag1":     "value1",
					"ENV":      "Prod",
					"Severity": "Critical",
				},
			},
		},
		{
			testName: "given set_tags action should return success with tags assigned from regex and tags from ExtraInfos",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetTags,
				Name:        "Tags",
				Description: "Test description",
				Value:       "Event.ExtraInfos.strparam",
			},
			event: types.Event{
				Entity: &types.Entity{},
				Tags: map[string]string{
					"tag0": "",
					"tag1": "value1",
				},
				ExtraInfos: map[string]interface{}{
					"strparam": "Prod ENV; Critical Severity;",
				},
			},
			regexMatches: eventfilter.RegexMatch{
				EventRegexMatches: match.EventRegexMatches{
					MatchedRegexp: func(s string) utils.RegexExpression {
						v, err := utils.NewRegexExpression(s)
						if err != nil {
							panic(err)
						}
						return v
					}(`(?P<value>[a-zA-Z]+)\s+(?P<name>[a-zA-Z]+);`),
					ExtraInfos: map[string]pattern.RegexMatches{
						"strparam": {
							"value": "Prod",
						},
					},
				},
			},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
				Tags: map[string]string{
					"tag0":     "",
					"tag1":     "value1",
					"ENV":      "Prod",
					"Severity": "Critical",
				},
				ExtraInfos: map[string]interface{}{
					"strparam": "Prod ENV; Critical Severity;",
				},
			},
		},
		{
			testName: "given set_tags_from_template action should return success with tag assigned from template",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetTagsFromTemplate,
				Name:        "ENV",
				Description: "Test description",
				ParsedValue: tplExecutor.Parse("{{.ExternalData.data_1}}"),
			},
			event: types.Event{
				Entity: &types.Entity{},
				Tags: map[string]string{
					"tag0": "",
					"tag1": "value1",
				},
			},
			regexMatches: eventfilter.RegexMatch{},
			externalData: map[string]interface{}{
				"data_1": "Prod",
			},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
				Tags: map[string]string{
					"tag0": "",
					"tag1": "value1",
					"ENV":  "Prod",
				},
			},
		},
		{
			testName: "given set_tags_from_template action should return success with tag assigned from template and pattern matched Output",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetTagsFromTemplate,
				Name:        "ENV",
				Description: "Test description",
				ParsedValue: tplExecutor.Parse("{{.Event.Output}}"),
			},
			event: types.Event{
				Entity: &types.Entity{},
				Tags: map[string]string{
					"tag0": "",
					"tag1": "value1",
				},
				Output: "Some text preceding tags. Prod ENV; Critical Severity;",
			},
			regexMatches: eventfilter.RegexMatch{
				EventRegexMatches: match.EventRegexMatches{
					MatchedRegexp: func(s string) utils.RegexExpression {
						v, err := utils.NewRegexExpression(s)
						if err != nil {
							panic(err)
						}
						return v
					}(`(?P<value>[a-zA-Z]+)\s+(?P<name>[a-zA-Z]+);`),
					Output: pattern.RegexMatches{
						"value": "Prod",
					},
				},
			},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
				Output: "Some text preceding tags. Prod ENV; Critical Severity;",
				Tags: map[string]string{
					"tag0":     "",
					"tag1":     "value1",
					"ENV":      "Prod",
					"Severity": "Critical",
				},
			},
		},
		{
			testName: "successful set_tags_from_template action set tags from template and pattern matched name-value pairs from ExtraInfos",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetTagsFromTemplate,
				Name:        "ENV",
				Description: "Test description",
				ParsedValue: tplExecutor.Parse("{{.Event.ExtraInfos.strparam}}"),
			},
			event: types.Event{
				Entity: &types.Entity{},
				Tags: map[string]string{
					"tag0": "",
					"tag1": "value1",
				},
				ExtraInfos: map[string]interface{}{
					"strparam": "Some text preceding tags. Prod ENV; Critical Severity;",
				},
			},
			regexMatches: eventfilter.RegexMatch{
				EventRegexMatches: match.EventRegexMatches{
					MatchedRegexp: func(s string) utils.RegexExpression {
						v, err := utils.NewRegexExpression(s)
						if err != nil {
							panic(err)
						}
						return v
					}(`(?P<value>[a-zA-Z]+)\s+(?P<name>[a-zA-Z]+);`),
					ExtraInfos: map[string]pattern.RegexMatches{
						"strparam": {
							"value": "Prod",
						},
					},
				},
			},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
				Tags: map[string]string{
					"tag0":     "",
					"tag1":     "value1",
					"ENV":      "Prod",
					"Severity": "Critical",
				},
				ExtraInfos: map[string]interface{}{
					"strparam": "Some text preceding tags. Prod ENV; Critical Severity;",
				},
			},
		},
		{
			testName: "given set_tags_from_template action should return success with tag assigned from template and named as ParsedAction.Name",
			action: eventfilter.ParsedAction{
				Type:        eventfilter.ActionSetTagsFromTemplate,
				Name:        "ENV",
				Description: "Test description",
				ParsedValue: tplExecutor.Parse("{{.Event.ExtraInfos.strparam}}"),
			},
			event: types.Event{
				Entity: &types.Entity{},
				Tags: map[string]string{
					"tag0": "",
					"tag1": "value1",
				},
				ExtraInfos: map[string]interface{}{
					"strparam": "Prod ENV; Critical Severity;",
				},
			},
			regexMatches: eventfilter.RegexMatch{
				EventRegexMatches: match.EventRegexMatches{
					ExtraInfos: map[string]pattern.RegexMatches{
						"strparam": {
							"value": "Prod",
						},
					},
				},
			},
			expectedEvent: types.Event{
				Entity: &types.Entity{},
				Tags: map[string]string{
					"tag0": "",
					"tag1": "value1",
					"ENV":  "Prod ENV; Critical Severity;",
				},
				ExtraInfos: map[string]interface{}{
					"strparam": "Prod ENV; Critical Severity;",
				},
			},
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
				},
				ExtraInfos: map[string]any{
					"dict": map[string]any{
						"key1": "value1",
						"key2": []string{"value2", "value3"},
						"key3": float64(3),
					},
				},
			},
			expectedEntityUpdated: true,
			expectedError:         false,
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
				},
				ExtraInfos: map[string]any{
					"dict": map[string]any{
						"key1": "value1",
						"key2": []string{"value2", "value3"},
						"key3": float64(3),
					},
				},
			},
			expectedEntityUpdated: true,
			expectedError:         false,
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
				},
				ExtraInfos: map[string]any{},
			},
			expectedEntityUpdated: false,
			expectedError:         false,
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
			expectedEntityUpdated: false,
			expectedError:         true,
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
			expectedEntityUpdated: false,
			expectedError:         true,
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
			expectedEntityUpdated: false,
			expectedError:         true,
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
			expectedEntityUpdated: false,
			expectedError:         false,
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
		dataset := dataset
		t.Run(dataset.testName, func(t *testing.T) {
			resultEntityUpdated, resultErr := processor.Process(context.Background(), "test", dataset.action, &dataset.event,
				dataset.regexMatches, dataset.externalData)
			if diff := pretty.Compare(dataset.expectedEvent, dataset.event); diff != "" {
				t.Errorf("unexpected event: %s", diff)
			}

			if dataset.expectedError && resultErr == nil {
				t.Error("expected an error")
			}

			if dataset.expectedEntityUpdated != resultEntityUpdated {
				t.Errorf("expected an entityUpdated = %v, but got %v", dataset.expectedEntityUpdated, resultEntityUpdated)
			}

			if !dataset.expectedError && resultErr != nil {
				t.Errorf("expected no error but got %+v", resultErr)
			}
			if resultErr != nil {
				t.Logf("%s returned: %s", dataset.testName, resultErr)
			}
		})
	}
}
