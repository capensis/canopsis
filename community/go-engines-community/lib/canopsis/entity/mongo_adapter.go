package entity

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoAdapter ...
type mongoAdapter struct {
	dbClient     mongo.DbClient
	dbCollection mongo.DbCollection

	bulk []mongodriver.WriteModel
}

// NewAdapter gives the correct entity adapter. Give nil to the redis client
// and it will create a new redis.Client with the dedicated redis database for entities.
func NewAdapter(dbClient mongo.DbClient) Adapter {
	return &mongoAdapter{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(mongo.EntityMongoCollection),
	}
}

// Insert add a new entity.
func (a *mongoAdapter) Insert(entity types.Entity) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := a.dbCollection.InsertOne(ctx, entity)
	return err
}

// Update updates an existing entity or creates a new one in db.
func (a *mongoAdapter) Update(entity types.Entity) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := a.dbCollection.ReplaceOne(ctx, bson.M{"_id": entity.ID}, entity)
	return err
}

func (a *mongoAdapter) Bulk(ctx context.Context, models []mongodriver.WriteModel) error {
	_, err := a.dbCollection.BulkWrite(ctx, models)

	return err
}

func (a *mongoAdapter) Remove(entity types.Entity) error {
	panic("not implemented")
}

// BulkInsert insert entities in bulk.
// Thread safe.
func (a *mongoAdapter) BulkInsert(entity types.Entity) error {
	operation := mongodriver.NewInsertOneModel()
	operation.SetDocument(entity)

	a.bulk = append(a.bulk, operation)

	return nil
}

func (a *mongoAdapter) AddToBulkUpdate(entityId string, data interface{}) error {
	operation := mongodriver.NewUpdateOneModel()
	operation.SetFilter(bson.M{"_id": entityId})
	operation.SetUpdate(data)

	a.bulk = append(a.bulk, operation)

	return nil
}

func (a *mongoAdapter) BulkUpsert(entity types.Entity, newImpacts []string, newDepends []string) error {
	operation := mongodriver.NewUpdateOneModel()
	operation.SetFilter(bson.M{"_id": entity.ID})
	operation.SetUpdate(entity.GetUpsertMongoBson(newImpacts, newDepends))
	operation.SetUpsert(true)

	a.bulk = append(a.bulk, operation)

	return nil
}

// FlushBulk force all bulk caches to be written.
func (a *mongoAdapter) FlushBulk() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(a.bulk) > 0 {
		_, err := a.dbCollection.BulkWrite(ctx, a.bulk)
		if err != nil {
			return fmt.Errorf("entity adapter flushbulk update: %v", err)
		}

		a.bulk = nil
	}

	return nil
}

func (a *mongoAdapter) FlushBulkInsert() error {
	panic("not implemented")
}

func (a *mongoAdapter) FlushBulkUpdate() error {
	panic("not implemented")
}

// Get is the same as GetEntityByID
// Return True if the document has been found
func (a *mongoAdapter) Get(id string) (types.Entity, bool) {
	entity, err := a.GetEntityByID(id)
	entity.EnsureInitialized()

	if err == mongodriver.ErrNoDocuments {
		return entity, false
	} else if err != nil {
		return entity, false
	}

	return entity, true
}

// GetEntityByID finds an Entity from is eid
func (a *mongoAdapter) GetEntityByID(id string) (types.Entity, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var ent types.Entity

	res := a.dbCollection.FindOne(ctx, bson.M{"_id": id})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return ent, errt.NewNotFound(err)
		}

		return ent, err
	}

	err := res.Decode(&ent)

	return ent, err
}

func (a *mongoAdapter) Count() (int, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := a.dbCollection.CountDocuments(ctx, bson.M{})
	return int(res), err
}

func (a *mongoAdapter) GetIDs(filter map[string]interface{}, ids *[]interface{}) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cursor, err := a.dbCollection.Find(ctx, filter)
	if err != nil {
		return err
	}

	err = cursor.All(ctx, ids)
	if err != nil {
		return err
	}

	return cursor.Close(ctx)
}

func (a *mongoAdapter) RemoveAll() error {
	panic("not implemented")
}

