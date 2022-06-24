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

const (
	defaultConnector     = "taskhandler"
	defaultConnectorName = "task_importctx"
)

type worker struct {
	entityCollection        libmongo.DbCollection
	categoryCollection      libmongo.DbCollection
	alarmCollection         libmongo.DbCollection
	alarmResolvedCollection libmongo.DbCollection

	publisher         EventPublisher
	metricMetaUpdater metrics.MetaUpdater
}

type parseResult struct {
	writeModels []mongo.WriteModel

	serviceEvents     []types.Event
	basicEntityEvents []types.Event
}

type parseEntityResult struct {
	writeModels []mongo.WriteModel

	componentInfos map[string]map[string]interface{}

	serviceEvents     []types.Event
	basicEntityEvents []types.Event
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

	err = w.sendUpdateServiceEvents(ctx)
	if err != nil {
		return stats, err
	}

	w.metricMetaUpdater.UpdateAll(ctx)

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

	serviceCount, err := w.entityCollection.CountDocuments(ctx, bson.M{"type": types.EntityTypeService})
	if err != nil {
		return stats, err
	}

	if serviceCount == 0 {
		for _, event := range res.serviceEvents {
			err := w.publisher.SendEvent(event)
			if err != nil {
				return stats, err
			}
		}
	} else if len(res.serviceEvents)+len(res.basicEntityEvents) <= int(serviceCount) {
		for _, event := range res.serviceEvents {
			err := w.publisher.SendEvent(event)
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
				err = w.publisher.SendEvent(fixedEvent)
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

	w.metricMetaUpdater.UpdateAll(ctx)

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

		if t == "links" {
			t, err := decoder.Token()
			if err != nil {
				return res, fmt.Errorf("failed to parse links: %v", err)
			}

			if t != json.Delim('[') {
				return res, fmt.Errorf("links should be an array")
			}

			linkWriteModels, err := w.parseLinks(ctx, decoder, entityParseRes.componentInfos)
			if err != nil {
				return res, err
			}

			writeModels = append(writeModels, linkWriteModels...)
		}
	}

	res.writeModels = writeModels
	res.serviceEvents = entityParseRes.serviceEvents
	res.basicEntityEvents = entityParseRes.basicEntityEvents

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
	serviceEvents := make([]types.Event, 0)
	basicEntityEvents := make([]types.Event, 0)
	now := types.NewCpsTime()
	componentInfos := make(map[string]map[string]interface{})

	for decoder.More() {
		var ci ConfigurationItem
		err := decoder.Decode(&ci)
		if err != nil {
			return res, fmt.Errorf("failed to decode cis item: %v", err)
		}

		w.fillDefaultFields(&ci)
		err = w.validate(ci)
		if err != nil {
			return res, fmt.Errorf("ci = %s, validation error: %s", ci.ID, err.Error())
		}

		ci.ImportSource = source
		ci.Imported = now

		eventType := ""
		var oldEntity ConfigurationItem

		switch ci.Action {
		case ActionCreate:
			writeModels = append(writeModels, w.createEntity(ci))
			if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
			}
			if ci.Enabled != nil && *ci.Enabled {
				switch *ci.Type {
				case types.EntityTypeService:
					eventType = types.EventTypeRecomputeEntityService
				default:
					eventType = types.EventTypeEntityUpdated
				}
			}
		case ActionSet:
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil && err != mongo.ErrNoDocuments {
				return res, err
			}

			if err == mongo.ErrNoDocuments {
				writeModels = append(writeModels, w.createEntity(ci))
				if ci.Enabled != nil && *ci.Enabled {
					switch *ci.Type {
					case types.EntityTypeService:
						eventType = types.EventTypeRecomputeEntityService
					default:
						eventType = types.EventTypeEntityUpdated
					}
				}
				break
			}

			writeModels = append(writeModels, w.updateEntity(&ci, oldEntity, true))
			if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
				writeModels = append(writeModels, w.updateComponentInfosOnComponentUpdate(ci))
			}

			if oldEntity.Enabled != nil && *oldEntity.Enabled || ci.Enabled != nil && *ci.Enabled {
				switch *ci.Type {
				case types.EntityTypeService:
					eventType = types.EventTypeRecomputeEntityService
				default:
					if *oldEntity.Enabled && *ci.Enabled {
						eventType = types.EventTypeEntityUpdated
					} else {
						eventType = types.EventTypeEntityToggled
					}
				}
			}
		case ActionUpdate:
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return res, fmt.Errorf("failed to update an entity with _id = %s", ci.ID)
				}

				return res, err
			}

			writeModels = append(writeModels, w.updateEntity(&ci, oldEntity, false))
			if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
				writeModels = append(writeModels, w.updateComponentInfosOnComponentUpdate(ci))
			}

			if oldEntity.Enabled != nil && *oldEntity.Enabled || ci.Enabled != nil && *ci.Enabled {
				switch *ci.Type {
				case types.EntityTypeService:
					eventType = types.EventTypeRecomputeEntityService
				default:
					if *oldEntity.Enabled && *ci.Enabled {
						eventType = types.EventTypeEntityUpdated
					} else {
						eventType = types.EventTypeEntityToggled
					}
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

			writeModels = append(writeModels, w.deleteEntity(ci)...)
			if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = nil
				writeModels = append(writeModels, w.updateComponentInfosOnComponentDelete(ci))
			}

			if oldEntity.Enabled != nil && *oldEntity.Enabled && *ci.Type == types.EntityTypeService {
				eventType = types.EventTypeRecomputeEntityService
			}
		case ActionEnable:
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return res, fmt.Errorf("failed to enable an entity with _id = %s", ci.ID)
				}

				return res, err
			}

			writeModels = append(writeModels, w.changeState(ci.ID, true, source, now))
			if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
			}

			switch *ci.Type {
			case types.EntityTypeService:
				eventType = types.EventTypeRecomputeEntityService
			default:
				eventType = types.EventTypeEntityToggled
			}
		case ActionDisable:
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return res, fmt.Errorf("failed to disable an entity with _id = %s", ci.ID)
				}

				return res, err
			}

			writeModels = append(writeModels, w.changeState(ci.ID, false, source, now))
			if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
			}

			switch *ci.Type {
			case types.EntityTypeService:
				eventType = types.EventTypeRecomputeEntityService
			default:
				eventType = types.EventTypeEntityToggled
			}
		default:
			return res, fmt.Errorf("the action %s is not recognized", ci.Action)
		}

		if withEvents && eventType != "" {
			switch *ci.Type {
			case types.EntityTypeService:
				serviceEvents = append(serviceEvents, w.createServiceEvent(ci, eventType, now))
			default:
				event, err := w.createBasicEntityEvent(ctx, ci, oldEntity, eventType, now)
				if err != nil {
					return res, err
				}
				if event.EventType != "" {
					basicEntityEvents = append(basicEntityEvents, event)
				}
			}
		}
	}

	res.writeModels = writeModels
	res.componentInfos = componentInfos
	res.serviceEvents = serviceEvents
	res.basicEntityEvents = basicEntityEvents
	return res, nil
}

