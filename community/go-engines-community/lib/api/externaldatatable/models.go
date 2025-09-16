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
	ImportStatusPreviewSucceeded
	ImportStatusPreviewFailed
)

const (
	JobTypeImport = iota
	JobTypePreview
)

const (
	// MaxStringLenStr and MaxIDLenStr are strings to avoid conversion.
	MaxStringLenStr = "255"
	MaxIDLenStr     = "36" // uuid len

	MaxStringLen = 255
)

const tmpTablePrefix = "tmp_"

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string   `form:"sort_by" binding:"oneoforempty=_id name description type"`
	IDs    []string `form:"ids[]"`
}

type EditRequest struct {
	Type        *int   `json:"type" binding:"required,oneof=0 1"`
	Name        string `json:"name" binding:"required,table_name"`
	Description string `json:"description" binding:"max=500"`
	Author      string `json:"author" swaggerignore:"true"`
}

type UpdateRequest struct {
	EditRequest
	ID         string `json:"-"`
	ColumnTags []int  `json:"column_tags" binding:"required,dive,oneof=0 1 2"`
}

type ImportCompleteRequest struct {
	ColumnTags []int `json:"column_tags" binding:"required,dive,oneof=0 1 2"`
}

type ListPreviewRequest struct {
	pagination.Query
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

type ColumnConfig struct {
	BaseColumnConfig `bson:",inline"`
	Tag              *int `bson:"tag,omitempty" json:"tag,omitempty" binding:"omitempty,oneof=0 1 2"`
}

func (c *ColumnConfig) IsRegexp() bool {
	return c.Type == externaldata.ColumnTypeRegexp
}

type BaseColumnConfig struct {
	Name string `bson:"name" json:"name" binding:"required"`
	// Possible type values.
	//   * `1` - type string
	//   * `2` - type boolean
	//   * `3` - type number
	//   * `4` - type string_array
	//   * `5` - type datetime
	//   * `6` - type timestamp
	//   * `7` - type regexp
	Type int `bson:"type" json:"type" binding:"required,min=1,max=7"`
	// Possible thousands separator values.
	//   * `dot` - dot separator
	//   * `comma` - comma separator
	//   * `space` - space separator
	ThousandsSeparator string `bson:"thousands_separator,omitempty" json:"thousands_separator,omitempty" binding:"oneoforempty=dot comma space"`
	// Possible decimal separator values.
	//   * `dot` - dot separator
	//   * `comma` - comma separator
	DecimalSeparator string `bson:"decimal_separator,omitempty" json:"decimal_separator,omitempty" binding:"oneoforempty=dot comma"`
	// Possible string array types.
	//   * `1` - json array
	//   * `2` - custom separator array
	StringArrayType      int    `bson:"string_array_type,omitempty" json:"string_array_type,omitempty" binding:"required_if=Type 4,omitempty,oneof=1 2"`
	StringArraySeparator string `bson:"string_array_separator,omitempty" json:"string_array_separator,omitempty" binding:"required_if=StringArrayType 2"`
}

type Table struct {
	ID                string                      `bson:"_id" json:"_id"`
	Type              int                         `bson:"type" json:"type"`
	Name              string                      `bson:"name" json:"name"`
	Description       string                      `bson:"description" json:"description"`
	ColumnConfigs     []externaldata.ColumnConfig `bson:"column_configs" json:"column_configs"`
	FromConfig        bool                        `bson:"from_config" json:"from_config"`
	RemovedFromConfig bool                        `bson:"removed_from_config" json:"removed_from_config"`
	Created           datetime.CpsTime            `bson:"created" json:"created" swaggertype:"integer"`
	Updated           datetime.CpsTime            `bson:"updated" json:"updated" swaggertype:"integer"`

	LinkedRules map[string][]struct {
		ID   string `bson:"_id" json:"_id"`
		Name string `bson:"name" json:"name"`
	} `bson:"linked_rules,omitempty" json:"linked_rules,omitempty"`
}

func (t *Table) getDBTableName() string {
	if t.Type == externaldata.TypeMongoDB {
		return externaldata.GetMongoCollectionName(t.Name, t.FromConfig)
	}

	return externaldata.GetPostgresTableName(t.Name)
}

func (t *Table) getPostgresTableIdentifier() pgx.Identifier {
	if t.FromConfig {
		return pgx.Identifier{t.Name}
	}

	return externaldata.GetPostgresTableIdentifier(t.Name)
}

func (t *Table) getColumns() []string {
	addPriorityColumn := false

	columns := make([]string, len(t.ColumnConfigs))
	for i := range t.ColumnConfigs {
		columns[i] = t.ColumnConfigs[i].Name
		if t.ColumnConfigs[i].Type == externaldata.ColumnTypeRegexp {
			addPriorityColumn = true
		}
	}

	if addPriorityColumn {
		columns = append(columns, priorityColumnName)
	}

	return columns
}

type AggregationResult struct {
	Data       []Table `bson:"data" json:"data"`
	TotalCount int64   `bson:"total_count" json:"total_count"`
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
	Created           datetime.CpsTime  `bson:"created" json:"-"`
	LastPing          *datetime.CpsTime `bson:"last_ping" json:"-"`
	Retries           int64             `bson:"retries" json:"-"`
	JobType           int               `bson:"job_type" json:"-"`

	ColumnConfigs     []ColumnConfig `bson:"column_configs,omitempty" json:"column_configs,omitempty"`
	PrevColumnConfigs []ColumnConfig `bson:"prev_column_configs,omitempty" json:"-"`

	ErrorInfo map[string]ErrorInfo `bson:"error_info,omitempty" json:"error_info,omitempty"`

	FailReason string `bson:"fail_reason,omitempty" json:"fail_reason,omitempty"`
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
	Table    Table             `json:"table,omitempty" bson:"table,omitempty"`
	Select   map[string]string `json:"select,omitempty" bson:"select,omitempty"`
	Regexp   map[string]string `json:"regexp,omitempty" bson:"regexp,omitempty"`
	SortBy   string            `json:"sort_by,omitempty" bson:"sort_by,omitempty"`
	Sort     string            `json:"sort,omitempty" bson:"sort,omitempty" binding:"oneoforempty=asc desc"`
	Optional bool              `json:"optional,omitempty" bson:"optional,omitempty"`

	// are used in api external data
	Request *request.Parameters `bson:"request,omitempty" json:"request,omitempty"`
}

type PreviewRequest struct {
	ColumnConfigs []ColumnConfig `json:"column_configs" binding:"required,dive"`
}

type ErrorInfo struct {
	Rows     []int    `json:"rows,omitempty"`
	Messages []string `json:"messages,omitempty"`
}