func (a *mongoAdapter) UpsertMany(entities []types.Entity) (map[string]bool, error) {
	if len(entities) == 0 {
		return nil, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	insertModels := make([]mongodriver.WriteModel, 0)
	for _, entity := range entities {
		insert := bson.M{
			"_id":             entity.ID,
			"name":            entity.Name,
			"measurements":    entity.Measurements,
			"enabled":         entity.Enabled,
			"type":            entity.Type,
			"enable_history":  entity.EnableHistory,
			"category":        entity.Category,
			"impact_level":    entity.ImpactLevel,
			"impact":          entity.Impacts,
			"depends":         entity.Depends,
			"infos":           entity.Infos,
			"created":         entity.Created,
			"last_event_date": entity.LastEventDate,
		}
		if entity.Component != "" {
			insert["component"] = entity.Component
		}

		insertModels = append(insertModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": entity.ID}).
			SetUpdate(bson.M{"$setOnInsert": insert}).
			SetUpsert(true))
	}
	res, err := a.dbCollection.BulkWrite(ctx, insertModels)
	if err != nil {
		return nil, err
	}

	upsertedIDs := make(map[string]bool, len(res.UpsertedIDs))
	for _, v := range res.UpsertedIDs {
		upsertedIDs[v.(string)] = true
	}
	// Update only enabled entities.
	updateModels := make([]mongodriver.WriteModel, 0)
	for _, entity := range entities {
		if upsertedIDs[entity.ID] {
			continue
		}

		upsertedIDs[entity.ID] = false
		set := bson.M{}
		if entity.LastEventDate != nil {
			set["last_event_date"] = entity.LastEventDate
		}

		if len(entity.Infos) > 0 {
			for k, v := range entity.Infos {
				key := "infos." + k
				set[key] = v
			}
		}

		update := bson.M{
			"$addToSet": bson.M{
				"impact":  bson.M{"$each": entity.Impacts},
				"depends": bson.M{"$each": entity.Depends},
			},
		}
		if len(set) > 0 {
			update["$set"] = set
		}
		updateModels = append(updateModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": entity.ID, "enabled": true}).
			SetUpdate(update))
	}

	if len(updateModels) > 0 {
		_, err = a.dbCollection.BulkWrite(ctx, updateModels)
		if err != nil {
			return nil, err
		}
	}

	return upsertedIDs, nil
}

func (a *mongoAdapter) AddImpacts(ids []string, impacts []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	writeModels := make([]mongodriver.WriteModel, len(ids))
	for i, id := range ids {
		writeModels[i] = mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": id}).
			SetUpdate(bson.M{
				"$push": bson.M{
					"impact": bson.M{"$each": impacts},
				},
			})
	}

	_, err := a.dbCollection.BulkWrite(ctx, writeModels)

	return err
}

func (a *mongoAdapter) RemoveImpacts(ids []string, impacts []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	writeModels := make([]mongodriver.WriteModel, len(ids))
	for i, id := range ids {
		writeModels[i] = mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": id}).
			SetUpdate(bson.M{
				"$pullAll": bson.M{
					"impact": impacts,
				},
			})
	}

	_, err := a.dbCollection.BulkWrite(ctx, writeModels)

	return err
}

func (a *mongoAdapter) AddImpactByQuery(query interface{}, impact string) ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := a.dbCollection.Find(
		ctx,
		bson.M{"$and": []interface{}{
			bson.M{"enabled": true},
			query,
			bson.M{"impact": bson.M{"$nin": bson.A{impact}}},
		}},
		options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}

	newEntities := make([]types.Entity, 0)
	err = res.All(ctx, &newEntities)
	if err != nil {
		return nil, err
	}

	newIDs := make([]string, len(newEntities))
	for i := range newEntities {
		newIDs[i] = newEntities[i].ID
	}

	if len(newIDs) > 0 {
		_, err = a.dbCollection.UpdateMany(
			ctx,
			bson.M{"_id": bson.M{"$in": newIDs}},
			bson.M{"$addToSet": bson.M{"impact": impact}},
		)
		if err != nil {
			return nil, err
		}
	}

	return newIDs, nil
}

func (a *mongoAdapter) RemoveImpactByQuery(query interface{}, impact string) ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := a.dbCollection.Find(
		ctx,
		bson.M{"$and": []interface{}{
			bson.M{"enabled": true},
			query,
			bson.M{"impact": bson.M{"$in": bson.A{impact}}},
		}},
		options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}

	removedEntities := make([]types.Entity, 0)
	err = res.All(ctx, &removedEntities)
	if err != nil {
		return nil, err
	}

	removedIDs := make([]string, len(removedEntities))
	for i := range removedEntities {
		removedIDs[i] = removedEntities[i].ID
	}

	if len(removedIDs) > 0 {
		_, err = a.dbCollection.UpdateMany(
			ctx,
			bson.M{"_id": bson.M{"$in": removedIDs}},
			bson.M{"$pull": bson.M{"impact": impact}},
		)
		if err != nil {
			return nil, err
		}
	}

	return removedIDs, nil
}

