package mongo

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
)

// IndexService is used to implement mongo indexes creations.
// Base implementation uses config configDir as source of index options.
type IndexService interface {
	Create(ctx context.Context) error
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
	Keys    IndexConfigKeys        `yaml:"keys"`
	Options map[string]interface{} `yaml:"options"`
}

type IndexConfigKeys struct {
	Keys   []string
	Values []int
}

func (k *IndexConfigKeys) UnmarshalYAML(value *yaml.Node) error {
	var err error

	if len(value.Content)%2 != 0 {
		return fmt.Errorf("content length should be an even number")
	}

	halfLen := len(value.Content) / 2

	k.Keys = make([]string, halfLen)
	k.Values = make([]int, halfLen)

	for i := 0; i < halfLen; i++ {
		k.Keys[i] = value.Content[2*i].Value
		k.Values[i], err = strconv.Atoi(value.Content[2*i+1].Value)

		if err != nil {
			return fmt.Errorf("failed to convert %s to int, error = %s", value.Content[2*i+1].Value, err)
		}
	}

	return nil
}

func (s *baseIndexService) Create(ctx context.Context) error {
	ch, err := s.parseConfigDir()

	if err != nil {
		return err
	}

	for config := range ch {
		s.createIndexes(ctx, config)
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

	// Complex index options are unmarshalled to map[interface{}]interface{} which
	// mongo package cannot marshal to bson with error
	// "cannot transform type map[interface{}]interface {} to a BSON Document: unsupported key type: interface {}".
	// Transform map[interface{}]interface{} to map[string]interface{} to avoid this error.
	for collection := range config.Indexes {
		for indexName := range config.Indexes[collection] {
			for option := range config.Indexes[collection][indexName].Options {
				config.Indexes[collection][indexName].Options[option] = transformInterfaceMapKeyToString(config.Indexes[collection][indexName].Options[option])
			}
		}
	}

	s.logger.
		Debug().
		Str("path", path).
		Msg("read config file")

	return &config, nil
}

// createIndexes creates mongo indexes by config.
func (s *baseIndexService) createIndexes(ctx context.Context, config *Config) {
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

			keysLen := len(params.Keys.Keys)

			indexKeys := make(bson.D, keysLen)
			for i := 0; i < keysLen; i++ {
				indexKeys[i] = bson.E{Key: params.Keys.Keys[i], Value: params.Keys.Values[i]}
			}

			indexOptions.Name = &indexName
			for i := 0; i < 2; i++ {
				_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
					Keys:    indexKeys,
					Options: &indexOptions,
				})

				if err != nil {
					if i == 0 && isIndexKeySpecsConflict(err) {
						if _, err = collection.Indexes().DropOne(ctx, *indexOptions.Name); err != nil {
							s.logger.
								Warn().
								Err(err).
								Str("collection", collectionName).
								Str("index", indexName).
								Msg("cannot drop index with conflicted specs")
						}
						continue
					}
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
				break
			}
		}
	}
}

func isIndexKeySpecsConflict(err error) bool {
	return strings.HasPrefix(err.Error(), "(IndexKeySpecsConflict)")
}

// transformInterfaceMapKeyToString replaces map[interface{}]interface{} to
// map[string]interface{} included nested maps.
func transformInterfaceMapKeyToString(i interface{}) interface{} {
	if i == nil {
		return i
	}

	if m, ok := i.(map[interface{}]interface{}); ok {
		strKeyMap := make(map[string]interface{}, len(m))
		for k, v := range m {
			strKeyMap[fmt.Sprintf("%v", k)] = transformInterfaceMapKeyToString(v)
		}

		return strKeyMap
	}

	return i
}
