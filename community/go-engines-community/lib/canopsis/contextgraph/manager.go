package contextgraph

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	added   = 0
	removed = 1
)

type Report struct {
	// The check flags show if an entity should be included in a second transaction search
	// to check services and state settings.
	CheckResource  bool
	CheckComponent bool
	CheckConnector bool

	// IsNew is used only for event metric
	IsNew bool
}

func NewManager(
	adapter libentity.Adapter,
	dbClient libmongo.DbClient,
	storage EntityServiceStorage,
	stateSettingService statesetting.Assigner,
	logger zerolog.Logger,
) Manager {
	return &manager{
		adapter:             adapter,
		entityCollection:    dbClient.Collection(libmongo.EntityMongoCollection),
		storage:             storage,
		stateSettingService: stateSettingService,
		logger:              logger,
	}
}

type manager struct {
	adapter             libentity.Adapter
	storage             EntityServiceStorage
	entityCollection    libmongo.DbCollection
	services            []entityservice.EntityService
	stateSettingService statesetting.Assigner
	logger              zerolog.Logger
}

func (m *manager) InheritComponentFields(resource, component *types.Entity, commRegister libmongo.CommandsRegister) error {
	update := make(bson.M)

	if len(component.Infos) > 0 {
		update["component_infos"] = component.Infos
	}

	if component.StateInfo != nil {
		if component.StateInfo.InheritedPattern != nil {
			matched, err := match.MatchEntityPattern(*component.StateInfo.InheritedPattern, resource)
			if err != nil {
				return err
			}

			if matched && !resource.ComponentStateSettings {
				resource.ComponentStateSettings = true
				if !resource.ComponentStateSettingsToRemove {
					resource.ComponentStateSettingsToAdd = true
				}

				resource.ComponentStateSettingsToRemove = false

				update["component_state_settings"] = resource.ComponentStateSettings
				update["component_state_settings_to_add"] = resource.ComponentStateSettingsToAdd
				update["component_state_settings_to_remove"] = resource.ComponentStateSettingsToRemove
			}

			if !matched && resource.ComponentStateSettings {
				resource.ComponentStateSettings = false
				if !resource.ComponentStateSettingsToAdd {
					resource.ComponentStateSettingsToRemove = true
				}

				resource.ComponentStateSettingsToAdd = false

				update["component_state_settings"] = resource.ComponentStateSettings
				update["component_state_settings_to_add"] = resource.ComponentStateSettingsToAdd
				update["component_state_settings_to_remove"] = resource.ComponentStateSettingsToRemove
			}
		}
	}

	if len(update) > 0 {
		commRegister.RegisterUpdate(resource.ID, update)
	}

	return nil
}

func (m *manager) RefreshServices(ctx context.Context) error {
	var err error

	m.services, err = m.storage.GetAll(ctx)

	return err
}

