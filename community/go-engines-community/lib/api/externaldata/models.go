package externaldata

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"github.com/jackc/pgx/v5"
)

const (
	TypeMongoDB = iota
	TypePostgreSQL
)

const (
	ColumnTypeNoType = iota
	ColumnTypeFilter
	ColumnTypeContext
)

const (
	ImportStatusCreated = iota
	ImportStatusRunning
	ImportStatusSucceeded
	ImportStatusFailed
)

const (
	mongoPrefix              = "externaldata_"
	postgresSchema           = "externaldata"
	postgresDefaultColumnLen = 255
	tmpTablePrefix           = "tmp_"
	// Use "_id" because it's not possible to change primary field name in MongoDB.
	idColumnName = "_id"
	uuidLen      = 36
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id name description type"`
}

type CreateRequest struct {
	Type        *int   `json:"type" binding:"required,oneof=0 1"`
	Name        string `json:"name" binding:"required,table_name"`
	Description string `json:"description" binding:"max=255"`
	Author      string `json:"author" swaggerignore:"true"`
}

type ImportCompleteRequest struct {
	Columns map[string]int `json:"columns" binding:"required"`
}

type ListDataRequest struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by"`
}

type Response struct {
	ID          string           `bson:"_id" json:"_id"`
	Type        int              `bson:"type" json:"type"`
	Name        string           `bson:"name" json:"name"`
	Description string           `bson:"description" json:"description"`
	Columns     map[string]int   `bson:"columns" json:"columns"`
	Created     datetime.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Updated     datetime.CpsTime `bson:"updated" json:"updated" swaggertype:"integer"`
}

type Document struct {
	ID          string           `bson:"_id"`
	Type        int              `bson:"type"`
	Name        string           `bson:"name"`
	Description string           `bson:"description"`
	Columns     map[string]int   `bson:"columns"`
	Author      string           `bson:"author"`
	Created     datetime.CpsTime `bson:"created"`
	Updated     datetime.CpsTime `bson:"updated"`
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
	Data       []map[string]string `bson:"data" json:"data"`
	TotalCount int64               `bson:"total_count" json:"total_count"`
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
	Columns           []string          `bson:"columns" json:"-"`
	Created           datetime.CpsTime  `bson:"created" json:"-"`
	LastPing          *datetime.CpsTime `bson:"last_ping" json:"-"`
	Retries           int64             `bson:"retries" json:"-"`
}

func getCollectionName(name string) string {
	return mongoPrefix + name
}

func getPostgresTableName(name string) string {
	return getPostgresTableIdentifier(name).Sanitize()
}

func getPostgresTableIdentifier(name string) pgx.Identifier {
	return pgx.Identifier{postgresSchema, name}
}
