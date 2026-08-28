package contextgraph

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/sync/singleflight"
)

type Report struct {
	// The check flags show if an entity should be included in a second transaction search
	// to check services and state settings.
	CheckResource  bool
	CheckComponent bool
	CheckConnector bool
	CheckService   bool

	CheckInfoChanged bool

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

		sfGroup: &singleflight.Group{},
	}
}

type manager struct {
	adapter             libentity.Adapter
	storage             EntityServiceStorage
	entityCollection    libmongo.DbCollection
	stateSettingService statesetting.Assigner
	logger              zerolog.Logger

	sfGroup                  *singleflight.Group
	servicesMx               sync.RWMutex
	services                 map[string]EntityService
	servicesByComponentInfos map[string]map[string]EntityService
	servicesByInfos          map[string]map[string]EntityService
}

func (m *manager) InheritComponentFields(resource, component *types.Entity, commRegister libmongo.CommandsRegister) error {
	update := make(bson.M)
	var err error
	if component.StateInfo != nil {
		matched := true

		if len(component.StateInfo.InheritedPattern) > 0 {
			matched, err = match.MatchEntityPattern(component.StateInfo.InheritedPattern, resource)
			if err != nil {
				return err
			}
		}

		if matched && !resource.ComponentStateSettings {
			resource.ComponentStateSettings = true
			if !resource.ComponentStateSettingsToRemove {
				resource.ComponentStateSettingsToAdd = true
			} else {
				resource.ComponentStateSettingsToRemove = false
			}

			update["component_state_settings"] = resource.ComponentStateSettings
			update["component_state_settings_to_add"] = resource.ComponentStateSettingsToAdd
			update["component_state_settings_to_remove"] = resource.ComponentStateSettingsToRemove
		} else if !matched && resource.ComponentStateSettings {
			resource.ComponentStateSettings = false
			if !resource.ComponentStateSettingsToAdd {
				resource.ComponentStateSettingsToRemove = true
			} else {
				resource.ComponentStateSettingsToAdd = false
			}

			update["component_state_settings"] = resource.ComponentStateSettings
			update["component_state_settings_to_add"] = resource.ComponentStateSettingsToAdd
			update["component_state_settings_to_remove"] = resource.ComponentStateSettingsToRemove
		}
	}

	if len(update) > 0 {
		commRegister.RegisterUpdate(resource.ID, update)
	}

	return nil
}

func (m *manager) LoadServices(ctx context.Context) error {
	_, err, _ := m.sfGroup.Do("load_services", func() (any, error) {
		services, err := m.storage.GetAll(ctx)
		if err != nil {
			return nil, err
		}

		servicesByComponentInfos := make(map[string]map[string]EntityService, len(m.servicesByComponentInfos))
		servicesByInfos := make(map[string]map[string]EntityService, len(m.servicesByInfos))
		servicesMap := make(map[string]EntityService, len(services))

		for _, service := range services {
			for _, key := range service.EntityPattern.GetComponentInfosNames() {
				if servicesByComponentInfos[key] == nil {
					servicesByComponentInfos[key] = make(map[string]EntityService)
				}

				servicesByComponentInfos[key][service.ID] = service
			}

			// here we add keys from inherited pattern only for component infos, because we need to send recompute event
			// on component info change to recompute also inherited counters.
			for _, key := range service.InheritedPattern.GetComponentInfosNames() {
				if servicesByComponentInfos[key] == nil {
					servicesByComponentInfos[key] = make(map[string]EntityService)
				}

				servicesByComponentInfos[key][service.ID] = service
			}

			// for regular infos there is no need to check inherited pattern, because it will be checked in engine-axe
			// on counters update.
			for _, key := range service.EntityPattern.GetInfosNames() {
				if servicesByInfos[key] == nil {
					servicesByInfos[key] = make(map[string]EntityService)
				}

				servicesByInfos[key][service.ID] = service
			}

			servicesMap[service.ID] = service
		}

		m.servicesMx.Lock()
		m.servicesByComponentInfos = servicesByComponentInfos
		m.servicesByInfos = servicesByInfos
		m.services = servicesMap
		m.servicesMx.Unlock()

		return nil, nil
	})

	return err
}

