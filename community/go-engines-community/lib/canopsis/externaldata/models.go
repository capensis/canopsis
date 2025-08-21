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

	// MaxStringLenStr and MaxIdLenStr are strings to avoid conversion.
	MaxStringLenStr = "255"
	MaxIdLenStr     = "36" // uuid len

	MaxStringLen = 255
)

const (
	mongoPrefix    = "externaldata_"
	postgresSchema = "externaldata"
)

type ColumnConfig struct {
	BaseColumnConfig `bson:",inline"`
	Tag              *int `bson:"tag,omitempty" json:"tag,omitempty" binding:"omitempty,oneof=0 1 2"`
}

type BaseColumnConfig struct {
	Name string `bson:"name" json:"name" binding:"required"`
	// Possible type values.
	//   * `1` - type string
	//   * `2` - type boolean
	//   * `3` - type number
	//   * `4` - type string_array
	//   * `5` - type datetime
	//	 * `6` - type timestamp
	Type int `bson:"type" json:"type" binding:"required,oneof=1 2 3 4 5 6"`
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
	ID                string           `bson:"_id,omitempty"`
	Type              int              `bson:"type"`
	Name              string           `bson:"name"`
	Description       string           `bson:"description,omitempty"`
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
