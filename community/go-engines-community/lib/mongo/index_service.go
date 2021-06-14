package mongo

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// IndexService is used to implement mongo indexes creations.
// Base implementation uses config configDir as source of index options.
type IndexService interface {
	Create() error
}

// NewIndexService creates index service.
func NewIndexService(
	dbClient DbClient,
	configDir string,
	logger *zerolog.Logger,
) IndexService {
	return &baseIndexService{
		dbClient:  dbClient,
		configDir: configDir,
		logger:    logger,
	}
}

type baseIndexService struct {
	dbClient  DbClient
	configDir string
	logger    *zerolog.Logger
}

// Config represents format of config file.
type Config struct {
	// Indexes structure is [collectionName][indexName]indexConfig
	Indexes map[string]map[string]IndexConfig `yaml:"indexes"`
}

// IndexConfig represent format of index in config file.
type IndexConfig struct {
	Keys    map[string]int         `yaml:"keys"`
	Options map[string]interface{} `yaml:"options"`
}

func (s *baseIndexService) Create() error {
	ch, err := s.parseConfigDir()

	if err != nil {
		return err
	}

	for config := range ch {
		s.createIndexes(config)
	}

	return nil
}

// parseConfigDir discovers all .yml file in config dir and parses them.
func (s *baseIndexService) parseConfigDir() (<-chan *Config, error) {
	fileInfoList, err := ioutil.ReadDir(s.configDir)

	if err != nil {
		s.logger.
			Error().
			Err(err).
			Str("configDir", s.configDir).
			Msg("cannot read config dir")

		return nil, err
	}

	s.logger.
		Debug().
		Str("configDir", s.configDir).
		Msg("read config dir")
	ch := make(chan *Config)

	go func() {
		defer close(ch)
		for _, fileInfo := range fileInfoList {
			if fileInfo.Mode().IsRegular() && strings.HasSuffix(fileInfo.Name(), ".yml") {
				path := filepath.Join(s.configDir, fileInfo.Name())
				c, err := s.parseConfigFile(path)

				if err != nil {
					continue
				}

				ch <- c
			}
		}
	}()

	return ch, nil
}

// parseConfigFile reads config file and parses it's content.
func (s *baseIndexService) parseConfigFile(path string) (*Config, error) {
	buf, err := ioutil.ReadFile(path)

	if err != nil {
		s.logger.
			Error().
			Err(err).
			Str("path", path).
			Msg("cannot read config file")

		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(buf, &config)

	if err != nil {
		s.logger.
			Error().
			Err(err).
			Str("path", path).
			Msg("cannot parse config file")

		return nil, err
	}

	s.logger.
		Debug().
		Str("path", path).
		Msg("read config file")

	return &config, nil
}

// createIndexes creates mongo indexes by config.
func (s *baseIndexService) createIndexes(config *Config) {
	for collectionName, indexes := range config.Indexes {
		collection := s.dbClient.Collection(collectionName)

		for indexName, params := range indexes {
			indexOptions := mongooptions.IndexOptions{}
			err := mapstructure.Decode(params.Options, &indexOptions)

			if err != nil {
				s.logger.
					Error().
					Err(err).
					Str("collection", collectionName).
					Str("index", indexName).
					Msg("cannot parse index options")

				continue
			}

			indexOptions.Name = &indexName
			ctx := context.Background()
			_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys:    &params.Keys,
				Options: &indexOptions,
			})

			if err != nil {
				s.logger.
					Error().
					Err(err).
					Str("collection", collectionName).
					Str("index", indexName).
					Msg("cannot create index")
			} else {
				s.logger.
					Info().
					Str("collection", collectionName).
					Str("index", indexName).
					Msg("create index")
			}
		}
	}
}