func (m *manager) AssignServices(entity *types.Entity, commRegister libmongo.CommandsRegister) {
	// array's indexes: 0 - added services, 1 - removed services
	var entitiesData [2][]string

	servicesMap := make(map[string]struct{}, len(entity.Services))
	for _, id := range entity.Services {
		servicesMap[id] = struct{}{}
	}

	for _, serv := range m.services {
		serviceID := serv.ID
		_, found := servicesMap[serviceID]
		matched := false
		if len(serv.EntityPattern) > 0 {
			var err error
			matched, err = match.MatchEntityPattern(serv.EntityPattern, entity)
			if err != nil {
				m.logger.Err(err).Str("service", serviceID).Msgf("service has invalid pattern")
			}
		}

		if matched {
			if !found && entity.Enabled {
				entitiesData[added] = append(entitiesData[added], serviceID)
			}

			if found && !entity.Enabled {
				entitiesData[removed] = append(entitiesData[removed], serviceID)
			}
		} else if found {
			entitiesData[removed] = append(entitiesData[removed], serviceID)
		}
	}

	addedTo := entitiesData[added]
	removedFrom := entitiesData[removed]

	if len(addedTo) == 0 && len(removedFrom) == 0 {
		return
	}

	toAddMap := make(map[string]bool, len(addedTo))
	for _, impact := range addedTo {
		toAddMap[impact] = true
	}

	toRemoveMap := make(map[string]bool, len(removedFrom))
	for _, impact := range removedFrom {
		toRemoveMap[impact] = true
	}

	servicesToAddMap := make(map[string]bool, len(entity.ServicesToAdd))
	for _, id := range entity.ServicesToAdd {
		servicesToAddMap[id] = true
	}

	servicesToRemoveMap := make(map[string]bool, len(entity.ServicesToRemove))
	for _, id := range entity.ServicesToRemove {
		servicesToRemoveMap[id] = true
	}

	newServices := make([]string, 0, len(entity.Services))
	newServicesToAdd := make([]string, 0, len(entity.ServicesToAdd))
	newServicesToRemove := make([]string, 0, len(entity.ServicesToRemove))

	newServices = append(newServices, addedTo...)

	for idx := 0; idx < len(addedTo); idx++ {
		if !servicesToRemoveMap[addedTo[idx]] {
			newServicesToAdd = append(newServicesToAdd, addedTo[idx])
		}
	}

	for idx := 0; idx < len(removedFrom); idx++ {
		if !servicesToAddMap[removedFrom[idx]] {
			newServicesToRemove = append(newServicesToRemove, removedFrom[idx])
		}
	}

	for idx := 0; idx < len(entity.ServicesToAdd); idx++ {
		if !toRemoveMap[entity.ServicesToAdd[idx]] {
			newServicesToAdd = append(newServicesToAdd, entity.ServicesToAdd[idx])
		}
	}

	for idx := 0; idx < len(entity.ServicesToRemove); idx++ {
		if !toAddMap[entity.ServicesToRemove[idx]] {
			newServicesToRemove = append(newServicesToRemove, entity.ServicesToRemove[idx])
		}
	}

	for idx := 0; idx < len(entity.Services); idx++ {
		if !toRemoveMap[entity.Services[idx]] {
			newServices = append(newServices, entity.Services[idx])
		}
	}

	commRegister.RegisterUpdate(
		entity.ID,
		bson.M{
			"services_to_add":    newServicesToAdd,
			"services_to_remove": newServicesToRemove,
			"services":           newServices,
		},
	)

	entity.ServicesToAdd = newServicesToAdd
	entity.ServicesToRemove = newServicesToRemove
	entity.Services = newServices
}

func (m *manager) RecomputeService(ctx context.Context, serviceID string, commRegister libmongo.CommandsRegister) (types.Entity, error) {
	if serviceID == "" {
		return types.Entity{}, nil
	}

	service, err := m.storage.Get(ctx, serviceID)
	if err != nil {
		return types.Entity{}, err
	}

	if !service.Enabled || service.ID == "" {
		err := m.processDisabledService(ctx, serviceID, commRegister)
		if err != nil {
			return types.Entity{}, err
		}

		// todo: should be called to get fresh services from the db, should be removed when we do something with cache
		err = m.RefreshServices(ctx)
		if err != nil {
			return types.Entity{}, err
		}

		m.AssignServices(&service.Entity, commRegister)

		return service.Entity, nil
	}

	query, negativeQuery, err := service.GetMongoQueries()
	if err != nil {
		return types.Entity{}, err
	}

	if query == nil || negativeQuery == nil {
		return types.Entity{}, fmt.Errorf("can't get queries from patterns")
	}

	var entitiesToRemove []types.Entity

	cursor, err := m.entityCollection.Find(
		ctx,
		bson.M{"$and": bson.A{
			bson.M{"services": serviceID},
			negativeQuery,
		}},
	)
	if err != nil {
		return types.Entity{}, err
	}

	err = cursor.All(ctx, &entitiesToRemove)
	if err != nil {
		return types.Entity{}, err
	}

	for _, ent := range entitiesToRemove {
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

		commRegister.RegisterUpdate(ent.ID, bson.M{"services": ent.Services, "services_to_add": ent.ServicesToAdd})
	}

	var entitiesToAdd []types.Entity
	cursor, err = m.entityCollection.Find(
		ctx,
		bson.M{"$and": bson.A{
			bson.M{"services": bson.M{"$ne": serviceID}},
			query,
		}})
	if err != nil {
		return types.Entity{}, err
	}

	err = cursor.All(ctx, &entitiesToAdd)
	if err != nil {
		return types.Entity{}, err
	}

	for _, ent := range entitiesToAdd {
		for idx, impServ := range ent.ServicesToRemove {
			if impServ == serviceID {
				ent.ServicesToRemove = append(ent.ServicesToRemove[:idx], ent.ServicesToRemove[idx+1:]...)
				break
			}
		}

		ent.Services = append(ent.Services, serviceID)
		commRegister.RegisterUpdate(ent.ID, bson.M{"services": ent.Services, "services_to_remove": ent.ServicesToRemove})
	}

	_, err = m.AssignStateSetting(ctx, &service.Entity, commRegister)
	if err != nil {
		return types.Entity{}, err
	}

	// todo: should be called to get fresh services from the db, should be removed when we do something with cache
	err = m.RefreshServices(ctx)
	if err != nil {
		return types.Entity{}, err
	}

	m.AssignServices(&service.Entity, commRegister)

	return service.Entity, nil
}

