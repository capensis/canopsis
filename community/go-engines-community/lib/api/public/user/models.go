package user

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type ListRequest struct {
	pagination.Query
	Sort    string   `form:"sort"`
	Enabled *bool    `form:"enabled"`
	Author  string   `form:"author"`
	Sources []string `form:"source"`
	Roles   []string `form:"role"`
}

type CreateRequest struct {
	EditRequest
	Source     string `json:"source" binding:"required,max=10000"`
	ExternalID string `json:"external_id" binding:"required,max=10000"`
}

type EditRequest struct {
	Name      string   `json:"name" binding:"required,max=255"`
	Firstname string   `json:"firstname" binding:"max=255"`
	Lastname  string   `json:"lastname" binding:"max=255"`
	Email     string   `json:"email" binding:"required,email"`
	Roles     []string `json:"roles" binding:"required,notblank"`
	Enabled   *bool    `json:"enabled" binding:"required"`
	Author    string   `json:"author" swaggerignore:"true"`
}

type PatchRequest struct {
	Name      *string  `json:"name" binding:"omitempty,max=255"`
	Firstname *string  `json:"firstname" binding:"omitempty,max=255"`
	Lastname  *string  `json:"lastname" binding:"omitempty,max=255"`
	Email     *string  `json:"email" binding:"omitempty,email"`
	Roles     []string `json:"roles" binding:"omitempty,notblank"`
	Enabled   *bool    `json:"enabled"`
	Author    string   `json:"author" swaggerignore:"true"`
}

type BulkUpdateRequestItem struct {
	ID string `json:"_id" binding:"required"`
	EditRequest
}

type BulkPatchRequestItem struct {
	ID string `json:"_id" binding:"required"`
	PatchRequest
}

type Response struct {
	ID         string               `bson:"_id" json:"_id"`
	Name       string               `bson:"name" json:"name"`
	Lastname   string               `bson:"lastname" json:"lastname"`
	Firstname  string               `bson:"firstname" json:"firstname"`
	Email      string               `bson:"email" json:"email"`
	Roles      []response.Reference `bson:"roles" json:"roles"`
	Enabled    bool                 `bson:"enabled" json:"enabled"`
	ExternalID string               `bson:"external_id" json:"external_id"`
	Source     string               `bson:"source" json:"source"`
	Author     *response.Reference  `bson:"author" json:"author"`
	Created    *datetime.CpsTime    `bson:"created" json:"created" swaggertype:"integer"`
	Updated    *datetime.CpsTime    `bson:"updated" json:"updated" swaggertype:"integer"`
}
