package dbquery

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetCategoryLookup(prefixArg ...string) []bson.M {
	prefix := ""
	if len(prefixArg) > 0 && prefixArg[0] != "" {
		prefix = prefixArg[0] + "."
	}

	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityCategoryMongoCollection,
			"localField":   prefix + "category",
			"foreignField": "_id",
			"as":           prefix + "category",
		}},
		{"$unwind": bson.M{"path": "$" + prefix + "category", "preserveNullAndEmptyArrays": true}},
	}
}

func GetPbehaviorInfoLastCommentLookup(authorProvider author.Provider, prefixArg ...string) []bson.M {
	prefix := ""
	if len(prefixArg) > 0 && prefixArg[0] != "" {
		prefix = prefixArg[0] + "."
	}

	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.PbehaviorMongoCollection,
			"foreignField": "_id",
			"localField":   prefix + "pbehavior_info.id",
			"as":           prefix + "pbehavior",
		}},
		{"$unwind": bson.M{"path": "$" + prefix + "pbehavior", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			prefix + "pbehavior.last_comment": bson.M{"$arrayElemAt": bson.A{"$" + prefix + "pbehavior.comments", -1}},
		}},
	}
	pipeline = append(pipeline, authorProvider.PipelineForField(prefix+"pbehavior.last_comment.author")...)
	pipeline = append(pipeline,
		bson.M{"$addFields": bson.M{
			prefix + "pbehavior.last_comment": bson.M{
				"$cond": bson.M{
					"if":   "$" + prefix + "pbehavior.last_comment._id",
					"then": "$" + prefix + "pbehavior.last_comment",
					"else": "$$REMOVE",
				},
			},
		}},
		bson.M{"$addFields": bson.M{
			prefix + "pbehavior_info": bson.M{"$cond": bson.M{
				"if": "$" + prefix + "pbehavior_info",
				"then": bson.M{"$mergeObjects": bson.A{
					"$" + prefix + "pbehavior_info",
					bson.M{
						"last_comment": bson.M{"$cond": bson.M{
							"if": "$" + prefix + "pbehavior.last_comment",
							"then": bson.M{"$mergeObjects": bson.A{
								"$" + prefix + "pbehavior.last_comment",
								bson.M{"author": bson.M{"$cond": bson.M{
									"if":   "$" + prefix + "pbehavior.last_comment.origin",
									"then": "$" + prefix + "pbehavior.last_comment.origin",
									"else": "$" + prefix + "pbehavior.last_comment.author.display_name",
								}}},
							}},
							"else": nil,
						}},
					},
				}},
				"else": nil,
			}},
		}},
		bson.M{"$project": bson.M{
			prefix + "pbehavior": 0,
		}},
	)

	return pipeline
}

func GetDependsCountPipeline(prefixArg ...string) []bson.M {
	prefix := ""
	if len(prefixArg) > 0 && prefixArg[0] != "" {
		prefix = prefixArg[0] + "."
	}

	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   prefix + "_id",
			"foreignField": "services",
			"as":           prefix + "service_depends",
			"pipeline": []bson.M{
				{"$match": bson.M{"soft_deleted": bson.M{"$exists": false}}},
				{"$project": bson.M{"_id": 1}},
			},
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   prefix + "_id",
			"foreignField": "component",
			"as":           prefix + "component_depends",
			"pipeline": []bson.M{
				{"$match": bson.M{
					"type":         types.EntityTypeResource,
					"soft_deleted": bson.M{"$exists": false},
				}},
				{"$project": bson.M{"_id": 1}},
			},
		}},
		{"$addFields": bson.M{
			prefix + "depends_count": bson.M{"$cond": bson.A{
				bson.M{"$and": []bson.M{
					{"$eq": bson.A{"$" + prefix + "type", types.EntityTypeComponent}},
					{"$eq": bson.A{bson.M{"$type": "$" + prefix + "state_info._id"}, "string"}},
					{"$ne": bson.A{"$" + prefix + "state_info._id", ""}},
				}},
				bson.M{"$size": "$" + prefix + "component_depends"},
				bson.M{"$size": "$" + prefix + "service_depends"},
			}},
		}},
		{"$project": bson.M{prefix + "service_depends": 0, prefix + "component_depends": 0}},
	}
}

func GetImpactsCountPipeline(prefixArg ...string) []bson.M {
	prefix := ""
	if len(prefixArg) > 0 && prefixArg[0] != "" {
		prefix = prefixArg[0] + "."
	}

	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   prefix + "services",
			"foreignField": "_id",
			"as":           prefix + "service_impacts",
			"pipeline": []bson.M{
				{"$match": bson.M{"soft_deleted": bson.M{"$exists": false}}},
				{"$project": bson.M{"_id": 1}},
			},
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   prefix + "component",
			"foreignField": "_id",
			"as":           prefix + "component_impacts",
			"pipeline": []bson.M{
				{"$project": bson.M{
					"hasStateSettings": bson.M{
						"$cond": []any{bson.M{
							"$and": []bson.M{
								{"$eq": []any{
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
		{"$unwind": bson.M{"path": "$" + prefix + "component_impacts", "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			prefix + "impacts_count": bson.M{
				"$cond": []any{
					bson.M{"$and": []any{
						bson.M{"$eq": []string{
							"$" + prefix + "type", "resource"},
						},
						"$" + prefix + "component_impacts.hasStateSettings",
					}},
					bson.M{"$sum": bson.A{1, bson.M{"$size": "$" + prefix + "service_impacts"}}},
					bson.M{"$size": "$" + prefix + "service_impacts"},
				},
			},
		}},
		{"$project": bson.M{prefix + "service_impacts": 0, prefix + "component_impacts": 0}},
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

func GetDownstreamCountPipeline(prefixArg ...string) []bson.M {
	prefix := ""
	if len(prefixArg) > 0 && prefixArg[0] != "" {
		prefix = prefixArg[0] + "."
	}

	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   prefix + "_id",
			"foreignField": "upstream",
			"as":           prefix + "downstreams",
			"pipeline": []bson.M{
				{"$match": bson.M{
					"type":    bson.M{"$in": bson.A{types.EntityTypeResource, types.EntityTypeComponent}},
					"enabled": true,
				}},
				{"$project": bson.M{"_id": 1}},
			},
		}},
		{"$addFields": bson.M{
			prefix + "downstream_count": bson.M{"$cond": bson.M{
				"if": bson.M{"$and": []bson.M{
					{"$in": bson.A{"$" + prefix + "type", bson.A{types.EntityTypeResource, types.EntityTypeComponent}}},
					{"$eq": bson.A{"$" + prefix + "enabled", true}},
				}},
				"then": bson.M{"$size": "$" + prefix + "downstreams"},
				"else": 0,
			}},
		}},
		{"$project": bson.M{prefix + "downstreams": 0}},
	}
}
