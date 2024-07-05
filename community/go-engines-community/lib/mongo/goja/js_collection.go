package goja

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/dop251/goja"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type jsCollection struct {
	dbCollection mongo.DbCollection
	dbClient     mongo.DbClient
	vm           *goja.Runtime
}

func (c *jsCollection) CreateIndex(ctx context.Context, orderedKeys, opts, commitQuorum goja.Value) (string, error) {
	if orderedKeys == nil {
		return "", errors.New("no keys provided")
	}

	if t := orderedKeys.ExportType(); t == nil || t.Kind() != reflect.Map {
		return "", errors.New("invalid type for keys")
	}

	obj := orderedKeys.ToObject(c.vm)
	objKeys := obj.Keys()
	if len(objKeys) == 0 {
		return "", errors.New("no keys provided")
	}

	dbKeys := make(bson.D, len(objKeys))
	for i, k := range objKeys {
		dbKeys[i] = bson.E{Key: k, Value: obj.Get(k)}
	}

	dbIndexOpts := options.Index()
	err := transformOptions(c.vm, opts, dbIndexOpts, map[string]mappingFunc{
		"expireAfterSeconds": func(v any) error {
			i, ok := v.(int64)
			if !ok {
				return errors.New("invalid type for expireAfterSeconds")
			}

			dbIndexOpts.SetExpireAfterSeconds(int32(i))

			return nil
		},
		"textIndexVersion": func(v any) error {
			i, ok := v.(int64)
			if !ok {
				return errors.New("invalid type for textIndexVersion")
			}

			dbIndexOpts.SetTextVersion(int32(i))

			return nil
		},
		"2dsphereIndexVersion": func(v any) error {
			i, ok := v.(int64)
			if !ok {
				return errors.New("invalid type for 2dsphereIndexVersion")
			}

			dbIndexOpts.SetSphereVersion(int32(i))

			return nil
		},
		"bits": func(v any) error {
			i, ok := v.(int64)
			if !ok {
				return errors.New("invalid type for bits")
			}

			dbIndexOpts.SetBits(int32(i))

			return nil
		},
		"min": func(v any) error {
			if i, ok := v.(int64); ok {
				dbIndexOpts.SetMin(float64(i))

				return nil
			}

			if f, ok := v.(float64); ok {
				dbIndexOpts.SetMin(f)

				return nil
			}

			return errors.New("invalid type for min")
		},
		"max": func(v any) error {
			if i, ok := v.(int64); ok {
				dbIndexOpts.SetMax(float64(i))

				return nil
			}

			if f, ok := v.(float64); ok {
				dbIndexOpts.SetMax(f)

				return nil
			}

			return errors.New("invalid type for max")
		},
	})
	if err != nil {
		return "", fmt.Errorf("invalid create index options: %w", err)
	}

	dbIndexModel := mongodriver.IndexModel{
		Keys:    dbKeys,
		Options: dbIndexOpts,
	}
	dbOpts := options.CreateIndexes()
	dbOpts.CommitQuorum, err = transformValue(c.vm, commitQuorum)
	if err != nil {
		return "", fmt.Errorf("invalid commit quorum: %w", err)
	}

	name, err := c.dbCollection.Indexes().CreateOne(ctx, dbIndexModel, dbOpts)
	if err != nil {
		return "", fmt.Errorf("error creating index: %w", err)
	}

	return name, nil
}

func (c *jsCollection) DropIndex(ctx context.Context, name string) (map[string]int64, error) {
	res, err := c.dbCollection.Indexes().DropOne(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("error dropping index: %w", err)
	}

	return map[string]int64{
		"nIndexesWas": res.Lookup("nIndexesWas").AsInt64(),
		"ok":          1,
	}, nil
}

func (c *jsCollection) GetIndexes(ctx context.Context) (map[string]any, error) {
	cursor, err := c.dbCollection.Indexes().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing indexes: %w", err)
	}

	return (&jsCursor{
		dbCursor: cursor,
		vm:       c.vm,
	}).getMethods(ctx), nil
}

