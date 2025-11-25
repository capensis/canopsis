package entity

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// mongoAdapter provides MongoDB-backed implementation for entity Adapter.
type mongoAdapter struct {
	dbCollection mongo.DbCollection
}

// NewAdapter returns an entity Adapter backed by MongoDB.
func NewAdapter(dbClient mongo.DbClient) Adapter {
	return &mongoAdapter{
		dbCollection: dbClient.Collection(mongo.EntityMongoCollection),
	}
}

func (a *mongoAdapter) Bulk(ctx context.Context, models []mongodriver.WriteModel) error {
	_, err := a.dbCollection.BulkWrite(ctx, models)

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

func (a *mongoAdapter) FindToCheckPbehaviorInfo(ctx context.Context, idsWithPbehaviors, exceptIDs, serviceIDs []string) (mongo.Cursor, error) {
	filter := bson.M{
		"enabled": true,
	}
	if len(exceptIDs) > 0 {
		filter["_id"] = bson.M{"$nin": exceptIDs}
	}

	var orStmt []bson.M

	if len(idsWithPbehaviors) > 0 {
		orStmt = append(orStmt, bson.M{"_id": bson.M{"$in": idsWithPbehaviors}})
	}

	if len(serviceIDs) > 0 {
		orStmt = append(orStmt, bson.M{"services": bson.M{"$in": serviceIDs}})
	}

	if len(orStmt) > 0 {
		orStmt = append(orStmt, bson.M{"pbehavior_info.id": bson.M{"$ne": nil}})
		filter["$or"] = orStmt
	} else {
		filter["pbehavior_info.id"] = bson.M{"$ne": nil}
	}

	return a.dbCollection.Find(ctx, filter)
}
