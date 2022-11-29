package importcontextgraph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type worker struct {
	entityCollection        libmongo.DbCollection
	categoryCollection      libmongo.DbCollection
	alarmCollection         libmongo.DbCollection
	alarmResolvedCollection libmongo.DbCollection

	publisher         EventPublisher
	metricMetaUpdater metrics.MetaUpdater
}

func NewWorker(
	dbClient libmongo.DbClient,
	publisher EventPublisher,
	metricMetaUpdater metrics.MetaUpdater,
) Worker {
	return &worker{
		entityCollection:        dbClient.Collection(libmongo.EntityMongoCollection),
		categoryCollection:      dbClient.Collection(libmongo.EntityCategoryMongoCollection),
		alarmCollection:         dbClient.Collection(libmongo.AlarmMongoCollection),
		alarmResolvedCollection: dbClient.Collection(libmongo.ResolvedAlarmMongoCollection),

		publisher:         publisher,
		metricMetaUpdater: metricMetaUpdater,
	}
}

func (w *worker) Work(ctx context.Context, filename, source string) (stats Stats, resErr error) {
	startTime := time.Now()
	defer func() {
		stats.ExecTime = time.Since(startTime)
	}()

	res, err := w.parseFile(ctx, filename, source, false)
	if err != nil {
		return stats, err
	}

	if len(res.writeModels) == 0 {
		return stats, fmt.Errorf("empty import")
	}

	stats.Updated, stats.Deleted, err = w.bulkWrite(ctx, res.writeModels, canopsis.DefaultBulkSize, canopsis.DefaultBulkBytesSize)
	if err != nil {
		return stats, err
	}

	w.metricMetaUpdater.UpdateAll(ctx)

	for _, event := range res.deleteEvents {
		err = w.publisher.SendEvent(ctx, event)
		if err != nil {
			return stats, err
		}
	}

	if stats.Updated > 0 || stats.Deleted > 0 {
		err = w.sendUpdateServiceEvents(ctx)
		if err != nil {
			return stats, err
		}
	}

	return stats, nil
}

func (w *worker) WorkPartial(ctx context.Context, filename, source string) (stats Stats, resErr error) {
	startTime := time.Now()
	defer func() {
		stats.ExecTime = time.Since(startTime)
	}()

	res, err := w.parseFile(ctx, filename, source, true)
	if err != nil {
		return stats, err
	}

	if len(res.writeModels) == 0 {
		return stats, fmt.Errorf("empty import")
	}

	stats.Updated, stats.Deleted, err = w.bulkWrite(ctx, res.writeModels, canopsis.DefaultBulkSize, canopsis.DefaultBulkBytesSize)
	if err != nil {
		return stats, err
	}

	if len(res.updatedIds) > 0 {
		w.metricMetaUpdater.UpdateById(ctx, res.updatedIds...)
	}
	if len(res.removedIds) > 0 {
		w.metricMetaUpdater.DeleteById(ctx, res.removedIds...)
	}

	serviceCount, err := w.entityCollection.CountDocuments(ctx, bson.M{"type": types.EntityTypeService})
	if err != nil {
		return stats, err
	}

	if serviceCount == 0 {
		for _, event := range res.serviceEvents {
			err := w.publisher.SendEvent(ctx, event)
			if err != nil {
				return stats, err
			}
		}
	} else if len(res.serviceEvents)+len(res.basicEntityEvents) <= int(serviceCount) {
		for _, event := range res.serviceEvents {
			err := w.publisher.SendEvent(ctx, event)
			if err != nil {
				return stats, err
			}
		}
		for _, event := range res.basicEntityEvents {
			fixedEvent, err := w.fillEventAfterLinksUpdate(ctx, event)
			if err != nil {
				return stats, err
			}
			if fixedEvent.EventType != "" {
				err = w.publisher.SendEvent(ctx, fixedEvent)
				if err != nil {
					return stats, err
				}
			}
		}
	} else {
		err = w.sendUpdateServiceEvents(ctx)
		if err != nil {
			return stats, err
		}
	}

	return stats, nil
}

