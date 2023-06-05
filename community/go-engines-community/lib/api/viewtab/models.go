package viewtab

import (
	"encoding/json"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widget"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
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

type Response struct {
	ID       string             `bson:"_id" json:"_id,omitempty"`
	Title    string             `bson:"title" json:"title"`
	View     string             `bson:"view" json:"-"`
	Position int64              `bson:"position" json:"-"`
	Widgets  *[]widget.Response `bson:"widgets" json:"widgets,omitempty"`
	Author   *author.Author     `bson:"author" json:"author,omitempty"`
	Created  *types.CpsTime     `bson:"created" json:"created,omitempty"`
	Updated  *types.CpsTime     `bson:"updated" json:"updated,omitempty"`
}
