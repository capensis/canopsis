package fixtures

import (
	"context"
	"fmt"
	"os"
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
	Clean(ctx context.Context) error
}

func NewLoader(
	client mongo.DbClient,
	dirs []string,
	parser Parser,
	logger zerolog.Logger,
) Loader {
	return &loader{
		client:   client,
		dirs:     dirs,
		parser:   parser,
		keepData: false,
		logger:   logger,
	}
}

func NewLoaderWithKeepData(
	client mongo.DbClient,
	dirs []string,
	parser Parser,
	logger zerolog.Logger,
) Loader {
	return &loader{
		client:   client,
		dirs:     dirs,
		parser:   parser,
		keepData: true,
		logger:   logger,
	}
}

type loader struct {
	client   mongo.DbClient
	dirs     []string
	parser   Parser
	keepData bool
	logger   zerolog.Logger

	collections []string
}

func (l *loader) Load(ctx context.Context) error {
	deleted := make(map[string]bool)

	for _, dir := range l.dirs {
		files, err := getFiles(dir)
		if err != nil {
			return err
		}

		for _, filename := range files {
			content, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("cannot read file %q: %w", filename, err)
			}

			docsByCollection, err := l.parser.Parse(content)
			if err != nil {
				return fmt.Errorf("cannot parse file %q: %w", filename, err)
			}

			for collectionName, docs := range docsByCollection {
				if !l.keepData {
					if !deleted[collectionName] {
						_, err := l.client.Collection(collectionName).DeleteMany(ctx, bson.M{})
						if err != nil {
							return fmt.Errorf("cannot delete collection %q: %w", collectionName, err)
						}

						deleted[collectionName] = true
						l.collections = append(l.collections, collectionName)
					}
				}

				if len(docs) > 0 {
					collection := l.client.Collection(collectionName)
					_, err = collection.InsertMany(ctx, docs)
					if err != nil {
						return fmt.Errorf("cannot save documents from file %q: %w", filename, err)
					}
				}
			}
		}
	}

	return nil
}

func (l *loader) Clean(ctx context.Context) error {
	for _, collectionName := range l.collections {
		_, err := l.client.Collection(collectionName).DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("cannot delete collection %q: %w", collectionName, err)
		}
	}

	return nil
}

func getFiles(dir string) ([]string, error) {
	stat, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read dir %q: %w", dir, err)
	}

	if stat.IsDir() {
		return filepath.Glob(filepath.Join(dir, filePattern))
	}

	return []string{dir}, nil
}
