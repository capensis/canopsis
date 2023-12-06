package entity

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("entity is not found")

// mongoAdapter ...
type mongoAdapter struct {
	dbCollection mongo.DbCollection
}

// NewAdapter gives the correct entity adapter. Give nil to the redis client
// and it will create a new redis.Client with the dedicated redis database for entities.
func NewAdapter(dbClient mongo.DbClient) Adapter {
	return &mongoAdapter{
		dbCollection: dbClient.Collection(mongo.EntityMongoCollection),
	}
}

// Insert add a new entity.
func (a *mongoAdapter) Insert(ctx context.Context, entity types.Entity) error {
	_, err := a.dbCollection.InsertOne(ctx, entity)
	return err
}

// Update updates an existing entity or creates a new one in db.
func (a *mongoAdapter) Update(ctx context.Context, entity types.Entity) error {
	_, err := a.dbCollection.ReplaceOne(ctx, bson.M{"_id": entity.ID}, entity)
	return err
}

func (a *mongoAdapter) Bulk(ctx context.Context, models []mongodriver.WriteModel) error {
	_, err := a.dbCollection.BulkWrite(ctx, models)

	return err
}

// Get is the same as GetEntityByID
// Return True if the document has been found
func (a *mongoAdapter) Get(ctx context.Context, id string) (types.Entity, bool) {
	entity, err := a.GetEntityByID(ctx, id)
	entity.EnsureInitialized()

	if errors.Is(err, ErrNotFound) {
		return entity, false
	} else if err != nil {
		return entity, false
	}

	return entity, true
}

// GetEntityByID finds an Entity from is eid
func (a *mongoAdapter) GetEntityByID(ctx context.Context, id string) (types.Entity, error) {
	var ent types.Entity

	res := a.dbCollection.FindOne(ctx, bson.M{"_id": id})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return ent, ErrNotFound
		}

		return ent, err
	}

	return ent, res.Decode(&ent)
}

func (a *mongoAdapter) Count(ctx context.Context) (int, error) {
	res, err := a.dbCollection.CountDocuments(ctx, bson.M{})
	return int(res), err
}

func (a *mongoAdapter) GetIDs(ctx context.Context, filter map[string]interface{}, ids *[]interface{}) error {
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

func (a *mongoAdapter) UpsertMany(ctx context.Context, entities []types.Entity) (map[string]bool, error) {
	if len(entities) == 0 {
		return nil, nil
	}

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
			"infos":           entity.Infos,
			"created":         entity.Created,
			"last_event_date": entity.LastEventDate,
		}
		if entity.Healthcheck {
			insert["healthcheck"] = true
		}
		if entity.Connector != "" {
			insert["connector"] = entity.Connector
		}
		if entity.Component != "" {
			insert["component"] = entity.Component
		}

		insertModels = append(insertModels, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": entity.ID}).
			SetUpdate(bson.M{"$setOnInsert": insert, "$unset": bson.M{"soft_deleted": ""}}).
			SetUpsert(true))
	}
	res, err := a.dbCollection.BulkWrite(ctx, insertModels)
	if err != nil {
		return nil, err
	}

	upsertedIDs := make(map[string]bool, len(res.UpsertedIDs))
	for _, v := range res.UpsertedIDs {
		if s, ok := v.(string); ok {
			upsertedIDs[s] = true
		}
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

		if len(set) > 0 {
			updateModels = append(updateModels, mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": entity.ID, "enabled": true}).
				SetUpdate(bson.M{"$set": set}))
		}
	}

	if len(updateModels) > 0 {
		_, err = a.dbCollection.BulkWrite(ctx, updateModels)
		if err != nil {
			return nil, err
		}
	}

	return upsertedIDs, nil
}

