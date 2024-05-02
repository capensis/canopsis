package cli

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

func NewUpCmd(
	path, to string,
	client mongo.DbClient,
	scriptExecutor mongo.ScriptExecutor,
	logger zerolog.Logger,
) Cmd {
	return &upCmd{
		path:           path,
		to:             to,
		collection:     client.Collection(collectionName),
		scriptExecutor: scriptExecutor,
		logger:         logger,
	}
}

type upCmd struct {
	path, to       string
	collection     mongo.DbCollection
	scriptExecutor mongo.ScriptExecutor
	logger         zerolog.Logger
}

func (c *upCmd) Exec(ctx context.Context) error {
	files, err := filepath.Glob(filepath.Join(c.path, "*"+fileNameSuffixUp))
	if err != nil {
		return fmt.Errorf("cannot read directory %q: %w", c.path, err)
	}

	ids := make([]string, 0)
	found := false

	for _, file := range files {
		id := strings.TrimSuffix(filepath.Base(file), fileNameSuffixUp)
		ids = append(ids, id)

		if c.to != "" && id == c.to {
			found = true
			break
		}
	}

	if c.to != "" && !found {
		return fmt.Errorf("unknown migration %q", c.to)
	}

	if len(ids) == 0 {
		if len(files) != 0 {
			return fmt.Errorf("no migration files found in %q, total files %d", c.path, len(files))
		}
		return nil
	}

	exists, err := c.findMigrations(ctx, ids)
	if err != nil {
		return err
	}

	for _, id := range ids {
		if exists[id] {
			continue
		}

		file := filepath.Join(c.path, id+fileNameSuffixUp)
		err = c.scriptExecutor.Exec(file)
		if err != nil {
			return err
		}

		_, err = c.collection.InsertOne(ctx, bson.M{"_id": id})
		if err != nil {
			return fmt.Errorf("cannot update migration history: %w", err)
		}

		c.logger.Info().Str("file", filepath.Base(file)).Msg("up migration script executed")
	}

	return nil
}

func (c *upCmd) findMigrations(ctx context.Context, ids []string) (map[string]bool, error) {
	data := struct {
		ID string `bson:"_id"`
	}{}
	cursor, err := c.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, fmt.Errorf("cannot fetch migrations: %w", err)
	}

	res := make(map[string]bool)
	for cursor.Next(ctx) {
		err := cursor.Decode(&data)
		if err != nil {
			return nil, fmt.Errorf("cannot decode migration: %w", err)
		}

		res[data.ID] = true
	}

	return res, nil
}
