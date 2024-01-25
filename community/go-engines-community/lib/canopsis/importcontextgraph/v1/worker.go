package v1

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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	publisher         importcontextgraph.EventPublisher
	metricMetaUpdater metrics.MetaUpdater
}

type parseResult struct {
	writeModels []mongo.WriteModel

	updatedIds        []string
	removedIds        []string
	serviceEvents     []types.Event
	basicEntityEvents []types.Event
}

type parseEntityResult struct {
	writeModels []mongo.WriteModel

	componentInfos map[string]map[string]interface{}

	entityTypes map[string]string

	updatedIds        []string
	removedIds        []string
	serviceEvents     []types.Event
	basicEntityEvents []types.Event
}

// left for backward compatibility
func NewWorker(
	dbClient libmongo.DbClient,
	publisher importcontextgraph.EventPublisher,
	metricMetaUpdater metrics.MetaUpdater,
) importcontextgraph.Worker {
	return &worker{
		entityCollection:        dbClient.Collection(libmongo.EntityMongoCollection),
		categoryCollection:      dbClient.Collection(libmongo.EntityCategoryMongoCollection),
		alarmCollection:         dbClient.Collection(libmongo.AlarmMongoCollection),
		alarmResolvedCollection: dbClient.Collection(libmongo.ResolvedAlarmMongoCollection),

		publisher:         publisher,
		metricMetaUpdater: metricMetaUpdater,
	}
}

func (w *worker) Work(ctx context.Context, filename, source string) (stats importcontextgraph.Stats, resErr error) {
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

func (w *worker) WorkPartial(ctx context.Context, filename, source string) (stats importcontextgraph.Stats, resErr error) {
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
			err = w.publisher.SendEvent(ctx, event)
			if err != nil {
				return stats, err
			}
		}
	} else if len(res.serviceEvents)+len(res.basicEntityEvents) <= int(serviceCount) {
		for _, event := range res.serviceEvents {
			err = w.publisher.SendEvent(ctx, event)
			if err != nil {
				return stats, err
			}
		}
		for _, event := range res.basicEntityEvents {
			err = w.publisher.SendEvent(ctx, event)
			if err != nil {
				return stats, err
			}
		}
	} else {
		err = w.sendUpdateServiceEvents(ctx)
		if err != nil {
			return stats, err
		}
	}

	if len(res.updatedIds) > 0 {
		w.metricMetaUpdater.UpdateById(ctx, res.updatedIds...)
	}
	if len(res.removedIds) > 0 {
		w.metricMetaUpdater.DeleteById(ctx, res.removedIds...)
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
			if errors.Is(err, io.EOF) {
				break
			}

			return res, err
		}

		if t == "cis" {
			t, err := decoder.Token()
			if err != nil {
				return res, fmt.Errorf("failed to parse cis: %w", err)
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
				return res, fmt.Errorf("failed to parse links: %w", err)
			}

			if t != json.Delim('[') {
				return res, fmt.Errorf("links should be an array")
			}

			linkWriteModels, err := w.parseLinks(ctx, decoder, entityParseRes.componentInfos, entityParseRes.entityTypes)
			if err != nil {
				return res, err
			}

			writeModels = append(writeModels, linkWriteModels...)
		}
	}

	res.writeModels = writeModels
	res.updatedIds = entityParseRes.updatedIds
	res.removedIds = entityParseRes.removedIds
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
	updatedIds := make([]string, 0)
	removedIds := make([]string, 0)
	serviceEvents := make([]types.Event, 0)
	basicEntityEvents := make([]types.Event, 0)
	entityTypes := make(map[string]string)
	now := datetime.NewCpsTime()
	componentInfos := make(map[string]map[string]interface{})

	for decoder.More() {
		var ci importcontextgraph.ConfigurationItem
		err := decoder.Decode(&ci)
		if err != nil {
			return res, fmt.Errorf("failed to decode cis item: %w", err)
		}

		w.fillDefaultFields(&ci)
		err = w.validate(ci)
		if err != nil {
			return res, fmt.Errorf("ci = %s, validation error: %s", ci.ID, err.Error())
		}

		ci.ImportSource = source
		ci.Imported = now

		eventType := ""
		var oldEntity importcontextgraph.ConfigurationItem

		if ci.Type != nil {
			entityTypes[ci.ID] = *ci.Type
		}

		switch ci.Action {
		case importcontextgraph.ActionCreate:
			updatedIds = append(updatedIds, ci.ID)
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
		case importcontextgraph.ActionSet:
			updatedIds = append(updatedIds, ci.ID)
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				return res, err
			}

			if errors.Is(err, mongo.ErrNoDocuments) {
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
		case importcontextgraph.ActionUpdate:
			updatedIds = append(updatedIds, ci.ID)
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
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
		case importcontextgraph.ActionDelete:
			removedIds = append(removedIds, ci.ID)
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
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
		case importcontextgraph.ActionEnable:
			updatedIds = append(updatedIds, ci.ID)
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
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
		case importcontextgraph.ActionDisable:
			updatedIds = append(updatedIds, ci.ID)
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
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
				event := w.createBasicEntityEvent(ci, oldEntity, eventType, now)
				if event.EventType != "" {
					basicEntityEvents = append(basicEntityEvents, event)
				}
			}
		}
	}

	res.updatedIds = updatedIds
	res.removedIds = removedIds
	res.writeModels = writeModels
	res.componentInfos = componentInfos
	res.serviceEvents = serviceEvents
	res.basicEntityEvents = basicEntityEvents
	res.entityTypes = entityTypes
	return res, nil
}