func (m *manager) HandleResource(ctx context.Context, event *types.Event, commRegister libmongo.CommandsRegister) (Report, error) {
	var report Report

	var resource *types.Entity
	var componentExist bool
	var connectorExist bool
	var err error

	componentID := event.Component
	resourceID := event.Resource + "/" + componentID
	connectorName := event.ConnectorName
	connectorID := event.Connector + "/" + connectorName

	if !event.IsContextable() || event.IsOnlyServiceUpdate() {
		resource, err = m.getEntity(ctx, resourceID)
		if err != nil {
			return report, err
		}

		if resource == nil {
			return report, fmt.Errorf("resource %s doesn't exist", resourceID)
		}

		if event.IsOnlyServiceUpdate() {
			report.CheckResource = true // to check services and state settings.
		}

		event.Entity = resource

		return report, nil
	} else {
		resource, componentExist, connectorExist, err = m.getResourceEntities(ctx, resourceID, componentID, connectorID)
		if err != nil {
			return report, err
		}
	}

	if resource != nil && resource.SoftDeleted != nil {
		event.Entity = resource

		// clean report
		return Report{}, nil
	}

	now := datetime.NewCpsTime()

	if resource == nil {
		if !connectorExist {
			commRegister.RegisterInsert(&types.Entity{
				ID:            connectorID,
				Name:          connectorName,
				EnableHistory: []datetime.CpsTime{now},
				Enabled:       true,
				Type:          types.EntityTypeConnector,
				Infos:         map[string]types.Info{},
				ImpactLevel:   types.EntityDefaultImpactLevel,
				Created:       now,
				LastEventDate: &now,
				Healthcheck:   event.Healthcheck,
			})

			report.CheckConnector = true
		} else {
			commRegister.RegisterUpdate(connectorID, bson.M{"last_event_date": now})
		}

		if !componentExist {
			commRegister.RegisterInsert(&types.Entity{
				ID:            componentID,
				Name:          componentID,
				Connector:     connectorID,
				EnableHistory: []datetime.CpsTime{now},
				Enabled:       true,
				Type:          types.EntityTypeComponent,
				Component:     componentID,
				Infos:         map[string]types.Info{},
				ImpactLevel:   types.EntityDefaultImpactLevel,
				Created:       now,
				LastEventDate: &now,
				Healthcheck:   event.Healthcheck,
			})

			report.CheckComponent = true
		}

		resource = &types.Entity{
			ID:            resourceID,
			Name:          event.Resource,
			EnableHistory: []datetime.CpsTime{now},
			Enabled:       true,
			Type:          types.EntityTypeResource,
			Connector:     connectorID,
			Component:     event.Component,
			Infos:         map[string]types.Info{},
			ImpactLevel:   types.EntityDefaultImpactLevel,
			Created:       now,
			LastEventDate: &now,
			Healthcheck:   event.Healthcheck,
		}

		commRegister.RegisterInsert(resource)
		report.CheckResource = true
		report.IsNew = true

		event.Entity = resource

		return report, nil
	}

	if resource.Connector != connectorID && !connectorExist {
		resource.Connector = connectorID

		commRegister.RegisterUpdate(resourceID, bson.M{"connector": connectorID, "last_event_date": now})
		commRegister.RegisterInsert(&types.Entity{
			ID:            connectorID,
			Name:          connectorName,
			EnableHistory: []datetime.CpsTime{now},
			Enabled:       true,
			Type:          types.EntityTypeConnector,
			Infos:         map[string]types.Info{},
			ImpactLevel:   types.EntityDefaultImpactLevel,
			Created:       now,
			LastEventDate: &now,
			Healthcheck:   event.Healthcheck,
		})

		report.CheckResource = true
		report.CheckConnector = true
	} else {
		commRegister.RegisterUpdate(connectorID, bson.M{"last_event_date": now})
		commRegister.RegisterUpdate(resourceID, bson.M{"last_event_date": now})
	}

	resource.LastEventDate = &now
	event.Entity = resource

	return report, nil
}

