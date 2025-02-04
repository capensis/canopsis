package externaldata

import (
	"context"
	"fmt"

	apicommon "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const maxMongoDocsToCheck = 10

func SyncMongoCollections(
	ctx context.Context,
	client mongo.DbClient,
	collectionNames []string,
	refCollectionNames []string,
	logger zerolog.Logger,
) error {
	collection := client.Collection(mongo.ExternalDataTableCollection)
	delCount, delErrCount, err := deleteMissingTables(ctx, client, collection, collectionNames, refCollectionNames, logger)
	if err != nil {
		return err
	}

	if len(collectionNames) > 0 {
		_, err = collection.UpdateMany(ctx, bson.M{"name": bson.M{"$in": collectionNames}, "from_config": true},
			bson.M{"$unset": bson.M{"removed_from_config": ""}})
		if err != nil {
			return fmt.Errorf("failed to update external data tables: %w", err)
		}
	}

	unprocessedCollNames, nonConfErrCount, err := findUnprocessedTables(ctx, collection, collectionNames, logger)
	if err != nil {
		return err
	}

	newCount, newErrCount, err := insertNewTables(ctx, client, collection, unprocessedCollNames, logger)
	if err != nil {
		return err
	}

	logger.Info().
		Int("deleted", delCount).
		Int("deleted_errors", delErrCount).
		Int("created", newCount).
		Int("created_errors", newErrCount+nonConfErrCount).
		Int("unmodified", len(collectionNames)-len(unprocessedCollNames)-nonConfErrCount).
		Msg("external data tables successfully updated")

	return nil
}

func insertNewTables(
	ctx context.Context,
	client mongo.DbClient,
	collection mongo.DbCollection,
	unprocessedCollNames []string,
	logger zerolog.Logger,
) (int, int, error) {
	now := datetime.NewCpsTime()
	docs := make([]any, 0, len(unprocessedCollNames))
	errCount := 0
	for _, collName := range unprocessedCollNames {
		cursor, err := client.Collection(collName).Find(ctx, bson.M{}, options.Find().SetLimit(maxMongoDocsToCheck))
		if err != nil {
			return 0, 0, fmt.Errorf("failed to read collection %q: %w", collName, err)
		}

		var columns []string
		var columnTypes []int
		var hasCol map[string]bool
		invalidFields := make([]string, 0)
		nonStrFields := make([]string, 0)
		for cursor.Next(ctx) {
			var collDoc bson.D
			if err := cursor.Decode(&collDoc); err != nil {
				return 0, 0, fmt.Errorf("failed to decode doc: %w", err)
			}

			if columns == nil {
				columns = make([]string, 0, len(collDoc)-1)
				columnTypes = make([]int, 0, len(collDoc)-1)
				hasCol = make(map[string]bool, len(collDoc)-1)
			}

			for _, f := range collDoc {
				if f.Key == IDColumnName {
					continue
				}

				if !hasCol[f.Key] && !apicommon.IsTableName(f.Key) {
					invalidFields = append(invalidFields, f.Key)
					continue
				}

				if _, ok := f.Value.(string); !ok {
					nonStrFields = append(nonStrFields, f.Key)
					continue
				}

				if !hasCol[f.Key] {
					hasCol[f.Key] = true
					columns = append(columns, f.Key)
					columnTypes = append(columnTypes, ColumnTypeNoType)
				}
			}
		}

		if err = cursor.Err(); err != nil {
			return 0, 0, fmt.Errorf("failed to fetch docs from collection: %w", err)
		}

		if err = cursor.Close(ctx); err != nil {
			return 0, 0, fmt.Errorf("failed to close cursor: %w", err)
		}

		if len(invalidFields) > 0 {
			errCount++
			logger.Error().Str("collection_name", collName).Strs("fields", invalidFields).Msg("MongoDB collection contains fields with invalid names")
			continue
		}

		if len(nonStrFields) > 0 {
			errCount++
			logger.Error().Str("collection_name", collName).Strs("fields", nonStrFields).Msg("MongoDB collection contains non string fields")
			continue
		}

		if len(columnTypes) == 0 {
			errCount++
			logger.Error().Str("collection_name", collName).Msg("MongoDB collection does not exist or empty")
			continue
		}

		docs = append(docs, Table{
			ID:          utils.NewID(),
			Type:        TypeMongoDB,
			Name:        collName,
			Columns:     columns,
			ColumnTypes: columnTypes,
			FromConfig:  true,
			Created:     now,
			Updated:     now,
		})
	}

	count := 0
	if len(docs) > 0 {
		r, err := collection.InsertMany(ctx, docs)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to insert external data tables: %w", err)
		}

		count = len(r)
	}

	return count, errCount, nil
}

func findUnprocessedTables(ctx context.Context, collection mongo.DbCollection, names []string, logger zerolog.Logger) ([]string, int, error) {
	if len(names) == 0 {
		return nil, 0, nil
	}

	namesMap := make(map[string]struct{}, len(names))
	for _, name := range names {
		namesMap[name] = struct{}{}
	}

	cursor, err := collection.Find(ctx, bson.M{
		"name": bson.M{"$in": names},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find external data tables: %w", err)
	}

	nonConfCount := 0
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		v := Table{}
		if err = cursor.Decode(&v); err != nil {
			return nil, 0, fmt.Errorf("failed to decode external data table: %w", err)
		}

		delete(namesMap, v.Name)
		if !v.FromConfig {
			nonConfCount++
			logger.Error().Str("collection_name", v.Name).Msg("MongoDB collection cannot be added to external data list because it's already created by API")
		}
	}

	if err = cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to fetch external data tables: %w", err)
	}

	res := make([]string, len(namesMap))
	i := 0
	for v := range namesMap {
		res[i] = v
		i++
	}

	return res, nonConfCount, nil
}

