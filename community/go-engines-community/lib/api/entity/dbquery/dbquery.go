package dbquery

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
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
			"as":           "service_depends",
			"pipeline": []bson.M{
				{"$project": bson.M{"_id": 1}},
			},
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "_id",
			"foreignField": "component",
			"as":           "component_depends",
			"pipeline": []bson.M{
				{"$match": bson.M{"type": types.EntityTypeResource}},
				{"$project": bson.M{"_id": 1}},
			},
		}},
		{"$addFields": bson.M{
			"depends_count": bson.M{"$cond": bson.A{
				bson.M{"$and": []bson.M{
					{"$eq": bson.A{"$type", types.EntityTypeComponent}},
					{"$eq": bson.A{bson.M{"$type": "$state_info._id"}, "string"}},
					{"$ne": bson.A{"$state_info._id", ""}},
				}},
				bson.M{"$size": "$component_depends"},
				bson.M{"$size": "$service_depends"},
			}},
		}},
		{"$project": bson.M{"service_depends": 0, "component_depends": 0}},
	}
}

func GetImpactsCountPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "services",
			"foreignField": "_id",
			"as":           "service_impacts",
			"pipeline": []bson.M{
				{"$project": bson.M{"_id": 1}},
			},
		}},
		{"$unwind": bson.M{"path": "$entity_counters", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "component",
			"foreignField": "_id",
			"as":           "component_impacts",
			"pipeline": []bson.M{
				{"$project": bson.M{
					"hasStateSettings": bson.M{
						"$cond": []interface{}{bson.M{
							"$and": []bson.M{
								{"$eq": []interface{}{
									bson.M{"$type": "$state_info._id"},
									"string",
								}},
								{"$ne": []string{"$state_info._id", ""}},
							}},
							true,
							false,
						},
					},
				}},
			},
		}},
		{"$unwind": bson.M{"path": "$component_impacts", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			"impacts_count": bson.M{
				"$cond": []interface{}{
					bson.M{"$and": []interface{}{
						bson.M{"$eq": []string{
							"$type", "resource"},
						},
						"$component_impacts.hasStateSettings",
					}},
					bson.M{"$sum": bson.A{1, "$service_impacts"}},
					bson.M{"$size": "$service_impacts"},
				},
			},
		}},
		{"$project": bson.M{"service_impacts": 0, "component_impacts": 0}},
	}
}

func GetStateSettingPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.StateSettingsMongoCollection,
			"localField":   "state_info._id",
			"foreignField": "_id",
			"as":           "state_setting",
		}},
		{"$unwind": bson.M{"path": "$state_setting", "preserveNullAndEmptyArrays": true}},
	}
}
