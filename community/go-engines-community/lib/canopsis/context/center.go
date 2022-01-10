package context

import (
	"context"
	"errors"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const bulkMaxSize = 10000

func NewEnrichmentCenter(
	adapter libentity.Adapter,
	enableEnrich bool,
	entityServiceManager entityservice.Manager,
	metricMetaUpdater metrics.MetaUpdater,
	logger zerolog.Logger,
) EnrichmentCenter {
	return &center{
		adapter:              adapter,
		enableEnrich:         enableEnrich,
		entityServiceManager: entityServiceManager,
		metricMetaUpdater:    metricMetaUpdater,
		logger:               logger,
	}
}

type center struct {
	adapter              libentity.Adapter
	enableEnrich         bool
	entityServiceManager entityservice.Manager
	metricMetaUpdater    metrics.MetaUpdater
	logger               zerolog.Logger
}

func (c *center) Handle(ctx context.Context, event types.Event, fields EnrichFields) (*types.Entity, UpdatedEntityServices, error) {
	updatedServices := UpdatedEntityServices{}
	eventEntity, entities, err := c.createEntities(ctx, event, fields)
	if err != nil {
		return nil, updatedServices, err
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

	updatedEntities := make([]string, 0)
	metaUpdated := false
	for _, entity := range entities {
		if eventEntity.ID == entity.ID {
			// Update new event entity synchronously to update metrics in following engines.
			c.metricMetaUpdater.UpdateById(context.Background(), eventEntity.ID)
			metaUpdated = true
		} else {
			updatedEntities = append(updatedEntities, entity.ID)
		}
	}
	if !metaUpdated {
		updatedEntities = append(updatedEntities, eventEntity.ID)
	}
	updatedEntities = append(updatedEntities, resources...)
	if len(updatedEntities) > 0 {
		go c.metricMetaUpdater.UpdateById(context.Background(), updatedEntities...)
	}

	if !eventEntity.Enabled {
		return eventEntity, updatedServices, nil
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

	if entity == nil || !entity.Enabled {
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

	go c.metricMetaUpdater.UpdateById(context.Background(), updatedEntities...)

	return updatedServices, nil
}

func (c *center) UpdateImpactedServices(ctx context.Context) error {
	cursor, err := c.adapter.GetImpactedServicesInfo(ctx)
	if err != nil {
		return err
	}

	writeModels := make([]mongo.WriteModel, 0, bulkMaxSize)

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var info struct {
			ID               string   `bson:"_id"`
			ImpactedServices []string `bson:"impacted_services"`
		}
		err := cursor.Decode(&info)
		if err != nil {
			return err
		}

		if len(info.ImpactedServices) > 0 {
			writeModels = append(writeModels, mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": info.ID}).SetUpdate(bson.M{
				"$set": bson.M{"impacted_services": info.ImpactedServices},
			}))
		} else {
			writeModels = append(writeModels, mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": info.ID}).SetUpdate(bson.M{
				"$unset": bson.M{"impacted_services": ""},
			}))
		}

		if len(writeModels) == entityservice.BulkMaxSize {
			err := c.adapter.Bulk(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
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
func (c *center) createEntities(ctx context.Context, event types.Event, fields EnrichFields) (*types.Entity, []types.Entity, error) {
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
		Impacts:       []string{},
		Depends:       []string{componentID},
		EnableHistory: []types.CpsTime{now},
		Enabled:       true,
		Type:          types.EntityTypeConnector,
		Infos:         map[string]types.Info{},
		ImpactLevel:   types.EntityDefaultImpactLevel,
		Created:       now,
		LastEventDate: &now,
	}
	component := types.Entity{
		ID:            componentID,
		Name:          event.Component,
		Impacts:       []string{connectorID},
		Depends:       []string{},
		EnableHistory: []types.CpsTime{now},
		Enabled:       true,
		Type:          types.EntityTypeComponent,
		Component:     event.Component,
		Infos:         map[string]types.Info{},
		ImpactLevel:   types.EntityDefaultImpactLevel,
		Created:       now,
	}
	if resourceID == "" {
		component.LastEventDate = &now
	} else {
		connector.Impacts = append(connector.Impacts, resourceID)
		component.Depends = append(component.Depends, resourceID)
	}

	var entities []types.Entity

	if entity != nil {
		if event.SourceType == types.SourceTypeResource && !entity.HasDepend(connectorID) {
			entity.Depends = append(entity.Depends, connectorID)
			entities = []types.Entity{connector, component, *entity}
		}

		if event.SourceType == types.SourceTypeComponent && !entity.HasImpact(connectorID) {
			entity.Impacts = append(entity.Impacts, connectorID)
			entities = []types.Entity{connector, *entity}
		}
	} else {
		entity = &component
		if resourceID != "" {
			resource := types.Entity{
				ID:            resourceID,
				Name:          event.Resource,
				Impacts:       []string{componentID},
				Depends:       []string{connectorID},
				EnableHistory: []types.CpsTime{now},
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Component:     event.Component,
				Infos:         map[string]types.Info{},
				ImpactLevel:   types.EntityDefaultImpactLevel,
				IsNew:         true,
				Created:       now,
				LastEventDate: &now,
			}

			entity = &resource
		}

		entities = []types.Entity{connector, component, *entity}
	}

	if c.enableEnrich {
		set, err := c.setExtraInfos(event, entity, fields)
		if err != nil {
			return nil, nil, err
		}

		if set {
			found := false
			for _, v := range entities {
				if entity.ID == v.ID {
					found = true
					break
				}
			}

			if !found {
				entities = append(entities, *entity)
			}
		}
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

// setExtraInfos sets infos from event to entity if they are allowed by enrich fields.
func (c *center) setExtraInfos(event types.Event, entity *types.Entity, fields EnrichFields) (bool, error) {
	set := false
	for key, val := range event.ExtraInfos {
		if !fields.Allow(key) {
			continue
		}

		value, err := types.InterfaceToString(val)
		if err != nil {
			return false, err
		}

		info := types.NewInfo(key, "", value)
		entity.Infos[info.Name] = info
		set = true
	}

	return set, nil
}