func (w *worker) parseLinks(
	ctx context.Context,
	decoder *json.Decoder,
	componentInfos map[string]map[string]interface{},
	entityTypes map[string]string,
) ([]mongo.WriteModel, error) {
	writeModels := make([]mongo.WriteModel, 0)
	if entityTypes == nil {
		entityTypes = make(map[string]string)
	}

	for decoder.More() {
		var link importcontextgraph.Link
		err := decoder.Decode(&link)
		if err != nil {
			return nil, fmt.Errorf("failed to decode links item: %w", err)
		}

		ciTo := importcontextgraph.ConfigurationItem{}
		err = w.entityCollection.FindOne(ctx, bson.M{"_id": link.To}).Decode(&ciTo)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}

		if _, ok := entityTypes[link.To]; !ok && ciTo.Type != nil {
			entityTypes[link.To] = *ciTo.Type
		}
		for _, from := range link.From {
			fromType := entityTypes[from]
			if fromType == "" {
				ciFrom := importcontextgraph.ConfigurationItem{}
				err = w.entityCollection.FindOne(ctx, bson.M{"_id": from}).Decode(&ciFrom)
				if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
					return nil, err
				}
				if ciFrom.Type != nil {
					entityTypes[fromType] = *ciFrom.Type
				}
			}
		}

		switch link.Action {
		case importcontextgraph.ActionCreate:
			linkWriterModels, err := w.createLink(link, entityTypes)
			if err != nil {
				return nil, err
			}
			writeModels = append(writeModels, linkWriterModels...)

			if infos, ok := componentInfos[link.To]; ok {
				writeModels = append(writeModels, w.updateComponentInfosOnLinkCreate(link, infos))
			} else if ciTo.Type != nil && *ciTo.Type == types.EntityTypeComponent {
				writeModels = append(writeModels, w.updateComponentInfosOnLinkCreate(link, ciTo.Infos))
			}
		case importcontextgraph.ActionDelete:
			writeModels = append(writeModels, w.deleteLink(link, entityTypes)...)

			if _, ok := componentInfos[link.To]; ok {
				writeModels = append(writeModels, w.updateComponentInfosOnLinkDelete(link))
			} else if ciTo.Type != nil && *ciTo.Type == types.EntityTypeComponent {
				writeModels = append(writeModels, w.updateComponentInfosOnLinkDelete(link))
			}
		case importcontextgraph.ActionUpdate:
			//wasn't implemented in python code
			return nil, importcontextgraph.ErrNotImplemented
		case importcontextgraph.ActionEnable:
			//wasn't implemented in python code
			return nil, importcontextgraph.ErrNotImplemented
		case importcontextgraph.ActionDisable:
			//wasn't implemented in python code
			return nil, importcontextgraph.ErrNotImplemented
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

		err = w.publisher.SendEvent(ctx, types.Event{
			EventType:     types.EventTypeRecomputeEntityService,
			Connector:     defaultConnector,
			ConnectorName: defaultConnectorName,
			Component:     service.ID,
			Timestamp:     datetime.NewCpsTime(),
			Author:        canopsis.DefaultEventAuthor,
			Initiator:     types.InitiatorSystem,
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

func (w *worker) validate(ci importcontextgraph.ConfigurationItem) error {
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

	if *ci.Type != types.EntityTypeService && len(ci.EntityPattern) > 0 {
		return fmt.Errorf("contains entity pattern, but ci is not a service")
	}

	return nil
}

func (w *worker) fillDefaultFields(ci *importcontextgraph.ConfigurationItem) {
	if ci.Name == nil {
		ci.Name = &ci.ID
	}

	if ci.ImpactLevel == nil {
		def := new(int64)
		*def = 1

		ci.ImpactLevel = def
	}
}

func (w *worker) createLink(link importcontextgraph.Link, entityTypes map[string]string) ([]mongo.WriteModel, error) {
	if link.To == "" || len(link.From) == 0 {
		return nil, nil
	}

	switch entityTypes[link.To] {
	case types.EntityTypeConnector:
		return []mongo.WriteModel{
			mongo.NewUpdateManyModel().
				SetFilter(bson.M{
					"_id":  bson.M{"$in": link.From},
					"type": bson.M{"$in": bson.A{types.EntityTypeComponent, types.EntityTypeResource}},
				}).
				SetUpdate(bson.M{"$set": bson.M{"connector": link.To}}),
		}, nil
	case types.EntityTypeComponent:
		return []mongo.WriteModel{
			mongo.NewUpdateManyModel().
				SetFilter(bson.M{
					"_id":  bson.M{"$in": link.From},
					"type": types.EntityTypeResource,
				}).
				SetUpdate(bson.M{"$set": bson.M{"component": link.To}}),
		}, nil
	case types.EntityTypeResource:
		if len(link.From) > 1 {
			return nil, fmt.Errorf("an entity cannot be connected to more then 1 connector")
		}
		if entityTypes[link.From[0]] != types.EntityTypeConnector {
			return nil, fmt.Errorf("a resource can be connected only to connector")
		}

		return []mongo.WriteModel{
			mongo.NewUpdateManyModel().
				SetFilter(bson.M{
					"_id":  link.To,
					"type": types.EntityTypeResource,
				}).
				SetUpdate(bson.M{"$set": bson.M{"connector": link.From[0]}}),
		}, nil
	}

	return nil, nil
}

func (w *worker) deleteLink(link importcontextgraph.Link, entityTypes map[string]string) []mongo.WriteModel {
	if link.To == "" || len(link.From) == 0 {
		return nil
	}

	switch entityTypes[link.To] {
	case types.EntityTypeConnector:
		return []mongo.WriteModel{
			mongo.NewUpdateManyModel().
				SetFilter(bson.M{
					"_id":       bson.M{"$in": link.From},
					"type":      bson.M{"$in": bson.A{types.EntityTypeComponent, types.EntityTypeResource}},
					"connector": link.To,
				}).
				SetUpdate(bson.M{"$unset": bson.M{"connector": ""}}),
		}
	case types.EntityTypeComponent:
		return []mongo.WriteModel{
			mongo.NewUpdateManyModel().
				SetFilter(bson.M{
					"_id":       bson.M{"$in": link.From},
					"type":      types.EntityTypeResource,
					"component": link.To,
				}).
				SetUpdate(bson.M{"$unset": bson.M{"component": ""}}),
		}
	case types.EntityTypeResource:
		if entityTypes[link.From[0]] != types.EntityTypeConnector {
			return nil
		}

		return []mongo.WriteModel{
			mongo.NewUpdateManyModel().
				SetFilter(bson.M{
					"_id":       link.To,
					"type":      types.EntityTypeResource,
					"connector": link.From[0],
				}).
				SetUpdate(bson.M{"$unset": bson.M{"connector": ""}}),
		}
	}

	return nil
}

func (w *worker) createEntity(ci importcontextgraph.ConfigurationItem) mongo.WriteModel {
	ci.Services = []string{}
	ci.EnableHistory = make([]int64, 0)

	if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
		ci.Component = ci.ID
	}

	if ci.Infos == nil {
		ci.Infos = make(map[string]interface{})
	}

	if ci.Measurements == nil {
		ci.Measurements = make([]interface{}, 0)
	}

	now := datetime.NewCpsTime()

	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"_id": ci.ID}).
		SetUpdate(bson.M{"$set": ci, "$setOnInsert": bson.M{"created": now}}).
		SetUpsert(true)
}

