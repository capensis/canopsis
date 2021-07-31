package baggotrule

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGetRule(t *testing.T) {
	client, err := mongo.NewClient(context.Background(), 1, time.Minute)
	ShouldBeNil(err)
	collection := client.Collection(mongo.BaggotRuleMongoCollection)
	_, err = collection.DeleteMany(context.Background(), bson.M{})
	ShouldBeNil(err)

	now := time.Now()
	_, err = collection.InsertMany(context.Background(), []interface{}{
		Rule{
			ID: "rule_1",
			Duration: types.DurationWithUnit{
				Seconds: 10,
				Unit:    "s",
			},
			AlarmPatterns:  pattern.AlarmPatternList{},
			EntityPatterns: pattern.EntityPatternList{},
			Updated: &types.CpsTime{
				Time: now.Add(5 * time.Minute),
			},
			Priority: 2,
		},
		Rule{
			ID: "rule_2",
			Duration: types.DurationWithUnit{
				Seconds: 10,
				Unit:    "s",
			},
			AlarmPatterns:  pattern.AlarmPatternList{},
			EntityPatterns: pattern.EntityPatternList{},
			Updated: &types.CpsTime{
				Time: now,
			},
			Priority: 1,
		},
		Rule{
			ID: "rule_3",
			Duration: types.DurationWithUnit{
				Seconds: 10,
				Unit:    "s",
			},
			AlarmPatterns:  pattern.AlarmPatternList{},
			EntityPatterns: pattern.EntityPatternList{},
			Updated: &types.CpsTime{
				Time: now.Add(10 * time.Minute),
			},
			Priority: 2,
		},
	})
	ShouldBeNil(err)

	ruleAdpt := &mongoAdapter{collection: collection}
	rules, err := ruleAdpt.Get(context.Background())
	ShouldBeNil(err)
	ShouldHaveLength(rules, 3)
	ShouldEqual(rules[0].ID, "rule_2")
	ShouldEqual(rules[1].ID, "rule_3")
	ShouldEqual(rules[2].ID, "rule_1")
}
