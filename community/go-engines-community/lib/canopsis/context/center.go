package context

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewEnrichmentCenter(
	adapter libentity.Adapter,
	dbClient libmongo.DbClient,
	entityServiceManager entityservice.Manager,
	metricMetaUpdater metrics.MetaUpdater,
) EnrichmentCenter {
	return &center{
		dbClient:             dbClient,
		dbCollection:         dbClient.Collection(libmongo.EntityMongoCollection),
		adapter:              adapter,
		entityServiceManager: entityServiceManager,
		metricMetaUpdater:    metricMetaUpdater,
	}
}

type center struct {
	dbClient             libmongo.DbClient
	dbCollection         libmongo.DbCollection
	adapter              libentity.Adapter
	entityServiceManager entityservice.Manager
	metricMetaUpdater    metrics.MetaUpdater
}

func (c *center) Handle(ctx context.Context, event types.Event) (*types.Entity, UpdatedEntityServices, error) {
	updatedServices := UpdatedEntityServices{}
	var eventEntity *types.Entity
	var entities []types.Entity
	var err error
	if event.IsOnlyServiceUpdate() {
		eventEntity, err = c.findEntityByID(ctx, event.GetEID())
		if err != nil {
			return nil, updatedServices, err
		}
	} else {
		eventEntity, entities, err = c.createEntities(ctx, event)
		if err != nil {
			return nil, updatedServices, err
		}
	}

	if eventEntity == nil {
		return nil, updatedServices, nil
	}

	resources, err := c.updateComponentInfos(ctx, event, eventEntity)
	if err != nil {
		return nil, updatedServices, err
	}

	if len(resources) > 0 {
		has, err := c.entityServiceManager.HasEntityServiceByComponentInfos(ctx)
		if err != nil {
			return nil, updatedServices, err
		}
		if has {
			updatedServices.UpdatedComponentInfosResources = resources
		}
	}

	updatedEntityIds := make([]string, len(entities))
	for i, entity := range entities {
		updatedEntityIds[i] = entity.ID
	}
	updatedEntityIds = append(updatedEntityIds, resources...)
	if len(updatedEntityIds) > 0 {
		c.metricMetaUpdater.UpdateById(ctx, updatedEntityIds...)
	}

	found := false
	for k, v := range entities {
		if eventEntity.ID == v.ID {
			found = true

			//need to update entity inside entities because of component infos enrichment
			entities[k] = *eventEntity
			break
		}
	}

	if !found {
		entities = append(entities, *eventEntity)
	}

	added, removed, err := c.entityServiceManager.UpdateServices(ctx, entities)
	if err != nil {
		return nil, updatedServices, err
	}

	updatedServices.AddedTo = added[event.GetEID()]
	updatedServices.RemovedFrom = removed[event.GetEID()]

	return eventEntity, updatedServices, nil
}

func (c *center) HandleEntityServiceUpdate(ctx context.Context, serviceID string) (*UpdatedEntityServices, error) {
	isUpdated, removed, err := c.entityServiceManager.UpdateService(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	var updatedServices *UpdatedEntityServices
	if isUpdated {
		updatedServices = &UpdatedEntityServices{RemovedFrom: removed}
	}

	entity, err := c.findEntityByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return updatedServices, nil
	}

	addedTo, removedFrom, err := c.entityServiceManager.UpdateServices(ctx, []types.Entity{*entity})
	if err != nil {
		return nil, err
	}

	if added, ok := addedTo[serviceID]; ok {
		if updatedServices == nil {
			updatedServices = &UpdatedEntityServices{AddedTo: added}
		} else {
			updatedServices.AddedTo = append(updatedServices.AddedTo, added...)
		}
	}

	if removed, ok := removedFrom[serviceID]; ok {
		if updatedServices == nil {
			updatedServices = &UpdatedEntityServices{RemovedFrom: removed}
		} else {
			updatedServices.RemovedFrom = append(updatedServices.RemovedFrom, removed...)
		}
	}

	return updatedServices, nil
}

func (c *center) Get(ctx context.Context, event types.Event) (*types.Entity, error) {
	return c.findEntityByID(ctx, event.GetEID())
}

