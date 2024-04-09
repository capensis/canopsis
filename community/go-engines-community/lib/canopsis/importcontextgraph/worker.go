package importcontextgraph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitycategory"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type parseResult struct {
	writeModels []mongo.WriteModel

	updatedIds []string
	removedIds []string

	serviceEvents          []types.Event
	resourceEvents         []types.Event
	updatedComponentEvents []types.Event
	existedComponentEvents []types.Event
}

type parseEntityResult struct {
	writeModels []mongo.WriteModel

	updatedIds []string
	removedIds []string

	serviceEvents          []types.Event
	resourceEvents         []types.Event
	updatedComponentEvents []types.Event
	existedComponentEvents []types.Event
}

type worker struct {
	entityCollection        libmongo.DbCollection
	categoryCollection      libmongo.DbCollection
	alarmCollection         libmongo.DbCollection
	alarmResolvedCollection libmongo.DbCollection

	publisher         EventPublisher
	metricMetaUpdater metrics.MetaUpdater

	connector string

	logger zerolog.Logger
}

func NewWorker(
	dbClient libmongo.DbClient,
	publisher EventPublisher,
	metricMetaUpdater metrics.MetaUpdater,
	connector string,
	logger zerolog.Logger,
) Worker {
	return &worker{
		entityCollection:        dbClient.Collection(libmongo.EntityMongoCollection),
		categoryCollection:      dbClient.Collection(libmongo.EntityCategoryMongoCollection),
		alarmCollection:         dbClient.Collection(libmongo.AlarmMongoCollection),
		alarmResolvedCollection: dbClient.Collection(libmongo.ResolvedAlarmMongoCollection),

		publisher:         publisher,
		metricMetaUpdater: metricMetaUpdater,

		connector: connector,

		logger: logger,
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
		for _, event := range res.updatedComponentEvents {
			err = w.publisher.SendEvent(ctx, event)
			if err != nil {
				return stats, err
			}
		}

		for _, event := range res.existedComponentEvents {
			err = w.publisher.SendEvent(ctx, event)
			if err != nil {
				return stats, err
			}
		}
	} else if len(res.serviceEvents)+len(res.resourceEvents) <= int(serviceCount) {
		for _, event := range res.serviceEvents {
			err = w.publisher.SendEvent(ctx, event)
			if err != nil {
				return stats, err
			}
		}

		for _, event := range res.resourceEvents {
			err = w.publisher.SendEvent(ctx, event)
			if err != nil {
				return stats, err
			}
		}

		for _, event := range res.updatedComponentEvents {
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

		for _, event := range res.updatedComponentEvents {
			err = w.publisher.SendEvent(ctx, event)
			if err != nil {
				return stats, err
			}
		}

		for _, event := range res.existedComponentEvents {
			err = w.publisher.SendEvent(ctx, event)
			if err != nil {
				return stats, err
			}
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

	t, err := decoder.Token()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return res, nil
		}

		return res, err
	}

	if t != json.Delim('[') {
		return res, fmt.Errorf("cis should be an array")
	}

	entityParseRes, err = w.parseEntities(ctx, decoder, source, withEvents)
	if err != nil {
		return res, err
	}

	writeModels = append(writeModels, entityParseRes.writeModels...)

	res.writeModels = writeModels
	res.updatedIds = entityParseRes.updatedIds
	res.removedIds = entityParseRes.removedIds
	res.serviceEvents = entityParseRes.serviceEvents
	res.resourceEvents = entityParseRes.resourceEvents
	res.updatedComponentEvents = entityParseRes.updatedComponentEvents
	res.existedComponentEvents = entityParseRes.existedComponentEvents

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

	now := datetime.NewCpsTime()
	componentInfos := make(map[string]map[string]types.Info)

	componentsExist := make(map[string]bool)
	componentsToDelete := make(map[string]bool)
	componentsToDisable := make(map[string]bool)

	deletedResources := make(map[string]bool)
	disabledResources := make(map[string]bool)

	serviceEvents := make([]types.Event, 0)
	resourceEvents := make([]types.Event, 0)
	existedComponentEvents := make([]types.Event, 0)
	updatedComponentEvents := make([]types.Event, 0)

	for decoder.More() {
		var ci EntityConfiguration
		err := decoder.Decode(&ci)
		if err != nil {
			return res, fmt.Errorf("failed to decode cis item: %w", err)
		}

		err = w.validate(ci)
		if err != nil {
			return res, fmt.Errorf("ci = %s, validation error: %s", ci.ID, err.Error())
		}

		if ci.Type == types.EntityTypeService && !match.ValidateEntityPattern(ci.EntityPattern, common.GetForbiddenFieldsInEntityPattern(libmongo.EntityMongoCollection)) {
			w.logger.Warn().Str("entity_name", ci.Name).Msg("invalid entity pattern, skip")
			continue
		}

		if ci.CategoryName != "" {
			var category entitycategory.Category

			err = w.categoryCollection.FindOne(ctx, bson.M{"name": ci.CategoryName}).Decode(&category)
			if err != nil {
				if !errors.Is(err, mongo.ErrNoDocuments) {
					return res, fmt.Errorf("failed to find a category with name = %s: %w", category.Name, err)
				}

				category = entitycategory.Category{
					ID:      utils.NewID(),
					Name:    ci.CategoryName,
					Created: &now,
					Updated: &now,
				}

				_, err = w.categoryCollection.InsertOne(ctx, category)
				if err != nil {
					return res, fmt.Errorf("failed to create a category with name = %s: %w", category.Name, err)
				}
			}

			ci.CategoryID = category.ID
		}

		w.fillDefaultFields(&ci, source, now)

		eventType := ""
		var oldEntity struct {
			EntityConfiguration `bson:",inline"`
			Resources           []string `bson:"resources"`
		}

		findCriteria := bson.M{"soft_deleted": bson.M{"$exists": false}}
		if ci.Type == types.EntityTypeService {
			findCriteria["name"] = ci.Name
		} else {
			findCriteria["_id"] = ci.ID
		}

		cursor, err := w.entityCollection.Aggregate(ctx, []bson.M{
			{"$match": findCriteria},
			{"$graphLookup": bson.M{
				"from":                    libmongo.EntityMongoCollection,
				"startWith":               "$_id",
				"connectFromField":        "_id",
				"connectToField":          "component",
				"as":                      "resources",
				"restrictSearchWithMatch": bson.M{"type": types.EntityTypeResource},
				"maxDepth":                0,
			}},
			{"$addFields": bson.M{
				"resources": bson.M{"$map": bson.M{"input": "$resources", "in": "$$this._id"}},
			}},
		})
		if err != nil {
			return res, err
		}
		if cursor.Next(ctx) {
			err = cursor.Decode(&oldEntity)
			if err != nil {
				_ = cursor.Close(ctx)
				return res, err
			}
		}
		err = cursor.Close(ctx)
		if err != nil {
			return res, err
		}

		switch ci.Action {
		case ActionSet:
			if ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
				componentsExist[ci.ID] = true
			}

			if oldEntity.ID == "" {
				writeModels = append(writeModels, w.createEntity(ci))
				if ci.Type == types.EntityTypeResource {
					if _, ok := componentsExist[ci.Component]; !ok {
						componentsExist[ci.Component] = false
					}
				}

				updatedIds = append(updatedIds, ci.ID)
			} else {
				writeModels = append(writeModels, w.updateEntity(&ci, oldEntity.EntityConfiguration, true))
				if ci.Type == types.EntityTypeResource {
					componentsExist[ci.Component] = true
				}

				updatedIds = append(updatedIds, oldEntity.ID)
			}

			if oldEntity.Enabled || ci.Enabled {
				switch ci.Type {
				case types.EntityTypeService:
					eventType = types.EventTypeRecomputeEntityService
				default:
					if oldEntity.Enabled && ci.Enabled {
						eventType = types.EventTypeEntityUpdated
					} else {
						eventType = types.EventTypeEntityToggled
					}
				}
			}
		case ActionDelete:
			if oldEntity.ID == "" {
				if ci.Type == types.EntityTypeService {
					err = fmt.Errorf("failed to delete an entity service with name = %s", ci.Name)
				} else {
					err = fmt.Errorf("failed to delete an entity with _id = %s", ci.ID)
				}

				return res, err
			}

			if ci.Type == types.EntityTypeResource && !deletedResources[ci.ID] {
				deletedResources[ci.ID] = true
			} else if ci.Type == types.EntityTypeComponent {
				componentsToDelete[ci.ID] = true

				for _, resourceID := range oldEntity.Resources {
					if !deletedResources[resourceID] {
						deletedResources[resourceID] = true

						writeModels = append(writeModels, w.deleteEntity(resourceID, now)...)
					}
				}
			} else if ci.Type == types.EntityTypeService {
				eventType = types.EventTypeRecomputeEntityService
			}

			writeModels = append(writeModels, w.deleteEntity(oldEntity.ID, now)...)
			removedIds = append(removedIds, oldEntity.ID)
		case ActionEnable:
			if oldEntity.ID == "" {
				if ci.Type == types.EntityTypeService {
					err = fmt.Errorf("failed to enable an entity service with name = %s", ci.Name)
				} else {
					err = fmt.Errorf("failed to enable an entity with _id = %s", ci.ID)
				}

				return res, err
			}

			switch ci.Type {
			case types.EntityTypeService:
				eventType = types.EventTypeRecomputeEntityService
			default:
				if ci.Type == types.EntityTypeComponent {
					componentInfos[ci.ID] = ci.Infos
				}

				eventType = types.EventTypeEntityToggled
			}

			writeModels = append(writeModels, w.changeState(oldEntity.ID, true, source, now))
			updatedIds = append(updatedIds, oldEntity.ID)
		case ActionDisable:
			if oldEntity.ID == "" {
				if ci.Type == types.EntityTypeService {
					err = fmt.Errorf("failed to disable an entity service with name = %s", ci.Name)
				} else {
					err = fmt.Errorf("failed to disable an entity with _id = %s", ci.ID)
				}

				return res, err
			}

			if ci.Type == types.EntityTypeResource && !deletedResources[ci.ID] {
				disabledResources[ci.ID] = true
			} else if ci.Type == types.EntityTypeComponent {
				componentsToDisable[ci.ID] = true
				componentInfos[ci.ID] = ci.Infos

				for _, resourceID := range oldEntity.Resources {
					if !deletedResources[resourceID] {
						deletedResources[resourceID] = true
						writeModels = append(writeModels, w.changeState(resourceID, false, source, now))
						updatedIds = append(updatedIds, resourceID)
					}
				}
			}

			switch ci.Type {
			case types.EntityTypeService:
				eventType = types.EventTypeRecomputeEntityService
			default:
				eventType = types.EventTypeEntityToggled
			}

			writeModels = append(writeModels, w.changeState(ci.ID, false, source, now))
			updatedIds = append(updatedIds, oldEntity.ID)
		default:
			return res, fmt.Errorf("the action %s is not recognized", ci.Action)
		}

		if withEvents && eventType != "" {
			switch ci.Type {
			case types.EntityTypeService:
				serviceEvents = append(serviceEvents, w.createServiceEvent(oldEntity.EntityConfiguration, eventType, now))
			case types.EntityTypeResource:
				resourceEvents = append(resourceEvents, w.createResourceEvent(eventType, ci.Name, ci.Component, now))
			case types.EntityTypeComponent:
				updatedComponentEvents = append(updatedComponentEvents, w.createComponentEvent(eventType, ci.ID, now))
			}
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
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": componentName, "soft_deleted": bson.M{"$exists": false}}).Decode(&oldEntity)
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				return res, err
			}

			if errors.Is(err, mongo.ErrNoDocuments) {
				ci := EntityConfiguration{
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
				updatedComponentEvents = append(updatedComponentEvents, w.createComponentEvent(types.EventTypeEntityUpdated, componentName, now))
			} else {
				if !oldEntity.Enabled {
					return res, fmt.Errorf("can't create resource for disabled component")
				}

				componentInfos[componentName] = oldEntity.Infos
				existedComponentEvents = append(existedComponentEvents, w.createComponentEvent(types.EventTypeEntityUpdated, componentName, now))
			}
		} else {
			existedComponentEvents = append(existedComponentEvents, w.createComponentEvent(types.EventTypeEntityUpdated, componentName, now))
		}

		if len(componentInfos[componentName]) > 0 {
			writeModels = append(writeModels, w.updateComponentInfos(componentName, componentInfos[componentName]))
		}
	}

	res.updatedIds = updatedIds
	res.removedIds = removedIds
	res.writeModels = writeModels
	res.serviceEvents = serviceEvents
	res.resourceEvents = resourceEvents
	res.updatedComponentEvents = updatedComponentEvents
	res.existedComponentEvents = existedComponentEvents

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
			Connector:     w.connector,
			ConnectorName: w.connector,
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

func (w *worker) validate(ci EntityConfiguration) error {
	switch ci.Type {
	case types.EntityTypeService:
		if len(ci.EntityPattern) == 0 && ci.Action == ActionSet {
			return fmt.Errorf("service %s contains empty pattern", ci.Name)
		}
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

func (w *worker) fillDefaultFields(ci *EntityConfiguration, source string, now datetime.CpsTime) {
	switch ci.Type {
	case types.EntityTypeService:
		ci.ID = ci.Name
	case types.EntityTypeResource:
		ci.ID = ci.Name + "/" + ci.Component
	case types.EntityTypeComponent:
		ci.ID = ci.Name
	}

	if ci.ImpactLevel == 0 {
		ci.ImpactLevel = 1
	}

	ci.ImportSource = source
	ci.Imported = now
}

func (w *worker) createEntity(ci EntityConfiguration) mongo.WriteModel {
	ci.Services = []string{}
	ci.EnableHistory = make([]int64, 0)

	if ci.Type == types.EntityTypeComponent {
		ci.Component = ci.ID
	}

	if ci.Infos == nil {
		ci.Infos = make(map[string]types.Info)
	}

	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"_id": ci.ID}).
		SetUpdate(bson.M{
			"$set":         ci,
			"$setOnInsert": bson.M{"created": datetime.NewCpsTime()},
			"$unset":       bson.M{"soft_deleted": ""},
		}).
		SetUpsert(true)
}

func (w *worker) updateEntity(ci *EntityConfiguration, oldEntity EntityConfiguration, mergeInfos bool) mongo.WriteModel {
	ci.EnableHistory = oldEntity.EnableHistory

	if ci.Type == types.EntityTypeComponent {
		ci.Component = ci.ID
	}

	if ci.Infos == nil {
		ci.Infos = make(map[string]types.Info)
	}

	if oldEntity.Infos == nil {
		oldEntity.Infos = make(map[string]types.Info)
	}

	if mergeInfos {
		for k, v := range ci.Infos {
			oldEntity.Infos[k] = v
		}

		ci.Infos = oldEntity.Infos
	}

	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"_id": oldEntity.ID}).
		SetUpdate(bson.M{
			"$set":         ci,
			"$setOnInsert": bson.M{"created": datetime.NewCpsTime()},
			"$unset":       bson.M{"soft_deleted": ""},
		}).
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

