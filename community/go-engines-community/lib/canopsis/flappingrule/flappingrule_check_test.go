package flappingrule

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

func TestSetDefaultFlappingCheck(t *testing.T) {
	dbClient, err := mongo.NewClient(context.Background(), 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	flappingCollection := dbClient.Collection(mongo.FlappingRuleMongoCollection)
	_, err = flappingCollection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		t.Fatal(err)
	}

	flappingAdapter := NewAdapter(dbClient)

	flapping := SetThenGetFlappingCheck(flappingAdapter, context.Background(), time.Second, log.NewLogger(false))
	Convey("get rules should work find", t, func() {
		rules := flapping.GetRules()
		So(rules, ShouldHaveLength, 0)

		_, err := flappingCollection.InsertOne(context.Background(), bson.M{
			"_id":         "test-flapping-rule-to-create-1",
			"description": "test create 1",
			"alarm_patterns": bson.A{
				bson.M{
					"v": bson.M{
						"component": "test-flapping-rule-to-create-1-pattern",
					},
				},
			},
			"flapping_interval": bson.M{
				"seconds": 10,
				"unit":    "s",
			},
			"flapping_freq_limit": 3,
			"priority":            5,
			"updated":             1629875485,
		})
		So(err, ShouldBeNil)
		So(flapping.GetRules(), ShouldHaveLength, 0)
		err = flapping.LoadRules(context.Background())
		So(err, ShouldBeNil)
		So(flapping.GetRules(), ShouldHaveLength, 1)
		_, err = flappingCollection.InsertOne(context.Background(), bson.M{
			"_id":         "test-flapping-rule-to-create-2",
			"description": "test create 1",
			"alarm_patterns": bson.A{
				bson.M{
					"v": bson.M{
						"component": "test-flapping-rule-to-create-1-pattern",
					},
				},
			},
			"flapping_interval": bson.M{
				"seconds": 10,
				"unit":    "s",
			},
			"flapping_freq_limit": 3,
			"priority":            5,
			"updated":             1629875495,
		})
		_, err = flappingCollection.UpdateOne(context.Background(), bson.M{"_id": "test-flapping-rule-to-create-1"},
			bson.M{"$set": bson.M{
				"description": "test create 1",
				"alarm_patterns": bson.A{
					bson.M{
						"v": bson.M{
							"component": "test-flapping-rule-to-create-1-pattern",
						},
					},
				},
				"flapping_interval": bson.M{
					"seconds": 10,
					"unit":    "s",
				},
				"flapping_freq_limit": 31,
				"priority":            15,
				"updated":             1629875695,
			},
			})
		So(err, ShouldBeNil)
		So(flapping.GetRules(), ShouldHaveLength, 1)
		err = flapping.LoadRules(context.Background())
		So(err, ShouldBeNil)
		So(flapping.GetRules(), ShouldHaveLength, 2)
	})
}
