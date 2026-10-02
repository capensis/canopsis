package widgettemplate

import (
	"bytes"
	"encoding/json"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id title type author.name author.display_name created updated"`
	Type   string `json:"type" form:"type"`
}

type EditRequest struct {
	ID          string                  `json:"-"`
	Title       string                  `json:"title" binding:"required,max=255"`
	Type        string                  `json:"type" binding:"required,oneof=alarm_columns entity_columns alarm_more_infos alarm_export_to_pdf alarm_quick_actions alarm_mass_quick_actions alarm_sort_columns weather_item weather_modal weather_entity ticket_formats"`
	Columns     []view.WidgetColumn     `json:"columns" binding:"dive"`
	Content     string                  `json:"content"`
	Actions     []string                `json:"actions" binding:"omitempty,dive,notblank,max=255"`
	SortColumns []view.WidgetSortColumn `json:"sort_columns" binding:"dive"`
	// Formats preserves ticket format strings verbatim; regexp syntax is not validated by the backend.
	// Empty strings and an empty list are allowed. A null list, null elements, and non-string entries are rejected; omission means an empty list.
	Formats ticketFormats `json:"formats,omitempty" binding:"max=100,dive,max=255" maxLength:"255" swaggertype:"array,string"`
	Author  string        `json:"author" swaggerignore:"true"`
}

type ticketFormats []string

func (f *ticketFormats) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if bytes.Equal(trimmed, []byte("null")) {
		return validation.NewSingleError("value_string_array", "Formats", "formats", nil)
	}

	var values []json.RawMessage
	if err := json.Unmarshal(trimmed, &values); err != nil {
		return validation.NewSingleError("value_string_array", "Formats", "formats", nil)
	}
	formats := make(ticketFormats, len(values))
	for i, value := range values {
		value = bytes.TrimSpace(value)
		if len(value) == 0 || value[0] != '"' {
			return validation.NewSingleError("value_string", "Formats", fmt.Sprintf("formats.%d", i), nil)
		}
		if err := json.Unmarshal(value, &formats[i]); err != nil {
			return validation.NewSingleError("value_string", "Formats", fmt.Sprintf("formats.%d", i), nil)
		}
	}

	*f = formats
	return nil
}

func normalizeFormatsResponse(r *Response) {
	if r.Type == view.WidgetTemplateTypeTicketFormats && r.Formats == nil {
		r.Formats = []string{}
	}
}

type Response struct {
	ID          string                  `bson:"_id" json:"_id"`
	Title       string                  `bson:"title" json:"title"`
	Type        string                  `bson:"type" json:"type"`
	Columns     []view.WidgetColumn     `bson:"columns" json:"columns,omitempty"`
	Content     string                  `bson:"content" json:"content,omitempty"`
	Actions     []string                `bson:"actions" json:"actions,omitempty"`
	SortColumns []view.WidgetSortColumn `bson:"sort_columns" json:"sort_columns,omitempty"`
	// Formats contains the exact strings stored by a ticket_formats template.
	Formats []string          `bson:"formats" json:"formats,omitzero"`
	Author  *author.Author    `bson:"author" json:"author"`
	Created *datetime.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Updated *datetime.CpsTime `bson:"updated" json:"updated" swaggertype:"integer"`
}

type AggregationResult struct {
	Data       []Response `bson:"data" json:"data"`
	TotalCount int64      `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() any {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
