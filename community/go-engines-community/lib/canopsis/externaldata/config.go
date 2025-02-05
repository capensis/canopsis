package externaldata

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const maxMongoDocsToCheck = 10

func SyncMongoCollections(ctx context.Context, client mongo.DbClient, collectionNames []string, logger zerolog.Logger) error {
	dbCollection := client.Collection(mongo.ExternalDataTableCollection)
	d, err := dbCollection.DeleteMany(ctx, bson.M{
		"name":        bson.M{"$nin": collectionNames},
		"from_config": true,
	})
	if err != nil {
		return fmt.Errorf("could not remove from list external data tables: %w", err)
	}

	unprocessedCollNames, err := findUnprocessedTables(ctx, dbCollection, collectionNames)
	if err != nil {
		return err
	}

	now := datetime.NewCpsTime()
	docs := make([]any, 0, len(unprocessedCollNames))
	for _, collName := range unprocessedCollNames {
		cursor, err := client.Collection(collName).Find(ctx, bson.M{}, options.Find().SetLimit(maxMongoDocsToCheck))
		if err != nil {
			return fmt.Errorf("failed to read collection %q: %w", collName, err)
		}

		columnTypes := make(map[string]int)
		nonStrFields := make([]string, 0)
		for cursor.Next(ctx) {
			var collDoc map[string]any
			if err := cursor.Decode(&collDoc); err != nil {
				return fmt.Errorf("failed to decode doc: %w", err)
			}

			for f, v := range collDoc {
				if f == IDColumnName {
					continue
				}

				if _, ok := v.(string); !ok {
					nonStrFields = append(nonStrFields, f)
					continue
				}

				columnTypes[f] = ColumnTypeNoType
			}
		}

		if err = cursor.Err(); err != nil {
			return fmt.Errorf("failed to fetch docs from collection: %w", err)
		}

		if err = cursor.Close(ctx); err != nil {
			return fmt.Errorf("failed to close cursor: %w", err)
		}

		if len(nonStrFields) > 0 {
			logger.Error().Str("collection_name", collName).Strs("fields", nonStrFields).Msg("MongoDB collection contains non string fields")
			continue
		}

		if len(columnTypes) == 0 {
			logger.Error().Str("collection_name", collName).Msg("MongoDB collection does not exist or empty")
			continue
		}

		docs = append(docs, Table{
			ID:          utils.NewID(),
			Type:        TypeMongoDB,
			Name:        collName,
			ColumnTypes: columnTypes,
			FromConfig:  true,
			Created:     now,
			Updated:     now,
		})
	}

	if len(docs) > 0 {
		_, err = dbCollection.InsertMany(ctx, docs)
		if err != nil {
			return fmt.Errorf("failed to insert external data tables: %w", err)
		}
	}

	logger.Info().
		Int64("deleted", d).
		Int("created", len(docs)).
		Int("unmodified", len(collectionNames)-len(unprocessedCollNames)).
		Msg("external data tables successfully updated")

	return nil
}

func findUnprocessedTables(ctx context.Context, dbCollection mongo.DbCollection, names []string) ([]string, error) {
	namesMap := make(map[string]struct{}, len(names))
	for _, name := range names {
		namesMap[name] = struct{}{}
	}

	cursor, err := dbCollection.Find(ctx, bson.M{
		"name":        bson.M{"$in": names},
		"from_config": true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find external data tables: %w", err)
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		v := Table{}
		if err = cursor.Decode(&v); err != nil {
			return nil, fmt.Errorf("failed to decode external data table: %w", err)
		}

		delete(namesMap, v.Name)
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("failed to fetch external data tables: %w", err)
	}

	res := make([]string, len(namesMap))
	i := 0
	for v := range namesMap {
		res[i] = v
		i++
	}

	return res, nil
}