func (w *worker) parseFile(ctx context.Context, filename, source string, withEvents bool) (_ parseResult, resErr error) {
	res := parseResult{}
	file, err := os.Open(filename)
	if err != nil {
		return res, err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			if resErr == nil {
				resErr = err
			}
		} else {
			err = os.Remove(filename)
			if err != nil && resErr == nil {
				resErr = err
			}
		}
	}()

	writeModels := make([]mongo.WriteModel, 0)
	var entityParseRes parseEntityResult
	decoder := json.NewDecoder(file)

	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}

			return res, err
		}

		if t == "cis" {
			t, err := decoder.Token()
			if err != nil {
				return res, fmt.Errorf("failed to parse cis: %v", err)
			}

			if t != json.Delim('[') {
				return res, fmt.Errorf("cis should be an array")
			}

			entityParseRes, err = w.parseEntities(ctx, decoder, source, withEvents)
			if err != nil {
				return res, err
			}

			writeModels = append(writeModels, entityParseRes.writeModels...)
		}
	}

	res.writeModels = writeModels
	res.updatedIds = entityParseRes.updatedIds
	res.removedIds = entityParseRes.removedIds
	res.serviceEvents = entityParseRes.serviceEvents
	res.basicEntityEvents = entityParseRes.basicEntityEvents
	res.deleteEvents = entityParseRes.deleteEvents

	return res, nil
}

func (w *worker) parseEntities(
	ctx context.Context,
	decoder *json.Decoder,
	source string,
	withEvents bool,
) (parseEntityResult, error) {
	res := parseEntityResult{}
	writeModels := make([]mongo.WriteModel, 0)
	updatedIds := make([]string, 0)
	removedIds := make([]string, 0)

	now := types.NewCpsTime()
	componentInfos := make(map[string]map[string]types.Info)

	componentsExist := make(map[string]bool)
	componentsToDelete := make(map[string]bool)
	componentsToDisable := make(map[string]bool)

	createLinks := make(map[string][]string)
	deletedResources := make(map[string]bool)
	deletedEvents := make([]types.Event, 0)

	for decoder.More() {
		var ci ConfigurationItem
		err := decoder.Decode(&ci)
		if err != nil {
			return res, fmt.Errorf("failed to decode cis item: %v", err)
		}

		err = w.validate(ci)
		if err != nil {
			return res, fmt.Errorf("ci = %s, validation error: %s", ci.ID, err.Error())
		}

		w.fillDefaultFields(&ci, source, now)

		var oldEntity ConfigurationItem

		switch ci.Action {
		case ActionSet:
			err = w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil && err != mongo.ErrNoDocuments {
				return res, err
			}

			if ci.Type == types.EntityTypeResource {
				createLinks[ci.Component] = append(createLinks[ci.Component], ci.ID)
			}

			if ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
				componentsExist[ci.ID] = true
			}

			if err == mongo.ErrNoDocuments {
				writeModels = append(writeModels, w.createEntity(ci))
				if ci.Type == types.EntityTypeResource {
					if _, ok := componentsExist[ci.Component]; !ok {
						componentsExist[ci.Component] = false
					}
				}
			} else {
				writeModels = append(writeModels, w.updateEntity(&ci, oldEntity, true))
				if ci.Type == types.EntityTypeResource {
					componentsExist[ci.Component] = true
				}
			}
		case ActionDelete:
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return res, fmt.Errorf("failed to delete an entity with _id = %s", ci.ID)
				}

				return res, err
			}

			if ci.Type == types.EntityTypeResource && !deletedResources[ci.ID] {
				writeModels = append(writeModels, w.deleteEntity(ci.ID)...)
				deletedEvents = append(deletedEvents, w.createResolveDeletedEvent(types.EntityTypeResource, ci.ID, now))

				deletedResources[ci.ID] = true
			}

			if ci.Type == types.EntityTypeComponent {
				componentsToDelete[ci.ID] = true

				for _, resourceID := range oldEntity.Depends {
					if !deletedResources[resourceID] {
						writeModels = append(writeModels, w.deleteEntity(resourceID)...)
						deletedEvents = append(deletedEvents, w.createResolveDeletedEvent(types.EntityTypeResource, resourceID, now))

						deletedResources[resourceID] = true
					}
				}

				writeModels = append(writeModels, w.deleteEntity(ci.ID)...)
				deletedEvents = append(deletedEvents, w.createResolveDeletedEvent(types.EntityTypeResource, ci.ID, now))
			}
		case ActionEnable:
			updatedIds = append(updatedIds, ci.ID)
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return res, fmt.Errorf("failed to enable an entity with _id = %s", ci.ID)
				}

				return res, err
			}

			writeModels = append(writeModels, w.changeState(ci.ID, true, source, now))
			if ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
			}
		case ActionDisable:
			updatedIds = append(updatedIds, ci.ID)
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return res, fmt.Errorf("failed to disable an entity with _id = %s", ci.ID)
				}

				return res, err
			}

			writeModels = append(writeModels, w.changeState(ci.ID, false, source, now))
			if ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
				componentsToDisable[ci.ID] = true
			}
		default:
			return res, fmt.Errorf("the action %s is not recognized", ci.Action)
		}
	}

	for componentName, exists := range componentsExist {
		if componentsToDelete[componentName] {
			return res, fmt.Errorf("can't create and delete component simutaneously")
		}

		if componentsToDisable[componentName] {
			return res, fmt.Errorf("can't create and disable component simutaneously")
		}

		if !exists {
			var oldEntity types.Entity
			updatedIds = append(updatedIds, componentName)
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": componentName}).Decode(&oldEntity)
			if err != nil && err != mongo.ErrNoDocuments {
				return res, err
			}

			if err == mongo.ErrNoDocuments {
				ci := ConfigurationItem{
					ID:           componentName,
					Name:         componentName,
					Component:    componentName,
					ImpactLevel:  1,
					ImportSource: source,
					Imported:     now,
					Type:         types.EntityTypeComponent,
					Enabled:      true,
				}

				writeModels = append(writeModels, w.createEntity(ci))
			} else {
				if !oldEntity.Enabled {
					return res, fmt.Errorf("can't create and disable component simutaneously")
				}

				componentInfos[componentName] = oldEntity.Infos
			}
		}

		resourceIDs, ok := createLinks[componentName]
		if ok && len(resourceIDs) > 0 {
			writeModels = append(writeModels, w.createLink(resourceIDs, componentName)...)
		}

		if len(componentInfos[componentName]) > 0 {
			writeModels = append(writeModels, w.updateComponentInfosOnComponentUpdate(componentName, componentInfos[componentName]))
		}
	}

	res.updatedIds = updatedIds
	res.removedIds = removedIds
	res.writeModels = writeModels
	res.deleteEvents = deletedEvents

	return res, nil
}

