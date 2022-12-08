package entitybasic

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	GetOneBy(ctx context.Context, id string) (*Entity, error)
	Update(ctx context.Context, r EditRequest) (*Entity, bool, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	dbClient          mongo.DbClient
	dbCollection      mongo.DbCollection
	alarmDbCollection mongo.DbCollection
	basicTypes        []string
}

func NewStore(db mongo.DbClient) Store {
	return &store{
		dbClient:          db,
		dbCollection:      db.Collection(mongo.EntityMongoCollection),
		alarmDbCollection: db.Collection(mongo.AlarmMongoCollection),
		basicTypes:        []string{types.EntityTypeResource, types.EntityTypeComponent, types.EntityTypeConnector},
	}
}

func (s *store) GetOneBy(ctx context.Context, id string) (*Entity, error) {
	pipeline := []bson.M{
		{"$match": bson.M{
			"_id":          id,
			"type":         bson.M{"$in": s.basicTypes},
			"soft_deleted": bson.M{"$in": bson.A{false, nil}},
		}},
		// Find category
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
		// Find links to basic entities
		{"$graphLookup": bson.M{
			"from":                    mongo.EntityMongoCollection,
			"startWith":               "$impact",
			"connectFromField":        "impact",
			"connectToField":          "_id",
			"as":                      "changeable_impact",
			"restrictSearchWithMatch": bson.M{"type": bson.M{"$in": s.basicTypes}, "soft_deleted": bson.M{"$exists": false}},
			"maxDepth":                0,
		}},
		{"$addFields": bson.M{
			"changeable_impact": bson.M{"$map": bson.M{"input": "$changeable_impact", "as": "each", "in": "$$each._id"}},
		}},
		{"$project": bson.M{"_id": "$_id", "changeable_impact": "$changeable_impact", "data": "$$ROOT", "depends": "$depends"}},
		{"$project": bson.M{"data.changeable_impact": 0}},
		{"$graphLookup": bson.M{
			"from":                    mongo.EntityMongoCollection,
			"startWith":               "$depends",
			"connectFromField":        "depends",
			"connectToField":          "_id",
			"as":                      "changeable_depends",
			"restrictSearchWithMatch": bson.M{"type": bson.M{"$in": s.basicTypes}, "soft_deleted": bson.M{"$exists": false}},
			"maxDepth":                0,
		}},
		{"$addFields": bson.M{
			"changeable_depends": bson.M{"$map": bson.M{"input": "$changeable_depends", "as": "each", "in": "$$each._id"}},
		}},
		{"$project": bson.M{"depends": 0}},
		{"$replaceRoot": bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.M{
					"changeable_impact": bson.M{"$map": bson.M{
						"input": "$changeable_impact",
						"as":    "i",
						"in":    "$$i",
					}},
					"changeable_depends": bson.M{"$map": bson.M{
						"input": "$changeable_depends",
						"as":    "d",
						"in":    "$$d",
					}},
				},
			}},
		}},
	}
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if cursor.Next(ctx) {
		res := Entity{}
		err := cursor.Decode(&res)
		if err != nil {
			return nil, fmt.Errorf("cannot decode: %w", err)
		}

		return &res, nil
	}

	return nil, nil
}

