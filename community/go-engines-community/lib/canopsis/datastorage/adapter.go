package datastorage

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type adapter struct {
	collection mongo.DbCollection
}

func NewAdapter(client mongo.DbClient) Adapter {
	return &adapter{
		collection: client.Collection(mongo.ConfigurationMongoCollection),
	}
}

func (a *adapter) Get(ctx context.Context) (DataStorage, error) {
	data := DataStorage{}
	err := a.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&data)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return data, nil
		}

		return data, err
	}

	return data, nil
}

func (a *adapter) UpdateHistoryJunit(ctx context.Context, t types.CpsTime) error {
	res, err := a.collection.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{
		"$set": bson.M{
			"history.junit": t,
		},
	})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("cannot find configuration _id=%s", ID)
	}

	return nil
}

func (a *adapter) UpdateHistoryRemediation(ctx context.Context, t types.CpsTime) error {
	res, err := a.collection.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{
		"$set": bson.M{
			"history.remediation": t,
		},
	})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("cannot find configuration _id=%s", ID)
	}

	return nil
}

func (a *adapter) UpdateHistoryAlarm(ctx context.Context, history HistoryWithCount) error {
	res, err := a.collection.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{
		"$set": bson.M{
			"history.alarm": history,
		},
	})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("cannot find configuration _id=%s", ID)
	}

	return nil
}

func (a *adapter) UpdateHistoryEntity(ctx context.Context, history HistoryWithCount) error {
	res, err := a.collection.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{
		"$set": bson.M{
			"history.entity": history,
		},
	})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("cannot find configuration _id=%s", ID)
	}

	return nil
}

func (a *adapter) UpdateHistoryPbehavior(ctx context.Context, t types.CpsTime) error {
	res, err := a.collection.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{
		"$set": bson.M{
			"history.pbehavior": t,
		},
	})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("cannot find configuration _id=%s", ID)
	}

	return nil
}

func (a *adapter) UpdateHistoryHealthCheck(ctx context.Context, t types.CpsTime) error {
	res, err := a.collection.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{
		"$set": bson.M{
			"history.health_check": t,
		},
	})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("cannot find configuration _id=%s", ID)
	}

	return nil
}

func (a *adapter) UpdateHistoryWebhook(ctx context.Context, t types.CpsTime) error {
	res, err := a.collection.UpdateOne(ctx, bson.M{"_id": ID}, bson.M{
		"$set": bson.M{
			"history.webhook": t,
		},
	})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("cannot find configuration _id=%s", ID)
	}

	return nil
}
