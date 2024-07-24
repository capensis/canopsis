package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDownCmd(
	path, to string,
	client mongo.DbClient,
	scriptExecutor mongo.ScriptExecutor,
	logger zerolog.Logger,
) Cmd {
	return &downCmd{
		path:           path,
		to:             to,
		collection:     client.Collection(collectionName),
		scriptExecutor: scriptExecutor,
		logger:         logger,
	}
}

type downCmd struct {
	path, to       string
	collection     mongo.DbCollection
	scriptExecutor mongo.ScriptExecutor
	logger         zerolog.Logger
}

func (c *downCmd) Exec(ctx context.Context) error {
	ids, err := c.findMigrations(ctx, c.to)
	if err != nil {
		return err
	}

	if c.to != "" && ids[len(ids)-1] != c.to {
		return fmt.Errorf("unknown migration %q", c.to)
	}

	for _, id := range ids {
		file := filepath.Join(c.path, id+fileNameSuffixDown)
		_, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				c.logger.Error().Msgf("not found down migration script %q", file)

				_, err = c.collection.DeleteOne(ctx, bson.M{"_id": id})
				if err != nil {
					return fmt.Errorf("cannot update migration history: %w", err)
				}

				continue
			}

			return fmt.Errorf("cannot check file exist: %w", err)
		}

		err = c.scriptExecutor.Exec(ctx, file)
		if err != nil {
			return err
		}

		_, err = c.collection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			return fmt.Errorf("cannot update migration history: %w", err)
		}

		c.logger.Info().Str("file", filepath.Base(file)).Msg("down migration script executed")
	}

	return nil
}

func (c *downCmd) findMigrations(ctx context.Context, id string) ([]string, error) {
	data := struct {
		ID string `bson:"_id"`
	}{}
	filter := bson.M{}
	if id != "" {
		filter = bson.M{"_id": bson.M{"$gte": id}}
	}
	cursor, err := c.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"_id": -1}))
	if err != nil {
		return nil, fmt.Errorf("cannot fetch migrations: %w", err)
	}

	res := make([]string, 0)
	for cursor.Next(ctx) {
		err := cursor.Decode(&data)
		if err != nil {
			return nil, fmt.Errorf("cannot decode migration: %w", err)
		}

		res = append(res, data.ID)
	}

	return res, nil
}