func (m *manager) AssignServices(entity *types.Entity, commRegister libmongo.CommandsRegister) {
	toAddMap := make(map[string]bool)
	toRemoveMap := make(map[string]bool)

	servicesMap := make(map[string]bool, len(entity.Services))
	for _, id := range entity.Services {
		servicesMap[id] = true
	}

	m.assignAllServices(entity, servicesMap, toAddMap, toRemoveMap)
	m.applyAssignedServices(entity, toAddMap, toRemoveMap, commRegister)
}

func (m *manager) AssignServicesByInfoNames(entity *types.Entity, infoNames []string, commRegister libmongo.CommandsRegister) {
	if len(infoNames) == 0 {
		return
	}

	toAddMap := make(map[string]bool)
	toRemoveMap := make(map[string]bool)

	m.assignServicesByInfos(entity, infoNames, toAddMap, toRemoveMap)
	m.applyAssignedServices(entity, toAddMap, toRemoveMap, commRegister)
}

func (m *manager) AssignServicesByComponentInfoNames(entity *types.Entity, componentInfoNames []string, commRegister libmongo.CommandsRegister) map[string]bool {
	if len(componentInfoNames) == 0 {
		return nil
	}

	toAddMap := make(map[string]bool)
	toRemoveMap := make(map[string]bool)
	inheritedMap := make(map[string]bool)

	m.assignServicesByComponentInfos(entity, componentInfoNames, toAddMap, toRemoveMap, inheritedMap)
	m.applyAssignedServices(entity, toAddMap, toRemoveMap, commRegister)

	// gather all maps to have map with affected services ids
	maps.Copy(toAddMap, toRemoveMap)
	maps.Copy(toAddMap, inheritedMap)

	return toAddMap
}

func (m *manager) assignAllServices(
	entity *types.Entity,
	servicesMap map[string]bool,
	toAddMap map[string]bool,
	toRemoveMap map[string]bool,
) {
	m.servicesMx.RLock()
	defer m.servicesMx.RUnlock()

	for _, service := range m.services {
		m.assignService(entity, service, servicesMap, toAddMap, toRemoveMap)
	}
}

func (m *manager) assignServicesByInfos(
	entity *types.Entity,
	infoUpdates []string,
	toAddMap, toRemoveMap map[string]bool,
) {
	m.servicesMx.RLock()
	defer m.servicesMx.RUnlock()

	servicesMap := make(map[string]bool, len(entity.Services))
	for _, id := range entity.Services {
		servicesMap[id] = true
	}

	seen := make(map[string]bool, len(infoUpdates))
	for _, infoName := range infoUpdates {
		for _, service := range m.servicesByInfos[infoName] {
			if seen[service.ID] {
				continue
			}

			seen[service.ID] = true
			m.assignService(entity, service, servicesMap, toAddMap, toRemoveMap)
		}
	}
}

func (m *manager) assignServicesByComponentInfos(
	entity *types.Entity,
	componentInfoUpdates []string,
	toAddMap, toRemoveMap, inheritedMap map[string]bool,
) {
	m.servicesMx.RLock()
	defer m.servicesMx.RUnlock()

	servicesMap := make(map[string]bool, len(entity.Services))
	for _, id := range entity.Services {
		servicesMap[id] = true
	}

	inheritedServicesMap := make(map[string]bool, len(entity.InheritedServices))
	for _, id := range entity.InheritedServices {
		inheritedServicesMap[id] = true
	}

	seen := make(map[string]bool, len(componentInfoUpdates))

	for _, infoName := range componentInfoUpdates {
		for _, service := range m.servicesByComponentInfos[infoName] {
			if seen[service.ID] {
				continue
			}

			seen[service.ID] = true

			if !m.assignService(entity, service, servicesMap, toAddMap, toRemoveMap) {
				matched, err := match.MatchEntityPattern(service.InheritedPattern, entity)
				if err != nil {
					m.logger.Err(err).Str("service", service.ID).Msgf("service has invalid inherited pattern")
					continue
				}

				if inheritedServicesMap[service.ID] != matched {
					inheritedMap[service.ID] = true
				}
			}
		}
	}
}

