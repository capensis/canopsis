package contextgraph

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeFormat = "15:04:05"

type mongoReporter struct {
	client     libmongo.DbClient
	collection libmongo.DbCollection
}

func NewMongoStatusReporter(client libmongo.DbClient) StatusReporter {
	return &mongoReporter{
		client:     client,
		collection: client.Collection(libmongo.ImportJobMongoCollection),
	}
}

func (r *mongoReporter) ReportCreate(ctx context.Context, job *ImportJob) error {
	job.ID = utils.NewID()

	_, err := r.collection.InsertOne(ctx, job)
	return err
}

func (r *mongoReporter) GetFirst(ctx context.Context, abandonedInterval time.Duration) (ImportJob, error) {
	job := ImportJob{}
	err := r.client.WithTransaction(ctx, func(ctx context.Context) error {
		job = ImportJob{}

		ongoingJob := ImportJob{}
		err := r.collection.FindOne(ctx, bson.M{"status": StatusOngoing}).Decode(&ongoingJob)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
		if ongoingJob.ID != "" {
			if ongoingJob.LastPing != nil && ongoingJob.LastPing.After(time.Now().Add(-abandonedInterval)) {
				return nil
			}

			now := time.Now()
			ongoingJob.LastPing = &now
			_, err = r.collection.UpdateOne(ctx, bson.M{"_id": ongoingJob.ID}, bson.M{"$set": bson.M{
				"last_ping": ongoingJob.LastPing,
			}})
			if err != nil {
				return err
			}

			job = ongoingJob
			return nil
		}

		err = r.collection.FindOneAndUpdate(
			ctx,
			bson.M{"status": StatusPending},
			bson.M{"$set": bson.M{
				"status":    StatusOngoing,
				"last_ping": time.Now(),
			}},
			options.FindOneAndUpdate().
				SetSort(bson.M{"creation": 1}).
				SetReturnDocument(options.After),
		).Decode(&job)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}

		return nil
	})

	return job, err
}

func (r *mongoReporter) ReportOngoing(ctx context.Context, job ImportJob) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":    job.ID,
		"status": StatusOngoing,
	}, bson.M{"$set": bson.M{
		"last_ping": time.Now(),
	}})

	return err
}

func (r *mongoReporter) ReportDone(ctx context.Context, job ImportJob, stats importcontextgraph.Stats) (bool, error) {
	t := time.Time{}
	execTime := t.Add(stats.ExecTime).Format(timeFormat)
	res, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":    job.ID,
		"status": StatusOngoing,
	}, bson.M{"$set": bson.M{
		"status":    StatusDone,
		"exec_time": execTime,
		"stats":     stats,
	}})
	if err != nil {
		return false, err
	}

	return res.ModifiedCount > 0, nil
}

func (r *mongoReporter) ReportError(ctx context.Context, job ImportJob, execDuration time.Duration, err error) (bool, error) {
	t := time.Time{}
	execTime := t.Add(execDuration).Format(timeFormat)

	res, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":    job.ID,
		"status": StatusOngoing,
	}, bson.M{"$set": bson.M{
		"status":    StatusFailed,
		"exec_time": execTime,
		"info":      err.Error(),
	}})
	if err != nil {
		return false, err
	}

	return res.ModifiedCount > 0, nil
}

func (r *mongoReporter) GetStatus(ctx context.Context, id string) (ImportJob, error) {
	var status ImportJob
	res := r.collection.FindOne(ctx, bson.M{"_id": id})
	err := res.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = ErrNotFound
		}
	} else {
		err = res.Decode(&status)
	}

	return status, err
}

func (r *mongoReporter) Clean(ctx context.Context, interval time.Duration) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{
		"status": bson.M{"$in": bson.A{StatusDone, StatusFailed}},
		"$or": []bson.M{
			{"creation": bson.M{"$lt": time.Now().Add(-interval)}},
			{"creation": bson.M{"$type": "string"}},
		}},
	)
	return err
}
