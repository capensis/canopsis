package closedelay

//go:generate go tool go.uber.org/mock/mockgen -destination=../../../mocks/lib/canopsis/closedelay/service.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/closedelay Service

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
)

// PeriodsUntilResend defines the number of periodical intervals to wait
// before the task needs to be resent.
const PeriodsUntilResend = 5

type Service interface {
	Process(ctx context.Context) ([]types.Event, error)
}

type service struct {
	collection     mongo.DbCollection
	eventGenerator libevent.Generator
	resendDelay    time.Duration
}

func NewService(client mongo.DbClient, generator libevent.Generator, periodicalInterval time.Duration) Service {
	return &service{
		collection:     client.Collection(mongo.CloseDelayJobCollection),
		eventGenerator: generator,
		resendDelay:    PeriodsUntilResend * periodicalInterval,
	}
}

func (s *service) Process(ctx context.Context) ([]types.Event, error) {
	now := datetime.NewCpsTime()

	cursor, err := s.collection.Find(ctx, bson.M{
		"exec_time": bson.M{"$lte": now},
		"$or": bson.A{
			bson.M{
				"resend_at": bson.M{"$exists": false},
			},
			bson.M{
				"resend_at": bson.M{"$lte": now},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find close delay jobs: %w", err)
	}

	defer cursor.Close(ctx)

	var events []types.Event
	writeModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)

	for cursor.Next(ctx) {
		var job Job

		err := cursor.Decode(&job)
		if err != nil {
			return nil, fmt.Errorf("failed to decode close delay job: %w", err)
		}

		event, err := s.eventGenerator.Generate(types.Entity{
			Type:      job.Type,
			Name:      job.Name,
			Component: job.Component,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to generate close delay event: %w", err)
		}

		event.AlarmID = job.ID
		event.EventType = types.EventTypeCheck
		event.Timestamp = now
		event.State = types.AlarmStateOK
		event.Output = fmt.Sprintf("closed after %d seconds delay", job.Delay)
		event.IsCloseDelayJob = true

		events = append(events, event)

		writeModels = append(
			writeModels,
			mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": job.ID}).
				SetUpdate(bson.M{"$set": bson.M{"resend_at": now.Add(s.resendDelay).Unix()}}),
		)

		if len(writeModels) == canopsis.DefaultBulkSize {
			_, err := s.collection.BulkWrite(ctx, writeModels)
			if err != nil {
				return nil, fmt.Errorf("failed to bulk update close delay jobs: %w", err)
			}

			writeModels = writeModels[:0]
		}
	}

	if len(writeModels) > 0 {
		_, err := s.collection.BulkWrite(ctx, writeModels)
		if err != nil {
			return nil, fmt.Errorf("failed to bulk update close delay jobs: %w", err)
		}
	}

	err = cursor.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to process close delay jobs cursor: %w", err)
	}

	return events, nil
}
