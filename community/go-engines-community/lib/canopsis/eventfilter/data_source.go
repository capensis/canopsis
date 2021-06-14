package eventfilter

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/globalsign/mgo/bson"
)

// DataSourceFactory is an interface for plugins that provide a DataSource.
type DataSourceFactory interface {
	// Create returns a new empty DataSourceGetter.
	Create(parameters map[string]interface{}) (DataSourceGetter, error)
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
	Get(parameters DataSourceGetterParameters) (interface{}, error)
}

// DataSource is a type that represents an external data source that can be
// used in an event filter rule.
type DataSource struct {
	DataSourceBase
	DataSourceGetter
}

// SetBSON unmarshals a BSON value into a DataSource.
func (d *DataSource) SetBSON(raw bson.Raw) error {
	err := raw.Unmarshal(&d.DataSourceBase)
	return err
}