func deleteMissingTables(
	ctx context.Context,
	client mongo.DbClient,
	collection mongo.DbCollection,
	collectionNames []string,
	refCollectionNames []string,
	logger zerolog.Logger,
) (int, int, error) {
	match := bson.M{
		"from_config": true,
	}
	if len(collectionNames) > 0 {
		match["name"] = bson.M{"$nin": collectionNames}
	}

	cursor, err := collection.Find(ctx, match, options.Find().SetProjection(bson.M{"_id": 1, "name": 1}))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to find external data tables to removed: %w", err)
	}

	defer cursor.Close(ctx)
	ids := make([]string, 0)
	names := make(map[string]string)
	for cursor.Next(ctx) {
		t := Table{}
		err := cursor.Decode(&t)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to decode external data table: %w", err)
		}

		ids = append(ids, t.ID)
		names[t.ID] = t.Name
	}

	if err = cursor.Err(); err != nil {
		return 0, 0, fmt.Errorf("failed to fetch external data tables to removed: %w", err)
	}

	if len(ids) == 0 {
		return 0, 0, nil
	}

	linked := make(map[string]bool, len(ids))
	linked, err = findWidgetLinkedTables(ctx, client.Collection(mongo.WidgetMongoCollection), ids, linked)
	if err != nil {
		return 0, 0, err
	}

	for _, c := range refCollectionNames {
		linked, err = findRuleLinkedTables(ctx, client.Collection(c), ids, linked)
		if err != nil {
			return 0, 0, err
		}
	}

	i := 0
	blockedIDs := make([]string, 0)
	for _, id := range ids {
		if linked[id] {
			logger.Error().Str("collection_name", names[id]).Msg("MongoDB collection cannot be removed from external data list because it's used in rules")
			blockedIDs = append(blockedIDs, id)
			continue
		}

		ids[i] = id
		i++
	}

	ids = ids[:i]
	count := 0
	if len(ids) > 0 {
		d, err := collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
		if err != nil {
			return 0, 0, fmt.Errorf("failed to delete external data tables: %w", err)
		}

		count = int(d)
	}

	errCount := 0
	if len(blockedIDs) > 0 {
		r, err := collection.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": blockedIDs}}, bson.M{"$set": bson.M{
			"removed_from_config": true,
		}})
		if err != nil {
			return 0, 0, fmt.Errorf("failed to update external data tables: %w", err)
		}

		errCount = int(r.MatchedCount)
	}

	return count, errCount, nil
}

func findWidgetLinkedTables(ctx context.Context, dbWidgetCollection mongo.DbCollection, ids []string, linked map[string]bool) (map[string]bool, error) {
	if len(ids) == 0 {
		return linked, nil
	}

	cursor, err := dbWidgetCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"type":             view.WidgetTypeExternalData,
			"parameters.table": bson.M{"$in": ids},
		}},
		{"$group": bson.M{
			"_id": "$parameters.table",
		}},
	})
	if err != nil {
		return linked, fmt.Errorf("failed to find linked rules: %w", err)
	}

	for cursor.Next(ctx) {
		r := struct {
			ID string `bson:"_id"`
		}{}
		err := cursor.Decode(&r)
		if err != nil {
			return linked, fmt.Errorf("failed to decode linked rule: %w", err)
		}

		linked[r.ID] = true
	}

	if err = cursor.Err(); err != nil {
		return linked, fmt.Errorf("failed to fetch linked rules: %w", err)
	}

	if err = cursor.Close(ctx); err != nil {
		return linked, fmt.Errorf("failed to close rule cursor: %w", err)
	}

	return linked, nil
}

func findRuleLinkedTables(ctx context.Context, dbRuleCollection mongo.DbCollection, ids []string, linked map[string]bool) (map[string]bool, error) {
	unlinkedIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		if !linked[id] {
			unlinkedIDs = append(unlinkedIDs, id)
		}
	}

	if len(unlinkedIDs) == 0 {
		return linked, nil
	}

	cursor, err := dbRuleCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"external_data.table": bson.M{"$in": unlinkedIDs}}},
		{"$unwind": "$external_data"},
		{"$match": bson.M{"external_data.table": bson.M{"$in": unlinkedIDs}}},
		{"$group": bson.M{
			"_id": "$external_data.table",
		}},
	})
	if err != nil {
		return linked, fmt.Errorf("failed to find linked rules: %w", err)
	}

	for cursor.Next(ctx) {
		r := struct {
			ID string `bson:"_id"`
		}{}
		err := cursor.Decode(&r)
		if err != nil {
			return linked, fmt.Errorf("failed to decode linked rule: %w", err)
		}

		linked[r.ID] = true
	}

	if err = cursor.Err(); err != nil {
		return linked, fmt.Errorf("failed to fetch linked rules: %w", err)
	}

	if err = cursor.Close(ctx); err != nil {
		return linked, fmt.Errorf("failed to close rule cursor: %w", err)
	}

	return linked, nil
}
