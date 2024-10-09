package correlation

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MetaAlarmStateService interface {
	// GetMetaAlarmState returns current alarm state.
	GetMetaAlarmState(ctx context.Context, ruleID string) (MetaAlarmState, error)
	// UpdateOpenedState updates opened metaalarm state by increasing its version and replaces alarm groups.
	UpdateOpenedState(ctx context.Context, state MetaAlarmState, previousVersion int64, previousState int, upsert bool) (bool, error)
	// ArchiveState moves state to a separate document, needed for late create metaalarm events to get their state instead of new one.
	ArchiveState(ctx context.Context, state MetaAlarmState) (bool, error)
	// DeleteState deletes current alarm state.
	DeleteState(ctx context.Context, ruleID string) (bool, error)
	// SwitchStateToReady switch state status to ready, should be used only after metaalarm is triggered.
	SwitchStateToReady(ctx context.Context, state MetaAlarmState, previousVersion int64, previousState int, upsert bool) (bool, error)
	// SwitchStateToCreated switch state status to created, should be used only after or during the metaalarm creation.
	SwitchStateToCreated(ctx context.Context, stateID string) (bool, error)
	// PushChild pushes child to the child group in the DB, should be used only when group is gathered and metaalarm is triggered.
	PushChild(ctx context.Context, state MetaAlarmState, previousVersion int64, previousState int, entityID string, alarmLastUpdate int64) (bool, error)
	// RefreshMetaAlarmStateGroup removes resolved alarms and current event entity from the group.
	RefreshMetaAlarmStateGroup(ctx context.Context, state *MetaAlarmState, entityID string, timeInterval int64) error
}

type metaAlarmStateService struct {
	alarmCollection           mongo.DbCollection
	metaAlarmStatesCollection mongo.DbCollection
}

func NewMetaAlarmStateService(dbClient mongo.DbClient) MetaAlarmStateService {
	return &metaAlarmStateService{
		alarmCollection:           dbClient.Collection(mongo.AlarmMongoCollection),
		metaAlarmStatesCollection: dbClient.Collection(mongo.MetaAlarmStatesCollection),
	}
}

func (a *metaAlarmStateService) GetMetaAlarmState(ctx context.Context, ruleID string) (MetaAlarmState, error) {
	var metaAlarmState MetaAlarmState

	err := a.metaAlarmStatesCollection.FindOne(ctx, bson.M{"_id": ruleID}).Decode(&metaAlarmState)
	if err != nil && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return metaAlarmState, err
	}

	return metaAlarmState, nil
}

