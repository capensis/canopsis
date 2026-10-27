package widgettemplate

import (
	"encoding/json"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestEditRequestUnmarshalFormats(t *testing.T) {
	for _, tc := range []struct {
		name    string
		payload string
		want    []string
		wantErr bool
	}{
		{name: "preserves strings", payload: `{"formats":["^INC-[0-9]+$","","  exact  "]}`, want: []string{"^INC-[0-9]+$", "", "  exact  "}},
		{name: "allows empty list", payload: `{"formats":[]}`, want: []string{}},
		{name: "rejects null list", payload: `{"formats":null}`, wantErr: true},
		{name: "rejects null entry", payload: `{"formats":[null]}`, wantErr: true},
		{name: "rejects case variant null entry", payload: `{"Formats":[null]}`, wantErr: true},
		{name: "rejects malformed duplicate key", payload: `{"formats":[null],"Formats":["INC"]}`, wantErr: true},
		{name: "rejects number entry", payload: `{"formats":[1]}`, wantErr: true},
		{name: "rejects object entry", payload: `{"formats":[{}]}`, wantErr: true},
		{name: "rejects non-array", payload: `{"formats":"INC"}`, wantErr: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var request EditRequest
			err := json.Unmarshal([]byte(tc.payload), &request)
			if (err != nil) != tc.wantErr {
				t.Fatalf("json.Unmarshal() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err == nil && !reflect.DeepEqual([]string(request.Formats), tc.want) {
				t.Fatalf("formats = %#v, want %#v", request.Formats, tc.want)
			}
		})
	}
}

func TestWidgetTemplateBSONFormats(t *testing.T) {
	t.Run("empty ticket formats are persisted as an array", func(t *testing.T) {
		result := marshalBSONDocument(t, transformModelToDocument(transformEditRequestToModel(EditRequest{Type: view.WidgetTemplateTypeTicketFormats})))
		formats, ok := result["formats"].(bson.A)
		if !ok || len(formats) != 0 {
			t.Fatalf("formats = %#v, want empty BSON array", result["formats"])
		}
	})

	t.Run("non-ticket formats omit populated inner formats", func(t *testing.T) {
		model := view.WidgetTemplate{
			ID:      "template-id",
			Title:   "More info",
			Type:    view.WidgetTemplateTypeAlarmMoreInfos,
			Content: "body",
			Formats: []string{"stale"},
		}
		result := marshalBSONDocument(t, transformModelToDocument(model))
		if _, ok := result["formats"]; ok {
			t.Fatalf("BSON document exposed formats: %#v", result["formats"])
		}
		if result["_id"] != "template-id" || result["title"] != "More info" || result["type"] != view.WidgetTemplateTypeAlarmMoreInfos || result["content"] != "body" {
			t.Fatalf("BSON document lost template fields: %#v", result)
		}
	})

	t.Run("ticket formats persist exact values and template fields", func(t *testing.T) {
		created := datetime.NewCpsTime(1700000000)
		updated := datetime.NewCpsTime(1700000100)
		model := view.WidgetTemplate{
			ID:          "ticket-template-id",
			Title:       "Ticket formats",
			Type:        view.WidgetTemplateTypeTicketFormats,
			Columns:     []view.WidgetColumn{{Value: "name", Label: "Name"}},
			Content:     "content",
			Actions:     []string{"open"},
			SortColumns: []view.WidgetSortColumn{{SortBy: "name", Sort: "asc"}},
			Formats:     []string{"^INC-[0-9]+$", "", "  exact  "},
			Author:      "author-id",
			Created:     created,
			Updated:     updated,
		}
		result := marshalBSONDocument(t, transformModelToDocument(model))
		formats, ok := result["formats"].(bson.A)
		if !ok || !reflect.DeepEqual([]any(formats), []any{"^INC-[0-9]+$", "", "  exact  "}) {
			t.Fatalf("formats = %#v, want exact BSON array", result["formats"])
		}
		if result["_id"] != model.ID || result["title"] != model.Title || result["type"] != model.Type || result["author"] != model.Author ||
			!reflect.DeepEqual(result["columns"], marshalBSONValue(t, model.Columns)) || result["content"] != model.Content ||
			!reflect.DeepEqual(result["actions"], marshalBSONValue(t, model.Actions)) || !reflect.DeepEqual(result["sort_columns"], marshalBSONValue(t, model.SortColumns)) ||
			result["created"] != marshalBSONValue(t, model.Created) || result["updated"] != marshalBSONValue(t, model.Updated) {
			t.Fatalf("BSON document lost template fields: %#v", result)
		}
	})

	t.Run("outer formats override embedded formats", func(t *testing.T) {
		formats := []string{"outer"}
		document := widgetTemplateDocument{
			WidgetTemplate: view.WidgetTemplate{
				Type:    view.WidgetTemplateTypeTicketFormats,
				Formats: []string{"embedded"},
			},
			Formats: &formats,
		}
		result := marshalBSONDocument(t, document)
		if !reflect.DeepEqual(result["formats"], bson.A{"outer"}) {
			t.Fatalf("formats = %#v, want outer formats", result["formats"])
		}
	})
}

func TestTransformEditRequestToUpdateFormats(t *testing.T) {
	t.Run("ticket formats update writes empty array", func(t *testing.T) {
		update := transformEditRequestToUpdate(transformEditRequestToModel(EditRequest{Type: view.WidgetTemplateTypeTicketFormats}))
		set, ok := update["$set"].(widgetTemplateDocument)
		if !ok {
			t.Fatalf("$set = %#v, want widgetTemplateDocument", update["$set"])
		}
		document := marshalBSONDocument(t, set)
		formats, ok := document["formats"].(bson.A)
		if !ok || len(formats) != 0 {
			t.Fatalf("formats = %#v, want empty BSON array", document["formats"])
		}
	})

	t.Run("other types unset formats", func(t *testing.T) {
		update := transformEditRequestToUpdate(transformEditRequestToModel(EditRequest{Type: view.WidgetTemplateTypeAlarmMoreInfos}))
		unset, ok := update["$unset"].(bson.M)
		if !ok || unset["formats"] == nil {
			t.Fatalf("$unset = %#v, want formats", update["$unset"])
		}
	})
}

func marshalBSONDocument(t *testing.T, value any) bson.M {
	t.Helper()
	data, err := bson.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var result bson.M
	if err := bson.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func marshalBSONValue(t *testing.T, value any) any {
	t.Helper()
	data, err := bson.Marshal(bson.M{"value": value})
	if err != nil {
		t.Fatal(err)
	}
	var result bson.M
	if err := bson.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	return result["value"]
}

func TestEditRequestUnmarshalPreservesID(t *testing.T) {
	request := EditRequest{ID: "existing-template-id", Author: "existing-author"}
	if err := json.Unmarshal([]byte(`{"title":"Updated","type":"ticket_formats","formats":[]}`), &request); err != nil {
		t.Fatal(err)
	}
	if request.ID != "existing-template-id" {
		t.Fatalf("ID = %q, want existing-template-id", request.ID)
	}
	if request.Author != "existing-author" {
		t.Fatalf("Author = %q, want existing-author", request.Author)
	}
}

func TestValidateEditRequestFormats(t *testing.T) {
	v := validator.New()
	v.RegisterStructValidation(ValidateEditRequest, EditRequest{})

	t.Run("allows empty and arbitrary format strings", func(t *testing.T) {
		err := v.Struct(EditRequest{
			Title:   "Ticket formats",
			Type:    view.WidgetTemplateTypeTicketFormats,
			Formats: []string{"", "not a regexp [", "  keep spaces  "},
		})
		if err != nil {
			t.Fatalf("validator rejected formats: %v", err)
		}
	})

	t.Run("allows empty formats and rejects incompatible fields", func(t *testing.T) {
		validErr := v.Struct(EditRequest{
			Title: "Ticket formats",
			Type:  view.WidgetTemplateTypeTicketFormats,
		})
		if validErr != nil {
			t.Fatalf("validator rejected empty formats: %v", validErr)
		}

		invalidErr := v.Struct(EditRequest{
			Title:   "Ticket formats",
			Type:    view.WidgetTemplateTypeTicketFormats,
			Content: "not allowed",
		})
		if invalidErr == nil {
			t.Fatal("validator accepted content on a ticket_formats template")
		}
	})

	t.Run("rejects formats on other template types", func(t *testing.T) {
		err := v.Struct(EditRequest{
			Title:   "More info",
			Type:    view.WidgetTemplateTypeAlarmMoreInfos,
			Content: "template",
			Formats: []string{"INC"},
		})
		if err == nil {
			t.Fatal("validator accepted formats on an alarm_more_infos template")
		}
	})
}

func TestResponseFormatsJSON(t *testing.T) {
	t.Run("empty ticket formats are returned as an array", func(t *testing.T) {
		response := Response{Type: view.WidgetTemplateTypeTicketFormats}
		normalizeFormatsResponse(&response)
		data, err := json.Marshal(response)
		if err != nil {
			t.Fatal(err)
		}
		var result map[string]any
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatal(err)
		}
		formats, ok := result["formats"].([]any)
		if !ok || len(formats) != 0 {
			t.Fatalf("formats = %#v, want empty array", result["formats"])
		}
	})

	t.Run("other template types do not expose formats", func(t *testing.T) {
		data, err := json.Marshal(Response{
			Type: view.WidgetTemplateTypeAlarmMoreInfos,
		})
		if err != nil {
			t.Fatal(err)
		}
		var result map[string]any
		if err := json.Unmarshal(data, &result); err != nil {
			t.Fatal(err)
		}
		if _, ok := result["formats"]; ok {
			t.Fatalf("response exposed formats: %s", data)
		}
	})
}