func (c *center) UpdateEntityInfos(ctx context.Context, entity *types.Entity) (UpdatedEntityServices, error) {
	if !entity.Enabled {
		return UpdatedEntityServices{}, errors.New("cannot update infos for disabled entity")
	}
	updatedServices := UpdatedEntityServices{}
	ok, err := c.adapter.AddInfos(ctx, entity.ID, entity.Infos)
	if err != nil {
		return updatedServices, err
	}

	if !ok {
		return updatedServices, err
	}

	entities := []types.Entity{*entity}
	addedTo, removedFrom, err := c.entityServiceManager.UpdateServices(ctx, entities)
	if err != nil {
		return updatedServices, err
	}

	updatedServices.AddedTo = addedTo[entity.ID]
	updatedServices.RemovedFrom = removedFrom[entity.ID]

	updatedEntities := []string{entity.ID}
	// Update component infos for related resource entities
	if entity.Type == types.EntityTypeComponent {
		resources, err := c.adapter.UpdateComponentInfosByComponent(ctx, entity.ID)
		if err != nil {
			return updatedServices, err
		}

		updatedEntities = append(updatedEntities, resources...)

		if len(resources) > 0 {
			has, err := c.entityServiceManager.HasEntityServiceByComponentInfos(ctx)
			if err != nil {
				return updatedServices, err
			}
			if has {
				updatedServices.UpdatedComponentInfosResources = resources
			}
		}
	}

	c.metricMetaUpdater.UpdateById(ctx, updatedEntities...)

	return updatedServices, nil
}

func (c *center) UpdateImpactedServices(ctx context.Context) error {
	cursor, err := c.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"connector": bson.M{"$nin": bson.A{nil, ""}}}},
		{"$unwind": "$services"},
		{"$group": bson.M{
			"_id":               "$connector",
			"impacted_services": bson.M{"$addToSet": "$services"},
		}},
	})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	var newModel mongo.WriteModel
	writeModels := make([]mongo.WriteModel, 0, canopsis.DefaultBulkSize)
	bulkBytesSize := 0

	for cursor.Next(ctx) {
		var info struct {
			ID               string   `bson:"_id"`
			ImpactedServices []string `bson:"impacted_services"`
		}

		err = cursor.Decode(&info)
		if err != nil {
			return err
		}

		if len(info.ImpactedServices) > 0 {
			newModel = mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": info.ID}).SetUpdate(bson.M{
				"$set": bson.M{"impacted_services": info.ImpactedServices},
			})
		} else {
			newModel = mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": info.ID}).SetUpdate(bson.M{
				"$unset": bson.M{"impacted_services": ""},
			})
		}

		b, err := bson.Marshal(newModel)
		if err != nil {
			return err
		}

		newModelLen := len(b)
		if bulkBytesSize+newModelLen > canopsis.DefaultBulkBytesSize {
			err := c.adapter.Bulk(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}

		bulkBytesSize += newModelLen
		writeModels = append(writeModels, newModel)

		if len(writeModels) == canopsis.DefaultBulkSize {
			err := c.adapter.Bulk(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}
	}

	if len(writeModels) > 0 {
		err = c.adapter.Bulk(ctx, writeModels)
	}

	return err
}

func (c *center) ReloadService(ctx context.Context, serviceID string) error {
	return c.entityServiceManager.ReloadService(ctx, serviceID)
}

func (c *center) LoadServices(ctx context.Context) error {
	return c.entityServiceManager.LoadServices(ctx)
}

func (c *center) findEntityByID(ctx context.Context, id string) (*types.Entity, error) {
	entity, err := c.adapter.GetEntityByID(ctx, id)
	if err != nil {
		if _, ok := err.(errt.NotFound); ok {
			return nil, nil
		}
		return nil, err
	}

	return &entity, nil
}

