package db

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func AggregateIDs(
	ctx context.Context,
	dbCollection mongo.DbCollection,
	pipeline any,
	opts ...options.Lister[options.AggregateOptions],
) (ids []string, err error) {
	cursor, err := dbCollection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var doc struct {
			ID string `bson:"_id"`
		}
		err = cursor.Decode(&doc)
		if err != nil {
			return nil, err
		}

		ids = append(ids, doc.ID)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func AggregateOne[T any](
	ctx context.Context,
	dbCollection mongo.DbCollection,
	pipeline any,
	opts ...options.Lister[options.AggregateOptions],
) (res T, err error) {
	cursor, err := dbCollection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return res, err
	}

	if !cursor.Next(ctx) {
		return res, httperror.ErrNotFound
	}

	err = cursor.Decode(&res)
	if err != nil {
		return res, err
	}

	if err = cursor.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func AggregateAll[T any](
	ctx context.Context,
	dbCollection mongo.DbCollection,
	pipeline any,
	opts ...options.Lister[options.AggregateOptions],
) (res []T, err error) {
	cursor, err := dbCollection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return res, err
	}

	err = cursor.All(ctx, &res)

	return res, err
}
