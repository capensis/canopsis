package parser

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestParser_Parse(t *testing.T) {
	p := NewParser()
	for _, data := range getParseDataSets() {
		output, err := p.Parse(data.Expression)

		if err != nil {
			if data.Err == nil || data.Err.Error() != err.Error() {
				t.Errorf("%s: expected no error but got %s %+v", data.Expression, err, data.Err)
			}
			continue
		}

		if !reflect.DeepEqual(output.Query(), data.Expected) {
			expectedMsg, _ := json.MarshalIndent(data.Expected, "", "  ")
			outputMsg, _ := json.MarshalIndent(output.Query(), "", "  ")
			t.Errorf("%s: expected\n%s\nbut got\n%s", data.Expression, expectedMsg, outputMsg)
		}
	}
}

type dateSet struct {
	Expression string
	Expected   interface{}
	Err        error
}

func getParseDataSets() []dateSet {
	return []dateSet{
		{
			Expression: "connector=\"test_connector\"",
			Expected:   bson.M{"connector": bson.M{"$eq": "test_connector"}},
		},
		{
			Expression: "connector!=\"test_connector\"",
			Expected:   bson.M{"connector": bson.M{"$ne": "test_connector"}},
		},
		{
			Expression: "connector<=\"test_connector\"",
			Expected:   bson.M{"connector": bson.M{"$lte": "test_connector"}},
		},
		{
			Expression: "connector>=\"test_connector\"",
			Expected:   bson.M{"connector": bson.M{"$gte": "test_connector"}},
		},
		{
			Expression: "connector<\"test_connector\"",
			Expected:   bson.M{"connector": bson.M{"$lt": "test_connector"}},
		},
		{
			Expression: "connector>\"test_connector\"",
			Expected:   bson.M{"connector": bson.M{"$gt": "test_connector"}},
		},
		{
			Expression: "d=\"test_component/test_resource\"",
			Expected:   bson.M{"d": bson.M{"$eq": "test_component/test_resource"}},
		},
		{
			Expression: "connector LIKE \"test_connector\"",
			Expected:   bson.M{"connector": bson.M{"$regex": "test_connector"}},
		},
		{
			Expression: "name LIKE \"criticité\"",
			Expected:   bson.M{"name": bson.M{"$regex": "criticité"}},
		},
		{
			Expression: "name LIKE 10",
			Expected:   bson.M{"name": bson.M{"$regex": "10"}},
		},
		{
			Expression: "connector NOT LIKE \"test_connector\"",
			Expected:   bson.M{"connector": bson.M{"$not": bson.M{"$regex": "test_connector"}}},
		},
		{
			Expression: "connector NOT LIKE 10",
			Expected:   bson.M{"connector": bson.M{"$not": bson.M{"$regex": "10"}}},
		},
		{
			Expression: "children CONTAINS \"test_connector\"",
			Expected:   bson.M{"children": bson.M{"$in": []interface{}{"test_connector"}}},
		},
		{
			Expression: "children NOT CONTAINS \"test_connector\"",
			Expected:   bson.M{"children": bson.M{"$nin": []interface{}{"test_connector"}}},
		},
		{
			Expression: "NOT connector=\"test_connector1\"",
			Expected:   bson.M{"connector": bson.M{"$not": bson.M{"$eq": "test_connector1"}}},
		},
		{
			Expression: "connector=\"test_connector1\" AND connector=\"test_connector2\"",
			Expected: bson.M{"$and": []bson.M{
				{"connector": bson.M{"$eq": "test_connector1"}},
				{"connector": bson.M{"$eq": "test_connector2"}},
			}},
		},
		{
			Expression: "connector=\"test_connector1\" OR connector=\"test_connector2\"",
			Expected: bson.M{"$or": []bson.M{
				{"connector": bson.M{"$eq": "test_connector1"}},
				{"connector": bson.M{"$eq": "test_connector2"}},
			}},
		},
		{
			Expression: "connector=\"test_connector\" AND resource=\"test_resource1\" OR resource=\"test_resource2\"",
			Expected: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"connector": bson.M{"$eq": "test_connector"}},
					{"resource": bson.M{"$eq": "test_resource1"}},
				}},
				{"resource": bson.M{"$eq": "test_resource2"}},
			}},
		},
		{
			Expression: "connector=\"test_connector\" OR resource=\"test_resource1\" AND resource=\"test_resource2\"",
			Expected: bson.M{"$or": []bson.M{
				{"connector": bson.M{"$eq": "test_connector"}},
				{"$and": []bson.M{
					{"resource": bson.M{"$eq": "test_resource1"}},
					{"resource": bson.M{"$eq": "test_resource2"}},
				}},
			}},
		},
		{
			Expression: "v.connector=\"test_connector\"",
			Expected:   bson.M{"v.connector": bson.M{"$eq": "test_connector"}},
		},
		{
			Expression: "v.state.val=0",
			Expected:   bson.M{"v.state.val": bson.M{"$eq": 0}},
		},
		{
			Expression: "v.state.val=1",
			Expected:   bson.M{"v.state.val": bson.M{"$eq": 1}},
		},
		{
			Expression: "v.state.val>1.5",
			Expected:   bson.M{"v.state.val": bson.M{"$gt": 1.5}},
		},
		{
			Expression: "v.state.val=v.status.val",
			Expected:   bson.M{"v.state.val": bson.M{"$eq": "v.status.val"}},
		},
		{
			Expression: "v.state.val=NULL",
			Expected:   bson.M{"v.state.val": bson.M{"$eq": nil}},
		},
		{
			Expression: "v.meta=TRUE",
			Expected:   bson.M{"v.meta": bson.M{"$eq": true}},
		},
		{
			Expression: "v.meta=FALSE",
			Expected:   bson.M{"v.meta": bson.M{"$eq": false}},
		},
		{
			Expression: "NOT connector=\"test_connector1\" OR connector=\"test_connector2\" AND connector LIKE \"test_connector3\"",
			Expected: bson.M{"$or": []bson.M{
				{"connector": bson.M{"$not": bson.M{"$eq": "test_connector1"}}},
				{"$and": []bson.M{
					{"connector": bson.M{"$eq": "test_connector2"}},
					{"connector": bson.M{"$regex": "test_connector3"}},
				}},
			}},
		},
		{
			Expression: "ressete",
			Err:        fmt.Errorf("comparison not found"),
		},
	}
}
