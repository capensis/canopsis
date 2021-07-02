package importcontextgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type worker struct {
	entityCollection   libmongo.DbCollection
	categoryCollection libmongo.DbCollection
	publisher          EventPublisher
}

func NewWorker(
	dbClient libmongo.DbClient,
	publisher EventPublisher,
) Worker {
	return &worker{
		entityCollection:   dbClient.Collection(libmongo.EntityMongoCollection),
		categoryCollection: dbClient.Collection(libmongo.EntityCategoryMongoCollection),
		publisher:          publisher,
	}
}

func (w *worker) Work(ctx context.Context, filename string) (stats Stats, resErr error) {
	startTime := time.Now()
	defer func() {
		stats.ExecTime = time.Since(startTime)
	}()

	writeModels, err := w.parseFile(ctx, filename)
	if err != nil {
		return stats, err
	}

	if len(writeModels) == 0 {
		return stats, fmt.Errorf("empty import")
	}

	stats.Updated, stats.Deleted, err = w.bulkWrite(ctx, writeModels, canopsis.DefaultBulkSize)
	if err != nil {
		return stats, err
	}

	err = w.sendUpdateServiceEvents(ctx)
	if err != nil {
		return stats, err
	}

	return stats, nil
}

func (w *worker) parseFile(ctx context.Context, filename string) (writeModels []mongo.WriteModel, resErr error) {
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

			entityWriteModels, err := w.parseEntities(ctx, decoder)
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

			linkWriteModels, err := w.parseLinks(decoder)
			if err != nil {
				return nil, err
			}

			writeModels = append(writeModels, linkWriteModels...)
		}
	}

	return writeModels, nil
}

func (w *worker) parseEntities(ctx context.Context, decoder *json.Decoder) ([]mongo.WriteModel, error) {
	writeModels := make([]mongo.WriteModel, 0)

	for decoder.More() {
		var ci ConfigurationItem
		err := decoder.Decode(&ci)
		if err != nil {
			return nil, fmt.Errorf("failed to decode cis item: %v", err)
		}

		w.fillDefaultFields(&ci)
		err = w.validate(ci)
		if err != nil {
			return nil, fmt.Errorf("ci = %s, validation error: %s", ci.ID, err.Error())
		}

		switch ci.Action {
		case ActionCreate:
			writeModels = append(writeModels, w.createEntity(ci))
		case ActionSet:
			var oldEntity ConfigurationItem
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil && err != mongo.ErrNoDocuments {
				return nil, err
			}

			if err == mongo.ErrNoDocuments {
				writeModels = append(writeModels, w.createEntity(ci))

				break
			}

			writeModels = append(writeModels, w.updateEntity(ci, oldEntity, true))
		case ActionUpdate:
			var oldEntity ConfigurationItem
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Decode(&oldEntity)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return nil, fmt.Errorf("failed to update an entity with _id = %s", ci.ID)
				}

				return nil, err
			}

			writeModels = append(writeModels, w.updateEntity(ci, oldEntity, false))
		case ActionDelete:
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Err()
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return nil, fmt.Errorf("failed to delete an entity with _id = %s", ci.ID)
				}

				return nil, err
			}

			writeModels = append(writeModels, w.deleteEntity(ci)...)
		case ActionEnable:
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Err()
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return nil, fmt.Errorf("failed to enable an entity with _id = %s", ci.ID)
				}

				return nil, err
			}

			writeModels = append(writeModels, w.changeState(ci.ID, true))
		case ActionDisable:
			err := w.entityCollection.FindOne(ctx, bson.M{"_id": ci.ID}).Err()
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return nil, fmt.Errorf("failed to disable an entity with _id = %s", ci.ID)
				}

				return nil, err
			}

			writeModels = append(writeModels, w.changeState(ci.ID, false))
		default:
			return nil, fmt.Errorf("the action %s is not recognized", ci.Action)
		}
	}

	return writeModels, nil
}

func (w *worker) parseLinks(decoder *json.Decoder) ([]mongo.WriteModel, error) {
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
		case ActionDelete:
			writeModels = append(writeModels, w.deleteLink(link)...)
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

func (w *worker) bulkWrite(ctx context.Context, writeModels []mongo.WriteModel, limit int) (int64, int64, error) {
	var updated, deleted int64

	for i := 0; i < len(writeModels); i += limit {
		end := i + limit
		if i+limit > len(writeModels) {
			end = len(writeModels)
		}

		p := writeModels[i:end]
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
	ci.EnableHistory = make([]string, 0)

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

func (w *worker) updateEntity(ci ConfigurationItem, oldEntity ConfigurationItem, mergeInfos bool) mongo.WriteModel {
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

func (w *worker) changeState(id string, enabled bool) mongo.WriteModel {
	return mongo.NewUpdateManyModel().
		SetFilter(bson.M{"_id": id}).
		SetUpdate(bson.M{"$set": bson.M{"enabled": enabled}})
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
