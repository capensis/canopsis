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

func NewSkipCmd(
	path string,
	version string,
	client mongo.DbClient,
	logger zerolog.Logger,
) Cmd {
	return &skipCmd{
		path:       path,
		version:    version,
		collection: client.Collection(collectionName),
		logger:     logger,
	}
}

type skipCmd struct {
	path       string
	version    string
	collection mongo.DbCollection
	logger     zerolog.Logger
}

func (c *skipCmd) Exec(ctx context.Context) error {
	files, err := filepath.Glob(filepath.Join(c.path, "*"+fileNameSuffixUp))
	if err != nil {
		return fmt.Errorf("cannot read directory %q: %w", c.path, err)
	}

	ids := make([]string, 0)
	found := false

	for _, file := range files {
		id := strings.TrimSuffix(filepath.Base(file), fileNameSuffixUp)
		ids = append(ids, id)

		if c.version != "" && id == c.version {
			found = true
			break
		}
	}

	if c.version != "" && !found {
		return fmt.Errorf("unknown migration %q", c.version)
	}

	if len(ids) == 0 {
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
		_, err = c.collection.InsertOne(ctx, bson.M{"_id": id})
		if err != nil {
			return fmt.Errorf("cannot update migration history: %w", err)
		}

		c.logger.Info().Str("file", filepath.Base(file)).Msg("up migration script skipped")
	}

	return nil
}

func (c *skipCmd) findMigrations(ctx context.Context, ids []string) (map[string]bool, error) {
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
