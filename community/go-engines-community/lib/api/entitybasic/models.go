package entitybasic

import (
	"encoding/json"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Entity struct {
	entity.Entity `bson:",inline"`
	Description   string `bson:"description" json:"description"`
	SliAvailState int64  `bson:"sli_avail_state" json:"sli_avail_state"`
}

type IdRequest struct {
	ID string `form:"_id" binding:"required"`
}

type EditRequest struct {
	ID            string        `json:"-"`
	Description   string        `json:"description" binding:"max=255"`
	Enabled       *bool         `json:"enabled" binding:"required"`
	Category      string        `json:"category"`
	ImpactLevel   int64         `json:"impact_level" binding:"required,min=1,max=10"`
	Infos         []InfoRequest `json:"infos" binding:"dive"`
	SliAvailState *int64        `json:"sli_avail_state" binding:"required,min=0,max=3"`

	Coordinates *types.Coordinates `json:"coordinates"`
}

type InfoRequest struct {
	Name        string      `json:"name" binding:"required,max=255"`
	Description string      `json:"description" binding:"max=255"`
	Value       interface{} `json:"value"`
}

func (r *InfoRequest) UnmarshalJSON(b []byte) error {
	type Alias InfoRequest
	tmp := Alias{}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	*r = InfoRequest(tmp)

	if r.Value != nil {
		switch v := r.Value.(type) {
		case float64, float32, int, int64, int32, bool, string:
			// do nothing
		case []interface{}:
			for _, item := range v {
				if item != nil {
					switch item.(type) {
					case float64, float32, int, int64, int32, bool, string:
						// do nothing
					default:
						return fmt.Errorf("invalid type of array item: %T", item)
					}
				}
			}
		default:
			return fmt.Errorf("invalid value type: %T", r.Value)
		}
	}

	return nil
}
