package pbehavior

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
	"sort"
	"testing"
)

func TestEntityMatcher_MatchAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	firstMockCursor := mock_mongo.NewMockCursor(ctrl)
	secondMockCursor := mock_mongo.NewMockCursor(ctrl)
	firstMockCursor.EXPECT().Next(gomock.Any()).Return(true)
	secondMockCursor.EXPECT().Next(gomock.Any()).Return(true)
	firstMockCursor.EXPECT().Close(gomock.Any())
	secondMockCursor.EXPECT().Close(gomock.Any())
	firstMockCursor.
		EXPECT().
		Decode(gomock.Any()).
		Do(func(val map[string][]string) {
			val["pbh1"] = []string{"1"}
			val["pbh2"] = []string{}
			val["pbh3"] = []string{}
			val["pbh4"] = []string{}
			val["pbh5"] = []string{}
		}).
		Return(nil)
	secondMockCursor.
		EXPECT().
		Decode(gomock.Any()).
		Do(func(val map[string][]string) {
			val["pbh6"] = []string{"1"}
			val["pbh7"] = []string{}
			val["pbh8"] = []string{"1"}
		}).
		Return(nil)
	entityID := "testid"
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	firstCall := mockDbCollection.
		EXPECT().
		Aggregate(gomock.Any(), gomock.Eq([]bson.M{
			{"$match": bson.M{"_id": entityID}},
			{"$facet": bson.M{
				"pbh1": []bson.M{{"$match": bson.M{"name": "resource1"}}},
				"pbh2": []bson.M{{"$match": bson.M{"name": "resource2"}}},
				"pbh3": []bson.M{{"$match": bson.M{"name": "resource3"}}},
				"pbh4": []bson.M{{"$match": bson.M{"name": "resource4"}}},
				"pbh5": []bson.M{{"$match": bson.M{"name": "resource5"}}},
			}},
			{"$addFields": bson.M{
				"ids": bson.M{
					"$arrayToObject": bson.M{
						"$map": bson.M{
							"input": bson.M{"$objectToArray": "$$ROOT"},
							"as":    "each",
							"in": bson.M{
								"k": "$$each.k",
								"v": bson.M{"$map": bson.M{
									"input": "$$each.v",
									"as":    "e",
									"in":    "$$e._id",
								}},
							},
						},
					},
				},
			}},
			{"$replaceRoot": bson.M{"newRoot": "$ids"}},
		})).
		Return(firstMockCursor, nil)
	secondCall := mockDbCollection.
		EXPECT().
		Aggregate(gomock.Any(), gomock.Eq([]bson.M{
			{"$match": bson.M{"_id": entityID}},
			{"$facet": bson.M{
				"pbh6": []bson.M{{"$match": bson.M{"name": "resource6"}}},
				"pbh7": []bson.M{{"$match": bson.M{"name": "resource7"}}},
				"pbh8": []bson.M{{"$match": bson.M{"name": "resource8"}}},
			}},
			{"$addFields": bson.M{
				"ids": bson.M{
					"$arrayToObject": bson.M{
						"$map": bson.M{
							"input": bson.M{"$objectToArray": "$$ROOT"},
							"as":    "each",
							"in": bson.M{
								"k": "$$each.k",
								"v": bson.M{"$map": bson.M{
									"input": "$$each.v",
									"as":    "e",
									"in":    "$$e._id",
								}},
							},
						},
					},
				},
			}},
			{"$replaceRoot": bson.M{"newRoot": "$ids"}},
		})).
		Return(secondMockCursor, nil)
	gomock.InOrder(firstCall, secondCall)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.
		EXPECT().
		Collection(mongo.EntityMongoCollection).
		Return(mockDbCollection)
	filters := map[string]interface{}{
		"pbh1": bson.M{"name": "resource1"},
		"pbh2": bson.M{"name": "resource2"},
		"pbh3": bson.M{"name": "resource3"},
		"pbh4": bson.M{"name": "resource4"},
		"pbh5": bson.M{"name": "resource5"},
		"pbh6": bson.M{"name": "resource6"},
		"pbh7": bson.M{"name": "resource7"},
		"pbh8": bson.M{"name": "resource8"},
	}
	expected := []string{
		"pbh1",
		"pbh6",
		"pbh8",
	}

	m := NewEntityMatcher(mockDbClient, 5)
	res, err := m.MatchAll(context.Background(), entityID, filters)

	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	sort.Strings(res)

	if diff := pretty.Compare(expected, res); diff != "" {
		t.Errorf("unexpected result: %s", diff)
	}
}
