package fixtures

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type LoadConfig struct {
	CollectionName string
	File           string
	Data           interface{}
}

func LoadFixtures(ctx context.Context, client mongo.DbClient, loadConfigs ...LoadConfig) error {
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
