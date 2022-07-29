package eventfilter

//go:generate mockgen -destination=../../../mocks/lib/canopsis/eventfilter/eventfilter.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter Service,Adapter

import (
	"context"
	libcontext "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// Service is an interface for the service that manages the event filter.
type Service interface {
	// LoadRules loads the event filter rules from the database, and adds them
	// to the service.
	LoadRules(ctx context.Context) error

	// LoadDataSourceFactories loads the data source factories and adds them to
	// the service.
	LoadDataSourceFactories(enrichmentCenter libcontext.EnrichmentCenter, dataSourceDirectory string) error

	// ProcessEvent processes an event with the rules of the event filter. It
	// returns a DropError if the event should be dropped by the eventfilter.
	ProcessEvent(ctx context.Context, event types.Event) (types.Event, Report, error)
}

// Adapter is a type that provides access to the MongoDB collection containing
// the event filter's rules
type Adapter interface {
	List(ctx context.Context) ([]Rule, error)
}