func (w *worker) sendUpdateServiceEvents(ctx context.Context) error {
	cursor, err := w.entityCollection.Find(ctx, bson.M{"type": types.EntityTypeService})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var service types.Entity

		err := cursor.Decode(&service)
		if err != nil {
			return err
		}

		err = w.publisher.SendEvent(ctx, types.Event{
			EventType:     types.EventTypeRecomputeEntityService,
			Connector:     types.ConnectorEngineService,
			ConnectorName: types.ConnectorEngineService,
			Component:     service.ID,
			Timestamp:     types.CpsTime{Time: time.Now()},
			Author:        canopsis.DefaultEventAuthor,
			SourceType:    types.SourceTypeService,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *worker) bulkWrite(ctx context.Context, writeModels []mongo.WriteModel, limit, limitBytes int) (int64, int64, error) {
	var updated, deleted int64

	start := 0
	end := 0
	for {
		if end == len(writeModels) {
			break
		}

		start = end
		end = start + limit

		if end > len(writeModels) {
			end = len(writeModels)
		}

		bulkSize := 0
		for i := start; i < end; i++ {
			b, err := bson.Marshal(writeModels[i])
			if err != nil {
				return 0, 0, err
			}

			l := len(b)
			if l+bulkSize >= limitBytes {
				if i > start {
					end = i
				} else {
					end = start + 1
				}
				break
			}

			bulkSize += l
		}

		p := writeModels[start:end]
		result, err := w.entityCollection.BulkWrite(ctx, p)
		if err != nil {
			return 0, 0, err
		}

		updated += result.UpsertedCount + result.ModifiedCount
		deleted += result.DeletedCount
	}

	return updated, deleted, nil
}

func (w *worker) validate(ci ConfigurationItem) error {
	switch ci.Type {
	case types.EntityTypeService:
	case types.EntityTypeResource:
		if ci.Component == "" {
			return fmt.Errorf("resource %s contains empty component", ci.Name)
		}
	case types.EntityTypeComponent:
	default:
		return fmt.Errorf("type is not valid %s", ci.Type)
	}

	if ci.Name == "" {
		return fmt.Errorf("empty name is not allowed")
	}

	if ci.Type != types.EntityTypeService && len(ci.EntityPattern) > 0 {
		return fmt.Errorf("contains entity pattern, but ci is not a service")
	}

	return nil
}

func (w *worker) fillDefaultFields(ci *ConfigurationItem, source string, now types.CpsTime) {
	ci.ID = ci.Name
	if ci.Type == types.EntityTypeResource {
		ci.ID = ci.Name + "/" + ci.Component
	}

	if ci.ImpactLevel == 0 {
		ci.ImpactLevel = 1
	}

	ci.ImportSource = source
	ci.Imported = now
}

func (w *worker) createLink(from []string, to string) []mongo.WriteModel {
	updateTo := bson.M{"$addToSet": bson.M{"depends": bson.M{"$each": from}}}
	updateFrom := bson.M{"$addToSet": bson.M{"impact": to}}

	return []mongo.WriteModel{
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"_id": to}).
			SetUpdate(updateTo),
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"_id": bson.M{"$in": from}}).
			SetUpdate(updateFrom),
	}
}