func (w *worker) updateEntity(ci *importcontextgraph.ConfigurationItem, oldEntity importcontextgraph.ConfigurationItem, mergeInfos bool) mongo.WriteModel {
	ci.EnableHistory = oldEntity.EnableHistory

	if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
		ci.Component = ci.ID
	}

	if ci.Infos == nil {
		ci.Infos = make(map[string]interface{})
	}

	if oldEntity.Infos == nil {
		oldEntity.Infos = make(map[string]interface{})
	}

	if mergeInfos {
		for k, v := range ci.Infos {
			oldEntity.Infos[k] = v
		}

		ci.Infos = oldEntity.Infos
	}

	now := datetime.NewCpsTime()

	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"_id": ci.ID}).
		SetUpdate(bson.M{"$set": ci, "$setOnInsert": bson.M{"created": now}}).
		SetUpsert(true)
}

func (w *worker) changeState(id string, enabled bool, importSource string, imported datetime.CpsTime) mongo.WriteModel {
	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"_id": id}).
		SetUpdate(bson.M{"$set": bson.M{
			"enabled":       enabled,
			"import_source": importSource,
			"imported":      imported,
		}})
}

func (w *worker) deleteEntity(ci importcontextgraph.ConfigurationItem) []mongo.WriteModel {
	return []mongo.WriteModel{
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"connector": ci.ID}).
			SetUpdate(bson.M{"$unset": bson.M{"connector": ""}}),
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"component": ci.ID}).
			SetUpdate(bson.M{"$unset": bson.M{"component": ""}}),
		mongo.NewUpdateManyModel().
			SetFilter(bson.M{"services": ci.ID}).
			SetUpdate(bson.M{"$pull": bson.M{"services": ci.ID}}),
		mongo.NewDeleteOneModel().
			SetFilter(bson.M{"_id": ci.ID}),
	}
}

