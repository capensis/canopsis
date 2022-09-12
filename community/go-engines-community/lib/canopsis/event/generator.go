package event

import (
	"context"
	"fmt"
	"strings"

	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Generator interface {
	Generate(
		ctx context.Context,
		entity types.Entity,
	) (types.Event, error)
}

func NewGenerator(entityAdapter libentity.Adapter) Generator {
	return &generator{
		entityAdapter: entityAdapter,

		connectors: make(map[string]types.Entity),
		components: make(map[string]types.Entity),
	}
}

type generator struct {
	entityAdapter libentity.Adapter

	connectors map[string]types.Entity
	components map[string]types.Entity
}

func (g *generator) Generate(
	ctx context.Context,
	entity types.Entity,
) (types.Event, error) {
	event := types.Event{
		Initiator: types.InitiatorSystem,
	}

	switch entity.Type {
	case types.EntityTypeConnector:
		event.Connector = strings.TrimSuffix(entity.ID, "/"+entity.Name)
		event.ConnectorName = entity.Name
	case types.EntityTypeComponent:
		if entity.Connector != "" {
			event.Connector, event.ConnectorName, _ = strings.Cut(entity.Connector, "/")
		} else {
			connector, err := g.findConnectorForComponent(ctx, entity)
			if err != nil {
				return event, err
			}
			if connector == nil {
				return event, fmt.Errorf("cannot generate event for entity %v : not found any alarm and not found linked connector", entity.ID)
			}
			event.Connector = strings.TrimSuffix(connector.ID, "/"+connector.Name)
			event.ConnectorName = connector.Name
		}
		event.Component = entity.Name
	case types.EntityTypeResource:
		if entity.Connector != "" {
			event.Connector, event.ConnectorName, _ = strings.Cut(entity.Connector, "/")
		} else {
			connector, err := g.findConnectorForResource(ctx, entity)
			if err != nil {
				return event, err
			}
			if connector == nil {
				return event, fmt.Errorf("cannot generate event for entity %v : not found any alarm and not found linked connector", entity.ID)
			}
			event.Connector = strings.TrimSuffix(connector.ID, "/"+connector.Name)
			event.ConnectorName = connector.Name
		}
		if entity.Component != "" {
			event.Component = entity.Component
		} else {
			component, err := g.findComponentForResource(ctx, entity)
			if err != nil {
				return event, err
			}
			if component == nil {
				return event, fmt.Errorf("cannot generate event for resource %v : not found any alarm and not found linked component", entity.ID)
			}
			event.Component = component.ID
		}
		event.Resource = entity.Name
	case types.EntityTypeService:
		event.Connector = types.ConnectorEngineService
		event.ConnectorName = types.ConnectorEngineService
		event.Component = entity.ID
	default:
		return event, fmt.Errorf("unknown entity type %v", entity.Type)
	}

	event.SourceType = event.DetectSourceType()

	return event, nil
}

func (g *generator) findConnectorForComponent(ctx context.Context, entity types.Entity) (*types.Entity, error) {
	for _, id := range entity.Impacts {
		if connector, ok := g.connectors[id]; ok {
			return &connector, nil
		}
	}

	connector, err := g.entityAdapter.FindConnectorForComponent(ctx, entity.ID)
	if err != nil || connector == nil {
		return nil, err
	}

	g.connectors[connector.ID] = *connector

	return connector, nil
}

func (g *generator) findConnectorForResource(ctx context.Context, entity types.Entity) (*types.Entity, error) {
	for _, id := range entity.Depends {
		if connector, ok := g.connectors[id]; ok {
			return &connector, nil
		}
	}

	connector, err := g.entityAdapter.FindConnectorForResource(ctx, entity.ID)
	if err != nil || connector == nil {
		return nil, err
	}

	g.connectors[connector.ID] = *connector

	return connector, nil
}

func (g *generator) findComponentForResource(ctx context.Context, entity types.Entity) (*types.Entity, error) {
	for _, id := range entity.Impacts {
		if component, ok := g.components[id]; ok {
			return &component, nil
		}
	}

	component, err := g.entityAdapter.FindComponentForResource(ctx, entity.ID)
	if err != nil || component == nil {
		return nil, err
	}

	g.components[component.ID] = *component

	return component, nil
}
