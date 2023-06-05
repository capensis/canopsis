package contextgraph

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	added   = 0
	removed = 1
)

func NewManager(
	adapter libentity.Adapter,
	dbClient libmongo.DbClient,
	storage EntityServiceStorage,
	metricMetaUpdater metrics.MetaUpdater,
	logger zerolog.Logger,
) Manager {
	return &manager{
		adapter:           adapter,
		collection:        dbClient.Collection(libmongo.EntityMongoCollection),
		storage:           storage,
		metricMetaUpdater: metricMetaUpdater,
		logger:            logger,
	}
}

type manager struct {
	adapter           libentity.Adapter
	storage           EntityServiceStorage
	collection        libmongo.DbCollection
	metricMetaUpdater metrics.MetaUpdater
	logger            zerolog.Logger
}

func (m *manager) UpdateEntities(ctx context.Context, eventEntityID string, entities []types.Entity) (types.Entity, error) {
	writeModels := make([]mongo.WriteModel, 0, canopsis.DefaultBulkSize)

	var eventEntity types.Entity

	bulkBytesSize := 0
	var newModel mongo.WriteModel
	for _, ent := range entities {
		set := bson.M{}

		if ent.ID == eventEntityID {
			eventEntity = ent
			if ent.LastEventDate != nil {
				set["last_event_date"] = ent.LastEventDate
			}
		}

		if ent.IsNew {
			newModel = mongo.NewInsertOneModel().SetDocument(ent)
		} else {
			if ent.Infos != nil {
				set["infos"] = ent.Infos
			}

			if ent.ComponentInfos != nil {
				set["component_infos"] = ent.ComponentInfos
			}

			if ent.Services != nil {
				set["services"] = ent.Services
			}

			if ent.ServicesToAdd != nil {
				set["services_to_add"] = ent.ServicesToAdd
			}

			if ent.ServicesToRemove != nil {
				set["services_to_remove"] = ent.ServicesToRemove
			}

			newModel = mongo.NewUpdateOneModel().
				SetFilter(bson.M{"_id": ent.ID}).
				SetUpdate(bson.M{"$set": set})
		}

		b, err := bson.Marshal(newModel)
		if err != nil {
			return types.Entity{}, err
		}

		newModelLen := len(b)
		if bulkBytesSize+newModelLen > canopsis.DefaultBulkBytesSize {
			_, err := m.collection.BulkWrite(ctx, writeModels)
			if err != nil {
				return types.Entity{}, err
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}

		bulkBytesSize += newModelLen
		writeModels = append(
			writeModels,
			newModel,
		)

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err := m.collection.BulkWrite(ctx, writeModels)
			if err != nil {
				return types.Entity{}, err
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}
	}

	if len(writeModels) != 0 {
		_, err := m.collection.BulkWrite(ctx, writeModels)
		if err != nil {
			return types.Entity{}, err
		}
	}

	return eventEntity, nil
}

func (m *manager) CheckServices(ctx context.Context, entities []types.Entity) ([]types.Entity, error) {
	if len(entities) == 0 {
		return nil, nil
	}

	services, err := m.storage.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	entitiesData := make(map[string][2][]string) // array's indexes: 0 - added to impact, 1 - removed from impact
	servicesData := make(map[string][2][]string) // array's indexes: 0 - added to depends, 1 - removed from depends

	for _, ent := range entities {
		entityID := ent.ID

		impactMap := make(map[string]struct{}, len(ent.Services))
		for _, impactedService := range ent.Services {
			impactMap[impactedService] = struct{}{}
		}

		for _, serv := range services {
			serviceID := serv.ID

			_, found := impactMap[serviceID]

			match := false
			if len(serv.EntityPattern) > 0 {
				var err error
				match, _, err = serv.EntityPattern.Match(ent)
				if err != nil {
					m.logger.Err(err).Str("service", serv.ID).Msgf("service has invalid pattern")
				}
			} else if serv.OldEntityPatterns.IsSet() {
				if serv.OldEntityPatterns.IsValid() {
					match = serv.OldEntityPatterns.Matches(&ent)
				} else {
					m.logger.Err(pattern.ErrInvalidOldEntityPattern).Str("service", serv.ID).Msgf("service has invalid pattern")
				}
			}

			if match {
				if !found && ent.Enabled {
					entData := entitiesData[entityID]
					entData[added] = append(entData[added], serviceID)
					entitiesData[entityID] = entData

					servData := servicesData[serviceID]
					servData[added] = append(servData[added], entityID)
					servicesData[serviceID] = servData
				}

				if found && !ent.Enabled {
					entData := entitiesData[entityID]
					entData[removed] = append(entData[removed], serviceID)
					entitiesData[entityID] = entData

					servData := servicesData[serviceID]
					servData[removed] = append(servData[removed], entityID)
					servicesData[serviceID] = servData
				}
			} else if found {
				entData := entitiesData[entityID]
				entData[removed] = append(entData[removed], serviceID)
				entitiesData[entityID] = entData

				servData := servicesData[serviceID]
				servData[removed] = append(servData[removed], entityID)
				servicesData[serviceID] = servData
			}
		}
	}

	updatedEntities := make([]types.Entity, 0, len(entities)+len(servicesData))
	for _, ent := range entities {
		data, ok := entitiesData[ent.ID]
		if !ok {
			updatedEntities = append(updatedEntities, ent)
			continue
		}

		addedTo := data[0]
		removedFrom := data[1]

		toAddMap := make(map[string]bool, len(addedTo))
		for _, impact := range addedTo {
			toAddMap[impact] = true
		}

		toRemoveMap := make(map[string]bool, len(removedFrom))
		for _, impact := range removedFrom {
			toRemoveMap[impact] = true
		}

		for idx := 0; idx < len(ent.ServicesToAdd) && len(removedFrom) != 0; idx++ {
			if toRemoveMap[ent.ServicesToAdd[idx]] {
				ent.ServicesToAdd = append(ent.ServicesToAdd[:idx], ent.ServicesToAdd[idx+1:]...)
				idx--
			}
		}

		for idx := 0; idx < len(ent.ServicesToRemove) && len(addedTo) != 0; idx++ {
			if toAddMap[ent.ServicesToRemove[idx]] {
				ent.ServicesToRemove = append(ent.ServicesToRemove[:idx], ent.ServicesToRemove[idx+1:]...)
				idx--
			}
		}

		for idx := 0; idx < len(ent.Services); idx++ {
			if toRemoveMap[ent.Services[idx]] {
				ent.Services = append(ent.Services[:idx], ent.Services[idx+1:]...)
				idx--
			}
		}

		ent.ServicesToAdd = append(ent.ServicesToAdd, addedTo...)
		ent.ServicesToRemove = append(ent.ServicesToRemove, removedFrom...)
		ent.Services = append(ent.Services, addedTo...)

		updatedEntities = append(updatedEntities, ent)
	}

	for _, serv := range services {
		if _, ok := servicesData[serv.ID]; ok {
			ent := serv

			updatedEntities = append(updatedEntities, types.Entity{
				ID:      ent.ID,
				Enabled: true,
				Type:    types.EntityTypeService,
			})
		}
	}

	return updatedEntities, nil
}

func (m *manager) RecomputeService(ctx context.Context, serviceID string) (types.Entity, []types.Entity, error) {
	if serviceID == "" {
		return types.Entity{}, nil, nil
	}

	service, err := m.storage.Get(ctx, serviceID)
	if err != nil {
		return types.Entity{}, nil, err
	}

	if !service.Enabled || service.ID == "" {
		var dependedEntities []types.Entity
		cursor, err := m.collection.Find(ctx, bson.M{"services": serviceID})
		if err != nil {
			return types.Entity{}, nil, err
		}

		err = cursor.All(ctx, &dependedEntities)
		if err != nil {
			return types.Entity{}, nil, err
		}

		for idx, ent := range dependedEntities {
			for impIdx, impServ := range ent.Services {
				if impServ == serviceID {
					ent.Services = append(ent.Services[:impIdx], ent.Services[impIdx+1:]...)
					break
				}
			}

			for impIdx, impServ := range ent.ServicesToAdd {
				if impServ == serviceID {
					ent.ServicesToAdd = append(ent.ServicesToAdd[:impIdx], ent.ServicesToAdd[impIdx+1:]...)
					break
				}
			}

			dependedEntities[idx] = ent
		}

		var impactedEntities []types.Entity
		cursor, err = m.collection.Find(ctx, bson.M{"depends": serviceID})
		if err != nil {
			return types.Entity{}, nil, err
		}

		err = cursor.All(ctx, &impactedEntities)
		if err != nil {
			return types.Entity{}, nil, err
		}

		for idx, ent := range impactedEntities {
			service.Services = append(service.Services, ent.ID)
			impactedEntities[idx] = ent
		}

		if service.Entity.ID != "" {
			impactedEntities = append(impactedEntities, service.Entity)
		}

		return service.Entity, append(dependedEntities, impactedEntities...), nil
	}

	var updatedEntities []types.Entity

	query, negativeQuery, err := getServiceQueries(service)
	if err != nil {
		return types.Entity{}, nil, err
	}

	if query == nil || negativeQuery == nil {
		return types.Entity{}, nil, fmt.Errorf("can't get queries from patterns")
	}

	var entitiesToRemove []types.Entity

	cursor, err := m.collection.Find(
		ctx,
		bson.M{"$and": bson.A{
			bson.M{"services": serviceID},
			negativeQuery,
		}},
	)
	if err != nil {
		return types.Entity{}, nil, err
	}

	err = cursor.All(ctx, &entitiesToRemove)
	if err != nil {
		return types.Entity{}, nil, err
	}

	entitiesToRemoveMap := make(map[string]bool, len(entitiesToRemove))
	for _, ent := range entitiesToRemove {
		entitiesToRemoveMap[ent.ID] = true

		for idx, impServ := range ent.Services {
			if impServ == serviceID {
				ent.Services = append(ent.Services[:idx], ent.Services[idx+1:]...)
				break
			}
		}

		for idx, impServ := range ent.ServicesToAdd {
			if impServ == serviceID {
				ent.ServicesToAdd = append(ent.ServicesToAdd[:idx], ent.ServicesToAdd[idx+1:]...)
				break
			}
		}

		updatedEntities = append(updatedEntities, ent)
	}

	var entitiesToAdd []types.Entity
	cursor, err = m.collection.Find(
		ctx,
		bson.M{"$and": bson.A{
			bson.M{"services": bson.M{"$ne": serviceID}},
			query,
		}})
	if err != nil {
		return types.Entity{}, nil, err
	}

	err = cursor.All(ctx, &entitiesToAdd)
	if err != nil {
		return types.Entity{}, nil, err
	}

	for _, ent := range entitiesToAdd {
		for idx, impServ := range ent.ServicesToRemove {
			if impServ == serviceID {
				ent.ServicesToRemove = append(ent.ServicesToRemove[:idx], ent.ServicesToRemove[idx+1:]...)
				break
			}
		}

		ent.Services = append(ent.Services, serviceID)
		updatedEntities = append(updatedEntities, ent)
	}

	return service.Entity, append(updatedEntities, types.Entity{
		ID:      service.ID,
		Enabled: service.Enabled,
	}), nil
}

func (m *manager) HandleEvent(ctx context.Context, event types.Event) (types.Entity, []types.Entity, error) {
	eventEntity, err := m.adapter.GetEntityByID(ctx, event.GetEID())
	isNew := errors.Is(err, libentity.ErrNotFound)
	if err != nil && !isNew {
		return types.Entity{}, nil, err
	}

	if !event.IsContextable() || event.IsOnlyServiceUpdate() {
		if isNew {
			return types.Entity{}, nil, fmt.Errorf("entity %s doesn't exist", event.GetEID())
		}

		return eventEntity, nil, nil
	}

	if event.SourceType == types.SourceTypeService {
		return eventEntity, nil, nil
	}

	if event.SourceType != types.SourceTypeResource && event.SourceType != types.SourceTypeComponent {
		return types.Entity{}, nil, nil
	}

	now := types.CpsTime{Time: time.Now()}
	if event.EventType == types.EventTypeCheck {
		eventEntity.LastEventDate = &now
	}

	connectorName := event.ConnectorName
	connectorID := event.Connector + "/" + connectorName

	var contextGraphEntities []types.Entity

	if event.SourceType == types.SourceTypeComponent {
		if isNew {
			exist, err := m.entityExist(ctx, connectorID)
			if err != nil {
				return types.Entity{}, nil, err
			}

			if !exist {
				contextGraphEntities = []types.Entity{
					{
						ID:            connectorID,
						Name:          connectorName,
						EnableHistory: []types.CpsTime{now},
						Enabled:       true,
						Type:          types.EntityTypeConnector,
						Infos:         map[string]types.Info{},
						ImpactLevel:   types.EntityDefaultImpactLevel,
						Created:       now,
						LastEventDate: &now,
						IsNew:         true,
					},
				}
			} else {
				_, err := m.collection.UpdateOne(
					ctx,
					bson.M{"_id": connectorID},
					bson.M{"$set": bson.M{"last_event_date": now}},
				)
				if err != nil {
					return types.Entity{}, nil, err
				}
			}

			return types.Entity{
				ID:            event.Component,
				Name:          event.Component,
				Connector:     connectorID,
				EnableHistory: []types.CpsTime{now},
				Enabled:       true,
				Type:          types.EntityTypeComponent,
				Component:     event.Component,
				Infos:         map[string]types.Info{},
				ImpactLevel:   types.EntityDefaultImpactLevel,
				Created:       now,
				IsNew:         true,
				LastEventDate: &now,
			}, contextGraphEntities, nil
		}

		// if component isn't new, then check if connector exists, if not upsert it.
		// if connector exists update last_event_date
		if eventEntity.Connector != connectorID {
			exist, err := m.entityExist(ctx, connectorID)
			if err != nil {
				return types.Entity{}, nil, err
			}

			if !exist {
				contextGraphEntities = []types.Entity{
					{
						ID:            connectorID,
						Name:          connectorName,
						EnableHistory: []types.CpsTime{now},
						Enabled:       true,
						Type:          types.EntityTypeConnector,
						Infos:         map[string]types.Info{},
						ImpactLevel:   types.EntityDefaultImpactLevel,
						Created:       now,
						LastEventDate: &now,
					},
				}
			} else {
				_, err := m.collection.UpdateOne(
					ctx,
					bson.M{"_id": connectorID},
					bson.M{"$set": bson.M{"last_event_date": now}},
				)
				if err != nil {
					return types.Entity{}, nil, err
				}
			}
		} else {
			_, err := m.collection.UpdateOne(
				ctx,
				bson.M{"_id": connectorID},
				bson.M{"$set": bson.M{"last_event_date": now}},
			)
			if err != nil {
				return types.Entity{}, nil, err
			}
		}

		return eventEntity, contextGraphEntities, err
	}

	//if resource is new, then upsert connector and component
	if isNew {
		exist, err := m.entityExist(ctx, connectorID)
		if err != nil {
			return types.Entity{}, nil, err
		}

		if !exist {
			contextGraphEntities = append(contextGraphEntities, types.Entity{
				ID:            connectorID,
				Name:          connectorName,
				EnableHistory: []types.CpsTime{now},
				Enabled:       true,
				Type:          types.EntityTypeConnector,
				Infos:         map[string]types.Info{},
				ImpactLevel:   types.EntityDefaultImpactLevel,
				Created:       now,
				LastEventDate: &now,
				IsNew:         true,
			})
		} else {
			_, err := m.collection.UpdateOne(
				ctx,
				bson.M{"_id": connectorID},
				bson.M{"$set": bson.M{"last_event_date": now}},
			)
			if err != nil {
				return types.Entity{}, nil, err
			}
		}

		componentInfosDoc := struct {
			ComponentInfos map[string]types.Info `bson:"infos"`
		}{}

		cursor, err := m.collection.Aggregate(
			ctx,
			[]bson.M{
				{
					"$match": bson.M{
						"_id": event.Component,
					},
				},
				{
					"$project": bson.M{
						"infos": 1,
					},
				},
			},
		)
		if err != nil {
			return types.Entity{}, nil, err
		}

		exist = false
		if cursor.Next(ctx) {
			exist = true
			err = cursor.Decode(&componentInfosDoc)
			if err != nil {
				return types.Entity{}, nil, err
			}
		}

		if !exist {
			contextGraphEntities = append(contextGraphEntities, types.Entity{
				ID:            event.Component,
				Name:          event.Component,
				Connector:     connectorID,
				EnableHistory: []types.CpsTime{now},
				Enabled:       true,
				Type:          types.EntityTypeComponent,
				Component:     event.Component,
				Infos:         map[string]types.Info{},
				ImpactLevel:   types.EntityDefaultImpactLevel,
				Created:       now,
				IsNew:         true,
				LastEventDate: &now,
			})
		}

		return types.Entity{
			ID:             event.Resource + "/" + event.Component,
			Name:           event.Resource,
			EnableHistory:  []types.CpsTime{now},
			Enabled:        true,
			Type:           types.EntityTypeResource,
			Connector:      connectorID,
			Component:      event.Component,
			Infos:          map[string]types.Info{},
			ComponentInfos: componentInfosDoc.ComponentInfos,
			ImpactLevel:    types.EntityDefaultImpactLevel,
			IsNew:          true,
			Created:        now,
			LastEventDate:  &now,
		}, contextGraphEntities, nil
	}

	//if resource isn't new, then check if component or connector exists, if not upsert them.
	if eventEntity.Connector != connectorID {
		exist, err := m.entityExist(ctx, connectorID)
		if err != nil {
			return types.Entity{}, nil, err
		}

		if !exist {
			contextGraphEntities = append(contextGraphEntities, types.Entity{
				ID:            connectorID,
				Name:          connectorName,
				EnableHistory: []types.CpsTime{now},
				Enabled:       true,
				Type:          types.EntityTypeConnector,
				Infos:         map[string]types.Info{},
				ImpactLevel:   types.EntityDefaultImpactLevel,
				Created:       now,
				LastEventDate: &now,
				IsNew:         true,
			})
		} else {
			_, err := m.collection.UpdateOne(
				ctx,
				bson.M{"_id": connectorID},
				bson.M{"$set": bson.M{"last_event_date": now}},
			)
			if err != nil {
				return types.Entity{}, nil, err
			}
		}
	} else {
		_, err := m.collection.UpdateOne(
			ctx,
			bson.M{"_id": connectorID},
			bson.M{"$set": bson.M{"last_event_date": now}},
		)
		if err != nil {
			return types.Entity{}, nil, err
		}
	}

	return eventEntity, contextGraphEntities, nil
}

func (m *manager) UpdateImpactedServicesFromDependencies(ctx context.Context) error {
	cursor, err := m.collection.Aggregate(ctx, []bson.M{
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
			err := m.adapter.Bulk(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}

		bulkBytesSize += newModelLen
		writeModels = append(writeModels, newModel)

		if len(writeModels) == canopsis.DefaultBulkSize {
			err := m.adapter.Bulk(ctx, writeModels)
			if err != nil {
				return err
			}

			writeModels = writeModels[:0]
			bulkBytesSize = 0
		}
	}

	if len(writeModels) > 0 {
		err = m.adapter.Bulk(ctx, writeModels)
	}

	return err
}

func (m *manager) FillResourcesWithInfos(ctx context.Context, component types.Entity) ([]types.Entity, error) {
	if component.Type != types.EntityTypeComponent {
		return nil, nil
	}

	cursor, err := m.collection.Find(ctx, bson.M{"component": component.ID})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	resources := make([]types.Entity, 0)
	for cursor.Next(ctx) {
		var resource types.Entity

		err = cursor.Decode(&resource)
		if err != nil {
			return nil, err
		}

		resource.ComponentInfos = component.Infos
		resources = append(resources, resource)
	}

	return resources, nil
}

func (m *manager) UpdateLastEventDate(ctx context.Context, eventType string, entityID string, timestamp types.CpsTime) error {
	if eventType != types.EventTypeCheck {
		return nil
	}

	_, err := m.collection.UpdateOne(
		ctx,
		bson.M{"_id": entityID},
		bson.M{"$set": bson.M{"last_event_date": timestamp}},
		options.Update().SetUpsert(true),
	)

	return err
}

func (m *manager) entityExist(ctx context.Context, id string) (bool, error) {
	err := m.collection.FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"_id": 1})).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func getServiceQueries(service entityservice.EntityService) (interface{}, interface{}, error) {
	var query, negativeQuery interface{}
	var err error

	if len(service.EntityPattern) > 0 {
		query, err = service.EntityPattern.ToMongoQuery("")
		if err != nil {
			return nil, nil, err
		}

		negativeQuery, err = service.EntityPattern.ToNegativeMongoQuery("")
		if err != nil {
			return nil, nil, err
		}
	} else if service.OldEntityPatterns.IsSet() {
		if !service.OldEntityPatterns.IsValid() {
			return nil, nil, pattern.ErrInvalidOldEntityPattern
		}
		query = service.OldEntityPatterns.AsMongoDriverQuery()
		negativeQuery = service.OldEntityPatterns.AsNegativeMongoDriverQuery()
	}

	return query, negativeQuery, nil
}
