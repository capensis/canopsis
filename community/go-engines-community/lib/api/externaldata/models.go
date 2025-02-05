package externaldata

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"github.com/jackc/pgx/v5"
)

const (
	ImportStatusCreated = iota
	ImportStatusRunning
	ImportStatusSucceeded
	ImportStatusFailed
)

const tmpTablePrefix = "tmp_"

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id name description type"`
}

type EditRequest struct {
	ID          string `json:"-"`
	Type        *int   `json:"type" binding:"required,oneof=0 1"`
	Name        string `json:"name" binding:"required,table_name"`
	Description string `json:"description" binding:"max=255"`
	Author      string `json:"author" swaggerignore:"true"`
}

type ImportCompleteRequest struct {
	ColumnTypes map[string]int `json:"column_types" binding:"required"`
}

type ListDataRequest struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by"`
}

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
}

type Response struct {
	ID            string           `bson:"_id" json:"_id"`
	Type          int              `bson:"type" json:"type"`
	Name          string           `bson:"name" json:"name"`
	Description   string           `bson:"description" json:"description"`
	ColumnTypes   map[string]int   `bson:"column_types" json:"column_types"`
	ColumnLengths map[string]int   `bson:"column_lengths" json:"-"`
	FromConfig    bool             `bson:"from_config" json:"from_config"`
	Created       datetime.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Updated       datetime.CpsTime `bson:"updated" json:"updated" swaggertype:"integer"`
}

func (r *Response) getDBTableName() string {
	if r.Type == externaldata.TypeMongoDB {
		return externaldata.GetMongoCollectionName(r.Name, r.FromConfig)
	}

	return externaldata.GetPostgresTableName(r.Name)
}

func (r *Response) getPostgresTableIdentifier() pgx.Identifier {
	if r.FromConfig {
		return pgx.Identifier{r.Name}
	}

	return externaldata.GetPostgresTableIdentifier(r.Name)
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

type AggregationDataResult struct {
	Data       []map[string]any `bson:"data" json:"data"`
	TotalCount int64            `bson:"total_count" json:"total_count"`
}

func (r *AggregationDataResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationDataResult) GetTotal() int64 {
	return r.TotalCount
}

type ImportJob struct {
	ID                string            `bson:"_id" json:"_id"`
	Status            int               `bson:"status" json:"status"`
	Type              int               `bson:"type" json:"-"`
	Table             string            `bson:"table" json:"-"`
	ExternalDataTable string            `bson:"exdt" json:"-"`
	Delimiter         rune              `bson:"delimiter" json:"-"`
	Filepath          string            `bson:"filepath" json:"-"`
	ColumnLengths     map[string]int    `bson:"column_lengths" json:"-"`
	Created           datetime.CpsTime  `bson:"created" json:"-"`
	LastPing          *datetime.CpsTime `bson:"last_ping" json:"-"`
	Retries           int64             `bson:"retries" json:"-"`
}

func (j *ImportJob) getDBTableName() string {
	if j.Type == externaldata.TypeMongoDB {
		return externaldata.GetMongoCollectionName(j.Table, false)
	}

	return externaldata.GetPostgresTableName(j.Table)
}
