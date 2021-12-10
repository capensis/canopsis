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

type EditPositionRequest struct {
	Items []string `json:"items" binding:"required,notblank,unique"`
}

func (r EditPositionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Items)
}

func (r *EditPositionRequest) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &r.Items)
}