func (w *worker) deleteEntity(id string, now datetime.CpsTime) []mongo.WriteModel {
	return []mongo.WriteModel{
		mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": id}).
			SetUpdate(bson.M{"$set": bson.M{"enabled": false, "soft_deleted": now}}),
	}
}

func (w *worker) updateComponentInfos(componentID string, infos map[string]types.Info) mongo.WriteModel {
	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"type": types.EntityTypeResource, "component": componentID}).
		SetUpdate(bson.M{"$set": bson.M{"component_infos": infos}})
}

func (w *worker) createServiceEvent(ci EntityConfiguration, eventType string, now datetime.CpsTime) types.Event {
	return types.Event{
		EventType:     eventType,
		Timestamp:     now,
		Author:        canopsis.DefaultEventAuthor,
		Connector:     w.connector,
		ConnectorName: w.connector,
		Component:     ci.ID,
		SourceType:    types.SourceTypeService,
		Initiator:     types.InitiatorSystem,
	}
}

func (w *worker) createResourceEvent(eventType, name, component string, now datetime.CpsTime) types.Event {
	return types.Event{
		Connector:     w.connector,
		ConnectorName: w.connector,
		EventType:     eventType,
		Timestamp:     now,
		Author:        canopsis.DefaultEventAuthor,
		Initiator:     types.InitiatorSystem,
		Resource:      name,
		Component:     component,
		SourceType:    types.SourceTypeResource,
	}
}

func (w *worker) createComponentEvent(eventType, name string, now datetime.CpsTime) types.Event {
	return types.Event{
		Connector:     w.connector,
		ConnectorName: w.connector,
		EventType:     eventType,
		Timestamp:     now,
		Author:        canopsis.DefaultEventAuthor,
		Initiator:     types.InitiatorSystem,
		Component:     name,
		SourceType:    types.SourceTypeComponent,
	}
}
