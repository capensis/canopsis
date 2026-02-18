package pbehaviorreason

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type ListRequest struct {
	pagination.Query
	Search string `form:"search"`
	Sort   string `form:"sort" binding:"oneoforempty=asc desc"`
	SortBy string `form:"sort_by,default:name" binding:"oneoforempty=name created updated"`
}

type EditRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description" binding:"required,max=255"`
	Author      string `json:"author" swaggerignore:"true"`
}

type Response struct {
	ID          string            `bson:"_id" json:"_id"`
	Name        string            `bson:"name" json:"name"`
	Description string            `bson:"description" json:"description"`
	Author      *author.Author    `bson:"author,omitempty" json:"author,omitempty"`
	Created     *datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated     *datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
}