func (m *manager) HandleComponent(ctx context.Context, event *types.Event, commRegister libmongo.CommandsRegister) (Report, error) {
	var report Report

	var component *types.Entity
	var connectorExist bool
	var err error

	componentID := event.Component
	connectorName := event.ConnectorName
	connectorID := event.Connector + "/" + connectorName

	if !event.IsContextable() || event.IsOnlyServiceUpdate() {
		component, err = m.getEntity(ctx, componentID)
		if err != nil {
			return report, err
		}

		if component == nil {
			return report, fmt.Errorf("component %s doesn't exist", componentID)
		}

		if event.IsOnlyServiceUpdate() {
			report.CheckComponent = true // to process state setting and component_infos for resources
		}

		event.Entity = component

		return report, nil
	} else {
		component, connectorExist, err = m.getComponentEntities(ctx, componentID, connectorID)
		if err != nil {
			return report, err
		}
	}

	if component != nil && component.SoftDeleted != nil {
		event.Entity = component

		// clean report
		return Report{}, nil
	}

	now := datetime.NewCpsTime()

	if component == nil {
		if !connectorExist {
			commRegister.RegisterInsert(&types.Entity{
				ID:            connectorID,
				Name:          connectorName,
				EnableHistory: []datetime.CpsTime{now},
				Enabled:       true,
				Type:          types.EntityTypeConnector,
				Infos:         map[string]types.Info{},
				ImpactLevel:   types.EntityDefaultImpactLevel,
				Created:       now,
				LastEventDate: &now,
				Healthcheck:   event.Healthcheck,
			})

			report.CheckConnector = true
		} else {
			commRegister.RegisterUpdate(connectorID, bson.M{"last_event_date": now})
		}

		component = &types.Entity{
			ID:            componentID,
			Name:          componentID,
			EnableHistory: []datetime.CpsTime{now},
			Enabled:       true,
			Type:          types.EntityTypeComponent,
			Connector:     connectorID,
			Component:     componentID,
			Infos:         map[string]types.Info{},
			ImpactLevel:   types.EntityDefaultImpactLevel,
			Created:       now,
			LastEventDate: &now,
			Healthcheck:   event.Healthcheck,
		}

		commRegister.RegisterInsert(component)
		report.CheckComponent = true
		report.IsNew = true

		event.Entity = component

		return report, nil
	}

	if component.Connector != connectorID && !connectorExist {
		component.Connector = connectorID

		commRegister.RegisterUpdate(componentID, bson.M{"connector": connectorID, "last_event_date": now})
		commRegister.RegisterInsert(&types.Entity{
			ID:            connectorID,
			Name:          connectorName,
			EnableHistory: []datetime.CpsTime{now},
			Enabled:       true,
			Type:          types.EntityTypeConnector,
			Infos:         map[string]types.Info{},
			ImpactLevel:   types.EntityDefaultImpactLevel,
			Created:       now,
			LastEventDate: &now,
			Healthcheck:   event.Healthcheck,
		})

		report.CheckComponent = true
		report.CheckConnector = true
	} else {
		commRegister.RegisterUpdate(connectorID, bson.M{"last_event_date": now})
		commRegister.RegisterUpdate(componentID, bson.M{"last_event_date": now})
	}

	component.LastEventDate = &now
	event.Entity = component

	return report, nil
}

func (m *manager) HandleService(ctx context.Context, event *types.Event, commRegister libmongo.CommandsRegister) (Report, error) {
	var service *types.Entity
	var err error

	serviceID := event.Component

	service, err = m.getEntity(ctx, serviceID)
	if err != nil {
		return Report{}, err
	}

	if service == nil {
		return Report{}, fmt.Errorf("service %s doesn't exist", serviceID)
	} else if service.SoftDeleted != nil {
		event.Entity = service

		return Report{}, nil
	}

	now := datetime.NewCpsTime()
	commRegister.RegisterUpdate(serviceID, bson.M{"last_event_date": now})
	service.LastEventDate = &now
	event.Entity = service

	return Report{}, nil
}

func (m *manager) HandleConnector(ctx context.Context, event *types.Event, commRegister libmongo.CommandsRegister) (Report, error) {
	var connector *types.Entity
	var err error

	connectorName := event.ConnectorName
	connectorID := event.Connector + "/" + connectorName

	connector, err = m.getEntity(ctx, connectorID)
	if err != nil {
		return Report{}, err
	}

	if connector == nil {
		return Report{}, fmt.Errorf("connector %s doesn't exist", connectorID)
	}

	event.Entity = connector
	commRegister.RegisterUpdate(connectorID, bson.M{"last_event_date": datetime.NewCpsTime()})

	return Report{}, nil
}