func (w *worker) updateComponentInfosOnComponentUpdate(ci importcontextgraph.ConfigurationItem) mongo.WriteModel {
	update := bson.M{"component": ci.ID}
	if len(ci.Infos) > 0 {
		update["component_infos"] = ci.Infos
	}

	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"type": types.EntityTypeResource, "component": ci.ID}).
		SetUpdate(bson.M{"$set": update})
}

func (w *worker) updateComponentInfosOnComponentDelete(ci importcontextgraph.ConfigurationItem) mongo.WriteModel {
	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"type": types.EntityTypeResource, "component": ci.ID}).
		SetUpdate(bson.M{"$unset": bson.M{"component_infos": "", "component": ""}})
}

func (w *worker) updateComponentInfosOnLinkCreate(link importcontextgraph.Link, infos map[string]interface{}) mongo.WriteModel {
	update := bson.M{"component": link.To}
	if len(infos) > 0 {
		update["component_infos"] = infos
	}

	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"_id": bson.M{"$in": link.From}, "type": types.EntityTypeResource}).
		SetUpdate(bson.M{"$set": update})
}

func (w *worker) updateComponentInfosOnLinkDelete(link importcontextgraph.Link) mongo.WriteModel {
	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"_id": bson.M{"$in": link.From}, "type": types.EntityTypeResource}).
		SetUpdate(bson.M{"$unset": bson.M{"component_infos": "", "component": ""}})
}

func (w *worker) createServiceEvent(ci importcontextgraph.ConfigurationItem, eventType string, now datetime.CpsTime) types.Event {
	return types.Event{
		EventType:     eventType,
		Timestamp:     now,
		Author:        canopsis.DefaultEventAuthor,
		Initiator:     types.InitiatorSystem,
		Connector:     defaultConnector,
		ConnectorName: defaultConnectorName,
		Component:     ci.ID,
		SourceType:    types.SourceTypeService,
	}
}

func (w *worker) createBasicEntityEvent(ci, oldEntity importcontextgraph.ConfigurationItem, eventType string, now datetime.CpsTime) types.Event {
	event := types.Event{
		EventType:     eventType,
		Timestamp:     now,
		Author:        canopsis.DefaultEventAuthor,
		Initiator:     types.InitiatorSystem,
		Connector:     defaultConnector,
		ConnectorName: defaultConnectorName,
	}

	name := ""
	if ci.Name != nil {
		name = *ci.Name
	} else if oldEntity.Name != nil {
		name = *oldEntity.Name
	}
	switch *ci.Type {
	case types.EntityTypeConnector:
		if name == "" {
			return types.Event{}
		}

		event.Connector = strings.TrimSuffix(ci.ID, "/"+name)
		event.ConnectorName = name
		event.SourceType = types.SourceTypeConnector
	case types.EntityTypeComponent:
		event.Component = ci.ID
		event.SourceType = types.SourceTypeComponent
	case types.EntityTypeResource:
		event.Resource = name
		event.Component = oldEntity.Component
		event.SourceType = types.SourceTypeResource
	}

	return event
}
