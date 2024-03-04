package icon

import (
	"mime/multipart"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type EditRequest struct {
	ID       string                `form:"-"`
	Title    string                `form:"title" binding:"required,max=255"`
	File     *multipart.FileHeader `form:"file" binding:"required"`
	MimeType string                `form:"-"`
}

type PatchRequest struct {
	ID       string                `form:"-"`
	Title    string                `form:"title" binding:"max=255"`
	File     *multipart.FileHeader `form:"file"`
	MimeType string                `form:"-"`
}

type Response struct {
	ID      string           `bson:"_id" json:"_id"`
	Title   string           `bson:"title" json:"title"`
	Created datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`

	Storage  string `bson:"storage" json:"-"`
	Etag     string `bson:"etag" json:"-"`
	MimeType string `bson:"mime_type" json:"-"`

	Content string `bson:"-" json:"content,omitempty"`
}

type AggregationResult struct {
	Data       []Response `bson:"data" json:"data"`
	TotalCount int64      `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
