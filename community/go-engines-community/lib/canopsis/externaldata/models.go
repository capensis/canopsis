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
	ColumnTypeNoType = iota
	ColumnTypeFilter
	ColumnTypeContext
)

const (
	// IDColumnName uses "_id" because it's not possible to change primary field name in MongoDB.
	IDColumnName             = "_id"
	PostgresDefaultColumnLen = 255
	PostgresIDColumnLen      = 36 // uuid len
)

const (
	mongoPrefix    = "externaldata_"
	postgresSchema = "externaldata"
)

type Table struct {
	ID            string           `bson:"_id,omitempty"`
	Type          int              `bson:"type"`
	Name          string           `bson:"name"`
	Description   string           `bson:"description,omitempty"`
	ColumnTypes   map[string]int   `bson:"column_types,omitempty"`
	ColumnLengths map[string]int   `bson:"column_lengths,omitempty"`
	FromConfig    bool             `bson:"from_config,omitempty"`
	Author        string           `bson:"author,omitempty"`
	Created       datetime.CpsTime `bson:"created,omitempty"`
	Updated       datetime.CpsTime `bson:"updated"`
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
