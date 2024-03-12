package widget

import (
	"encoding/json"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgetfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
)

type EditRequest struct {
	Title          string                 `json:"title" binding:"max=255"`
	Type           string                 `json:"type" binding:"required,max=255"`
	GridParameters map[string]interface{} `json:"grid_parameters"`
	Parameters     view.Parameters        `json:"parameters"`
	Author         string                 `json:"author" swaggerignore:"true"`

	Filters []FilterRequest `json:"filters" binding:"dive"`
}

type CreateRequest struct {
	EditRequest
	Tab string `json:"tab" binding:"required"`
}

type UpdateRequest struct {
	EditRequest
	ID string `json:"-"`
}

type FilterRequest struct {
	widgetfilter.BaseEditRequest
	ID string `json:"_id"`
}

type EditGridPositionItemRequest struct {
	ID             string                 `json:"_id"`
	GridParameters map[string]interface{} `json:"grid_parameters"`
}

type EditGridPositionRequest struct {
	Items []EditGridPositionItemRequest `json:"items" binding:"required,notblank"`
}

func (r EditGridPositionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Items)
}

func (r *EditGridPositionRequest) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &r.Items)
}

type Response struct {
	ID             string                 `bson:"_id" json:"_id,omitempty"`
	Title          string                 `bson:"title" json:"title"`
	Tab            string                 `bson:"tab" json:"-"`
	Type           string                 `bson:"type" json:"type"`
	GridParameters map[string]interface{} `bson:"grid_parameters" json:"grid_parameters"`
	Parameters     view.Parameters        `bson:"parameters" json:"parameters"`
	Author         *author.Author         `bson:"author" json:"author,omitempty"`
	Created        *datetime.CpsTime      `bson:"created" json:"created,omitempty" swaggertype:"integer"`
	Updated        *datetime.CpsTime      `bson:"updated" json:"updated,omitempty" swaggertype:"integer"`

	Filters []widgetfilter.Response `bson:"filters" json:"filters"`

	IsPrivate bool `bson:"is_private" json:"is_private"`
}
