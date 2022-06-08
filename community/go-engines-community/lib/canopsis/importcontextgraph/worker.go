package importcontextgraph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type worker struct {
	entityCollection   libmongo.DbCollection
	categoryCollection libmongo.DbCollection
	publisher          EventPublisher
	metricMetaUpdater  metrics.MetaUpdater
}

func NewWorker(
	dbClient libmongo.DbClient,
	publisher EventPublisher,
	metricMetaUpdater metrics.MetaUpdater,
) Worker {
	return &worker{
		entityCollection:   dbClient.Collection(libmongo.EntityMongoCollection),
		categoryCollection: dbClient.Collection(libmongo.EntityCategoryMongoCollection),
		publisher:          publisher,
		metricMetaUpdater:  metricMetaUpdater,
	}
}

func (w *worker) Work(ctx context.Context, filename, source string) (stats Stats, resErr error) {
	startTime := time.Now()
	defer func() {
		stats.ExecTime = time.Since(startTime)
	}()

	writeModels, err := w.parseFile(ctx, filename, source)
	if err != nil {
		return stats, err
	}

	if len(writeModels) == 0 {
		return stats, fmt.Errorf("empty import")
	}

	stats.Updated, stats.Deleted, err = w.bulkWrite(ctx, writeModels, canopsis.DefaultBulkSize, canopsis.DefaultBulkBytesSize)
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

func (w *worker) parseFile(ctx context.Context, filename, source string) (writeModels []mongo.WriteModel, resErr error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
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

	writeModels = make([]mongo.WriteModel, 0)
	decoder := json.NewDecoder(file)
	var componentInfos map[string]map[string]interface{}

	for {
		t, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		if t == "cis" {
			t, err := decoder.Token()
			if err != nil {
				return nil, fmt.Errorf("failed to parse cis: %v", err)
			}

			if t != json.Delim('[') {
				return nil, fmt.Errorf("cis should be an array")
			}

			var entityWriteModels []mongo.WriteModel
			entityWriteModels, componentInfos, err = w.parseEntities(ctx, decoder, source)
			if err != nil {
				return nil, err
			}

			writeModels = append(writeModels, entityWriteModels...)
		}

		if t == "links" {
			t, err := decoder.Token()
			if err != nil {
				return nil, fmt.Errorf("failed to parse links: %v", err)
			}

			if t != json.Delim('[') {
				return nil, fmt.Errorf("links should be an array")
			}

			linkWriteModels, err := w.parseLinks(ctx, decoder, componentInfos)
			if err != nil {
				return nil, err
			}

			writeModels = append(writeModels, linkWriteModels...)
		}
	}

	return writeModels, nil
}

func (w *worker) parseEntities(
	ctx context.Context,
	decoder *json.Decoder,
	source string,
) ([]mongo.WriteModel, map[string]map[string]interface{}, error) {
	writeModels := make([]mongo.WriteModel, 0)
	now := types.CpsTime{Time: time.Now()}
	componentInfos := make(map[string]map[string]interface{})

	for decoder.More() {
		var ci ConfigurationItem
		err := decoder.Decode(&ci)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decode cis item: %v", err)
		}

		w.fillDefaultFields(&ci)
		err = w.validate(ci)
		if err != nil {
			return nil, nil, fmt.Errorf("ci = %s, validation error: %s", ci.ID, err.Error())
		}

		ci.ImportSource = source
		ci.Imported = now

		switch ci.Action {
		case ActionCreate:
			writeModels = append(writeModels, w.createEntity(ci))
			if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
			}
		case ActionSet:
			var oldEntity ConfigurationItem
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil && err != mongo.ErrNoDocuments {
				return nil, nil, err
			}

			if err == mongo.ErrNoDocuments {
				writeModels = append(writeModels, w.createEntity(ci))

				break
			}

			writeModels = append(writeModels, w.updateEntity(&ci, oldEntity, true))
			if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
				writeModels = append(writeModels, w.updateComponentInfosOnComponentUpdate(ci))
			}
		case ActionUpdate:
			var oldEntity ConfigurationItem
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return nil, nil, fmt.Errorf("failed to update an entity with _id = %s", ci.ID)
				}

				return nil, nil, err
			}

			writeModels = append(writeModels, w.updateEntity(&ci, oldEntity, false))
			if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
				writeModels = append(writeModels, w.updateComponentInfosOnComponentUpdate(ci))
			}
		case ActionDelete:
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Err()
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return nil, nil, fmt.Errorf("failed to delete an entity with _id = %s", ci.ID)
				}

				return nil, nil, err
			}

			writeModels = append(writeModels, w.deleteEntity(ci)...)
			if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = nil
				writeModels = append(writeModels, w.updateComponentInfosOnComponentDelete(ci))
			}
		case ActionEnable:
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Err()
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return nil, nil, fmt.Errorf("failed to enable an entity with _id = %s", ci.ID)
				}

				return nil, nil, err
			}

			writeModels = append(writeModels, w.changeState(ci.ID, true, source, now))
			if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
			}
		case ActionDisable:
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Err()
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return nil, nil, fmt.Errorf("failed to disable an entity with _id = %s", ci.ID)
				}

				return nil, nil, err
			}

			writeModels = append(writeModels, w.changeState(ci.ID, false, source, now))
			if ci.Type != nil && *ci.Type == types.EntityTypeComponent {
				componentInfos[ci.ID] = ci.Infos
			}
		default:
			return nil, nil, fmt.Errorf("the action %s is not recognized", ci.Action)
		}
	}

	return writeModels, componentInfos, nil
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

		err = w.publisher.SendUpdateEntityServiceEvent(service.ID)
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

	if *ci.Type != types.EntityTypeService && len(ci.EntityPattern) > 0 {
		return fmt.Errorf("contains entity pattern, but ci is not a service")
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