func (w *worker) deleteLink(from []string, to string) []mongo.WriteModel {
	updateTo := bson.M{"$pull": bson.M{"depends": bson.M{"$in": from}}}
	updateFrom := bson.M{"$pull": bson.M{"impact": to}}

	return []mongo.WriteModel{
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"_id": to}).
			SetUpdate(updateTo),
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"_id": bson.M{"$in": from}}).
			SetUpdate(updateFrom),
	}
}

func (w *worker) createEntity(ci ConfigurationItem) mongo.WriteModel {
	ci.Depends = []string{}
	ci.Impact = []string{}
	ci.EnableHistory = make([]int64, 0)

	if ci.Type == types.EntityTypeComponent {
		ci.Component = ci.ID
	}

	if ci.Infos == nil {
		ci.Infos = make(map[string]types.Info)
	}

	if ci.Measurements == nil {
		ci.Measurements = make([]interface{}, 0)
	}

	now := types.CpsTime{Time: time.Now()}

	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"_id": ci.ID}).
		SetUpdate(bson.M{"$set": ci, "$setOnInsert": bson.M{"created": now}}).
		SetUpsert(true)
}

func (w *worker) updateEntity(ci *ConfigurationItem, oldEntity ConfigurationItem, mergeInfos bool) mongo.WriteModel {
	ci.Depends = oldEntity.Depends
	ci.Impact = oldEntity.Impact
	ci.EnableHistory = oldEntity.EnableHistory

	if ci.Type == types.EntityTypeComponent {
		ci.Component = ci.ID
	}

	if ci.Infos == nil {
		ci.Infos = make(map[string]types.Info)
	}

	if mergeInfos {
		for k, v := range ci.Infos {
			oldEntity.Infos[k] = v
		}

		ci.Infos = oldEntity.Infos
	}

	now := types.CpsTime{Time: time.Now()}

	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"_id": ci.ID}).
		SetUpdate(bson.M{"$set": ci, "$setOnInsert": bson.M{"created": now}}).
		SetUpsert(true)
}

func (w *worker) changeState(id string, enabled bool, importSource string, imported types.CpsTime) mongo.WriteModel {
	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"_id": id}).
		SetUpdate(bson.M{"$set": bson.M{
			"enabled":       enabled,
			"import_source": importSource,
			"imported":      imported,
		}})
}

func (w *worker) deleteEntity(id string) []mongo.WriteModel {
	return []mongo.WriteModel{
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"impact": id}).
			SetUpdate(bson.M{"$pull": bson.M{"impact": id}}),
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"depends": id}).
			SetUpdate(bson.M{"$pull": bson.M{"depends": id}}),
		mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": id}).
			SetUpdate(bson.M{"$set": bson.M{"soft_deleted": true}}),
	}
}

func (w *worker) createResolveDeletedEvent(t, id string, now types.CpsTime) types.Event {
	event := types.Event{
		EventType:     types.EventTypeResolveDeleted,
		Timestamp:     now,
		Author:        canopsis.DefaultEventAuthor,
		Connector:     "import",
		ConnectorName: "import",
	}

	switch t {
	case types.EntityTypeComponent:
		event.Component = id
		event.SourceType = types.SourceTypeComponent
	case types.EntityTypeResource:
		idParts := strings.Split(id, "/")
		if len(idParts) != 2 {
			return types.Event{}
		}
		event.Resource = idParts[0]
		event.Component = idParts[1]
		event.SourceType = types.SourceTypeResource
	}

	return event
}

func (w *worker) updateComponentInfosOnComponentUpdate(componentID string, infos map[string]types.Info) mongo.WriteModel {
	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"type": types.EntityTypeResource, "impact": componentID}).
		SetUpdate(bson.M{"$set": bson.M{"component_infos": infos}})
}

