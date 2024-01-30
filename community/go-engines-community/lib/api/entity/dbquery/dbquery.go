package dbquery

import (
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func GetCategoryLookup() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   "category",
			"foreignField": "_id",
			"as":           "category",
		}},
		{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
	}
}

func GetPbehaviorInfoTypeLookup() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorTypeMongoCollection,
			"foreignField": "_id",
			"localField":   "pbehavior_info.type",
			"as":           "pbehavior_info_type",
		}},
		{"$unwind": bson.M{"path": "$pbehavior_info_type", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"pbehavior_info": bson.M{"$cond": bson.M{
				"if": "$pbehavior_info",
				"then": bson.M{"$mergeObjects": bson.A{
					"$pbehavior_info",
					bson.M{"icon_name": "$pbehavior_info_type.icon_name"},
				}},
				"else": nil,
			}},
		}},
		{"$project": bson.M{"pbehavior_info_type": 0}},
	}
}

func GetDependsCountPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "_id",
			"foreignField": "services",
			"as":           "depends",
			"pipeline": []bson.M{
				{"$project": bson.M{"_id": 1}},
			},
		}},
		{"$addFields": bson.M{
			"depends_count": bson.M{"$size": "$depends"},
		}},
		{"$project": bson.M{"depends": 0}},
	}
}

func GetStateSettingLookup(entityPrefix ...string) []bson.M {
	prefix := ""
	if len(entityPrefix) > 0 {
		prefix = strings.TrimRight(entityPrefix[0], ".") + "."
	}
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.StateSettingsMongoCollection,
			"localField":   prefix + "state_info._id",
			"foreignField": "_id",
			"as":           "state_setting",
		}},
		{"$unwind": bson.M{"path": "$state_setting", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{prefix + "state_setting": bson.M{
			"$cond": bson.M{"else": "$$REMOVE", "if": "$state_setting", "then": bson.M{
				"title":                    "$state_setting.title",
				"method":                   "$state_setting.method",
				"inherited_entity_pattern": "$state_setting.inherited_entity_pattern",
				"state_thresholds":         "$state_setting.state_thresholds",
			}},
		}}},
	}
}
