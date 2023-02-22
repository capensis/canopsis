package contextgraph

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/importcontextgraph"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const timeFormat = "15:04:05"

type mongoReporter struct {
	collection libmongo.DbCollection
}

func NewMongoStatusReporter(client libmongo.DbClient) StatusReporter {
	return &mongoReporter{
		collection: client.Collection(libmongo.ImportJobMongoCollection),
	}
}

func (r *mongoReporter) ReportCreate(ctx context.Context, job *ImportJob) error {
	job.ID = utils.NewID()

	_, err := r.collection.InsertOne(ctx, job)
	return err
}

func (r *mongoReporter) ReportOngoing(ctx context.Context, job ImportJob) (bool, error) {
	res, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":    job.ID,
		"status": bson.M{"$in": bson.A{StatusPending, StatusOngoing}},
	}, bson.M{"$set": bson.M{
		"status":   StatusOngoing,
		"launched": time.Now(),
	}})
	if err != nil {
		return false, err
	}

	return res.MatchedCount > 0, nil
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

	return res.MatchedCount > 0, nil
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

	return res.MatchedCount > 0, nil
}

func (r *mongoReporter) GetStatus(ctx context.Context, id string) (ImportJob, error) {
	var status ImportJob
	res := r.collection.FindOne(ctx, bson.M{"_id": id})
	err := res.Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = ErrNotFound
		}
	} else {
		err = res.Decode(&status)
	}

	return status, err
}

func (r *mongoReporter) GetAbandoned(ctx context.Context, createdInterval, launchedInterval time.Duration) ([]ImportJob, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"$or": []bson.M{
		{
			"status":   StatusPending,
			"creation": bson.M{"$lt": time.Now().Add(-createdInterval)},
		},
		{
			"status":   StatusOngoing,
			"launched": bson.M{"$lt": time.Now().Add(-launchedInterval)},
		},
	}})
	if err != nil {
		return nil, err
	}

	res := make([]ImportJob, 0)
	err = cursor.All(ctx, &res)
	return res, err
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
