package fixtures

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type LoadConfig struct {
	CollectionName string
	File           string
	Data           interface{}
}

func Load(ctx context.Context, client mongo.DbClient, dirs []string) error {
	loadConfigs, err := getConfigs(dirs)
	if err != nil {
		return err
	}

	for _, loadConfig := range loadConfigs {
		_, err := client.Collection(loadConfig.CollectionName).DeleteMany(ctx, bson.M{})
		if err != nil {
			return err
		}

		filename := loadConfig.File
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			// skip import when no fixture found
			continue
		}
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		var data []interface{}

		err = json.Unmarshal(file, &data)
		if err != nil {
			return err
		}

		if len(data) > 0 {
			_, err = client.Collection(loadConfig.CollectionName).InsertMany(ctx, data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getConfigs(dirs []string) ([]LoadConfig, error) {
	configs := make([]LoadConfig, 0)
	re := regexp.MustCompile(`^([a-z_]+)\.json$`)

	for _, dirPath := range dirs {
		files, err := ioutil.ReadDir(dirPath)
		if err != nil {
			return nil, err
		}

		for _, fileInfo := range files {
			filename := fileInfo.Name()
			matches := re.FindStringSubmatch(filename)
			if len(matches) < 2 {
				continue
			}

			collection := matches[1]
			configs = append(configs, LoadConfig{
				CollectionName: collection,
				File:           filepath.Join(dirPath, filename),
			})
		}
	}

	return configs, nil
}
