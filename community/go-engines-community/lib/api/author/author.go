package author

import (
	"encoding/json"
	"regexp"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Author struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
}

type Role struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
}

func Pipeline() []bson.M {
	return PipelineForField("author")
}

func PipelineForField(field string) []bson.M {
	tmpField := getTempField(field)
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.RightsMongoCollection,
			"localField":   field,
			"foreignField": "_id",
			"as":           tmpField,
		}},
		{"$unwind": bson.M{"path": "$" + tmpField, "preserveNullAndEmptyArrays": true}},
		// keep only _id and name of author struct
		{"$addFields": bson.M{
			field + "._id":  "$" + tmpField + "._id",
			field + ".name": "$" + tmpField + ".crecord_name",
		}},
		{"$project": bson.M{tmpField: 0}},
		{"$addFields": bson.M{
			field: bson.M{"$cond": bson.M{
				"if":   "$" + field + "._id",
				"then": "$" + field,
				"else": "$$REMOVE",
			}},
		}},
	}
}

// getTempField returns a field name with "_" prefix and 3 random characters
func getTempField(s string) string {
	prefix := "_" + utils.RandString(3)

	if fs := strings.Split(s, "."); len(fs) > 1 {
		return strings.Join(fs[:len(fs)-1], ".") + "." + prefix + fs[len(fs)-1]
	}
	return prefix + s
}

// StripAuthorRandomPrefix removes random prefix from author field in pipeline
// This is required to make unit test's comparison
func StripAuthorRandomPrefix(pipeline []bson.M) []bson.M {
	b, err := json.Marshal(pipeline)
	if err != nil {
		return pipeline
	}
	re := regexp.MustCompile(`\._[0-9a-z]{3}author\b`)
	b = re.ReplaceAll(b, []byte(`.author`))
	var pipeline2 []bson.M
	err = json.Unmarshal(b, &pipeline2)
	if err != nil {
		return pipeline
	}
	return pipeline2
}
