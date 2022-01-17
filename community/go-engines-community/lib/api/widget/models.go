package widget

import (
	"encoding/json"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
)

type EditRequest struct {
	ID             string                 `json:"-"`
	Tab            string                 `json:"tab" binding:"required"`
	Title          string                 `json:"title" binding:"required,max=255"`
	Type           string                 `json:"type" binding:"required,max=255"`
	GridParameters map[string]interface{} `json:"grid_parameters"`
	Parameters     view.Parameters        `json:"parameters"`
	Author         string                 `json:"author" swaggerignore:"true"`
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

type CopyRequest struct {
	Tab    string `json:"tab" binding:"required"`
	Author string `json:"author" swaggerignore:"true"`
}
