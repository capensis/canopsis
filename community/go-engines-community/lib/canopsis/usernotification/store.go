package usernotification

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	CreateForInstructionApprove(ctx context.Context, time datetime.CpsTime, id, name, author, comment, userID, roleID string) error
	RemoveForInstructionApprove(ctx context.Context, id string) error

	CreateForInstructionRate(ctx context.Context, time datetime.CpsTime, id, name, userID string) error
	GetWriteModelForInstructionRate(time datetime.CpsTime, id, name, userID string) mongodriver.WriteModel
	RemoveForInstructionRate(ctx context.Context, id, userID string) error

	CreateForEventFilterFailure(ctx context.Context, time datetime.CpsTime, id, description string, roleIDs []string) error
	GetWriteModelForEventFilterFailure(time datetime.CpsTime, id, description string, roleIDs []string) mongodriver.WriteModel
	RemoveForEventFilterFailure(ctx context.Context, id string) error
}

type store struct {
	collection mongo.DbCollection
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		collection: dbClient.Collection(mongo.UserNotificationCollection),
	}
}

func (s *store) CreateForInstructionApprove(ctx context.Context, time datetime.CpsTime, id, name, author, comment, userID, roleID string) error {
	t := TypeInstructionApprove
	roleIDs := make([]string, 0)
	if roleID != "" {
		roleIDs = append(roleIDs, roleID)
	}

	_, err := s.collection.UpdateOne(ctx,
		bson.M{
			"type":     t,
			"rule._id": id,
		},
		bson.M{
			"$set": bson.M{
				"time":    time,
				"user":    userID,
				"roles":   roleIDs,
				"comment": comment,
				"rule": bson.M{
					"_id":  id,
					"name": name,
				},
				"author": author,
			},
			"$setOnInsert": bson.M{
				"_id":  utils.NewID(),
				"type": t,
			},
		},
		options.Update().SetUpsert(true),
	)

	return err
}

func (s *store) RemoveForInstructionApprove(ctx context.Context, id string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"type": TypeInstructionApprove, "rule._id": id})

	return err
}

func (s *store) CreateForInstructionRate(ctx context.Context, time datetime.CpsTime, id, name, userID string) error {
	wm := s.GetWriteModelForInstructionRate(time, id, name, userID)
	_, err := s.collection.BulkWrite(ctx, []mongodriver.WriteModel{wm})

	return err
}

func (s *store) GetWriteModelForInstructionRate(time datetime.CpsTime, id, name, userID string) mongodriver.WriteModel {
	t := TypeInstructionRate

	return mongodriver.NewUpdateOneModel().
		SetFilter(bson.M{
			"type":     t,
			"rule._id": id,
			"user":     userID,
		}).
		SetUpdate(bson.M{
			"$set": bson.M{
				"time":  time,
				"user":  userID,
				"roles": []string{},
				"rule": bson.M{
					"_id":  id,
					"name": name,
				},
			},
			"$setOnInsert": bson.M{
				"_id":  utils.NewID(),
				"type": t,
			},
		}).
		SetUpsert(true)
}

func (s *store) RemoveForInstructionRate(ctx context.Context, id, userID string) error {
	if userID == "" {
		_, err := s.collection.DeleteMany(ctx, bson.M{
			"type":     TypeInstructionRate,
			"rule._id": id,
		})

		return err
	}

	_, err := s.collection.DeleteOne(ctx, bson.M{
		"type":     TypeInstructionRate,
		"user":     userID,
		"rule._id": id,
	})

	return err
}

func (s *store) CreateForEventFilterFailure(
	ctx context.Context,
	time datetime.CpsTime,
	id, description string,
	roleIDs []string,
) error {
	wm := s.GetWriteModelForEventFilterFailure(time, id, description, roleIDs)
	_, err := s.collection.BulkWrite(ctx, []mongodriver.WriteModel{wm})

	return err
}

func (s *store) GetWriteModelForEventFilterFailure(time datetime.CpsTime, id, description string, roleIDs []string) mongodriver.WriteModel {
	t := TypeEventFilterFailure

	return mongodriver.NewUpdateOneModel().
		SetFilter(bson.M{
			"type":     t,
			"rule._id": id,
		}).
		SetUpdate(bson.M{
			"$set": bson.M{
				"time":  time,
				"user":  "",
				"roles": roleIDs,
				"rule": bson.M{
					"_id":  id,
					"name": description,
				},
			},
			"$setOnInsert": bson.M{
				"_id":  utils.NewID(),
				"type": t,
			},
		}).
		SetUpsert(true)
}

func (s *store) RemoveForEventFilterFailure(ctx context.Context, id string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"type": TypeEventFilterFailure, "rule._id": id})

	return err
}
