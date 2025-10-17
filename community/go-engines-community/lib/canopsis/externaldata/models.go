package externaldata

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"github.com/jackc/pgx/v5"
)

const (
	TypeMongoDB = iota
	TypePostgreSQL
)

const (
	ColumnTagNoTag = iota
	ColumnTagFilter
	ColumnTagContext
)

const (
	ColumnTypeUnknown = iota
	ColumnTypeString
	ColumnTypeBoolean
	ColumnTypeNumber
	ColumnTypeStringArray
	ColumnTypeDateTime
	ColumnTypeTimestamp
)

const (
	// IDColumnName uses "_id" because it's not possible to change the primary field name in MongoDB.
	IDColumnName = "_id"
)

const (
	mongoPrefix    = "externaldata_"
	postgresSchema = "externaldata"
)

type ColumnConfig struct {
	Name string `bson:"name" json:"name"`
	Type int    `bson:"type" json:"type"`
	Tag  *int   `bson:"tag,omitempty" json:"tag,omitempty"`
}

type Table struct {
	ID                string           `bson:"_id,omitempty"`
	Type              int              `bson:"type"`
	Name              string           `bson:"name"`
	Description       string           `bson:"description"`
	ColumnConfigs     []ColumnConfig   `bson:"column_configs,omitempty"`
	FromConfig        bool             `bson:"from_config,omitempty"`
	RemovedFromConfig bool             `bson:"removed_from_config,omitempty"`
	Author            string           `bson:"author,omitempty"`
	Created           datetime.CpsTime `bson:"created,omitempty"`
	Updated           datetime.CpsTime `bson:"updated"`
}

func (t *Table) GetDBName() string {
	if t.Type == TypeMongoDB {
		return GetMongoCollectionName(t.Name, t.FromConfig)
	}

	return GetPostgresTableName(t.Name)
}

func GetMongoCollectionName(name string, fromConfig bool) string {
	if fromConfig {
		return name
	}

	return mongoPrefix + name
}

func GetPostgresTableName(name string) string {
	return GetPostgresTableIdentifier(name).Sanitize()
}

func GetPostgresTableIdentifier(name string) pgx.Identifier {
	return pgx.Identifier{postgresSchema, name}
}
