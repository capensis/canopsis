package db

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func AggregateOne[T any](
	ctx context.Context,
	dbCollection mongo.DbCollection,
	pipeline any,
	opts ...options.Lister[options.AggregateOptions],
) (T, error) {
	var res T
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
