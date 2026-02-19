package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Cursor interface {
	All(ctx context.Context, results any) error
	Next(ctx context.Context) bool
	Close(ctx context.Context) error
	Decode(val any) error
	Err() error
}

type cursor struct {
	mongoCursor *mongo.Cursor
}

func (curs *cursor) All(ctx context.Context, results any) error {
	return curs.mongoCursor.All(ctx, results)
}

func (curs *cursor) Next(ctx context.Context) bool {
	return curs.mongoCursor.Next(ctx)
}

func (curs *cursor) Close(ctx context.Context) error {
	return curs.mongoCursor.Close(ctx)
}

func (curs *cursor) Decode(val any) error {
	return curs.mongoCursor.Decode(val)
}

func (curs *cursor) Err() error {
	return curs.mongoCursor.Err()
}