func (m *manager) UpdateImpactedServicesFromDependencies(ctx context.Context) error {
	cursor, err := m.entityCollection.Aggregate(ctx, []bson.M{
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

func (m *manager) ProcessComponentDependencies(ctx context.Context, component *types.Entity, commRegister libmongo.CommandsRegister) ([]string, error) {
	if component.Type != types.EntityTypeComponent {
		return nil, nil
	}

	cursor, err := m.entityCollection.Find(
		ctx,
		bson.M{"_id": bson.M{"$ne": component.ID}, "component": component.ID},
	)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var ids []string
	for cursor.Next(ctx) {
		var resource types.Entity
		update := make(bson.M)

		err = cursor.Decode(&resource)
		if err != nil {
			return nil, err
		}

		ids = append(ids, resource.ID)
		update["component_infos"] = component.Infos

		matched := true

		if component.StateInfo != nil {
			if component.StateInfo.InheritedPattern != nil {
				matched, err = match.MatchEntityPattern(*component.StateInfo.InheritedPattern, &resource)
				if err != nil {
					return nil, err
				}
			}
		} else {
			matched = false
		}

		if matched {
			update["component_state_settings"] = true
			update["component_state_settings_to_add"] = false
			update["component_state_settings_to_remove"] = false
		}

		if !matched {
			update["component_state_settings"] = false
			update["component_state_settings_to_add"] = false
			update["component_state_settings_to_remove"] = false
		}

		commRegister.RegisterUpdate(resource.ID, update)
	}

	return ids, nil
}

func (m *manager) UpdateLastEventDate(ctx context.Context, event *types.Event, updateConnectorLastEventDate bool) error {
	if event.EventType != types.EventTypeCheck || event.Entity.LastEventDate == nil {
		return nil
	}

	var query bson.M
	if updateConnectorLastEventDate {
		query = bson.M{"_id": bson.M{"$in": bson.A{event.Entity.ID, event.Entity.Connector}}}
	} else {
		query = bson.M{"_id": event.Entity.ID}
	}

	_, err := m.entityCollection.UpdateMany(
		ctx,
		query,
		bson.M{"$set": bson.M{"last_event_date": event.Entity.LastEventDate}},
	)

	return err
}

func (m *manager) AssignStateSetting(ctx context.Context, entity *types.Entity, commRegister libmongo.CommandsRegister) (bool, error) {
	return m.stateSettingService.AssignStateSetting(ctx, entity, commRegister)
}

func (m *manager) getResourceEntities(ctx context.Context, resourceID, componentID, connectorID string) (*types.Entity, bool, bool, error) {
	var resource *types.Entity
	var componentExist bool
	var connectorExist bool

	cursor, err := m.entityCollection.Find(ctx, bson.M{"_id": bson.M{"$in": bson.A{resourceID, componentID, connectorID}}})
	if err != nil {
		return nil, componentExist, connectorExist, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var ent types.Entity

		err = cursor.Decode(&ent)
		if err != nil {
			return nil, componentExist, connectorExist, err
		}

		if ent.Type == types.EntityTypeResource {
			resource = &ent
		} else if ent.Type == types.EntityTypeComponent {
			componentExist = ent.ID != ""
		} else {
			connectorExist = ent.ID != ""
		}
	}

	return resource, componentExist, connectorExist, nil
}

func (m *manager) getComponentEntities(ctx context.Context, componentID, connectorID string) (*types.Entity, bool, error) {
	var component *types.Entity
	var connectorExist bool

	cursor, err := m.entityCollection.Find(ctx, bson.M{"_id": bson.M{"$in": bson.A{componentID, connectorID}}})
	if err != nil {
		return nil, connectorExist, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var ent types.Entity

		err = cursor.Decode(&ent)
		if err != nil {
			return nil, connectorExist, err
		}

		if ent.Type == types.EntityTypeComponent {
			component = &ent
		} else {
			connectorExist = ent.ID != ""
		}
	}

	return component, connectorExist, nil
}

func (m *manager) getEntity(ctx context.Context, id string) (*types.Entity, error) {
	var eventEntity types.Entity

	err := m.entityCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&eventEntity)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &eventEntity, nil
}

func (m *manager) processDisabledService(ctx context.Context, serviceID string, commRegister libmongo.CommandsRegister) error {
	var dependedEntities []types.Entity
	cursor, err := m.entityCollection.Find(ctx, bson.M{"services": serviceID})
	if err != nil {
		return err
	}

	err = cursor.All(ctx, &dependedEntities)
	if err != nil {
		return err
	}

	for _, ent := range dependedEntities {
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

		commRegister.RegisterUpdate(ent.ID, bson.M{"services": ent.Services, "services_to_add": ent.ServicesToAdd})
	}

	return nil
}
