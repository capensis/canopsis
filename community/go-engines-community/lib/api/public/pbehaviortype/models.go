package pbehaviortype

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
)

type ListRequest struct {
	pagination.Query
	Search      string   `form:"search"`
	Sort        string   `form:"sort" binding:"oneoforempty=asc desc"`
	SortBy      string   `form:"sort_by,default:name" binding:"oneoforempty=name priority created updated"`
	OnlyDefault bool     `form:"default"`
	Types       []string `form:"types[]"`
}

type EditRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description" binding:"required,max=255"`
	Type        string `json:"type" binding:"required,oneof=active inactive maintenance pause"`
	Priority    int64  `json:"priority" binding:"required,min=1"`
	Color       string `json:"color" binding:"required,iscolor"`
	Author      string `json:"author" swaggerignore:"true"`
}

type CreateRequest struct {
	EditRequest
	IconName string `json:"icon_name" binding:"required,max=255"`
}

type UpdateRequest struct {
	EditRequest
	IconName string `json:"icon_name" binding:"max=255"`
}

type Response struct {
	pbehavior.Type `bson:",inline"`
	Author         *author.Author `bson:"author" json:"author"`
	Default        *bool          `bson:"default,omitempty" json:"default,omitempty"`
}
