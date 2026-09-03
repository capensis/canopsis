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

// ErrLeaseLost reports that another worker owns the job's current lease.
var ErrLeaseLost = errors.New("async job lease lost")

// Store creates and updates the jobs of a single type in their own collection.
// T is the job document, which is expected to embed Base inline.
type Store[T Job] struct {
	collection libmongo.DbCollection
}

func NewStore[T Job](client libmongo.DbClient, collection string) *Store[T] {
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

// Claim moves a pending job to ongoing and returns its post-update document.
// A false return means another worker already took the job - on a redelivery,
// for instance - and the caller must not process it.
//
// LeaseGeneration is incremented atomically, so it can be used to fence workers
// which lost their lease and were superseded by another worker.
func (s *Store[T]) Claim(ctx context.Context, id string) (T, bool, error) {
	var job T
	err := s.collection.FindOneAndUpdate(ctx,
		bson.M{"_id": id, "status": StatusPending},
		bson.M{"$set": bson.M{
			"status":    StatusOngoing,
			"last_ping": datetime.NewCpsTime(),
		}, "$inc": bson.M{
			"lease_generation": 1,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&job)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return job, false, nil
		}

		return job, false, err
	}

	return job, true, nil
}

// Ping refreshes last_ping, so a long running job can be told apart from a job left ongoing by a
// worker which died. Handlers are expected to call it periodically. generation
// must be the one returned by Claim.
func (s *Store[T]) Ping(ctx context.Context, id string, generation int64) error {
	return s.updateClaimed(ctx, id, generation, bson.M{
		"$set": bson.M{"last_ping": datetime.NewCpsTime()},
	})
}

// SetDone marks the job as successfully processed. set holds the type-specific result fields, it may be nil.
func (s *Store[T]) SetDone(ctx context.Context, id string, generation int64, d time.Duration, set bson.M) error {
	return s.complete(ctx, id, generation, StatusDone, "", d, set)
}

// SetFailed marks the job as failed, errorMessage being the reason reported to the user.
func (s *Store[T]) SetFailed(ctx context.Context, id string, generation int64, errorMessage string, d time.Duration) error {
	return s.complete(ctx, id, generation, StatusFailed, errorMessage, d, nil)
}

func (s *Store[T]) complete(
	ctx context.Context,
	id string,
	generation int64,
	status, errorMessage string,
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

	return s.updateClaimed(ctx, id, generation, bson.M{"$set": update})
}

func (s *Store[T]) updateClaimed(ctx context.Context, id string, generation int64, update bson.M) error {
	res, err := s.collection.UpdateOne(ctx,
		bson.M{
			"_id":              id,
			"status":           StatusOngoing,
			"lease_generation": generation,
		},
		update,
	)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return ErrLeaseLost
	}

	return nil
}

// DeleteCompleted deletes one job only if it is still eligible for cleanup.
func (s *Store[T]) DeleteCompleted(ctx context.Context, id string, t datetime.CpsTime) (bool, error) {
	count, err := s.collection.DeleteOne(ctx, bson.M{
		"_id":       id,
		"status":    bson.M{"$in": bson.A{StatusDone, StatusFailed}},
		"completed": bson.M{"$lt": t},
	})
	return count > 0, err
}

// FindCompletedBefore returns completed jobs eligible for cleanup.
func (s *Store[T]) FindCompletedBefore(ctx context.Context, t datetime.CpsTime) ([]T, error) {
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
