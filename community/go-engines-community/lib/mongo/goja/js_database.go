package goja

import (
	"context"
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/dop251/goja"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type jsDatabase struct {
	dbClient mongo.DbClient
	vm       *goja.Runtime
}

func (c *jsDatabase) CreateCollection(ctx context.Context, name string, opts goja.Value) (map[string]int, error) {
	dbOpts := options.CreateCollection()
	err := transformOptions(c.vm, opts, dbOpts, map[string]mappingFunc{
		"size": func(v any) error {
			i, ok := v.(int64)
			if !ok {
				return errors.New("invalid type for size")
			}

			dbOpts.SetSizeInBytes(i)

			return nil
		},
		"max": func(v any) error {
			i, ok := v.(int64)
			if !ok {
				return errors.New("invalid type for max")
			}

			dbOpts.SetMaxDocuments(i)

			return nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("invalid create collection options: %w", err)
	}

	err = c.dbClient.CreateCollection(ctx, name, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("error creating collection %q: %w", name, err)
	}

	return map[string]int{"ok": 1}, nil
}

func (c *jsDatabase) GetCollectionNames(ctx context.Context, filter, opts goja.Value) ([]string, error) {
	dbFilter, err := transformValue(c.vm, filter)
	if err != nil {
		return nil, fmt.Errorf("invalid filter: %w", err)
	}

	if dbFilter == nil {
		dbFilter = bson.D{}
	}

	dbOpts := options.ListCollections()
	err = transformOptions(c.vm, opts, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("invalid list collections options: %w", err)
	}

	collectionNames, err := c.dbClient.ListCollectionNames(ctx, dbFilter, dbOpts)
	if err != nil {
		return nil, fmt.Errorf("error listing collections: %w", err)
	}

	return collectionNames, nil
}

func (c *jsDatabase) GetCollection(ctx context.Context, collectionName string) map[string]any {
	return (&jsCollection{
		dbClient:     c.dbClient,
		dbCollection: c.dbClient.Collection(collectionName),
		vm:           c.vm,
	}).getMethods(ctx)
}

func (c *jsDatabase) RunCommand(ctx context.Context, command goja.Value) (map[string]any, error) {
	dbCommand, err := transformValue(c.vm, command)
	if err != nil {
		return nil, fmt.Errorf("invalid command: %w", err)
	}

	res := make(map[string]any)
	err = c.dbClient.RunCommand(ctx, dbCommand).Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("error running command: %w", err)
	}

	return res, nil
}

func (c *jsDatabase) getMethods(ctx context.Context, collectionNames []string) map[string]any {
	methods := map[string]any{
		"createCollection": func(name string, opts goja.Value) (map[string]int, error) {
			return c.CreateCollection(ctx, name, opts)
		},
		"getCollectionNames": func(filter, opts goja.Value) ([]string, error) {
			return c.GetCollectionNames(ctx, filter, opts)
		},
		"getCollection": func(collectionName string) map[string]any {
			return c.GetCollection(ctx, collectionName)
		},
		"runCommand": func(command goja.Value) (map[string]any, error) {
			return c.RunCommand(ctx, command)
		},
	}
	for _, name := range collectionNames {
		methods[name] = c.GetCollection(ctx, name)
	}

	return methods
}