func (c *jsCollection) DeleteOne(ctx context.Context, filter, opts goja.Value) (map[string]int64, error) {
	dbFilter, err := transformValue(c.vm, filter)
	if err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	if dbFilter == nil {
		dbFilter = bson.M{}
	}

	dbOpts := options.Delete()
	err = transformOptions(c.vm, opts, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("invalid delete options: %w", err)
	}

	res, err := c.dbCollection.DeleteOne(ctx, dbFilter, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("error deleting document: %w", err)
	}

	return map[string]int64{"deletedCount": res}, nil
}

func (c *jsCollection) DeleteMany(ctx context.Context, filter, opts goja.Value) (map[string]int64, error) {
	dbFilter, err := transformValue(c.vm, filter)
	if err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	if dbFilter == nil {
		dbFilter = bson.M{}
	}

	dbOpts := options.Delete()
	err = transformOptions(c.vm, opts, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("invalid delete options: %w", err)
	}

	res, err := c.dbCollection.DeleteMany(ctx, dbFilter, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("error deleting documents: %w", err)
	}

	return map[string]int64{"deletedCount": res}, nil
}

func (c *jsCollection) Drop(ctx context.Context) (bool, error) {
	err := c.dbCollection.Drop(ctx)
	if err != nil {
		return false, fmt.Errorf("error dropping collection: %w", err)
	}

	return true, nil
}

