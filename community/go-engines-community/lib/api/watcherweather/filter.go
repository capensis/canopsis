package watcherweather

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
)

func ParseFilter(filter string) (bson.M, error) {
	if filter == "" {
		return nil, nil
	}

	var parsedFilter bson.M
	err := json.Unmarshal([]byte(filter), &parsedFilter)
	if err != nil {
		return nil, err
	}

	return parsedFilter, nil
}