func (a *metaAlarmStateService) UpdateOpenedState(ctx context.Context, state MetaAlarmState, previousVersion int64, previousState int, upsert bool) (bool, error) {
	res, err := a.metaAlarmStatesCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":     state.ID,
			"version": previousVersion,
			"state":   previousState,
		},
		bson.M{
			"$inc": bson.M{
				"version": 1,
			},
			"$set": bson.M{
				"children_entity_ids": state.ChildrenEntityIDs,
				"children_timestamps": state.ChildrenTimestamps,
				"parents_entity_ids":  state.ParentsEntityIDs,
				"parents_timestamps":  state.ParentsTimestamps,
				"meta_alarm_name":     state.MetaAlarmName,
				"state":               Opened,
			},
		},
		options.Update().SetUpsert(upsert),
	)
	if err != nil || res.MatchedCount == 0 && res.UpsertedCount == 0 {
		if mongodriver.IsDuplicateKeyError(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (a *metaAlarmStateService) ArchiveState(ctx context.Context, state MetaAlarmState) (bool, error) {
	state.ID = state.ID + "-" + state.MetaAlarmName
	state.CreatedAt = time.Now()

	_, err := a.metaAlarmStatesCollection.InsertOne(ctx, state)
	if err != nil {
		if mongodriver.IsDuplicateKeyError(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (a *metaAlarmStateService) DeleteState(ctx context.Context, ruleID string) (bool, error) {
	deleted, err := a.metaAlarmStatesCollection.DeleteOne(ctx, bson.M{"_id": ruleID})

	return deleted > 0, err
}

func (a *metaAlarmStateService) SwitchStateToReady(ctx context.Context, state MetaAlarmState, previousVersion int64, previousState int, upsert bool) (bool, error) {
	res, err := a.metaAlarmStatesCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":     state.ID,
			"version": previousVersion,
			"state":   previousState,
		},
		bson.M{
			"$set": bson.M{
				"state":               Ready,
				"children_entity_ids": state.ChildrenEntityIDs,
				"children_timestamps": state.ChildrenTimestamps,
				"parents_entity_ids":  state.ParentsEntityIDs,
				"parents_timestamps":  state.ParentsTimestamps,
				"meta_alarm_name":     state.MetaAlarmName,
				"version":             0,
			},
		},
		options.Update().SetUpsert(upsert),
	)
	if err != nil || res.MatchedCount == 0 && res.UpsertedCount == 0 {
		if mongodriver.IsDuplicateKeyError(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (a *metaAlarmStateService) PushChild(ctx context.Context, state MetaAlarmState, previousVersion int64, previousState int, entityID string, alarmLastUpdate int64) (bool, error) {
	res, err := a.metaAlarmStatesCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":     state.ID,
			"version": previousVersion,
			"state":   previousState,
		},
		bson.M{
			"$push": bson.M{
				"children_entity_ids": entityID,
				"children_timestamps": alarmLastUpdate,
			},
		},
	)
	if err != nil || res.MatchedCount == 0 {
		return false, err
	}

	return true, nil
}

func (a *metaAlarmStateService) getAlarmsByIDs(entityIDs []string, excludeID string) []bson.M {
	return []bson.M{
		{
			"$match": bson.M{
				"$and": bson.A{
					bson.M{"_id": bson.M{"$in": entityIDs}},
					bson.M{"d": bson.M{"$ne": excludeID}},
				},
				"v.resolved": nil,
			},
		},
		{
			"$lookup": bson.M{
				"from":         mongo.EntityMongoCollection,
				"localField":   "d",
				"foreignField": "_id",
				"as":           "entity",
			},
		},
		{
			"$unwind": "$entity",
		},
		{
			"$project": bson.M{
				"_id":         "$d",
				"last_update": "$v.last_update_date",
			},
		},
	}
}

func (a *metaAlarmStateService) RefreshMetaAlarmStateGroup(ctx context.Context, state *MetaAlarmState, entityID string, timeInterval int64) error {
	if len(state.ChildrenEntityIDs) > 0 {
		cursor, err := a.alarmCollection.Aggregate(ctx, a.getAlarmsByIDs(state.ChildrenEntityIDs, entityID))
		if err != nil {
			return err
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var doc struct {
				ID         string `bson:"_id"`
				LastUpdate int64  `bson:"last_update"`
			}

			err = cursor.Decode(&doc)
			if err != nil {
				return err
			}

			state.PushChild(doc.ID, doc.LastUpdate, timeInterval)
		}
	}

	if len(state.ParentsEntityIDs) > 0 {
		cursor, err := a.alarmCollection.Aggregate(ctx, a.getAlarmsByIDs(state.ParentsEntityIDs, entityID))
		if err != nil {
			return err
		}

		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var doc struct {
				ID         string `bson:"_id"`
				LastUpdate int64  `bson:"last_update"`
			}

			err = cursor.Decode(&doc)
			if err != nil {
				return err
			}

			state.PushParent(doc.ID, doc.LastUpdate, timeInterval)
		}
	}

	return nil
}

func (a *metaAlarmStateService) SwitchStateToCreated(ctx context.Context, stateID string) (bool, error) {
	res, err := a.metaAlarmStatesCollection.UpdateOne(ctx, bson.M{"_id": stateID, "state": Ready}, bson.M{
		"$set": bson.M{
			"state": Created,
		},
	})
	if err != nil || res.MatchedCount == 0 {
		return false, err
	}

	return true, nil
}