func (m *manager) assignService(
	entity *types.Entity,
	service EntityService,
	servicesMap map[string]bool,
	toAddMap map[string]bool,
	toRemoveMap map[string]bool,
) bool {
	serviceChanged := false
	serviceID := service.ID

	found := servicesMap[serviceID]
	matched := false
	if len(service.EntityPattern) > 0 {
		var err error
		matched, err = match.MatchEntityPattern(service.EntityPattern, entity)
		if err != nil {
			m.logger.Err(err).Str("service", serviceID).Msgf("service has invalid pattern")
		}
	}

	if matched {
		if !found && entity.Enabled {
			toAddMap[serviceID] = true
			serviceChanged = true
		}

		if found && !entity.Enabled {
			toRemoveMap[serviceID] = true
			serviceChanged = true
		}
	} else if found {
		toRemoveMap[serviceID] = true
		serviceChanged = true
	}

	return serviceChanged
}

func (m *manager) applyAssignedServices(
	entity *types.Entity,
	toAddMap map[string]bool,
	toRemoveMap map[string]bool,
	commRegister libmongo.CommandsRegister,
) {
	if len(toAddMap) == 0 && len(toRemoveMap) == 0 {
		return
	}

	servicesToAddMap := make(map[string]bool, len(entity.ServicesToAdd))
	for _, id := range entity.ServicesToAdd {
		servicesToAddMap[id] = true
	}

	servicesToRemoveMap := make(map[string]bool, len(entity.ServicesToRemove))
	for _, id := range entity.ServicesToRemove {
		servicesToRemoveMap[id] = true
	}

	newServices := make([]string, 0, len(toAddMap)+len(entity.Services)-len(toRemoveMap))
	newServicesToAdd := make([]string, 0, max(len(entity.ServicesToAdd), len(toAddMap)))
	newServicesToRemove := make([]string, 0, max(len(entity.ServicesToRemove), len(toRemoveMap)))

	for id := range toAddMap {
		newServices = append(newServices, id)
		if !servicesToRemoveMap[id] {
			newServicesToAdd = append(newServicesToAdd, id)
		}
	}

	for id := range toRemoveMap {
		if !servicesToAddMap[id] {
			newServicesToRemove = append(newServicesToRemove, id)
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

	if commRegister != nil {
		commRegister.RegisterUpdate(
			entity.ID,
			bson.M{
				"services_to_add":    newServicesToAdd,
				"services_to_remove": newServicesToRemove,
				"services":           newServices,
			},
		)
	}

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
			return types.Entity{}, fmt.Errorf("recompute service %s: failed to process disabled service: %w", serviceID, err)
		}

		// todo: should be called to get fresh services from the db, should be removed when we do something with cache
		err = m.LoadServices(ctx)
		if err != nil {
			return types.Entity{}, fmt.Errorf("recompute service %s: failed to load services: %w", serviceID, err)
		}

		m.AssignServices(&service.Entity, commRegister)

		return service.Entity, nil
	}

	query, negativeQuery, err := service.GetMongoQueries()
	if err != nil {
		return types.Entity{}, fmt.Errorf("recompute service %s: failed to get mongo queries: %w", serviceID, err)
	}

	if query == nil || negativeQuery == nil {
		return types.Entity{}, fmt.Errorf("recompute service %s: can't get queries from patterns", serviceID)
	}

	var entitiesToRemove []types.Entity

	cursor, err := m.entityCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{
					"$and": bson.A{
						negativeQuery,
						// here we need to check by service ID not only in services but in all auxiliary fields
						// to avoid extra counters increments/decrements in the engine axe counters calculator.
						bson.M{"$or": bson.A{
							bson.M{"services": serviceID},
							bson.M{"inherited_services": serviceID},
							bson.M{"services_to_add": serviceID},
							bson.M{"services_to_remove": serviceID},
						}},
					},
				},
			},
			{
				"$project": bson.M{
					"_id":                1,
					"services":           1,
					"inherited_services": 1,
					"services_to_add":    1,
					"services_to_remove": 1,
				},
			},
			{
				"$addFields": bson.M{
					"services": bson.M{
						"$setDifference": bson.A{
							"$services",
							bson.A{serviceID},
						},
					},
				},
			},
			{
				"$addFields": bson.M{
					"inherited_services": bson.M{
						"$setDifference": bson.A{
							"$inherited_services",
							bson.A{serviceID},
						},
					},
				},
			},
			{
				"$addFields": bson.M{
					"services_to_add": bson.M{
						"$setDifference": bson.A{
							"$services_to_add",
							bson.A{serviceID},
						},
					},
				},
			},
			{
				"$addFields": bson.M{
					"services_to_remove": bson.M{
						"$setDifference": bson.A{
							"$services_to_remove",
							bson.A{serviceID},
						},
					},
				},
			},
		},
	)
	if err != nil {
		return types.Entity{}, fmt.Errorf("recompute service %s: failed to query entities to remove: %w", serviceID, err)
	}

	err = cursor.All(ctx, &entitiesToRemove)
	if err != nil {
		return types.Entity{}, fmt.Errorf("recompute service %s: failed to decode entities to remove: %w", serviceID, err)
	}

	for _, ent := range entitiesToRemove {
		commRegister.RegisterUpdate(ent.ID, bson.M{
			"services":           ent.Services,
			"inherited_services": ent.InheritedServices,
			"services_to_add":    ent.ServicesToAdd,
			"services_to_remove": ent.ServicesToRemove,
		})
	}

	_, err = m.AssignStateSetting(ctx, &service.Entity, commRegister)
	if err != nil {
		return types.Entity{}, err
	}

	var matchQuery bson.M

	if service.StateInfo != nil && service.StateInfo.InheritedPattern != nil {
		// here we match all entities regardless of services or inherited_services because we need to check them all.
		matchQuery = bson.M{
			"$match": bson.M{
				"$or": bson.A{
					query,
					bson.M{"services_to_add": serviceID},
					bson.M{"services_to_remove": serviceID},
				},
			},
		}
	} else {
		matchQuery = bson.M{
			"$match": bson.M{
				"$and": bson.A{
					query,
					bson.M{"$or": bson.A{
						bson.M{"services": bson.M{"$ne": serviceID}},
						bson.M{"inherited_services": serviceID},
						bson.M{"services_to_add": serviceID},
						bson.M{"services_to_remove": serviceID},
					}},
				},
			},
		}
	}

	cursor, err = m.entityCollection.Aggregate(
		ctx,
		[]bson.M{
			matchQuery,
			{
				"$addFields": bson.M{
					"services": bson.M{
						"$setDifference": bson.A{
							"$services",
							bson.A{serviceID},
						},
					},
				},
			},
			{
				"$addFields": bson.M{
					"inherited_services": bson.M{
						"$setDifference": bson.A{
							"$inherited_services",
							bson.A{serviceID},
						},
					},
				},
			},
			{
				"$addFields": bson.M{
					"services_to_add": bson.M{
						"$setDifference": bson.A{
							"$services_to_add",
							bson.A{serviceID},
						},
					},
				},
			},
			{
				"$addFields": bson.M{
					"services_to_remove": bson.M{
						"$setDifference": bson.A{
							"$services_to_remove",
							bson.A{serviceID},
						},
					},
				},
			},
		},
	)
	if err != nil {
		return types.Entity{}, fmt.Errorf("recompute service %s: failed to query entities to add: %w", serviceID, err)
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var ent types.Entity
		err = cursor.Decode(&ent)
		if err != nil {
			return types.Entity{}, fmt.Errorf("recompute service %s: failed to decode entity to add: %w", serviceID, err)
		}

		matched, err := match.MatchEntityPattern(service.EntityPattern, &ent)
		if err != nil {
			return types.Entity{}, fmt.Errorf("recompute service %s: failed to match entity pattern: %w", serviceID, err)
		}

		if matched {
			ent.Services = append(ent.Services, serviceID)
			if service.StateInfo != nil {
				inheritedMatched, err := match.MatchEntityPattern(service.StateInfo.InheritedPattern, &ent)
				if err != nil {
					return types.Entity{}, fmt.Errorf("recompute service %s: failed to match inherited pattern: %w", serviceID, err)
				}

				if inheritedMatched {
					ent.InheritedServices = append(ent.InheritedServices, serviceID)
				}
			}
		}

		commRegister.RegisterUpdate(ent.ID, bson.M{
			"services":           ent.Services,
			"inherited_services": ent.InheritedServices,
			"services_to_add":    ent.ServicesToAdd,
			"services_to_remove": ent.ServicesToRemove,
		})
	}

	if err := cursor.Err(); err != nil {
		return types.Entity{}, fmt.Errorf("recompute service %s: failed to process entities to add cursor: %w", serviceID, err)
	}

	// todo: should be called to get fresh services from the db, should be removed when we do something with cache
	err = m.LoadServices(ctx)
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
	var componentInfos map[string]types.Info
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
	}

	resource, componentExist, componentInfos, connectorExist, err = m.getResourceEntities(ctx, resourceID, componentID, connectorID)
	if err != nil {
		return report, err
	}

	if resource != nil && resource.SoftDeleted != nil {
		event.Entity = resource

		// clean report
		return Report{}, nil
	}

	now := datetime.NewCpsTime()
	var lastEventDate *datetime.CpsTime
	if event.EventType == types.EventTypeCheck {
		lastEventDate = &now
	}

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
				LastEventDate: lastEventDate,
				Healthcheck:   event.Healthcheck,
			})

			report.CheckConnector = true
		} else if lastEventDate != nil {
			commRegister.RegisterUpdate(connectorID, bson.M{"last_event_date": *lastEventDate})
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
				LastEventDate: lastEventDate,
				Healthcheck:   event.Healthcheck,
			})

			report.CheckComponent = true
		}

		resource = &types.Entity{
			ID:             resourceID,
			Name:           event.Resource,
			EnableHistory:  []datetime.CpsTime{now},
			Enabled:        true,
			Type:           types.EntityTypeResource,
			Connector:      connectorID,
			Component:      event.Component,
			ComponentInfos: componentInfos,
			Infos:          map[string]types.Info{},
			ImpactLevel:    types.EntityDefaultImpactLevel,
			Created:        now,
			LastEventDate:  lastEventDate,
			Healthcheck:    event.Healthcheck,
		}
		if resource.ID != event.Upstream {
			resource.Upstream = event.Upstream
		}

		commRegister.RegisterInsert(resource)
		report.CheckResource = true
		report.IsNew = true

		event.Entity = resource

		return report, nil
	}

	if resource.Connector != connectorID && !connectorExist {
		resource.Connector = connectorID

		upd := bson.M{
			"connector": connectorID,
		}
		if lastEventDate != nil {
			upd["last_event_date"] = *lastEventDate
		}
		commRegister.RegisterUpdate(resourceID, upd)
		commRegister.RegisterInsert(&types.Entity{
			ID:            connectorID,
			Name:          connectorName,
			EnableHistory: []datetime.CpsTime{now},
			Enabled:       true,
			Type:          types.EntityTypeConnector,
			Infos:         map[string]types.Info{},
			ImpactLevel:   types.EntityDefaultImpactLevel,
			Created:       now,
			LastEventDate: lastEventDate,
			Healthcheck:   event.Healthcheck,
		})

		report.CheckResource = true
		report.CheckConnector = true
	} else if lastEventDate != nil {
		commRegister.RegisterUpdate(connectorID, bson.M{"last_event_date": *lastEventDate})
		commRegister.RegisterUpdate(resourceID, bson.M{"last_event_date": *lastEventDate})
	}

	if resource.Upstream != event.Upstream {
		if resourceID == event.Upstream {
			commRegister.RegisterUpdate(resourceID, bson.M{
				"upstream":            "",
				"is_upstream_changed": true,
			})
			resource.IsUpstreamChanged = true
			resource.Upstream = ""
		} else {
			commRegister.RegisterUpdate(resourceID, bson.M{
				"upstream":            event.Upstream,
				"is_upstream_changed": true,
			})
			resource.IsUpstreamChanged = true
			resource.Upstream = event.Upstream
		}
	}

	resource.LastEventDate = lastEventDate
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

	if !event.IsContextable() || event.IsOnlyServiceUpdate() || event.Initiator == types.InitiatorSystem {
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
	var lastEventDate *datetime.CpsTime
	if event.EventType == types.EventTypeCheck {
		lastEventDate = &now
	}

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
				LastEventDate: lastEventDate,
				Healthcheck:   event.Healthcheck,
			})

			report.CheckConnector = true
		} else if lastEventDate != nil {
			commRegister.RegisterUpdate(connectorID, bson.M{"last_event_date": *lastEventDate})
		}

		component = &types.Entity{
			ID:            componentID,
			Name:          componentID,
			EnableHistory: []datetime.CpsTime{now},
			Enabled:       true,
			Type:          types.EntityTypeComponent,
			Connector:     connectorID,
			Component:     componentID,
			Upstream:      event.Upstream,
			Infos:         map[string]types.Info{},
			ImpactLevel:   types.EntityDefaultImpactLevel,
			Created:       now,
			LastEventDate: lastEventDate,
			Healthcheck:   event.Healthcheck,
		}
		if component.ID != event.Upstream {
			component.Upstream = event.Upstream
		}

		commRegister.RegisterInsert(component)
		report.CheckComponent = true
		report.IsNew = true

		event.Entity = component

		return report, nil
	}

	if component.Connector != connectorID && !connectorExist {
		component.Connector = connectorID

		upd := bson.M{
			"connector": connectorID,
		}
		if lastEventDate != nil {
			upd["last_event_date"] = *lastEventDate
		}
		commRegister.RegisterUpdate(componentID, upd)
		commRegister.RegisterInsert(&types.Entity{
			ID:            connectorID,
			Name:          connectorName,
			EnableHistory: []datetime.CpsTime{now},
			Enabled:       true,
			Type:          types.EntityTypeConnector,
			Infos:         map[string]types.Info{},
			ImpactLevel:   types.EntityDefaultImpactLevel,
			Created:       now,
			LastEventDate: lastEventDate,
			Healthcheck:   event.Healthcheck,
		})

		report.CheckComponent = true
		report.CheckConnector = true
	} else if lastEventDate != nil {
		commRegister.RegisterUpdate(connectorID, bson.M{"last_event_date": *lastEventDate})
		commRegister.RegisterUpdate(componentID, bson.M{"last_event_date": *lastEventDate})
	}

	if component.Upstream != event.Upstream {
		report.CheckComponent = true
		if componentID == event.Upstream {
			commRegister.RegisterUpdate(componentID, bson.M{
				"upstream":            "",
				"is_upstream_changed": true,
			})
			component.IsUpstreamChanged = true
			component.Upstream = ""
		} else {
			commRegister.RegisterUpdate(componentID, bson.M{
				"upstream":            event.Upstream,
				"is_upstream_changed": true,
			})
			component.IsUpstreamChanged = true
			component.Upstream = event.Upstream
		}
	}

	component.LastEventDate = lastEventDate
	event.Entity = component

	return report, nil
}