func (w *worker) parseLinks(
	ctx context.Context,
	decoder *json.Decoder,
	componentInfos map[string]map[string]interface{},
) ([]mongo.WriteModel, error) {
	writeModels := make([]mongo.WriteModel, 0)

	for decoder.More() {
		var link Link
		err := decoder.Decode(&link)
		if err != nil {
			return nil, fmt.Errorf("failed to decode links item: %v", err)
		}

		switch link.Action {
		case ActionCreate:
			writeModels = append(writeModels, w.createLink(link)...)

			if infos, ok := componentInfos[link.To]; ok {
				writeModels = append(writeModels, w.updateComponentInfosOnLinkCreate(link, infos))
			} else {
				ci := ConfigurationItem{}
				err := w.entityCollection.
					FindOne(ctx, bson.M{"_id": link.To, "type": types.EntityTypeComponent}).
					Decode(&ci)
				if err == nil {
					writeModels = append(writeModels, w.updateComponentInfosOnLinkCreate(link, ci.Infos))
				} else if !errors.Is(err, mongo.ErrNoDocuments) {
					return nil, err
				}
			}
		case ActionDelete:
			writeModels = append(writeModels, w.deleteLink(link)...)

			if _, ok := componentInfos[link.To]; ok {
				writeModels = append(writeModels, w.updateComponentInfosOnLinkDelete(link))
			} else {
				ci := ConfigurationItem{}
				err := w.entityCollection.
					FindOne(ctx, bson.M{"_id": link.To, "type": types.EntityTypeComponent}).
					Decode(&ci)
				if err == nil {
					writeModels = append(writeModels, w.updateComponentInfosOnLinkDelete(link))
				} else if !errors.Is(err, mongo.ErrNoDocuments) {
					return nil, err
				}
			}
		case ActionUpdate:
			//wasn't implemented in python code
			return nil, ErrNotImplemented
		case ActionEnable:
			//wasn't implemented in python code
			return nil, ErrNotImplemented
		case ActionDisable:
			//wasn't implemented in python code
			return nil, ErrNotImplemented
		default:
			return nil, fmt.Errorf("the action %s is not recognized", link.Action)
		}
	}

	return writeModels, nil
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

		err = w.publisher.SendEvent(types.Event{
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
	if ci.ID == "" {
		return fmt.Errorf("_id is required")
	}

	if ci.Type == nil {
		return fmt.Errorf("type is required")
	}

	switch *ci.Type {
	case types.EntityTypeService:
	case types.EntityTypeResource:
	case types.EntityTypeComponent:
	case types.EntityTypeConnector:
	default:
		return fmt.Errorf("type is not valid %q", *ci.Type)
	}

	if *ci.Type != types.EntityTypeService && ci.EntityPatterns != nil {
		return fmt.Errorf("contains entity patterns, but ci is not a service")
	}

	return nil
}

func (w *worker) fillDefaultFields(ci *ConfigurationItem) {
	if ci.Name == nil {
		ci.Name = &ci.ID
	}

	if ci.ImpactLevel == nil {
		def := new(int64)
		*def = 1

		ci.ImpactLevel = def
	}
}

func (w *worker) createLink(link Link) []mongo.WriteModel {
	return []mongo.WriteModel{
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"_id": link.To}).
			SetUpdate(bson.M{"$addToSet": bson.M{"depends": bson.M{"$each": link.From}}}),
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"_id": bson.M{"$in": link.From}}).
			SetUpdate(bson.M{"$addToSet": bson.M{"impact": link.To}}),
	}
}

