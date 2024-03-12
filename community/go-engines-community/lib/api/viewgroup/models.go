package viewgroup

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type ListRequest struct {
	pagination.Query
	Search      string `form:"search" json:"search"`
	WithViews   bool   `form:"with_views" json:"with_views"`
	WithTabs    bool   `form:"with_tabs" json:"with_tabs"`
	WithWidgets bool   `form:"with_widgets" json:"with_widgets"`
	WithFlags   bool   `form:"with_flags" json:"with_flags"`
	WithPrivate bool   `form:"with_private" json:"with_private"`

	UserID string `form:"-" json:"-"`
}

type EditRequest struct {
	ID     string `json:"-"`
	Title  string `json:"title" binding:"required,max=255"`
	Author string `json:"author" swaggerignore:"true"`
}

type ViewGroup struct {
	ID        string            `bson:"_id" json:"_id,omitempty"`
	Title     string            `bson:"title" json:"title"`
	Author    *author.Author    `bson:"author" json:"author,omitempty"`
	Views     *[]view.Response  `bson:"views,omitempty" json:"views,omitempty"`
	Created   *datetime.CpsTime `bson:"created" json:"created,omitempty" swaggertype:"integer"`
	Updated   *datetime.CpsTime `bson:"updated" json:"updated,omitempty" swaggertype:"integer"`
	Deletable *bool             `bson:"deletable,omitempty" json:"deletable,omitempty"`

	IsPrivate bool `bson:"is_private" json:"is_private"`
}

type AggregationResult struct {
	Data       []ViewGroup `bson:"data" json:"data"`
	TotalCount int64       `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
