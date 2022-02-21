package fixtures

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	filePattern = "*.yml"
)

type Loader interface {
	Load(ctx context.Context) error
}

func NewLoader(
	client mongo.DbClient,
	dirs []string,
	deleteAllBefore bool,
	parser Parser,
	logger zerolog.Logger,
) Loader {
	return &loader{
		client:          client,
		dirs:            dirs,
		deleteAllBefore: deleteAllBefore,
		parser:          parser,
		logger:          logger,
	}
}

type loader struct {
	client          mongo.DbClient
	dirs            []string
	deleteAllBefore bool
	parser          Parser
	logger          zerolog.Logger
}

func (l *loader) Load(ctx context.Context) error {
	deleted := make(map[string]bool)

	for _, dir := range l.dirs {
		files, err := filepath.Glob(filepath.Join(dir, filePattern))
		if err != nil {
			return fmt.Errorf("cannot read dir %q: %w", dir, err)
		}

		for _, filename := range files {
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("cannot read file %q: %w", filename, err)
			}

			docsByCollection, err := l.parser.Parse(content)
			if err != nil {
				return fmt.Errorf("cannot parse file %q: %w", filename, err)
			}

			for collectionName, docs := range docsByCollection {
				if len(docs) == 0 {
					continue
				}

				deleted, err = l.deleteAll(ctx, docsByCollection, deleted)
				if err != nil {
					return err
				}

				collection := l.client.Collection(collectionName)
				_, err = collection.InsertMany(ctx, docs)
				if err != nil {
					return fmt.Errorf("cannot save documents from file %q: %w", filename, err)
				}
			}
		}
	}

	return nil
}

func (l *loader) deleteAll(ctx context.Context, docs map[string][]interface{}, deleted map[string]bool) (map[string]bool, error) {
	if !l.deleteAllBefore {
		return deleted, nil
	}

	for collectionName := range docs {
		if deleted[collectionName] {
			continue
		}

		_, err := l.client.Collection(collectionName).DeleteMany(ctx, bson.M{})
		if err != nil {
			return nil, fmt.Errorf("cannot delete collection %q: %w", collectionName, err)
		}

		deleted[collectionName] = true
	}

	return deleted, nil
}
