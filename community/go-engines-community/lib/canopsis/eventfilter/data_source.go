package eventfilter

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// DataSourceFactory is an interface for plugins that provide a DataSource.
type DataSourceFactory interface {
	// Create returns a new empty DataSourceGetter.
	Create(dbClient mongo.DbClient, parameters map[string]interface{}) (DataSourceGetter, error)
}

// DataSourceGetterParameters is a type containing the parameters that can be
// used by a data source.
type DataSourceGetterParameters struct {
	// Event is the event for which the data needs to be fetched.
	Event types.Event

	// RegexMatch contains the values of the sub-expressions of the regular
	// expressions used in the pattern.
	RegexMatch pattern.EventRegexMatches
}

// DataSourceBase is a type containing the fields that are common to all data
// sources.
type DataSourceBase struct {
	// Type is the string used to determine the DataSourceFactory that should
	// be used to get a DataSourceGetter.
	Type string `bson:"type"`

	// Parameters contains the parameters of the data source. They will be used
	// as arguments to DataSourceFactory.Create.
	Parameters map[string]interface{} `bson:",inline"`
}

// DataSourceGetter is an interface for an external data source that can be
// used in an event filter rule.
type DataSourceGetter interface {
	// Get returns the data corresponding to a set of parameters.
	Get(ctx context.Context, parameters DataSourceGetterParameters, report *Report) (interface{}, error)
}

// DataSource is a type that represents an external data source that can be
// used in an event filter rule.
type DataSource struct {
	DataSourceBase
	DataSourceGetter
}

func (d *DataSource) UnmarshalBSONValue(_ bsontype.Type, b []byte) error {
	return bson.Unmarshal(b, &d.DataSourceBase)
}
