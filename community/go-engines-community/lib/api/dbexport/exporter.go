package dbexport

import (
	"bytes"
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/valyala/fastjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

type Exporter interface {
	Export(ctx context.Context, collection string, r Request) ([]byte, error)
}

type exporter struct {
	dbClient mongo.DbClient
}

func NewExporter(dbClient mongo.DbClient) Exporter {
	return &exporter{
		dbClient: dbClient,
	}
}

func (e *exporter) Export(ctx context.Context, collection string, r Request) ([]byte, error) {
	cursor, err := e.dbClient.Collection(collection).Find(ctx, bson.M{"_id": bson.M{"$in": r.IDs}})
	if err != nil {
		return nil, fmt.Errorf("failed to find documents from %s with ids = %v: %w", collection, r.IDs, err)
	}

	defer cursor.Close(ctx)

	buf := new(bytes.Buffer)
	vw, err := bsonrw.NewExtJSONValueWriter(buf, true, false)
	if err != nil {
		panic(err)
	}

	encoder, err := bson.NewEncoder(vw)
	if err != nil {
		panic(err)
	}

	// bson.MarshalExtJSON doesn't support arrays, using fastjson to help with that.
	var arena fastjson.Arena
	arr := arena.NewArray()

	idx := 0

	for cursor.Next(ctx) {
		var doc bson.Raw

		err = cursor.Decode(&doc)
		if err != nil {
			return nil, fmt.Errorf("failed to decode document from %s: %w", collection, err)
		}

		err = encoder.Encode(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal raw document from %s to extended json: %w", collection, err)
		}

		obj, err := fastjson.ParseBytes(buf.Bytes())
		if err != nil {
			return nil, fmt.Errorf("failed to parse ext json with fastjson: %w", err)
		}

		arr.SetArrayItem(idx, obj)
		buf.Reset()

		idx++
	}

	err = cursor.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to process mognodb cursor correctly: %w", err)
	}

	return arr.MarshalTo(nil), nil
}
