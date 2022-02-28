package userpreferences

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgetfilter"

type EditRequest struct {
	Widget  string                 `json:"widget" binding:"required"`
	Content map[string]interface{} `json:"content" binding:"required"`
}

type Response struct {
	Widget  string                 `bson:"widget" json:"widget"`
	Content map[string]interface{} `bson:"content" json:"content"`

	Filters []widgetfilter.Response `bson:"filters" json:"filters"`
}
