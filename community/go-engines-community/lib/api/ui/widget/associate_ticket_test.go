package widget

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.uber.org/mock/gomock"
)

func TestValidateAssociateTicketParameters(t *testing.T) {
	validate := validator.New()
	validate.RegisterStructValidation(NewValidator(nil).ValidateEditRequest, EditRequest{})

	for _, tc := range []struct {
		name      string
		payload   string
		wantTag   string
		wantField string
	}{
		{name: "absent settings preserve existing requests", payload: `{"type":"AlarmsList"}`},
		{name: "empty and arbitrary custom strings are accepted", payload: `{"type":"AlarmsList","parameters":{"associateTicket":{"formats":["","not a regexp ["]}}}`},
		{name: "empty custom list is accepted", payload: `{"type":"AlarmsList","parameters":{"associateTicket":{"formats":[]}}}`},
		{name: "nonempty template ignores custom formats", payload: `{"type":"AlarmsList","parameters":{"associateTicket":{"formatsTemplate":"tpl","formats":[null]}}}`},
		{name: "other widget types do not validate associate ticket", payload: `{"type":"Other","parameters":{"associateTicket":[]}}`},
		{name: "node must be object", payload: `{"type":"AlarmsList","parameters":{"associateTicket":[]}}`, wantTag: "invalid", wantField: "parameters.associateTicket"},
		{name: "null node is invalid", payload: `{"type":"AlarmsList","parameters":{"associateTicket":null}}`, wantTag: "invalid", wantField: "parameters.associateTicket"},
		{name: "template reference must be string", payload: `{"type":"AlarmsList","parameters":{"associateTicket":{"formatsTemplate":null}}}`, wantTag: "value_string", wantField: "parameters.associateTicket.formatsTemplate"},
		{name: "custom list must be array", payload: `{"type":"AlarmsList","parameters":{"associateTicket":{"formats":"INC"}}}`, wantTag: "value_string_array", wantField: "parameters.associateTicket.formats"},
		{name: "null list is invalid", payload: `{"type":"AlarmsList","parameters":{"associateTicket":{"formats":null}}}`, wantTag: "value_string_array", wantField: "parameters.associateTicket.formats"},
		{name: "number item is invalid", payload: `{"type":"AlarmsList","parameters":{"associateTicket":{"formats":[1]}}}`, wantTag: "value_string", wantField: "parameters.associateTicket.formats.0"},
		{name: "object item is invalid", payload: `{"type":"AlarmsList","parameters":{"associateTicket":{"formats":[{}]}}}`, wantTag: "value_string", wantField: "parameters.associateTicket.formats.0"},
		{name: "null item is invalid", payload: `{"type":"AlarmsList","parameters":{"associateTicket":{"formats":[null]}}}`, wantTag: "value_string", wantField: "parameters.associateTicket.formats.0"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var request EditRequest
			if err := json.Unmarshal([]byte(tc.payload), &request); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			err := validate.Struct(request)
			if tc.wantTag == "" {
				if err != nil {
					t.Fatalf("validate.Struct() error = %v, want no error", err)
				}
				return
			}
			var validationErrs validator.ValidationErrors
			if !errors.As(err, &validationErrs) || len(validationErrs) != 1 {
				t.Fatalf("validate.Struct() error = %v, want one validation error", err)
			}
			fieldErr := validationErrs[0]
			if fieldErr.Tag() != tc.wantTag {
				t.Errorf("validation tag = %q, want %q", fieldErr.Tag(), tc.wantTag)
			}
			gotField := validation.NewError(validationErrs, request).TransformNamespace(fieldErr.StructNamespace())
			if gotField != tc.wantField {
				t.Errorf("validation field = %q, want %q", gotField, tc.wantField)
			}
		})
	}
}

func TestAssociateTicketFormatLimits(t *testing.T) {
	validate := validator.New()
	validate.RegisterStructValidation(NewValidator(nil).ValidateEditRequest, EditRequest{})
	tooMany := make([]any, view.TicketFormatsMaxCount+1)
	for i := range tooMany {
		tooMany[i] = ""
	}
	valid := make([]any, view.TicketFormatsMaxCount)
	for i := range valid {
		valid[i] = ""
	}
	validRequests := []EditRequest{
		{Type: view.WidgetTypeAlarmsList, Parameters: view.Parameters{RemainParameters: map[string]any{"associateTicket": map[string]any{"formats": valid}}}},
		{Type: view.WidgetTypeAlarmsList, Parameters: view.Parameters{RemainParameters: map[string]any{"associateTicket": map[string]any{"formats": []any{strings.Repeat("é", view.TicketFormatMaxLength)}}}}},
		{Type: view.WidgetTypeAlarmsList, Parameters: view.Parameters{RemainParameters: map[string]any{"associateTicket": map[string]any{"formatsTemplate": "selected", "formats": tooMany}}}},
	}
	for i, request := range validRequests {
		if err := validate.Struct(request); err != nil {
			t.Errorf("valid boundary request %d rejected: %v", i, err)
		}
	}

	tests := []struct {
		name      string
		formats   []any
		field     string
		wantParam string
	}{
		{
			name:      "too many custom formats",
			formats:   tooMany,
			field:     "parameters.associateTicket.formats",
			wantParam: "100",
		},
		{
			name:      "too many unicode characters",
			formats:   []any{strings.Repeat("é", view.TicketFormatMaxLength+1)},
			field:     "parameters.associateTicket.formats.0",
			wantParam: "255",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request := EditRequest{
				Type: view.WidgetTypeAlarmsList,
				Parameters: view.Parameters{RemainParameters: map[string]any{
					"associateTicket": map[string]any{"formats": tc.formats},
				}},
			}
			err := validate.Struct(request)
			var validationErr validator.ValidationErrors
			if !errors.As(err, &validationErr) {
				t.Fatalf("validator error = %v, want validation errors", err)
			}
			if len(validationErr) != 1 {
				t.Fatalf("validation errors = %v, want one", validationErr)
			}
			fieldErr := validationErr[0]
			if fieldErr.Tag() != "max" || fieldErr.Param() != tc.wantParam {
				t.Errorf("tag/param = %q/%q, want max/%q", fieldErr.Tag(), fieldErr.Param(), tc.wantParam)
			}
			gotField := validation.NewError(validationErr, request).TransformNamespace(fieldErr.StructNamespace())
			if gotField != tc.field {
				t.Errorf("validation field = %q, want %q", gotField, tc.field)
			}
		})
	}
}

