package export

import (
	"context"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/jackc/pgx/v5"
)

func NewMongoCursor(
	cursor mongo.Cursor,
	fields []string,
	transform func(k string, v any) any,
) DataCursor {
	return &mongoCursor{
		cursor:    cursor,
		fields:    fields,
		transform: transform,
	}
}

func NewPostgresCursor(
	cursor pgx.Rows,
	fields []string,
	transform func(k string, v any) any,
) DataCursor {
	return &postgresCursor{
		cursor:    cursor,
		fields:    fields,
		transform: transform,
	}
}

type mongoCursor struct {
	cursor    mongo.Cursor
	fields    []string
	transform func(k string, v any) any
}

func (c *mongoCursor) Next(ctx context.Context) bool {
	return c.cursor.Next(ctx)
}

func (c *mongoCursor) Scan(m *map[string]any) error {
	err := c.cursor.Decode(m)
	if err != nil {
		return err
	}

	*m = filterFields(*m, c.fields, c.transform)
	return nil
}

func (c *mongoCursor) Close(ctx context.Context) error {
	return c.cursor.Close(ctx)
}

type postgresCursor struct {
	cursor    pgx.Rows
	fields    []string
	transform func(k string, v any) any
}

func (c *postgresCursor) Next(_ context.Context) bool {
	return c.cursor.Next()
}

func (c *postgresCursor) Scan(m *map[string]any) error {
	*m = make(map[string]any, len(c.fields))
	v, err := c.cursor.Values()
	if err != nil {
		return err
	}

	for i, f := range c.fields {
		if c.transform == nil {
			(*m)[f] = v[i]
		} else {
			(*m)[f] = c.transform(f, v[i])
		}
	}

	return nil
}

func (c *postgresCursor) Close(_ context.Context) error {
	c.cursor.Close()

	return nil
}

func filterFields(
	m map[string]any,
	fields []string,
	transform func(k string, v any) any,
) map[string]any {
	res := make(map[string]any, len(fields))
	for _, field := range fields {
		v, ok := getNestedVal(m, strings.Split(field, "."))
		if !ok {
			continue
		}

		if transform == nil {
			res[field] = v
		} else {
			res[field] = transform(field, v)
		}
	}

	return res
}

func getNestedVal(m map[string]any, keys []string) (any, bool) {
	if len(keys) == 0 {
		return nil, false
	}

	if v, ok := m[keys[0]]; ok {
		if len(keys) == 1 {
			return v, true
		}

		if mv, ok := v.(map[string]any); ok {
			return getNestedVal(mv, keys[1:])
		}
	}

	return nil, false
}