func (a *mongoAdapter) AddInfos(id string, infos map[string]types.Info) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	set := bson.M{}
	for k, v := range infos {
		set["infos."+k] = v
	}

	res, err := a.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": set,
	})

	if err != nil {
		return false, err
	}

	return res.ModifiedCount > 0, nil
}

func (a *mongoAdapter) UpdateComponentInfos(id, componentID string) (map[string]types.Info, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res := a.dbCollection.FindOne(
		ctx,
		bson.M{
			"_id":     componentID,
			"type":    types.EntityTypeComponent,
			"enabled": true,
		},
		options.FindOne().SetProjection(bson.M{"infos": 1}),
	)
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	component := types.Entity{}
	err := res.Decode(&component)
	if err != nil {
		return nil, err
	}

	updateRes, err := a.dbCollection.UpdateOne(
		ctx,
		bson.M{"_id": id, "type": types.EntityTypeResource},
		bson.M{"$set": bson.M{"component_infos": component.Infos}},
	)
	if err != nil || updateRes.ModifiedCount == 0 {
		return nil, err
	}

	return component.Infos, nil
}

func (a *mongoAdapter) UpdateComponentInfosByComponent(componentID string) ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cursor, err := a.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"_id":     componentID,
			"type":    types.EntityTypeComponent,
			"enabled": true,
		}},
		{"$graphLookup": bson.M{
			"from":                    mongo.EntityMongoCollection,
			"startWith":               "$depends",
			"connectFromField":        "depends",
			"connectToField":          "_id",
			"restrictSearchWithMatch": bson.M{"type": types.EntityTypeResource},
			"as":                      "resources",
		}},
		{"$project": bson.M{
			"infos":   1,
			"depends": 1,
			"resources": bson.M{"$map": bson.M{
				"input": "$resources",
				"as":    "each",
				"in":    "$$each._id",
			}},
		}},
	})
	if err != nil {
		return nil, err
	}

	if !cursor.Next(ctx) {
		return nil, nil
	}

	component := struct {
		Infos     interface{} `bson:"infos"`
		Resources []string    `bson:"resources"`
	}{}
	err = cursor.Decode(&component)
	if err != nil {
		return nil, err
	}

	resUpdate, err := a.dbCollection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": component.Resources}},
		bson.M{"$set": bson.M{"component_infos": component.Infos}},
	)

	if err != nil {
		return nil, err
	}

	if resUpdate.ModifiedCount > 0 {
		return component.Resources, nil
	}

	return nil, nil
}

func (a *mongoAdapter) UpdateLastEventDate(ids []string, time types.CpsTime) error {
	if len(ids) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := a.dbCollection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": ids}},
		bson.M{"$set": bson.M{"last_event_date": time}},
	)

	return err
}

func (a *mongoAdapter) UpdateIdleFields(ctx context.Context, id string,
	idleSince *types.CpsTime, lastIdleRuleApply string) error {
	set := bson.M{}
	unset := bson.M{}

	if idleSince == nil {
		unset["idle_since"] = ""
	} else {
		set["idle_since"] = idleSince
	}

	if lastIdleRuleApply == "" {
		unset["last_idle_rule_apply"] = ""
	} else {
		set["last_idle_rule_apply"] = lastIdleRuleApply
	}

	update := bson.M{}
	if len(set) > 0 {
		update["$set"] = set
	}
	if len(unset) > 0 {
		update["$unset"] = unset
	}

	_, err := a.dbCollection.UpdateMany(ctx, bson.M{"_id": id}, update)

	return err
}

func (a *mongoAdapter) FindByIDs(ids []string) ([]types.Entity, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cursor, err := a.dbCollection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	entities := make([]types.Entity, 0)
	err = cursor.All(ctx, &entities)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (a *mongoAdapter) GetAllWithLastUpdateDateBefore(
	ctx context.Context,
	time types.CpsTime,
	exclude []string,
) (mongo.Cursor, error) {
	return a.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"_id":                  bson.M{"$nin": exclude},
			"enabled":              true,
			"last_idle_rule_apply": nil,
			"type": bson.M{"$in": []string{
				types.EntityTypeConnector,
				types.EntityTypeComponent,
				types.EntityTypeResource,
			}},
			"$or": []bson.M{
				{"$and": []bson.M{
					{"last_event_date": bson.M{"$ne": nil}},
					{"last_event_date": bson.M{"$lte": time}},
				}},
				{"$and": []bson.M{
					{"last_event_date": nil},
					{"created": bson.M{"$lte": time}},
				}},
				{"$and": []bson.M{
					{"last_event_date": nil},
					{"created": nil},
				}},
			},
		}},
	})
}