func (m *manager) HandleService(ctx context.Context, event *types.Event, commRegister libmongo.CommandsRegister) (Report, error) {
	report := Report{}

	serviceID := event.Component
	service, err := m.getEntity(ctx, serviceID)
	if err != nil {
		return report, err
	}

	if service == nil {
		return report, fmt.Errorf("service %s doesn't exist", serviceID)
	}

	if service.SoftDeleted != nil {
		event.Entity = service

		return report, nil
	}

	if event.IsContextable() && !event.IsOnlyServiceUpdate() && event.EventType == types.EventTypeCheck {
		now := datetime.NewCpsTime()
		commRegister.RegisterUpdate(serviceID, bson.M{"last_event_date": now})
		service.LastEventDate = &now
	}

	if event.IsOnlyServiceUpdate() {
		report.CheckService = true
	}

	event.Entity = service

	return report, nil
}

func (m *manager) HandleConnector(ctx context.Context, event *types.Event, commRegister libmongo.CommandsRegister) (Report, error) {
	report := Report{}

	connectorName := event.ConnectorName
	connectorID := event.Connector + "/" + connectorName
	connector, err := m.getEntity(ctx, connectorID)
	if err != nil {
		return report, err
	}

	if connector == nil {
		return report, fmt.Errorf("connector %s doesn't exist", connectorID)
	}

	if event.IsContextable() && !event.IsOnlyServiceUpdate() && event.EventType == types.EventTypeCheck {
		now := datetime.NewCpsTime()
		commRegister.RegisterUpdate(connectorID, bson.M{"last_event_date": now})
		connector.LastEventDate = &now
	}

	if event.IsOnlyServiceUpdate() {
		report.CheckConnector = true
	}

	event.Entity = connector

	return report, nil
}

