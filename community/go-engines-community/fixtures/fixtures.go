package fixtures

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type LoadConfig struct {
	CollectionName string
	File           string
	Data           interface{}
}

func LoadFixtures(loadConfigs ...LoadConfig) error {
	client, err := mongo.NewClient(0, 0)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
