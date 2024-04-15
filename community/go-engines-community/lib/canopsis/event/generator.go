package event

import (
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Generator interface {
	Generate(entity types.Entity) (types.Event, error)
}

func NewGenerator(connector, connectorName string) Generator {
	return &generator{
		connector:     connector,
		connectorName: connectorName,
	}
}

type generator struct {
	connector, connectorName string
}

func (g *generator) Generate(
	entity types.Entity,
) (types.Event, error) {
	event := types.Event{
		Author:    canopsis.DefaultEventAuthor,
		Initiator: types.InitiatorSystem,
	}

	switch entity.Type {
	case types.EntityTypeConnector:
		event.Connector = strings.TrimSuffix(entity.ID, "/"+entity.Name)
		event.ConnectorName = entity.Name
		event.SourceType = types.SourceTypeConnector
	case types.EntityTypeComponent:
		event.Connector = g.connector
		event.ConnectorName = g.connectorName
		event.Component = entity.Name
		event.SourceType = types.SourceTypeComponent
	case types.EntityTypeResource:
		event.Connector = g.connector
		event.ConnectorName = g.connectorName
		event.Component = entity.Component
		event.Resource = entity.Name
		event.SourceType = types.SourceTypeResource
	case types.EntityTypeService:
		event.Connector = g.connector
		event.ConnectorName = g.connectorName
		event.Component = entity.ID
		event.SourceType = types.SourceTypeService
	default:
		return event, fmt.Errorf("unknown entity type %v", entity.Type)
	}

	return event, nil
}