func (m *manager) UpdateImpactedServicesFromDependencies(ctx context.Context) error {
	cursor, err := m.entityCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"enabled":    true,
			"connector":  bson.M{"$gt": ""},
			"services.0": bson.M{"$exists": true},
		}},
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

	if err := cursor.Err(); err != nil {
		return err
	}

	if len(writeModels) > 0 {
		err = m.adapter.Bulk(ctx, writeModels)
	}

	return err
}

func (m *manager) ProcessComponentInfos(ctx context.Context, component *types.Entity, updatedComponentInfos []string, stateSettingUpdated bool, commRegister libmongo.CommandsRegister) ([]string, []string, error) {
	if component.Type != types.EntityTypeComponent || !stateSettingUpdated && len(updatedComponentInfos) == 0 {
		return nil, nil, nil
	}

	cursor, err := m.entityCollection.Find(
		ctx,
		bson.M{"_id": bson.M{"$ne": component.ID}, "component": component.ID},
	)
	if err != nil {
		return nil, nil, err
	}

	defer cursor.Close(ctx)

	var resourceIDs []string
	var servicesIDs []string
	servicesIDsMap := make(map[string]bool)

	for cursor.Next(ctx) {
		var resource types.Entity

		err = cursor.Decode(&resource)
		if err != nil {
			return nil, nil, err
		}

		update := bson.M{}

		if stateSettingUpdated {
			matched := component.StateInfo != nil
			if matched && len(component.StateInfo.InheritedPattern) > 0 {
				matched, err = match.MatchEntityPattern(component.StateInfo.InheritedPattern, &resource)
				if err != nil {
					return nil, nil, err
				}
			}

			if matched {
				update["component_state_settings"] = true
				update["component_state_settings_to_add"] = false
				update["component_state_settings_to_remove"] = false
			} else {
				update["component_state_settings"] = false
				update["component_state_settings_to_add"] = false
				update["component_state_settings_to_remove"] = false
			}
		}

		if len(updatedComponentInfos) > 0 {
			resource.ComponentInfos = component.Infos
			update["component_infos"] = component.Infos

			// do not pass commRegister services will be assigned by service recompute event.
			affectedServicesIDsMap := m.AssignServicesByComponentInfoNames(&resource, updatedComponentInfos, nil)
			for serviceID := range affectedServicesIDsMap {
				if !servicesIDsMap[serviceID] {
					servicesIDsMap[serviceID] = true
					servicesIDs = append(servicesIDs, serviceID)
				}
			}
		}

		commRegister.RegisterUpdate(resource.ID, update)
		resourceIDs = append(resourceIDs, resource.ID)
	}

	return resourceIDs, servicesIDs, nil
}

