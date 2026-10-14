package entity

import (
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type pipelineResult struct {
	BeforeLimit []bson.M
	AfterLimit  []bson.M
	Sort        bson.M
}

func TestListPipeline_ShouldCheckFullPipeline(t *testing.T) {
	// Spelled out in full, so a change inside any of the stage helpers fails
	// here instead of silently passing everywhere else.
	expected := pipelineResult{
		BeforeLimit: []bson.M{},
		AfterLimit: []bson.M{
			{"$addFields": bson.M{
				"raw": "$$ROOT",
			}},
			{"$lookup": bson.M{
				"from":         mongo.EntityCategoryMongoCollection,
				"localField":   "category",
				"foreignField": "_id",
				"as":           "category",
			}},
			{"$unwind": bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}},
			{"$addFields": bson.M{
				"pbehavior_info": bson.M{"$cond": bson.M{
					"if": "$pbehavior_info.id",
					"then": bson.M{
						"timestamp": "$pbehavior_info.timestamp",
						"rule": bson.M{
							"_id":  "$pbehavior_info.id",
							"name": "$pbehavior_info.name",
						},
						"reason": bson.M{
							"_id":  "$pbehavior_info.reason",
							"name": "$pbehavior_info.reason_name",
						},
						"type": bson.M{
							"_id":            "$pbehavior_info.type",
							"name":           "$pbehavior_info.type_name",
							"canonical_type": "$pbehavior_info.canonical_type",
						},
					},
					"else": nil,
				}},
			}},
			{"$project": bson.M{
				"services":  0,
				"perf_data": 0,
				"comments":  0,
				"imtags":    0,
			}},
		},
		Sort: bson.M{"$sort": bson.D{
			{Key: "name", Value: 1},
			{Key: "_id", Value: 1},
		}},
	}

	beforeLimit, afterLimit, sort, err := ListPipeline(ListRequest{})
	if err != nil {
		t.Fatalf("expected no error; got %v", err)
	}

	assertPipeline(t, pipelineResult{beforeLimit, afterLimit, sort}, expected)
}

func assertPipeline(t *testing.T, got, expected pipelineResult) {
	t.Helper()

	if diff := pretty.Compare(expected, got); diff != "" {
		t.Errorf("unexpected pipeline (-want +got)\n%s", diff)
	}
}
