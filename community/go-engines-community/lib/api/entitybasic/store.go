package entitybasic

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	GetOneBy(id string) (*Entity, error)
	Update(EditRequest) (*Entity, bool, error)
	Delete(id string) (bool, error)
}

type store struct {
	db                    mongo.DbClient
	dbCollection          mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
	basicTypes            []string
}

func NewStore(db mongo.DbClient) Store {
	return &store{
		db:           db,
		dbCollection: db.Collection(mongo.EntityMongoCollection),
		defaultSearchByFields: []string{
			"name", "type",
		},
		defaultSortBy: "name",
		basicTypes:    []string{types.EntityTypeResource, types.EntityTypeComponent, types.EntityTypeConnector},
	}
}

func (s *store) GetOneBy(id string) (*Entity, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pipeline := []bson.M{
		{"$match": bson.M{
			"_id":  id,
			"type": bson.M{"$in": s.basicTypes},
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
			"restrictSearchWithMatch": bson.M{"type": bson.M{"$in": s.basicTypes}},
			"maxDepth":                0,
		}},
		{"$unwind": bson.M{"path": "$changeable_impact", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"changeable_impact._id": 1}},
		{"$group": bson.M{
			"_id":               "$_id",
			"data":              bson.M{"$first": "$$ROOT"},
			"depends":           bson.M{"$first": "$depends"},
			"changeable_impact": bson.M{"$push": "$changeable_impact"},
		}},
		{"$graphLookup": bson.M{
			"from":                    mongo.EntityMongoCollection,
			"startWith":               "$depends",
			"connectFromField":        "depends",
			"connectToField":          "_id",
			"as":                      "changeable_depends",
			"restrictSearchWithMatch": bson.M{"type": bson.M{"$in": s.basicTypes}},
			"maxDepth":                0,
		}},
		{"$unwind": bson.M{"path": "$changeable_depends", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"changeable_depends._id": 1}},
		{"$group": bson.M{
			"_id":                "$_id",
			"data":               bson.M{"$first": "$data"},
			"changeable_impact":  bson.M{"$first": "$changeable_impact"},
			"changeable_depends": bson.M{"$push": "$changeable_depends"},
		}},
		{"$replaceRoot": bson.M{
			"newRoot": bson.M{"$mergeObjects": bson.A{
				"$data",
				bson.M{
					"changeable_impact": bson.M{"$map": bson.M{
						"input": "$changeable_impact",
						"as":    "i",
						"in":    "$$i._id",
					}},
					"changeable_depends": bson.M{"$map": bson.M{
						"input": "$changeable_depends",
						"as":    "d",
						"in":    "$$d._id",
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

func (s *store) Update(r EditRequest) (*Entity, bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	entity, err := s.GetOneBy(r.ID)
	if err != nil || entity == nil {
		return nil, false, err
	}

	res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID},
		bson.M{"$set": bson.M{
			"description":  r.Description,
			"enabled":      *r.Enabled,
			"category":     r.Category,
			"impact_level": r.ImpactLevel,
			"infos":        transformInfos(r),
		}},
	)
	if err != nil || res.MatchedCount == 0 {
		return nil, false, err
	}

	err = s.removeEntityLinks(*entity, r)
	if err != nil {
		return nil, false, err
	}
	err = s.addEntityLinks(*entity, r)
	if err != nil {
		return nil, false, err
	}

	updatedEntity, err := s.GetOneBy(r.ID)
	if err != nil {
		return nil, false, err
	}

	isToggled := updatedEntity.Enabled != entity.Enabled

	return updatedEntity, isToggled, nil
}

func (s *store) Delete(id string) (bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	alarmRes := s.db.Collection(mongo.AlarmMongoCollection).FindOne(ctx, bson.M{"d": entity.ID})
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

	return true, nil
}

func (s *store) removeEntityLinks(entity Entity, r EditRequest) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func (s *store) addEntityLinks(entity Entity, r EditRequest) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