func (a *mongoAdapter) FindConnectorForComponent(ctx context.Context, id string) (*types.Entity, error) {
	cursor, err := a.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"_id": id,
		}},
		{"$graphLookup": bson.M{
			"from":                    mongo.EntityMongoCollection,
			"startWith":               "$impact",
			"connectFromField":        "impact",
			"connectToField":          "_id",
			"as":                      "connector",
			"restrictSearchWithMatch": bson.M{"type": types.EntityTypeConnector},
			"maxDepth":                0,
		}},
		{"$unwind": "$connector"},
		{"$project": bson.M{
			"connector": 1,
		}},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		res := &struct {
			Connector types.Entity `bson:"connector"`
		}{}
		err := cursor.Decode(res)
		if err != nil {
			return nil, err
		}

		return &res.Connector, nil
	}

	return nil, nil
}

func (a *mongoAdapter) FindConnectorForResource(ctx context.Context, id string) (*types.Entity, error) {
	cursor, err := a.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"_id": id,
		}},
		{"$graphLookup": bson.M{
			"from":                    mongo.EntityMongoCollection,
			"startWith":               "$depends",
			"connectFromField":        "depends",
			"connectToField":          "_id",
			"as":                      "connector",
			"restrictSearchWithMatch": bson.M{"type": types.EntityTypeConnector},
			"maxDepth":                0,
		}},
		{"$unwind": "$connector"},
		{"$project": bson.M{
			"connector": 1,
		}},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		res := &struct {
			ID        string       `bson:"_id"`
			Connector types.Entity `bson:"connector"`
		}{}
		err := cursor.Decode(res)
		if err != nil {
			return nil, err
		}

		return &res.Connector, nil
	}

	return nil, nil
}

func (a *mongoAdapter) FindComponentForResource(ctx context.Context, id string) (*types.Entity, error) {
	cursor, err := a.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"_id": id,
		}},
		{"$graphLookup": bson.M{
			"from":                    mongo.EntityMongoCollection,
			"startWith":               "$impact",
			"connectFromField":        "impact",
			"connectToField":          "_id",
			"as":                      "component",
			"restrictSearchWithMatch": bson.M{"type": types.EntityTypeComponent},
			"maxDepth":                0,
		}},
		{"$unwind": "$component"},
		{"$project": bson.M{
			"component": 1,
		}},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		res := &struct {
			ID        string       `bson:"_id"`
			Component types.Entity `bson:"component"`
		}{}
		err := cursor.Decode(res)
		if err != nil {
			return nil, err
		}

		return &res.Component, nil
	}

	return nil, nil
}

func (a *mongoAdapter) GetWithIdleSince(ctx context.Context) (mongo.Cursor, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"idle_since", 1}})

	return a.dbCollection.Find(
		ctx,
		bson.M{
			"idle_since": bson.M{"$gt": 0},
			"type":       bson.M{"$in": []string{types.EntityTypeResource, types.EntityTypeComponent, types.EntityTypeConnector}},
			"enabled":    true,
		},
		findOptions,
	)
}

func (a *mongoAdapter) GetImpactedServicesInfo(ctx context.Context) (mongo.Cursor, error) {
	return a.dbCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"type": types.EntityTypeConnector,
			},
		},
		{
			"$addFields": bson.M{
				"dependencies": bson.M{
					"$concatArrays": bson.A{"$depends", "$impact"},
				},
			},
		},
		{
			"$project": bson.M{
				"dependencies": 1,
			},
		},
		{
			"$lookup": bson.M{
				"from": mongo.EntityMongoCollection,
				"let":  bson.M{"dependencies": "$dependencies"},
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"type": types.EntityTypeService,
						},
					},
					{
						"$addFields": bson.M{
							"inter": bson.M{
								"$setIntersection": bson.A{"$depends", "$$dependencies"},
							},
						},
					},
					{
						"$match": bson.M{
							"inter": bson.M{
								"$ne": bson.A{},
							},
						},
					},
				},
				"as": "impacted_services",
			},
		},
		{
			"$project": bson.M{
				"_id": 1,
				"impacted_services": bson.M{
					"$map": bson.M{
						"input": "$impacted_services",
						"as":    "item",
						"in":    "$$item._id",
					},
				},
			},
		},
	})
}