func (a *mongoAdapter) AddInfos(ctx context.Context, id string, infos map[string]types.Info) (bool, error) {
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

func (a *mongoAdapter) UpdateComponentInfos(ctx context.Context, id, componentID string) (map[string]types.Info, error) {
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
		if errors.Is(err, mongodriver.ErrNoDocuments) {
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

func (a *mongoAdapter) UpdateComponentInfosByComponent(ctx context.Context, componentID string) ([]string, error) {
	cursor, err := a.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"_id":     componentID,
			"type":    types.EntityTypeComponent,
			"enabled": true,
		}},
		{"$graphLookup": bson.M{
			"from":                    mongo.EntityMongoCollection,
			"startWith":               "$_id",
			"connectFromField":        "_id",
			"connectToField":          "component",
			"restrictSearchWithMatch": bson.M{"type": types.EntityTypeResource},
			"as":                      "resources",
		}},
		{"$project": bson.M{
			"infos": 1,
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

func (a *mongoAdapter) UpdateLastEventDate(ctx context.Context, ids []string, time datetime.CpsTime) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := a.dbCollection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": ids}},
		bson.M{"$set": bson.M{"last_event_date": time}},
	)

	return err
}

func (a *mongoAdapter) UpdateIdleFields(ctx context.Context, id string,
	idleSince *datetime.CpsTime, lastIdleRuleApply string) error {
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

func (a *mongoAdapter) FindByIDs(ctx context.Context, ids []string) ([]types.Entity, error) {
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
	time datetime.CpsTime,
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

func (a *mongoAdapter) GetWithIdleSince(ctx context.Context) (mongo.Cursor, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "idle_since", Value: 1}})

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

func (a *mongoAdapter) FindToCheckPbehaviorInfo(ctx context.Context, idsWithPbehaviors []string, exceptIds []string) (mongo.Cursor, error) {
	filter := bson.M{
		"enabled": true,
	}
	if len(exceptIds) > 0 {
		filter["_id"] = bson.M{"$nin": exceptIds}
	}
	if len(idsWithPbehaviors) > 0 {
		filter["$or"] = []bson.M{
			{"_id": bson.M{"$in": idsWithPbehaviors}},
			{"pbehavior_info": bson.M{"$ne": nil}},
		}
	} else {
		filter["pbehavior_info"] = bson.M{"$ne": nil}
	}

	opts := &options.FindOptions{}
	return a.dbCollection.Find(ctx, filter,
		opts.SetProjection(bson.M{"services": 0}))
}

func (a *mongoAdapter) UpdatePbehaviorInfo(ctx context.Context, id string, info types.PbehaviorInfo) error {
	if info.IsZero() {
		_, err := a.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$unset": bson.M{"pbehavior_info": ""}})
		return err
	}

	_, err := a.dbCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"pbehavior_info":      info,
		"last_pbehavior_date": info.Timestamp,
	}})
	return err
}

func (a *mongoAdapter) FindConnector(ctx context.Context, id string) (*types.Entity, error) {
	cursor, err := a.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"_id": id,
		}},
		{"$graphLookup": bson.M{
			"from":             mongo.EntityMongoCollection,
			"startWith":        "$connector",
			"connectFromField": "connector",
			"connectToField":   "_id",
			"as":               "connector",
			"maxDepth":         0,
		}},
		{"$unwind": "$connector"},
		{"$replaceRoot": bson.M{
			"newRoot": "$connector",
		}},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		res := types.Entity{}
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, nil
}

func (a *mongoAdapter) FindComponent(ctx context.Context, id string) (*types.Entity, error) {
	cursor, err := a.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"_id": id,
		}},
		{"$graphLookup": bson.M{
			"from":             mongo.EntityMongoCollection,
			"startWith":        "$component",
			"connectFromField": "component",
			"connectToField":   "_id",
			"as":               "component",
			"maxDepth":         0,
		}},
		{"$unwind": "$component"},
		{"$replaceRoot": bson.M{
			"newRoot": "$component",
		}},
	})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		res := types.Entity{}
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, nil
}