// updateComponentInfos updates component infos of resource entity if it's new resource event
// and component infos of all connected resource entities if it's component event.
func (c *center) updateComponentInfos(ctx context.Context, event types.Event, entity *types.Entity) ([]string, error) {
	// Update component infos for related resource entities
	if event.SourceType == types.SourceTypeComponent {
		resourceIDs, err := c.adapter.UpdateComponentInfosByComponent(ctx, event.Component)
		if err != nil {
			return nil, err
		}

		return resourceIDs, nil
	} else if event.SourceType == types.SourceTypeResource && entity != nil && entity.IsNew {
		infos, err := c.adapter.UpdateComponentInfos(ctx, entity.ID, event.Component)
		if err != nil {
			return nil, err
		}

		if infos != nil {
			entity.ComponentInfos = infos
		}
	}

	return nil, nil
}

// createEntities creates connection, component, resource entities if they don't exist.
func (c *center) createEntities(ctx context.Context, event types.Event) (*types.Entity, []types.Entity, error) {
	if event.SourceType != types.SourceTypeResource && event.SourceType != types.SourceTypeComponent {
		return nil, nil, nil
	}

	entity, err := c.Get(ctx, event)
	if err != nil {
		return nil, nil, err
	}

	connectorID := event.Connector + "/" + event.ConnectorName
	componentID := event.Component
	resourceID := ""
	if event.SourceType == types.SourceTypeResource {
		resourceID = event.Resource + "/" + event.Component
	}

	now := types.CpsTime{Time: time.Now()}
	connector := types.Entity{
		ID:            connectorID,
		Name:          event.ConnectorName,
		EnableHistory: []types.CpsTime{now},
		Enabled:       true,
		Type:          types.EntityTypeConnector,
		Services:      []string{},
		Infos:         map[string]types.Info{},
		ImpactLevel:   types.EntityDefaultImpactLevel,
		Created:       now,
		LastEventDate: &now,
	}
	component := types.Entity{
		ID:            componentID,
		Name:          event.Component,
		EnableHistory: []types.CpsTime{now},
		Enabled:       true,
		Type:          types.EntityTypeComponent,
		Services:      []string{},
		Connector:     connectorID,
		Component:     event.Component,
		Infos:         map[string]types.Info{},
		ImpactLevel:   types.EntityDefaultImpactLevel,
		Created:       now,
	}
	if resourceID == "" {
		component.LastEventDate = &now
	}

	var entities []types.Entity

	if entity != nil {
		if event.SourceType == types.SourceTypeResource && (entity.SoftDeleted != nil || entity.Connector != connectorID) {
			// save old connector for the component
			component.Connector = entity.Connector
			entity.Connector = connectorID
			entities = []types.Entity{connector, component, *entity}
		}

		if event.SourceType == types.SourceTypeComponent && (entity.SoftDeleted != nil || entity.Connector != connectorID) {
			entity.Connector = connectorID
			entities = []types.Entity{connector, *entity}
		}
	} else if resourceID == "" {
		entity = &component
		entities = []types.Entity{connector, *entity}
	} else {
		resource := types.Entity{
			ID:            resourceID,
			Name:          event.Resource,
			EnableHistory: []types.CpsTime{now},
			Enabled:       true,
			Type:          types.EntityTypeResource,
			Services:      []string{},
			Connector:     connectorID,
			Component:     componentID,
			Infos:         map[string]types.Info{},
			ImpactLevel:   types.EntityDefaultImpactLevel,
			IsNew:         true,
			Created:       now,
			LastEventDate: &now,
		}

		entity = &resource
		entities = []types.Entity{connector, component, *entity}
	}

	insertedIDs, err := c.adapter.UpsertMany(ctx, entities)
	if err != nil {
		return nil, nil, err
	}

	inserted := make([]types.Entity, 0)
	connectorUpdated := false
	entityUpdated := false
	for _, e := range entities {
		if e.ID == connectorID {
			connectorUpdated = true
		}
		if e.ID == entity.ID {
			entityUpdated = true
		}

		if insertedIDs[e.ID] {
			inserted = append(inserted, e)
		}
	}

	updateLastEventDate := make([]string, 0)
	if !entityUpdated {
		updateLastEventDate = append(updateLastEventDate, entity.ID)
		entity.LastEventDate = &now
	}
	if !connectorUpdated {
		updateLastEventDate = append(updateLastEventDate, connectorID)
	}

	err = c.adapter.UpdateLastEventDate(ctx, updateLastEventDate, now)
	if err != nil {
		return nil, nil, err
	}

	return entity, inserted, nil
}
