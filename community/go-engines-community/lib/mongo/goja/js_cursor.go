package goja

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/dop251/goja"
)

type jsCursor struct {
	dbCursor mongo.Cursor
	vm       *goja.Runtime
}

func (c *jsCursor) ForEach(ctx context.Context, f func(call goja.FunctionCall) goja.Value) error {
	for c.dbCursor.Next(ctx) {
		doc := make(map[string]any)
		err := c.dbCursor.Decode(&doc)
		if err != nil {
			return fmt.Errorf("cursor decoding failed: %w", err)
		}

		f(goja.FunctionCall{
			Arguments: []goja.Value{
				c.vm.ToValue(doc),
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