func TestTicketFormatsTemplateMapping(t *testing.T) {
	parameters := view.GetWidgetTemplateParameters()[view.WidgetTypeAlarmsList][view.WidgetTemplateTypeTicketFormats]
	if len(parameters) != 1 || parameters[0] != "associateTicket.formats" {
		t.Fatalf("ticket format parameters = %#v, want [associateTicket.formats]", parameters)
	}
}

func TestTransformTicketFormatsTemplate(t *testing.T) {
	t.Run("template values and title replace caller values", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		collection := mock_mongo.NewMockDbCollection(ctrl)
		result := mock_mongo.NewMockSingleResultHelper(ctrl)
		collection.EXPECT().FindOne(gomock.Any(), gomock.Eq(bson.M{
			"_id":  "ticket-template-id",
			"type": view.WidgetTemplateTypeTicketFormats,
		})).Return(result)
		result.EXPECT().Decode(gomock.Any()).DoAndReturn(func(value *view.WidgetTemplate) error {
			*value = view.WidgetTemplate{
				Title:   "Resolved ticket formats",
				Type:    view.WidgetTemplateTypeTicketFormats,
				Formats: []string{"^INC-[0-9]+$", ""},
			}
			return nil
		})

		request := EditRequest{
			Type: view.WidgetTypeAlarmsList,
			Parameters: view.Parameters{RemainParameters: map[string]any{
				"associateTicket": map[string]any{
					"formatsTemplate":      "ticket-template-id",
					"formatsTemplateTitle": "caller title",
					"formats":              []any{"custom value"},
				},
			}},
		}
		store := store{widgetTemplateCollection: collection}
		if err := store.transformTemplateFields(context.Background(), &request); err != nil {
			t.Fatal(err)
		}
		associateTicket, ok := request.Parameters.RemainParameters["associateTicket"].(map[string]any)
		if !ok {
			t.Fatalf("associateTicket is not a map[string]any")
		}
		if associateTicket["formatsTemplateTitle"] != "Resolved ticket formats" {
			t.Fatalf("formatsTemplateTitle = %#v", associateTicket["formatsTemplateTitle"])
		}
		got, ok := associateTicket["formats"].([]string)
		want := []string{"^INC-[0-9]+$", ""}
		if !ok || !slices.Equal(got, want) {
			t.Fatalf("formats = %#v, want %#v", got, want)
		}
	})

	t.Run("empty template values resolve to an empty list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		collection := mock_mongo.NewMockDbCollection(ctrl)
		result := mock_mongo.NewMockSingleResultHelper(ctrl)
		collection.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(result)
		result.EXPECT().Decode(gomock.Any()).DoAndReturn(func(value *view.WidgetTemplate) error {
			*value = view.WidgetTemplate{Type: view.WidgetTemplateTypeTicketFormats}
			return nil
		})
		request := EditRequest{
			Type: view.WidgetTypeAlarmsList,
			Parameters: view.Parameters{RemainParameters: map[string]any{
				"associateTicket": map[string]any{"formatsTemplate": "ticket-template-id"},
			}},
		}
		store := store{widgetTemplateCollection: collection}
		if err := store.transformTemplateFields(context.Background(), &request); err != nil {
			t.Fatal(err)
		}
		assocTicketParams, ok := request.Parameters.RemainParameters["associateTicket"].(map[string]any)
		if !ok {
			t.Fatalf("associateTicket is not a map[string]any")
		}
		formats, ok := assocTicketParams["formats"]
		if !ok {
			t.Fatalf("formats key is missing in associateTicket")
		}
		if got, ok := formats.([]string); !ok || len(got) != 0 {
			t.Fatalf("formats = %#v, want empty []string", formats)
		}
	})

	t.Run("custom settings clear supplied template title", func(t *testing.T) {
		request := EditRequest{
			Type: view.WidgetTypeAlarmsList,
			Parameters: view.Parameters{RemainParameters: map[string]any{
				"associateTicket": map[string]any{
					"formatsTemplate":      "",
					"formatsTemplateTitle": "stale title",
					"formats":              []any{"custom"},
				},
			}},
		}
		store := store{}
		if err := store.transformTemplateFields(context.Background(), &request); err != nil {
			t.Fatal(err)
		}
		associateTicket, ok := request.Parameters.RemainParameters["associateTicket"].(map[string]any)
		if !ok {
			t.Fatalf("associateTicket is not a map[string]any")
		}
		if _, ok = associateTicket["formatsTemplateTitle"]; ok {
			t.Fatalf("custom settings retained template title: %#v", associateTicket)
		}
	})
}
