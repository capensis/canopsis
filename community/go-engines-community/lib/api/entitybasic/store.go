package entitybasic

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	GetOneBy(ctx context.Context, id string) (*Entity, error)
	Update(ctx context.Context, r EditRequest) (*Entity, bool, error)
	Delete(ctx context.Context, id, userId string) (bool, error)
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
			"soft_deleted": bson.M{"$exists": false},
			"healthcheck":  bson.M{"$in": bson.A{nil, false}},
		}},
		// Find category
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
	}
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

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
	now := datetime.NewCpsTime()

	set := bson.M{
		"description":     r.Description,
		"enabled":         *r.Enabled,
		"category":        r.Category,
		"impact_level":    r.ImpactLevel,
		"infos":           entity.TransformInfosRequest(r.Infos),
		"sli_avail_state": r.SliAvailState,
		"author":          r.Author,
		"updated":         now,
	}
	update := bson.M{}
	if r.Coordinates == nil {
		update["$unset"] = bson.M{"coordinates": ""}
	} else {
		set["coordinates"] = r.Coordinates
	}
	update["$set"] = set

	var isToggled bool
	var updatedEntity *Entity

	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		isToggled = false
		updatedEntity = nil

		oldEntity := Entity{}
		cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
			{"$match": bson.M{"_id": r.ID}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$_id",
				"connectFromField":        "_id",
				"connectToField":          "component",
				"as":                      "resources",
				"restrictSearchWithMatch": bson.M{"type": types.EntityTypeResource},
				"maxDepth":                0,
			}},
			{"$project": bson.M{
				"_id":       1,
				"enabled":   1,
				"type":      1,
				"resources": bson.M{"$map": bson.M{"input": "$resources", "in": "$$this._id"}},
			}},
		})
		if err != nil {
			return err
		}

		defer cursor.Close(ctx)

		if cursor.Next(ctx) {
			err = cursor.Decode(&oldEntity)
			if err != nil {
				return err
			}
		} else {
			return nil
		}

		res, err := s.dbCollection.UpdateOne(ctx, bson.M{"_id": r.ID}, update)
		if err != nil || res.MatchedCount == 0 {
			return err
		}

		updatedEntity, err = s.GetOneBy(ctx, r.ID)
		if err != nil || updatedEntity == nil {
			return err
		}

		isToggled = updatedEntity.Enabled != oldEntity.Enabled
		updatedEntity.Resources = oldEntity.Resources
		return nil
	})

	if err != nil || updatedEntity == nil {
		return nil, false, err
	}

	if isToggled && !updatedEntity.Enabled && updatedEntity.Type == types.EntityTypeComponent {
		depLen := len(updatedEntity.Resources)
		from := 0

		for to := canopsis.DefaultBulkSize; to <= depLen; to += canopsis.DefaultBulkSize {
			_, err = s.dbCollection.UpdateMany(
				ctx,
				bson.M{"_id": bson.M{"$in": updatedEntity.Resources[from:to]}},
				bson.M{"$set": bson.M{
					"enabled": updatedEntity.Enabled,
					"author":  r.Author,
					"updated": now,
				}},
			)
			if err != nil {
				return nil, false, err
			}

			from = to
		}

		if from < depLen {
			_, err = s.dbCollection.UpdateMany(
				ctx,
				bson.M{"_id": bson.M{"$in": updatedEntity.Resources[from:depLen]}},
				bson.M{"$set": bson.M{
					"enabled": updatedEntity.Enabled,
					"author":  r.Author,
					"updated": now,
				}},
			)
			if err != nil {
				return nil, false, err
			}
		}
	}

	return updatedEntity, isToggled, nil
}

func (s *store) Delete(ctx context.Context, id, userId string) (bool, error) {
	res := false
	err := s.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		res = false
		entity := &types.Entity{}
		err := s.dbCollection.FindOneAndUpdate(ctx, bson.M{
			"_id":  id,
			"type": bson.M{"$in": s.basicTypes},
		}, bson.M{
			"$set": bson.M{
				"author": userId,
			},
		}).Decode(&entity)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}
			return err
		}

		if entity.Type == types.EntityTypeComponent {
			c, err := s.dbCollection.CountDocuments(ctx, bson.M{"component": entity.ID, "type": types.EntityTypeResource})
			if err != nil {
				return err
			}
			if c > 0 {
				return ErrComponent
			}
		}

		err = s.alarmDbCollection.FindOne(ctx, bson.M{
			"d":          entity.ID,
			"v.resolved": nil,
		}).Err()
		if err == nil {
			return ErrLinkedEntityToAlarm
		} else if !errors.Is(err, mongodriver.ErrNoDocuments) {
			return err
		}

		deleted, err := s.dbCollection.DeleteOne(ctx, bson.M{
			"_id":  id,
			"type": bson.M{"$in": s.basicTypes},
		})
		if err != nil || deleted == 0 {
			return err
		}

		if entity.Type == types.EntityTypeConnector {
			_, err = s.dbCollection.UpdateMany(ctx, bson.M{"connector": entity.ID},
				bson.M{
					"$set": bson.M{
						"author":  userId,
						"updated": datetime.NewCpsTime(),
					},
					"$unset": bson.M{"connector": ""},
				},
			)
			if err != nil {
				return err
			}
		}

		res = true
		return nil
	})

	return res, err
}
