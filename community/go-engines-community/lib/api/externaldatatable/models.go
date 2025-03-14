package externaldatatable

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
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
	Type        *int   `json:"type" binding:"required,oneof=0 1"`
	Name        string `json:"name" binding:"required,table_name"`
	Description string `json:"description" binding:"max=500"`
	Author      string `json:"author" swaggerignore:"true"`
}

type UpdateRequest struct {
	EditRequest
	ID          string `json:"-"`
	ColumnTypes []int  `json:"column_types" binding:"required,dive,oneof=0 1 2"`
}

type ImportCompleteRequest struct {
	ColumnTypes []int `json:"column_types" binding:"required,dive,oneof=0 1 2"`
}

type ListDataRequest struct {
	pagination.FilteredQuery
	SearchBy []string `json:"search_by" form:"search_by[]"`
	SortBy   string   `json:"sort_by" form:"sort_by"`
}

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
}

type ExportFetchParameters struct {
	ID       string   `json:"_id" swaggerignore:"true"`
	Search   string   `json:"search"`
	SearchBy []string `json:"search_by"`
}

type ExportRequest struct {
	ExportFetchParameters
	Fields    export.Fields `json:"fields"`
	Separator string        `json:"separator" binding:"oneoforempty=comma semicolon tab space"`
}

type Response struct {
	ID                string           `bson:"_id" json:"_id"`
	Type              int              `bson:"type" json:"type"`
	Name              string           `bson:"name" json:"name"`
	Description       string           `bson:"description" json:"description"`
	Columns           []string         `bson:"columns" json:"columns"`
	ColumnTypes       []int            `bson:"column_types" json:"column_types"`
	ColumnLengths     []int            `bson:"column_lengths" json:"-"`
	FromConfig        bool             `bson:"from_config" json:"from_config"`
	RemovedFromConfig bool             `bson:"removed_from_config" json:"removed_from_config"`
	Created           datetime.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Updated           datetime.CpsTime `bson:"updated" json:"updated" swaggertype:"integer"`

	LinkedRules map[string][]struct {
		ID   string `bson:"_id" json:"_id"`
		Name string `bson:"name" json:"name"`
	} `bson:"linked_rules,omitempty" json:"linked_rules,omitempty"`
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
	Separator         rune              `bson:"separator" json:"-"`
	Filepath          string            `bson:"filepath" json:"-"`
	Columns           []string          `bson:"columns" json:"-"`
	ColumnLengths     []int             `bson:"column_lengths" json:"-"`
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

type ExportResponse struct {
	ID     string `json:"_id"`
	Status int64  `json:"status"`
}

type RefResponse struct {
	Reference string `bson:"reference" json:"reference"`
	Type      string `json:"type" bson:"type"`

	// are used in db external data
	Table    Response          `json:"table,omitempty" bson:"table,omitempty"`
	Select   map[string]string `json:"select,omitempty" bson:"select,omitempty"`
	Regexp   map[string]string `json:"regexp,omitempty" bson:"regexp,omitempty"`
	SortBy   string            `json:"sort_by,omitempty" bson:"sort_by,omitempty"`
	Sort     string            `json:"sort,omitempty" bson:"sort,omitempty" binding:"oneoforempty=asc desc"`
	Optional bool              `json:"optional,omitempty" bson:"optional,omitempty"`

	// are used in api external data
	Request *request.Parameters `bson:"request,omitempty" json:"request,omitempty"`
}
