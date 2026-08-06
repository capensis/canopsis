package pbehaviortype

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type ListRequest struct {
	pagination.Query
	Sort    string   `form:"sort" binding:"oneoforempty=asc desc"`
	SortBy  string   `form:"sort_by,default:name" binding:"oneoforempty=_id name priority created updated"`
	Author  string   `form:"author"`
	Default *bool    `form:"default"`
	Types   []string `form:"type"`
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
	ID          string              `bson:"_id" json:"_id"`
	Name        string              `bson:"name" json:"name"`
	Description string              `bson:"description" json:"description"`
	Type        string              `bson:"type" json:"type"`
	Priority    int64               `bson:"priority" json:"priority"`
	IconName    string              `bson:"icon_name" json:"icon_name"`
	Color       string              `bson:"color" json:"color"`
	Created     *datetime.CpsTime   `bson:"created" json:"created" swaggertype:"integer"`
	Updated     *datetime.CpsTime   `bson:"updated" json:"updated" swaggertype:"integer"`
	Author      *response.Reference `bson:"author" json:"author"`
	Default     *bool               `bson:"default" json:"default"`
}