func (c *jsCollection) Find(ctx context.Context, filter, projection, opts goja.Value) (map[string]any, error) {
	dbFilter, err := transformValue(c.vm, filter)
	if err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	if dbFilter == nil {
		dbFilter = bson.M{}
	}

	dbOpts := options.Find()
	dbProjection, err := transformValue(c.vm, projection)
	if err != nil {
		return nil, fmt.Errorf("invalid projection: %w", err)
	}

	dbOpts.SetProjection(dbProjection)
	err = transformOptions(c.vm, opts, dbOpts, map[string]mappingFunc{
		"maxAwaitTimeMS": func(v any) error {
			return setMaxAwaitTimeMS[*options.FindOptions](v, dbOpts)
		},
		"maxTimeMS": func(v any) error {
			return setMaxTimeMS[*options.FindOptions](v, dbOpts)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("invalid find options: %w", err)
	}

	cursor, err := c.dbCollection.Find(ctx, dbFilter, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("error finding documents: %w", err)
	}

	return (&jsCursor{
		dbCursor: cursor,
		vm:       c.vm,
	}).getMethods(ctx), nil
}

func (c *jsCollection) Aggregate(ctx context.Context, pipeline, opts goja.Value) (map[string]any, error) {
	dbOpts := options.Aggregate()
	err := transformOptions(c.vm, opts, dbOpts, map[string]mappingFunc{
		"maxAwaitTimeMS": func(v any) error {
			return setMaxAwaitTimeMS[*options.AggregateOptions](v, dbOpts)
		},
		"maxTimeMS": func(v any) error {
			return setMaxTimeMS[*options.AggregateOptions](v, dbOpts)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("invalid aggregate options: %w", err)
	}

	dbPipeline, err := transformValue(c.vm, pipeline)
	if err != nil {
		return nil, fmt.Errorf("invalid pipeline: %w", err)
	}

	cursor, err := c.dbCollection.Aggregate(ctx, dbPipeline, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("error aggregating documents: %w", err)
	}

	return (&jsCursor{
		dbCursor: cursor,
		vm:       c.vm,
	}).getMethods(ctx), nil
}

func (c *jsCollection) FindOne(ctx context.Context, filter, opts goja.Value) (any, error) {
	dbFilter, err := transformValue(c.vm, filter)
	if err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	if dbFilter == nil {
		dbFilter = bson.M{}
	}

	dbOpts := options.FindOne()
	err = transformOptions(c.vm, opts, dbOpts, map[string]mappingFunc{
		"maxTimeMS": func(v any) error {
			return setMaxTimeMS[*options.FindOneOptions](v, dbOpts)
		},
	})
	if err != nil {
		return 0, fmt.Errorf("invalid find one options: %w", err)
	}

	doc := make(map[string]any)
	err = c.dbCollection.FindOne(ctx, dbFilter, dbOpts).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}

		return nil, fmt.Errorf("error finding document: %w", err)
	}

	return doc, nil
}

func (c *jsCollection) InsertOne(ctx context.Context, doc, opts goja.Value) (map[string]any, error) {
	dbDoc, err := transformValue(c.vm, doc)
	if err != nil {
		return nil, fmt.Errorf("invalid document: %w", err)
	}

	dbOpts := options.InsertOne()
	err = transformOptions(c.vm, opts, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("invalid insert one options: %w", err)
	}

	res, err := c.dbCollection.InsertOne(ctx, dbDoc, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("error inserting document: %w", err)
	}

	return map[string]any{"insertedId": res}, nil
}

func (c *jsCollection) InsertMany(ctx context.Context, docs, opts goja.Value) (map[string]any, error) {
	dbDocs, err := transformValue(c.vm, docs)
	if err != nil {
		return nil, fmt.Errorf("invalid documents: %w", err)
	}

	dbDocsSlice, ok := dbDocs.(bson.A)
	if !ok {
		return nil, errors.New("invalid document type")
	}

	dbOpts := options.InsertMany()
	err = transformOptions(c.vm, opts, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("invalid insert many options: %w", err)
	}

	res, err := c.dbCollection.InsertMany(ctx, dbDocsSlice, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("error inserting documents: %w", err)
	}

	return map[string]any{"insertedIds": res}, nil
}

func (c *jsCollection) UpdateOne(ctx context.Context, filter, update, opts goja.Value) (map[string]any, error) {
	dbFilter, err := transformValue(c.vm, filter)
	if err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	if dbFilter == nil {
		dbFilter = bson.M{}
	}

	dbUpdate, err := transformValue(c.vm, update)
	if err != nil {
		return nil, fmt.Errorf("invalid update: %w", err)
	}

	dbOpts := options.Update()
	err = transformOptions(c.vm, opts, dbOpts, map[string]mappingFunc{
		"arrayFilters": func(v any) error {
			return setArrayFilters[*options.UpdateOptions](v, dbOpts)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("invalid update options: %w", err)
	}

	res, err := c.dbCollection.UpdateOne(ctx, dbFilter, dbUpdate, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("error updating document: %w", err)
	}

	return map[string]any{
		"insertedId":    res.UpsertedID,
		"matchedCount":  res.MatchedCount,
		"modifiedCount": res.ModifiedCount,
		"upsertedCount": res.UpsertedCount,
	}, nil
}

func (c *jsCollection) UpdateMany(ctx context.Context, filter, update, opts goja.Value) (map[string]any, error) {
	dbFilter, err := transformValue(c.vm, filter)
	if err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	if dbFilter == nil {
		dbFilter = bson.M{}
	}

	dbUpdate, err := transformValue(c.vm, update)
	if err != nil {
		return nil, fmt.Errorf("invalid update: %w", err)
	}

	dbOpts := options.Update()
	err = transformOptions(c.vm, opts, dbOpts, map[string]mappingFunc{
		"arrayFilters": func(v any) error {
			return setArrayFilters[*options.UpdateOptions](v, dbOpts)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("invalid update many options: %w", err)
	}

	res, err := c.dbCollection.UpdateMany(ctx, dbFilter, dbUpdate, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("error updating documents: %w", err)
	}

	return map[string]any{
		"insertedId":    res.UpsertedID,
		"matchedCount":  res.MatchedCount,
		"modifiedCount": res.ModifiedCount,
		"upsertedCount": res.UpsertedCount,
	}, nil
}

func (c *jsCollection) CountDocuments(ctx context.Context, filter, opts goja.Value) (int64, error) {
	dbFilter, err := transformValue(c.vm, filter)
	if err != nil {
		return 0, fmt.Errorf("invalid filter: %w", err)
	}

	if dbFilter == nil {
		dbFilter = bson.M{}
	}

	dbOpts := options.Count()
	err = transformOptions(c.vm, opts, dbOpts, map[string]mappingFunc{
		"maxTimeMS": func(v any) error {
			return setMaxTimeMS[*options.CountOptions](v, dbOpts)
		},
	})
	if err != nil {
		return 0, fmt.Errorf("invalid count options: %w", err)
	}

	res, err := c.dbCollection.CountDocuments(ctx, dbFilter, dbOpts)
	if err != nil {
		return 0, fmt.Errorf("error counting documents: %w", err)
	}

	return res, nil
}

func (c *jsCollection) RenameCollection(ctx context.Context, target string, dropTarget bool) (map[string]int64, error) {
	// rename can only be run on the 'admin' database
	err := c.dbClient.RunAdminCommand(ctx, bson.D{
		{Key: "renameCollection", Value: c.dbClient.Name() + "." + c.dbCollection.Name()},
		{Key: "to", Value: c.dbClient.Name() + "." + target},
		{Key: "dropTarget", Value: dropTarget},
	}).Err()
	if err != nil {
		return nil, fmt.Errorf("error renaming collection: %w", err)
	}

	return map[string]int64{"ok": 1}, nil
}

func (c *jsCollection) BulkWrite(ctx context.Context, operations, opts goja.Value) (map[string]any, error) {
	dbOpts := options.BulkWrite()
	err := transformOptions(c.vm, opts, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("invalid bulk write options: %w", err)
	}

	if operations == nil {
		return nil, errors.New("no operations provided")
	}

	if t := operations.ExportType(); t == nil || t.Kind() != reflect.Array && t.Kind() != reflect.Slice {
		return nil, errors.New("invalid type for operations")
	}

	obj := operations.ToObject(c.vm)
	keys := obj.Keys()
	writeModels := make([]mongodriver.WriteModel, len(keys))
	for i, k := range keys {
		v := obj.Get(k)
		if t := v.ExportType(); t == nil || t.Kind() != reflect.Map {
			return nil, errors.New("invalid type for operation at " + strconv.Itoa(i))
		}

		vObj := v.ToObject(c.vm)
		vKeys := vObj.Keys()
		if len(vKeys) != 1 {
			return nil, errors.New("invalid keys length at " + strconv.Itoa(i))
		}

		params := vObj.Get(vKeys[0])
		var mapping map[string]mappingFunc
		switch vKeys[0] {
		case "insertOne":
			writeModels[i] = mongodriver.NewInsertOneModel()
		case "updateOne":
			updateModel := mongodriver.NewUpdateOneModel()
			writeModels[i] = updateModel
			mapping = map[string]mappingFunc{
				"arrayFilters": func(v any) error {
					err := setArrayFilters[*mongodriver.UpdateOneModel](v, updateModel)
					if err != nil {
						return fmt.Errorf("invalid update model at %d: %w", i, err)
					}

					return nil
				},
			}
		case "updateMany":
			updateModel := mongodriver.NewUpdateManyModel()
			writeModels[i] = updateModel
			mapping = map[string]mappingFunc{
				"arrayFilters": func(v any) error {
					err := setArrayFilters[*mongodriver.UpdateManyModel](v, updateModel)
					if err != nil {
						return fmt.Errorf("invalid update model at %d: %w", i, err)
					}

					return nil
				},
			}
		case "deleteOne":
			writeModels[i] = mongodriver.NewDeleteOneModel()
		case "deleteMany":
			writeModels[i] = mongodriver.NewDeleteManyModel()
		case "replaceOne":
			writeModels[i] = mongodriver.NewReplaceOneModel()
		default:
			return nil, errors.New("unknown operation " + vKeys[0] + " at " + strconv.Itoa(i))
		}

		err = transformOptions(c.vm, params, writeModels[i], mapping)
		if err != nil {
			return nil, fmt.Errorf("invalid %s options at %d: %w", vKeys[0], i, err)
		}
	}

	res, err := c.dbCollection.BulkWrite(ctx, writeModels)
	if err != nil {
		return nil, fmt.Errorf("error bulk write documents: %w", err)
	}

	return map[string]any{
		"insertedCount": res.InsertedCount,
		"matchedCount":  res.MatchedCount,
		"modifiedCount": res.ModifiedCount,
		"deletedCount":  res.DeletedCount,
		"upsertedCount": res.UpsertedCount,
		"upsertedIds":   res.UpsertedIDs,
	}, nil
}

func (c *jsCollection) getMethods(ctx context.Context) map[string]any {
	return map[string]any{
		"createIndex": func(orderedKeys, opts, commitQuorum goja.Value) (string, error) {
			return c.CreateIndex(ctx, orderedKeys, opts, commitQuorum)
		},
		"dropIndex": func(name string) (map[string]int64, error) {
			return c.DropIndex(ctx, name)
		},
		"getIndexes": func() (map[string]any, error) {
			return c.GetIndexes(ctx)
		},
		"deleteOne": func(filter, opts goja.Value) (map[string]int64, error) {
			return c.DeleteOne(ctx, filter, opts)
		},
		"deleteMany": func(filter, opts goja.Value) (map[string]int64, error) {
			return c.DeleteMany(ctx, filter, opts)
		},
		"drop": func() (bool, error) {
			return c.Drop(ctx)
		},
		"find": func(filter, projection, opts goja.Value) (map[string]any, error) {
			return c.Find(ctx, filter, projection, opts)
		},
		"aggregate": func(pipeline, opts goja.Value) (map[string]any, error) {
			return c.Aggregate(ctx, pipeline, opts)
		},
		"findOne": func(filter, opts goja.Value) (any, error) {
			return c.FindOne(ctx, filter, opts)
		},
		"insertOne": func(doc, opts goja.Value) (map[string]any, error) {
			return c.InsertOne(ctx, doc, opts)
		},
		"insertMany": func(docs, opts goja.Value) (map[string]any, error) {
			return c.InsertMany(ctx, docs, opts)
		},
		"updateOne": func(filter, update, opts goja.Value) (map[string]any, error) {
			return c.UpdateOne(ctx, filter, update, opts)
		},
		"updateMany": func(filter, update, opts goja.Value) (map[string]any, error) {
			return c.UpdateMany(ctx, filter, update, opts)
		},
		"countDocuments": func(filter, opts goja.Value) (int64, error) {
			return c.CountDocuments(ctx, filter, opts)
		},
		"renameCollection": func(target string, dropTarget bool) (map[string]int64, error) {
			return c.RenameCollection(ctx, target, dropTarget)
		},
		"bulkWrite": func(operations, opts goja.Value) (map[string]any, error) {
			return c.BulkWrite(ctx, operations, opts)
		},
	}
}

type arrayFiltersSetter[T any] interface {
	SetArrayFilters(filters options.ArrayFilters) T
}

func setArrayFilters[T any](v any, setter arrayFiltersSetter[T]) error {
	filters, ok := v.(bson.A)
	if !ok {
		return errors.New("invalid type for arrayFilters")
	}

	setter.SetArrayFilters(options.ArrayFilters{
		Filters: filters,
	})

	return nil
}

type maxAwaitTimeMSSetter[T any] interface {
	SetMaxAwaitTime(d time.Duration) T
}

func setMaxAwaitTimeMS[T any](v any, setter maxAwaitTimeMSSetter[T]) error {
	i, ok := v.(int64)
	if !ok {
		return errors.New("invalid type for maxAwaitTimeMS")
	}

	setter.SetMaxAwaitTime(time.Duration(i) * time.Millisecond)

	return nil
}

type maxTimeMSSetter[T any] interface {
	SetMaxTime(d time.Duration) T
}

func setMaxTimeMS[T any](v any, setter maxTimeMSSetter[T]) error {
	i, ok := v.(int64)
	if !ok {
		return errors.New("invalid type for maxTimeMS")
	}

	setter.SetMaxTime(time.Duration(i) * time.Millisecond)

	return nil
}