func (w *worker) deleteLink(link Link) []mongo.WriteModel {
	return []mongo.WriteModel{
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"_id": link.To}).
			SetUpdate(bson.M{"$pull": bson.M{"depends": bson.M{"$in": link.From}}}),
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"_id": bson.M{"$in": link.From}}).
			SetUpdate(bson.M{"$pull": bson.M{"impact": link.To}}),
	}
}

func (w *worker) createEntity(ci ConfigurationItem) mongo.WriteModel {
	ci.Depends = []string{}
	ci.Impact = []string{}
	ci.EnableHistory = make([]int64, 0)

	if ci.Infos == nil {
		ci.Infos = make(map[string]interface{})
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

	if ci.Infos == nil {
		ci.Infos = make(map[string]interface{})
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

func (w *worker) deleteEntity(ci ConfigurationItem) []mongo.WriteModel {
	return []mongo.WriteModel{
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"impact": ci.ID}).
			SetUpdate(bson.M{"$pull": bson.M{"impact": ci.ID}}),
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"depends": ci.ID}).
			SetUpdate(bson.M{"$pull": bson.M{"depends": ci.ID}}),
		mongo.NewDeleteOneModel().
			SetFilter(bson.M{"_id": ci.ID}),
	}
}

func (w *worker) updateComponentInfosOnComponentUpdate(ci ConfigurationItem) mongo.WriteModel {
	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"type": types.EntityTypeResource, "impact": ci.ID}).
		SetUpdate(bson.M{"$set": bson.M{"component_infos": ci.Infos, "component": ci.ID}})
}

func (w *worker) updateComponentInfosOnComponentDelete(ci ConfigurationItem) mongo.WriteModel {
	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"type": types.EntityTypeResource, "impact": ci.ID}).
		SetUpdate(bson.M{"$unset": bson.M{"component_infos": "", "component": ""}})
}

func (w *worker) updateComponentInfosOnLinkCreate(link Link, infos map[string]interface{}) mongo.WriteModel {
	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"_id": bson.M{"$in": link.From}, "type": types.EntityTypeResource}).
		SetUpdate(bson.M{"$set": bson.M{"component_infos": infos, "component": link.To}})
}

func (w *worker) updateComponentInfosOnLinkDelete(link Link) mongo.WriteModel {
	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"_id": bson.M{"$in": link.From}, "type": types.EntityTypeResource}).
		SetUpdate(bson.M{"$unset": bson.M{"component_infos": "", "component": ""}})
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
		name := ""
		if ci.Name != nil {
			name = *ci.Name
		} else if oldEntity.Name != nil {
			name = *oldEntity.Name
		}
		switch *ci.Type {
		case types.EntityTypeConnector:
			if name == "" {
				return types.Event{}, nil
			}

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
