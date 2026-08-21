package asyncjob

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Store creates and updates the jobs of a single type in their own collection.
// T is the job document, which is expected to embed Base inline.
type Store[T any] struct {
	collection libmongo.DbCollection
}

func NewStore[T any](client libmongo.DbClient, collection string) *Store[T] {
	return &Store[T]{
		collection: client.Collection(collection),
	}
}

// Create stores a pending job.
func (s *Store[T]) Create(ctx context.Context, job T) error {
	_, err := s.collection.InsertOne(ctx, job)

	return err
}

// Get returns the job, or httperror.ErrNotFound if it doesn't exist.
func (s *Store[T]) Get(ctx context.Context, id string) (T, error) {
	var job T
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&job)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = httperror.ErrNotFound
		}

		return job, err
	}

	return job, nil
}

// GetUncompleted returns a job which is still pending or ongoing, or httperror.ErrNotFound if it doesn't exist.
func (s *Store[T]) GetUncompleted(ctx context.Context) (T, error) {
	var job T
	err := s.collection.FindOne(ctx,
		bson.M{"status": bson.M{"$in": bson.A{StatusPending, StatusOngoing}}},
		options.FindOne().SetSort(bson.M{"created": 1}),
	).Decode(&job)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = httperror.ErrNotFound
		}

		return job, err
	}

	return job, nil
}

// SetOngoing moves a pending job to ongoing and reports whether it was still pending.
// A false return means another worker already took the job - on a redelivery,
// for instance - and the caller must not process it.
func (s *Store[T]) SetOngoing(ctx context.Context, id string) (bool, error) {
	res, err := s.collection.UpdateOne(ctx,
		bson.M{"_id": id, "status": StatusPending},
		bson.M{"$set": bson.M{
			"status":    StatusOngoing,
			"last_ping": datetime.NewCpsTime(),
		}},
	)
	if err != nil {
		return false, err
	}

	return res.MatchedCount > 0, nil
}

// Ping refreshes last_ping, so a long running job can be told apart from a job left ongoing by a
// worker which died. Handlers are expected to call it periodically.
func (s *Store[T]) Ping(ctx context.Context, id string) error {
	_, err := s.collection.UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"last_ping": datetime.NewCpsTime()}},
	)

	return err
}

// SetDone marks the job as successfully processed. set holds the type-specific result fields, it may be nil.
func (s *Store[T]) SetDone(ctx context.Context, id string, d time.Duration, set bson.M) error {
	return s.complete(ctx, id, StatusDone, "", d, set)
}

// SetFailed marks the job as failed, errorMessage being the reason reported to the user.
func (s *Store[T]) SetFailed(ctx context.Context, id, errorMessage string, d time.Duration) error {
	return s.complete(ctx, id, StatusFailed, errorMessage, d, nil)
}

func (s *Store[T]) complete(
	ctx context.Context,
	id, status, errorMessage string,
	d time.Duration,
	set bson.M,
) error {
	update := bson.M{
		"status":      status,
		"duration_ms": d.Milliseconds(),
		"completed":   datetime.NewCpsTime(),
	}

	if errorMessage != "" {
		update["error_message"] = errorMessage
	}

	for k, v := range set {
		update[k] = v
	}

	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})

	return err
}

// DeleteCompletedBefore deletes the jobs which completed before t and returns them,
// so the files they own can be removed too.
func (s *Store[T]) DeleteCompletedBefore(ctx context.Context, t datetime.CpsTime) ([]T, error) {
	filter := bson.M{
		"status":    bson.M{"$in": bson.A{StatusDone, StatusFailed}},
		"completed": bson.M{"$lt": t},
	}

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var jobs []T
	err = cursor.All(ctx, &jobs)
	if err != nil {
		return nil, err
	}

	if len(jobs) == 0 {
		return nil, nil
	}

	_, err = s.collection.DeleteMany(ctx, filter)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// claimStale moves back to pending one job left ongoing by a worker which died, ie. which hasn't pinged since t, and returns its id.
// It returns false when there's no such job left to claim.
func (s *Store[T]) claimStale(ctx context.Context, t datetime.CpsTime, maxAttempts int64) (string, bool, error) {
	var job Base
	err := s.collection.FindOneAndUpdate(ctx,
		bson.M{
			"status":    StatusOngoing,
			"last_ping": bson.M{"$lt": t},
			"attempts":  bson.M{"$lt": maxAttempts},
		},
		bson.M{
			"$set": bson.M{
				"status":    StatusPending,
				"last_ping": datetime.NewCpsTime(),
			},
			"$inc": bson.M{"attempts": 1},
		},
	).Decode(&job)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", false, nil
		}

		return "", false, err
	}

	return job.ID, true, nil
}

// claimUnsignaled claims one pending job which hasn't been heard of since t and returns its id,
// so it can be signaled again. It returns false when there's no such job left to claim.
func (s *Store[T]) claimUnsignaled(ctx context.Context, t datetime.CpsTime) (string, bool, error) {
	var job Base
	err := s.collection.FindOneAndUpdate(ctx,
		bson.M{
			"status":    StatusPending,
			"last_ping": bson.M{"$lt": t},
		},
		bson.M{"$set": bson.M{
			"last_ping": datetime.NewCpsTime(),
		}},
	).Decode(&job)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", false, nil
		}

		return "", false, err
	}

	return job.ID, true, nil
}

// failStale marks as failed the jobs which were left ongoing by a worker which died and ran out of attempts.
func (s *Store[T]) failStale(ctx context.Context, t datetime.CpsTime, maxAttempts int64) error {
	_, err := s.collection.UpdateMany(ctx,
		bson.M{
			"status":    StatusOngoing,
			"last_ping": bson.M{"$lt": t},
			"attempts":  bson.M{"$gte": maxAttempts},
		},
		bson.M{"$set": bson.M{
			"status":        StatusFailed,
			"completed":     datetime.NewCpsTime(),
			"error_message": "job was interrupted too many times",
		}},
	)

	return err
}
