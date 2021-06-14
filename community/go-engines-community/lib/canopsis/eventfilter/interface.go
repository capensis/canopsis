package eventfilter

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

// Service is an interface for the service that manages the event filter.
type Service interface {
	// LoadRules loads the event filter rules from the database, and adds them
	// to the service.
	LoadRules() error

	// LoadDataSourceFactories loads the data source factories and adds them to
	// the service.
	LoadDataSourceFactories(enrichmentCenter context.EnrichmentCenter, enrichFields context.EnrichFields, dataSourceDirectory string) error

	// ProcessEvent processes an event with the rules of the event filter. It
	// returns a DropError if the event should be dropped by the eventfilter.
	ProcessEvent(event types.Event) (types.Event, Report, error)
}
