package che

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	minInfoLength   = 2
	buildingTimeout = 5 * time.Minute
	// stopListId is ID of document with stop list of keys;
	// keys with more than InfosDictionaryLimit unique values are added to stop list
	// this need to prevent dictionary from growing with values which are not frequently used
	stopListId = "stop_list"
	// batchSize is an approximate number of entity infos to process in one iteration
	batchSize = 100000
	// entitiesLimit is 1/10 of batchSize to have about 10 infos per entity
	entitiesLimit = batchSize / 10
)

// A composite id is used, because it works faster with a lot of bulk upserts instead of filter and uuid
type infosDictID struct {
	Key   string `bson:"k"`
	Value string `bson:"v"`
}

type infosDictDoc struct {
	ID    infosDictID `bson:"_id"`
	EntID string      `bson:"ent_id"`
}

func NewInfosDictionaryPeriodicalWorker(
	client mongo.DbClient,
	periodicalInterval time.Duration,
	logger zerolog.Logger,
) engine.PeriodicalWorker {
	return &infosDictionaryPeriodicalWorker{
		entityCollection:          client.Collection(mongo.EntityMongoCollection),
		entityInfosDictCollection: client.Collection(mongo.EntityInfosDictionaryCollection),
		configCollection:          client.Collection(mongo.ConfigurationMongoCollection),
		periodicalInterval:        periodicalInterval,
		logger:                    logger,
	}
}

type infosDictionaryPeriodicalWorker struct {
	entityCollection          mongo.DbCollection
	entityInfosDictCollection mongo.DbCollection
	configCollection          mongo.DbCollection
	periodicalInterval        time.Duration
	logger                    zerolog.Logger
}

func (w *infosDictionaryPeriodicalWorker) GetInterval() time.Duration {
	return w.periodicalInterval
}

func (w *infosDictionaryPeriodicalWorker) Work(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, buildingTimeout)
	defer cancel()

	conf := config.CanopsisConf{}

	err := w.configCollection.FindOne(ctx, bson.M{"_id": config.ConfigKeyName}).Decode(&conf)
	if err != nil {
		w.logger.Error().Err(err).Msg("unable to get canopsis config")
		return
	}

	if !conf.Global.BuildEntityInfosDictionary {
		return
	}

	now := types.NewCpsTime()

	err = w.buildDictionary(ctx, now, conf.Global.InfosDictionaryLimit)
	if err != nil {
		w.logger.Err(err).Msg("failed to build entity infos dictionary")
		return
	}

	_, err = w.entityInfosDictCollection.DeleteMany(ctx, bson.M{"last_update": bson.M{"$lt": now}})
	if err != nil {
		w.logger.Err(err).Msg("unable to delete outdated entity infos dictionary documents")
		return
	}
}

type stopValuesRes struct {
	Values []string `bson:"stop_values"`
	Limit  int      `bson:"limit"`
}

