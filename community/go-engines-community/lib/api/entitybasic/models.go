package entitybasic

import (
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
	ID            string               `json:"-"`
	Description   string               `json:"description" binding:"max=255"`
	Enabled       *bool                `json:"enabled" binding:"required"`
	Category      string               `json:"category"`
	ImpactLevel   int64                `json:"impact_level" binding:"required,min=1,max=10"`
	Infos         []entity.InfoRequest `json:"infos" binding:"dive"`
	SliAvailState *int64               `json:"sli_avail_state" binding:"required,min=0,max=3"`
	Coordinates   *types.Coordinates   `json:"coordinates"`
	Author        string               `json:"author" swaggerignore:"true"`
}
