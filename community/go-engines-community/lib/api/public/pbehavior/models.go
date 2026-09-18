package pbehavior

import (
	"encoding/json"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/exdate"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

const (
	BulkConnectorActionCreate = "create"
	BulkConnectorActionDelete = "delete"
)

type ListRequest struct {
	pagination.Query
	Sort      string   `form:"sort"`
	Enabled   *bool    `form:"enabled"`
	Types     []string `form:"type"`
	Reasons   []string `form:"reason"`
	Author    string   `form:"author"`
	Inherited *bool    `form:"inherited"`
	Origins   []string `form:"origin"`
}

type EditRequest struct {
	Enabled       *bool             `json:"enabled" binding:"required"`
	Name          string            `json:"name" binding:"required,max=255"`
	Reason        string            `json:"reason" binding:"required"`
	Type          string            `json:"type" binding:"required"`
	EntityPattern pattern.Entity    `json:"entity_pattern" binding:"required"`
	Author        string            `json:"author" swaggerignore:"true"`
	RRule         string            `json:"rrule"`
	Timezone      string            `json:"timezone"`
	Start         *datetime.CpsTime `json:"tstart" binding:"required" swaggertype:"integer"`
	Stop          *datetime.CpsTime `json:"tstop" swaggertype:"integer"`
	Exdates       []ExdateRequest   `json:"exdates" binding:"dive"`
	Color         string            `json:"color" binding:"iscolororempty"`
	Inherited     bool              `json:"inherited"`
}

type ExdateRequest struct {
	exdate.Request
	Type string `json:"type" binding:"required"`
}

type BulkUpdateRequestItem struct {
	EditRequest
	ID string `json:"_id" binding:"required"`
}

type BulkPatchRequestItem struct {
	PatchRequest
	ID string `json:"_id" binding:"required"`
}

type PatchRequest struct {
	Enabled       *bool           `json:"enabled"`
	Name          *string         `json:"name" binding:"omitempty,max=255"`
	Reason        *string         `json:"reason"`
	Type          *string         `json:"type"`
	EntityPattern pattern.Entity  `json:"entity_pattern"`
	Author        string          `json:"author" swaggerignore:"true"`
	RRule         *string         `json:"rrule"`
	Timezone      *string         `json:"timezone"`
	Start         *int64          `json:"tstart" swaggertype:"integer"`
	Stop          NullableTime    `json:"tstop" swaggertype:"integer"`
	Exdates       []ExdateRequest `json:"exdates" binding:"dive"`
	Color         *string         `json:"color" binding:"omitempty,iscolororempty"`
	Inherited     *bool           `json:"inherited"`
}

type NullableTime struct {
	val   *int64
	isSet bool
}

func (t *NullableTime) UnmarshalJSON(data []byte) error {
	t.isSet = true
	if string(data) == "null" {
		return nil
	}

	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	t.val = &i

	return nil
}

type BulkConnectorCreateRequestItem struct {
	Author   string            `json:"author" swaggerignore:"true"`
	Entities []string          `json:"entities" binding:"required,notblank"`
	Origin   string            `json:"origin" binding:"required,max=255"`
	Start    *datetime.CpsTime `json:"tstart" binding:"required" swaggertype:"integer"`
	Stop     *datetime.CpsTime `json:"tstop" binding:"required" swaggertype:"integer"`
	Comment  string            `json:"comment" binding:"required"`
	Name     string            `json:"name" binding:"required,max=255"`
	Reason   string            `json:"reason" binding:"required"`
	Type     string            `json:"type" binding:"required"`
	Color    string            `json:"color" binding:"iscolororempty"`
}

type BulkConnectorDeleteRequestItem struct {
	Author   string            `json:"author" swaggerignore:"true"`
	Entities []string          `json:"entities" binding:"required,notblank"`
	Origin   string            `json:"origin" binding:"required"`
	Start    *datetime.CpsTime `json:"tstart" binding:"required" swaggertype:"integer"`
	Stop     *datetime.CpsTime `json:"tstop" binding:"required" swaggertype:"integer"`
	Comment  string            `json:"comment" binding:"required"`
}

type BulkConnectorEditRequestItem struct {
	Action   string            `json:"action" binding:"oneof=create delete"`
	Author   string            `json:"author" swaggerignore:"true"`
	Entities []string          `json:"entities" binding:"required,notblank"`
	Origin   string            `json:"origin" binding:"required,max=255"`
	Start    *datetime.CpsTime `json:"tstart" binding:"required" swaggertype:"integer"`
	Stop     *datetime.CpsTime `json:"tstop" binding:"required" swaggertype:"integer"`
	Comment  string            `json:"comment" binding:"required"`
	Name     string            `json:"name" binding:"max=255"`
	Reason   string            `json:"reason"`
	Type     string            `json:"type"`
	Color    string            `json:"color" binding:"iscolororempty"`
}

type Response struct {
	ID            string               `bson:"_id" json:"_id"`
	Enabled       bool                 `bson:"enabled" json:"enabled"`
	Name          string               `bson:"name" json:"name"`
	Reason        *response.Reference  `bson:"reason" json:"reason"`
	Type          *response.Reference  `bson:"type" json:"type"`
	EntityPattern pattern.Entity       `bson:"entity_pattern" json:"entity_pattern,omitempty"`
	Origin        string               `bson:"origin" json:"origin"`
	OriginComment string               `bson:"origin_comment" json:"origin_comment"`
	Author        *response.Reference  `bson:"author" json:"author"`
	RRule         string               `bson:"rrule" json:"rrule"`
	Timezone      string               `bson:"timezone" json:"timezone"`
	Start         *datetime.CpsTime    `bson:"tstart" json:"tstart" swaggertype:"integer"`
	Stop          *datetime.CpsTime    `bson:"tstop" json:"tstop" swaggertype:"integer"`
	Exdates       []ExdateResponse     `bson:"exdates" json:"exdates"`
	Exceptions    []response.Reference `bson:"exceptions" json:"exceptions"`
	Color         string               `bson:"color" json:"color"`
	Inherited     bool                 `bson:"inherited" json:"inherited"`
	Created       *datetime.CpsTime    `bson:"created" json:"created" swaggertype:"integer"`
	Updated       *datetime.CpsTime    `bson:"updated" json:"updated" swaggertype:"integer"`
}

type ExdateResponse struct {
	Begin datetime.CpsTime   `bson:"begin" json:"begin" swaggertype:"integer"`
	End   datetime.CpsTime   `bson:"end" json:"end" swaggertype:"integer"`
	Type  response.Reference `bson:"type" json:"type"`
}