func (s *store) Update(ctx context.Context, r EditRequest) (*Entity, bool, error) {
	entity, err := s.GetOneBy(ctx, r.ID)
	if err != nil || entity == nil {
		return nil, false, err
	}

	set := bson.M{
		"description":     r.Description,
		"enabled":         *r.Enabled,
		"category":        r.Category,
		"impact_level":    r.ImpactLevel,
		"infos":           transformInfos(r),
		"sli_avail_state": r.SliAvailState,
	}
	update := bson.M{}
	if r.Coordinates == nil {
		update["$unset"] = bson.M{"coordinates": ""}
	} else {
		set["coordinates"] = r.Coordinates
	}
	update["$set"] = set

	var updatedEntity *Entity

	err = s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedEntity = nil
		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, update)
		if err != nil || res.MatchedCount == 0 {
			return err
		}
		err = s.removeEntityLinks(ctx, *entity, r)
		if err != nil {
			return err
		}
		err = s.addEntityLinks(ctx, *entity, r)
		if err != nil {
			return err
		}

		updatedEntity, err = s.GetOneBy(ctx, r.ID)
		return err
	})

	if err != nil || updatedEntity == nil {
		return nil, false, err
	}

	isToggled := updatedEntity.Enabled != entity.Enabled

	return updatedEntity, isToggled, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	res := s.dbCollection.FindOne(ctx, bson.M{
		"_id":  id,
		"type": bson.M{"$in": s.basicTypes},
	})
	if err := res.Err(); err != nil {
		if err == mongodriver.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	entity := &types.Entity{}
	err := res.Decode(entity)
	if err != nil {
		return false, err
	}

	alarmRes := s.alarmDbCollection.FindOne(ctx, bson.M{
		"d":          entity.ID,
		"v.resolved": nil,
	})
	if err := alarmRes.Err(); err == nil {
		return false, ErrLinkedEntityToAlarm
	} else if err != mongodriver.ErrNoDocuments {
		return false, err
	}

	deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{
		"_id":  id,
		"type": bson.M{"$in": s.basicTypes},
	})
	if err != nil {
		return false, err
	}

	if deleted == 0 {
		return false, nil
	}

	_, err = s.dbCollection.UpdateMany(ctx, bson.M{"impact": id},
		bson.M{"$pull": bson.M{"impact": id}})
	if err != nil {
		return false, err
	}
	_, err = s.dbCollection.UpdateMany(ctx, bson.M{"depends": id},
		bson.M{"$pull": bson.M{"depends": id}})
	if err != nil {
		return false, err
	}
	_, err = s.dbCollection.UpdateMany(ctx, bson.M{"connector": id},
		bson.M{"$unset": bson.M{"connector": ""}})
	if err != nil {
		return false, err
	}
	_, err = s.dbCollection.UpdateMany(ctx, bson.M{"component": id},
		bson.M{"$unset": bson.M{"component": ""}})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) removeEntityLinks(ctx context.Context, entity Entity, r EditRequest) error {
	impacts := diffSlices(entity.ChangeableImpacts, r.Impacts)
	depends := diffSlices(entity.ChangeableDepends, r.Depends)
	models := make([]mongodriver.WriteModel, 0)
	pull := bson.M{}

	if len(impacts) > 0 {
		pull["impact"] = impacts
		for _, impact := range impacts {
			models = append(models, mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": impact}).
				SetUpdate(bson.M{"$pull": bson.M{"depends": entity.ID}}))
		}
	}
	if len(depends) > 0 {
		pull["depends"] = depends
		for _, depend := range depends {
			models = append(models, mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": depend}).
				SetUpdate(bson.M{"$pull": bson.M{"impact": entity.ID}}))
		}
	}
	if len(pull) > 0 {
		models = append(models, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": entity.ID}).
			SetUpdate(bson.M{"$pullAll": pull}))
	}

	if len(models) > 0 {
		_, err := s.dbCollection.BulkWrite(ctx, models)
		return err
	}

	return nil
}

func (s *store) addEntityLinks(ctx context.Context, entity Entity, r EditRequest) error {
	impacts := diffSlices(r.Impacts, entity.ChangeableImpacts)
	depends := diffSlices(r.Depends, entity.ChangeableDepends)
	models := make([]mongodriver.WriteModel, 0)

	addToSet := bson.M{}
	if len(impacts) > 0 {
		addToSet["impact"] = bson.M{"$each": impacts}
		for _, impact := range impacts {
			models = append(models, mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": impact}).
				SetUpdate(bson.M{"$addToSet": bson.M{"depends": entity.ID}}))
		}
	}
	if len(depends) > 0 {
		addToSet["depends"] = bson.M{"$each": depends}
		for _, depend := range depends {
			models = append(models, mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": depend}).
				SetUpdate(bson.M{"$addToSet": bson.M{"impact": entity.ID}}))
		}
	}
	if len(addToSet) > 0 {
		models = append(models, mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": entity.ID}).
			SetUpdate(bson.M{"$addToSet": addToSet}))
	}

	if len(models) > 0 {
		_, err := s.dbCollection.BulkWrite(ctx, models)
		return err
	}

	return nil
}

func transformInfos(request EditRequest) map[string]types.Info {
	infos := make(map[string]types.Info, len(request.Infos))
	for _, v := range request.Infos {
		infos[v.Name] = types.Info{
			Name:        v.Name,
			Description: v.Description,
			Value:       v.Value,
		}
	}

	return infos
}

func diffSlices(l, r []string) []string {
	diff := make([]string, 0)

	for _, lv := range l {
		found := false
		for _, rv := range r {
			if lv == rv {
				found = true
				break
			}
		}

		if !found {
			diff = append(diff, lv)
		}
	}

	return diff
}
