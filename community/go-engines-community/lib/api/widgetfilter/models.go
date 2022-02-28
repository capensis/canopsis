package widgetfilter

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EditRequest struct {
	ID        string `json:"-"`
	Widget    string `json:"widget" binding:"required"`
	Title     string `json:"title" binding:"required,max=255"`
	IsPrivate *bool  `json:"is_private" binding:"required"`
	Author    string `json:"author" swaggerignore:"true"`

	common.AlarmPatternFieldsRequest
	common.EntityPatternFieldsRequest
	common.PbehaviorPatternFieldsRequest
}

type Response struct {
	ID        string         `bson:"_id" json:"_id"`
	Widget    string         `bson:"widget" json:"-"`
	Title     string         `bson:"title" json:"title"`
	IsPrivate *bool          `bson:"is_private" json:"is_private,omitempty"`
	Author    string         `bson:"author" json:"author,omitempty"`
	Created   *types.CpsTime `bson:"created" json:"created,omitempty" swaggertype:"integer"`
	Updated   *types.CpsTime `bson:"updated" json:"updated,omitempty" swaggertype:"integer"`

	OldMongoQuery map[string]interface{} `bson:"old_mongo_query" json:"old_mongo_query,omitempty"`

	savedpattern.AlarmPatternFields     `bson:",inline"`
	savedpattern.EntityPatternFields    `bson:",inline"`
	savedpattern.PbehaviorPatternFields `bson:",inline"`
}