func (m *manager) AssignStateSetting(ctx context.Context, entity *types.Entity, commRegister libmongo.CommandsRegister) (bool, error) {
	return m.stateSettingService.AssignStateSetting(ctx, entity, commRegister)
}

func (m *manager) getResourceEntities(ctx context.Context, resourceID, componentID, connectorID string) (*types.Entity, bool, map[string]types.Info, bool, error) {
	var resource *types.Entity
	var componentExist bool
	var componentInfos map[string]types.Info
	var connectorExist bool

	cursor, err := m.entityCollection.Find(ctx, bson.M{"_id": bson.M{"$in": bson.A{resourceID, componentID, connectorID}}})
	if err != nil {
		return nil, componentExist, componentInfos, connectorExist, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var ent types.Entity

		err = cursor.Decode(&ent)
		if err != nil {
			return nil, componentExist, componentInfos, connectorExist, err
		}

		switch ent.Type {
		case types.EntityTypeResource:
			resource = &ent
		case types.EntityTypeComponent:
			componentExist = true
			componentInfos = ent.Infos
		default:
			connectorExist = true
		}
	}

	return resource, componentExist, componentInfos, connectorExist, nil
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
			connectorExist = true
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
	cursor, err := m.entityCollection.Aggregate(
		ctx,
		[]bson.M{
			{
				"$match": bson.M{"services": serviceID},
			},
			{
				"$project": bson.M{
					"_id":                1,
					"services":           1,
					"inherited_services": 1,
					"services_to_add":    1,
				},
			},
			{
				"$addFields": bson.M{
					"services": bson.M{
						"$setDifference": bson.A{
							"$services",
							bson.A{serviceID},
						},
					},
				},
			},
			{
				"$addFields": bson.M{
					"inherited_services": bson.M{
						"$setDifference": bson.A{
							"$inherited_services",
							bson.A{serviceID},
						},
					},
				},
			},
			{
				"$addFields": bson.M{
					"services_to_add": bson.M{
						"$setDifference": bson.A{
							"$services_to_add",
							bson.A{serviceID},
						},
					},
				},
			},
		},
	)
	if err != nil {
		return err
	}

	err = cursor.All(ctx, &dependedEntities)
	if err != nil {
		return err
	}

	for _, ent := range dependedEntities {
		commRegister.RegisterUpdate(ent.ID, bson.M{
			"services":           ent.Services,
			"inherited_services": ent.InheritedServices,
			"services_to_add":    ent.ServicesToAdd,
		})
	}

	return nil
}