func (w *infosDictionaryPeriodicalWorker) buildDictionary(ctx context.Context, t types.CpsTime, limit int) (err error) {
	stopValues := stopValuesRes{}
	err = w.entityInfosDictCollection.FindOne(ctx, bson.M{"_id": stopListId}).Decode(&stopValues)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return fmt.Errorf("unable to get stop list: %w", err)
	}
	stopList := stopValues.Values
	if stopList == nil || stopValues.Limit != limit {
		stopList = make([]string, 0)
	}
	stopListLen := len(stopList)

	defer func(ctx context.Context) {
		if len(stopList) != stopListLen {
			_, errUpd := w.entityInfosDictCollection.UpdateOne(
				ctx, bson.M{"_id": stopListId},
				bson.M{"$set": bson.M{"stop_values": stopList, "limit": limit}},
				options.Update().SetUpsert(true))
			if errUpd != nil && err == nil {
				err = fmt.Errorf("unable to update stop list: %w", errUpd)
			}
		}
	}(context.WithoutCancel(ctx))

	lastEntityID := ""
	writeModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)
	infoDictDocs := make([]infosDictDoc, 0, batchSize)

	// key, "" -> count of unique key-value pairs;
	// non-empty string value as 2nd key in map to skip already processed key-value pairs
	keysCounts := make(map[string]map[string]int)

	for {
		pipeline := []bson.M{
			{"$match": bson.M{"_id": bson.M{"$gt": lastEntityID}}},
			{"$limit": entitiesLimit},
			{"$addFields": bson.M{"Infos": bson.M{"$objectToArray": "$infos"}}},
			{"$project": bson.M{"_id": 1, "Infos": 1}},
			{"$sort": bson.M{"_id": 1}},
			{"$unwind": "$Infos"},
			{"$unwind": "$Infos.v.value"},
			{"$project": bson.M{"k": "$Infos.k", "v": "$Infos.v.value"}},

			{"$match": bson.M{"k": bson.M{"$nin": stopList}}},
			{"$project": bson.M{"_id": bson.M{"k": "$k", "v": bson.M{
				"$cond": bson.M{"if": bson.M{"$and": []bson.M{
					{"$eq": bson.A{bson.M{"$type": "$v"}, "string"}},
					{"$gt": bson.A{bson.M{"$strLenCP": "$v"}, minInfoLength}},
				}}, "then": "$v", "else": ""}}}, "ent_id": "$_id"}},
		}
		entCursor, err := w.entityCollection.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		infoDictDocs = infoDictDocs[:0]
		if err := entCursor.All(ctx, &infoDictDocs); err != nil {
			return fmt.Errorf("unable to decode entity infos data: %w", err)
		}

		if len(infoDictDocs) == 0 {
			break
		}

		writeModels = writeModels[:0]
		bulkBytesSize := 0
		modelsOrdered := false

		for i := range infoDictDocs {
			key, value := infoDictDocs[i].ID.Key, infoDictDocs[i].ID.Value
			if _, ok := slices.BinarySearch(stopList, key); ok {
				continue
			}

			if _, ok := keysCounts[key]; !ok {
				keysCounts[key] = make(map[string]int)
			}
			if _, ok := keysCounts[key][value]; ok {
				continue
			}
			keysCounts[key][""]++
			if value != "" {
				keysCounts[key][value] = 0
			}

			var newModel mongodriver.WriteModel
			if limit > 0 && keysCounts[key][""] > limit {
				pos := sort.SearchStrings(stopList, key)
				stopList = slices.Insert(stopList, pos, key)

				newModel = mongodriver.
					NewDeleteManyModel().
					SetFilter(bson.M{"_id.k": key})
				modelsOrdered = true
				delete(keysCounts, key)
			} else {
				newModel = getUpsertOneModel(infoDictDocs[i].ID, t)
			}
			writeModels, bulkBytesSize, err = w.appendWriteModel(ctx, newModel, writeModels, bulkBytesSize, modelsOrdered)
			if err != nil {
				return err
			}
		}

		if len(writeModels) > 0 {
			_, err = w.entityInfosDictCollection.BulkWrite(ctx, writeModels,
				options.BulkWrite().SetOrdered(modelsOrdered))
			if err != nil {
				return fmt.Errorf("unable to bulk write entity infos dictionary: %w", err)
			}
		}
		lastEntityID = infoDictDocs[len(infoDictDocs)-1].EntID
	}

	if len(stopList) != stopListLen {
		writeModels = writeModels[:0]
		bulkBytesSize := 0

		for _, key := range stopList {
			newModel := getUpsertOneModel(infosDictID{Key: key, Value: ""}, t)

			writeModels, bulkBytesSize, err = w.appendWriteModel(ctx, newModel, writeModels, bulkBytesSize, false)
			if err != nil {
				return err
			}
		}

		if len(writeModels) > 0 {
			_, err = w.entityInfosDictCollection.BulkWrite(ctx, writeModels,
				options.BulkWrite().SetOrdered(false))
			if err != nil {
				return fmt.Errorf("unable to bulk write entity infos dictionary: %w", err)
			}
		}
	}
	return nil
}

func (w *infosDictionaryPeriodicalWorker) appendWriteModel(ctx context.Context, newModel mongodriver.WriteModel, writeModels []mongodriver.WriteModel, bulkBytesSize int, isModelsOrdered bool) ([]mongodriver.WriteModel, int, error) {
	b, err := bson.Marshal(newModel)
	if err != nil {
		return writeModels, bulkBytesSize, err
	}

	newModelLen := len(b)
	if bulkBytesSize+newModelLen > canopsis.DefaultBulkBytesSize ||
		len(writeModels) == canopsis.DefaultBulkSize {
		_, err = w.entityInfosDictCollection.BulkWrite(ctx, writeModels,
			options.BulkWrite().SetOrdered(isModelsOrdered))
		if err != nil {
			return writeModels, bulkBytesSize, fmt.Errorf("unable to bulk write entity infos dictionary: %w", err)
		}

		writeModels = writeModels[:0]
		bulkBytesSize = 0
	}

	bulkBytesSize += newModelLen
	writeModels = append(writeModels, newModel)

	return writeModels, bulkBytesSize, nil
}

func getUpsertOneModel(id infosDictID, t types.CpsTime) mongodriver.WriteModel {
	return mongodriver.
		NewUpdateOneModel().
		SetFilter(bson.M{"_id": id}).
		SetUpdate(bson.M{"$set": bson.M{"last_update": t}}).
		SetUpsert(true)
}