func (w *worker) createServiceEvent(ci ConfigurationItem, eventType string, now types.CpsTime) types.Event {
	return types.Event{
		EventType:     eventType,
		Timestamp:     now,
		Author:        canopsis.DefaultEventAuthor,
		Connector:     types.ConnectorEngineService,
		ConnectorName: types.ConnectorEngineService,
		Component:     ci.ID,
		SourceType:    types.SourceTypeService,
	}
}

func (w *worker) createBasicEntityEvent(ctx context.Context, ci, oldEntity ConfigurationItem, eventType string, now types.CpsTime) (types.Event, error) {
	event := types.Event{
		EventType: eventType,
		Timestamp: now,
		Author:    canopsis.DefaultEventAuthor,
	}

	alarm := types.Alarm{}
	findOptions := options.FindOne().SetSort(bson.M{"t": -1}).SetProjection(bson.M{"v.steps": 0})
	err := w.alarmCollection.FindOne(ctx, bson.M{"d": ci.ID}, findOptions).Decode(&alarm)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return event, fmt.Errorf("failed to fetch an alarm: %w", err)
	}
	if alarm.ID == "" {
		err := w.alarmResolvedCollection.FindOne(ctx, bson.M{"d": ci.ID}, findOptions).Decode(&alarm)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return event, fmt.Errorf("failed to fetch an alarm: %w", err)
		}
	}
	if alarm.ID != "" {
		event.Connector = alarm.Value.Connector
		event.ConnectorName = alarm.Value.ConnectorName
		event.Component = alarm.Value.Component
		event.Resource = alarm.Value.Resource
		event.SourceType = event.DetectSourceType()
	} else {
		name := ci.Name
		if name == "" {
			name = oldEntity.Name
		}
		switch ci.Type {
		case types.EntityTypeConnector:
			event.Connector = strings.TrimSuffix(ci.ID, "/"+name)
			event.ConnectorName = name
			event.SourceType = types.SourceTypeConnector
		case types.EntityTypeComponent:
			event.Component = ci.ID
			event.SourceType = types.SourceTypeComponent
		case types.EntityTypeResource:
			event.Resource = ci.ID // use id to retrieve component and connector after links parsing
			event.SourceType = types.SourceTypeResource
		}
	}

	return event, nil
}

func (w *worker) fillEventAfterLinksUpdate(ctx context.Context, event types.Event) (types.Event, error) {
	switch event.SourceType {
	case types.SourceTypeComponent:
		if event.Connector == "" {
			connector := types.Entity{}
			err := w.entityCollection.FindOne(ctx, bson.M{
				"type":    types.EntityTypeConnector,
				"depends": event.Component,
			}, options.FindOne().SetProjection(bson.M{"impact": 0, "depends": 0})).Decode(&connector)
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				return types.Event{}, err
			}
			if connector.ID == "" {
				event.Connector = defaultConnector
				event.ConnectorName = defaultConnectorName
			} else {
				event.Connector = strings.TrimSuffix(connector.ID, "/"+connector.Name)
				event.ConnectorName = connector.Name
			}
		}
	case types.SourceTypeResource:
		if event.Connector == "" {
			connector := types.Entity{}
			err := w.entityCollection.FindOne(ctx, bson.M{
				"type":   types.EntityTypeConnector,
				"impact": event.Resource,
			}, options.FindOne().SetProjection(bson.M{"impact": 0, "depends": 0})).Decode(&connector)
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				return types.Event{}, err
			}
			component := types.Entity{}
			err = w.entityCollection.FindOne(ctx, bson.M{
				"type":    types.EntityTypeComponent,
				"depends": event.Resource,
			}, options.FindOne().SetProjection(bson.M{"impact": 0, "depends": 0})).Decode(&component)
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				return types.Event{}, err
			}

			if connector.ID == "" {
				event.Connector = defaultConnector
				event.ConnectorName = defaultConnectorName
			} else {
				event.Connector = strings.TrimSuffix(connector.ID, "/"+connector.Name)
				event.ConnectorName = connector.Name
			}
			if component.ID == "" {
				idParts := strings.Split(event.Resource, "/")
				if len(idParts) != 2 {
					return types.Event{}, nil
				}
				event.Resource = idParts[0]
				event.Component = idParts[1]
			} else {
				event.Component = component.Name
				event.Resource = strings.TrimSuffix(event.Resource, "/"+component.Name)
			}
		}
	}

	return event, nil
}
