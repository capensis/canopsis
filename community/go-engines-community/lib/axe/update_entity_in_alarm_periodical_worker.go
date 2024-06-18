package axe

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type updateEntityInAlarmWorker struct {
	AlarmCollection  mongo.DbCollection
	EntityCollection mongo.DbCollection
}

func (w *updateEntityInAlarmWorker) Work(ctx context.Context) error {
	stream, err := w.EntityCollection.Watch(
		ctx,
		[]bson.M{{"$match": bson.M{
			"operationType": "update",
		}}},
	)
	if err != nil {
		return fmt.Errorf("cannot create change stream: %w", err)
	}

	defer stream.Close(ctx)
	for stream.Next(ctx) {
		select {
		case <-ctx.Done():
			return nil
		default:
			var changeEvent struct {
				DocumentKey struct {
					ID string `bson:"_id"`
				} `bson:"documentKey"`
				UpdateDescription struct {
					UpdatedFields map[string]any `bson:"updatedFields"`
				} `bson:"updateDescription"`
			}
			err = stream.Decode(&changeEvent)
			if err != nil {
				return fmt.Errorf("cannot decode change event: %w", err)
			}

			updatedFields := make(map[string]any, len(changeEvent.UpdateDescription.UpdatedFields))
			for k, v := range changeEvent.UpdateDescription.UpdatedFields {
				updatedFields["entity."+k] = v
			}

			_, err = w.AlarmCollection.UpdateMany(ctx, bson.M{"d": changeEvent.DocumentKey.ID}, bson.M{"$set": updatedFields})
			if err != nil {
				return fmt.Errorf("cannot update alarms: %w", err)
			}
		}
	}

	return nil
}
