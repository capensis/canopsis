package viewtab

import (
	"encoding/json"
)

type EditRequest struct {
	ID     string `json:"-"`
	Title  string `json:"title" binding:"required,max=255"`
	View   string `json:"view" binding:"required"`
	Author string `json:"author" swaggerignore:"true"`
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

type CopyRequest struct {
	View   string `json:"view" binding:"required"`
	Author string `json:"author" swaggerignore:"true"`
}
