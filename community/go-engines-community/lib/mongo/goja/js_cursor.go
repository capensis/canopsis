package goja

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/dop251/goja"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type jsCursor struct {
	dbCursor mongo.Cursor
	vm       *goja.Runtime
}

func (c *jsCursor) ForEach(ctx context.Context, f func(call goja.FunctionCall) goja.Value) error {
	for c.dbCursor.Next(ctx) {
		var doc bson.M
		err := c.dbCursor.Decode(&doc)
		if err != nil {
			return fmt.Errorf("cursor decoding failed: %w", err)
		}

		arg := c.bsonToMap(doc)
		f(goja.FunctionCall{
			Arguments: []goja.Value{
				c.vm.ToValue(arg),
			},
		})
	}

	return nil
}

func (c *jsCursor) Close(ctx context.Context) error {
	return c.dbCursor.Close(ctx)
}

func (c *jsCursor) HasNext(ctx context.Context) bool {
	return c.dbCursor.Next(ctx)
}

func (c *jsCursor) Next() (map[string]any, error) {
	doc := make(map[string]any)
	err := c.dbCursor.Decode(&doc)
	if err != nil {
		return nil, fmt.Errorf("cursor decoding failed: %w", err)
	}

	return doc, nil
}

func (c *jsCursor) getMethods(ctx context.Context) map[string]any {
	return map[string]any{
		"forEach": func(f func(call goja.FunctionCall) goja.Value) error {
			return c.ForEach(ctx, f)
		},
		"close": func() error {
			return c.Close(ctx)
		},
		"hasNext": func() bool {
			return c.HasNext(ctx)
		},
		"next": func() (map[string]any, error) {
			return c.Next()
		},
	}
}

func (c *jsCursor) bsonToMap(b any) any {
	switch bv := b.(type) {
	case bson.M:
		m := make(map[string]any, len(bv))
		for k, v := range bv {
			m[k] = c.bsonToMap(v)
		}

		return m
	case bson.D:
		m := make(map[string]any, len(bv))
		for _, v := range bv {
			m[v.Key] = c.bsonToMap(v.Value)
		}

		return m
	case bson.A:
		s := make([]any, len(bv))
		for i, v := range bv {
			s[i] = c.bsonToMap(v)
		}

		return s
	default:
		return b
	}
}
