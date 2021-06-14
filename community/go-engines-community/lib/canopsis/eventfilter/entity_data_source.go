package eventfilter

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"strings"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/context"
)

// EntityDataSourceFactory is a factory for EntityDataSourceGetter. It
// implements the DataSourceFactory interface.
type EntityDataSourceFactory struct {
	enrichmentCenter context.EnrichmentCenter
	enrichFields     context.EnrichFields
}

// NewEntityDataSourceFactory creates a new EntityDataSourceFactory.
func NewEntityDataSourceFactory(enrichmentCenter context.EnrichmentCenter, enrichFields context.EnrichFields) EntityDataSourceFactory {
	return EntityDataSourceFactory{
		enrichmentCenter: enrichmentCenter,
		enrichFields:     enrichFields,
	}
}

// Create returns a new empty EntityDataSourceGetter.
func (p EntityDataSourceFactory) Create(parameters map[string]interface{}) (DataSourceGetter, error) {
	if len(parameters) != 0 {
		unexpectedParameters := make([]string, 0, len(parameters))
		for key := range parameters {
			unexpectedParameters = append(unexpectedParameters, key)
		}
		return nil, fmt.Errorf("unexpected parameters for entity data source: %s", strings.Join(unexpectedParameters, ", "))
	}

	return EntityDataSourceGetter{
		EnrichmentCenter: p.enrichmentCenter,
		EnrichFields:     p.enrichFields,
	}, nil
}

// EntityDataSourceGetter is a DataSourceGetter that gets the entity
// corresponding to an event, and creates the related entities (resource,
// component and connector) if they do not exist.
type EntityDataSourceGetter struct {
	EnrichmentCenter context.EnrichmentCenter
	EnrichFields     context.EnrichFields
}

// Get returns the entity corresponding to an event.
func (g EntityDataSourceGetter) Get(parameters DataSourceGetterParameters) (interface{}, error) {
	var entity *types.Entity
	var err error

	if parameters.Event.IsContextable() {
		entity = g.EnrichmentCenter.Handle(parameters.Event, g.EnrichFields)
		if entity == nil {
			return nil, fmt.Errorf("unable to get entity")
		}
	} else {
		entity, err = g.EnrichmentCenter.Get(parameters.Event)
		if err != nil {
			return nil, err
		}
	}

	if entity == nil {
		return nil, nil
	}

	return *entity, nil
}
