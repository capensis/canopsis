package eventfilter

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"strings"

	libcontext "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/context"
)

// EntityDataSourceFactory is a factory for EntityDataSourceGetter. It
// implements the DataSourceFactory interface.
type EntityDataSourceFactory struct {
	enrichmentCenter libcontext.EnrichmentCenter
	enrichFields     libcontext.EnrichFields
}

// NewEntityDataSourceFactory creates a new EntityDataSourceFactory.
func NewEntityDataSourceFactory(enrichmentCenter libcontext.EnrichmentCenter, enrichFields libcontext.EnrichFields) EntityDataSourceFactory {
	return EntityDataSourceFactory{
		enrichmentCenter: enrichmentCenter,
		enrichFields:     enrichFields,
	}
}

// Create returns a new empty EntityDataSourceGetter.
func (p EntityDataSourceFactory) Create(_ mongo.DbClient, parameters map[string]interface{}) (DataSourceGetter, error) {
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
	EnrichmentCenter libcontext.EnrichmentCenter
	EnrichFields     libcontext.EnrichFields
}

// Get returns the entity corresponding to an event.
func (g EntityDataSourceGetter) Get(ctx context.Context, parameters DataSourceGetterParameters, report *Report) (interface{}, error) {
	if parameters.Event.IsContextable() {
		entity, updatedServices, err := g.EnrichmentCenter.Handle(ctx, parameters.Event, g.EnrichFields)
		if err != nil {
			return nil, fmt.Errorf("unable to get entity: %w", err)
		}

		if report != nil {
			report.UpdatedEntityServices = report.UpdatedEntityServices.Add(updatedServices)
		}

		if entity == nil || !entity.Enabled {
			return nil, nil
		}

		return *entity, nil
	}

	entity, err := g.EnrichmentCenter.Get(ctx, parameters.Event)
	if err != nil || entity == nil || !entity.Enabled {
		return nil, err
	}

	return *entity, nil
}
