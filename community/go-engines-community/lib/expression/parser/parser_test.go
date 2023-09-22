package parser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
)

func TestParser_Parse(t *testing.T) {
	p := NewParser()
	for i, data := range getParseDataSets() {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			output, err := p.Parse(data.Expression)

			if err != nil {
				if data.Err == nil || data.Err.Error() != err.Error() {
					t.Errorf("%s: expected no error but got %s %+v", data.Expression, err, data.Err)
				}
				return
			}

			if diff := pretty.Compare(output.Query(), data.ExpectedQuery); diff != "" {
				t.Errorf("unexpected query: %s", diff)
			}
			if diff := pretty.Compare(output.ExprQuery(), data.ExpectedExprQuery); diff != "" {
				t.Errorf("unexpected query: %s", diff)
			}
		})
	}
}

type dateSet struct {
	Expression        string
	ExpectedQuery     bson.M
	ExpectedExprQuery bson.M
	Err               error
}

func getParseDataSets() []dateSet {
	return []dateSet{
		{
			Expression:        "connector=\"test_connector\"",
			ExpectedQuery:     bson.M{"connector": bson.M{"$eq": "test_connector"}},
			ExpectedExprQuery: bson.M{"$eq": bson.A{"$connector", "test_connector"}},
		},
		{
			Expression:        "connector!=\"test_connector\"",
			ExpectedQuery:     bson.M{"connector": bson.M{"$ne": "test_connector"}},
			ExpectedExprQuery: bson.M{"$ne": bson.A{"$connector", "test_connector"}},
		},
		{
			Expression:        "connector<=\"test_connector\"",
			ExpectedQuery:     bson.M{"connector": bson.M{"$lte": "test_connector"}},
			ExpectedExprQuery: bson.M{"$lte": bson.A{"$connector", "test_connector"}},
		},
		{
			Expression:        "connector>=\"test_connector\"",
			ExpectedQuery:     bson.M{"connector": bson.M{"$gte": "test_connector"}},
			ExpectedExprQuery: bson.M{"$gte": bson.A{"$connector", "test_connector"}},
		},
		{
			Expression:        "connector<\"test_connector\"",
			ExpectedQuery:     bson.M{"connector": bson.M{"$lt": "test_connector"}},
			ExpectedExprQuery: bson.M{"$lt": bson.A{"$connector", "test_connector"}},
		},
		{
			Expression:        "connector>\"test_connector\"",
			ExpectedQuery:     bson.M{"connector": bson.M{"$gt": "test_connector"}},
			ExpectedExprQuery: bson.M{"$gt": bson.A{"$connector", "test_connector"}},
		},
		{
			Expression:        "d=\"test_component/test_resource\"",
			ExpectedQuery:     bson.M{"d": bson.M{"$eq": "test_component/test_resource"}},
			ExpectedExprQuery: bson.M{"$eq": bson.A{"$d", "test_component/test_resource"}},
		},
		{
			Expression:        "connector LIKE \"test_connector\"",
			ExpectedQuery:     bson.M{"connector": bson.M{"$regex": "test_connector"}},
			ExpectedExprQuery: bson.M{"$regexMatch": bson.M{"input": "$connector", "regex": "test_connector"}},
		},
		{
			Expression:        "name LIKE \"criticité\"",
			ExpectedQuery:     bson.M{"name": bson.M{"$regex": "criticité"}},
			ExpectedExprQuery: bson.M{"$regexMatch": bson.M{"input": "$name", "regex": "criticité"}},
		},
		{
			Expression:        "name LIKE 10",
			ExpectedQuery:     bson.M{"name": bson.M{"$regex": "10"}},
			ExpectedExprQuery: bson.M{"$regexMatch": bson.M{"input": "$name", "regex": "10"}},
		},
		{
			Expression:        "connector NOT LIKE \"test_connector\"",
			ExpectedQuery:     bson.M{"connector": bson.M{"$not": bson.M{"$regex": "test_connector"}}},
			ExpectedExprQuery: bson.M{"$not": bson.M{"$regexMatch": bson.M{"input": "$connector", "regex": "test_connector"}}},
		},
		{
			Expression:        "connector NOT LIKE 10",
			ExpectedQuery:     bson.M{"connector": bson.M{"$not": bson.M{"$regex": "10"}}},
			ExpectedExprQuery: bson.M{"$not": bson.M{"$regexMatch": bson.M{"input": "$connector", "regex": "10"}}},
		},
		{
			Expression:        "children CONTAINS \"test_connector\"",
			ExpectedQuery:     bson.M{"children": bson.M{"$in": []interface{}{"test_connector"}}},
			ExpectedExprQuery: bson.M{"$in": bson.A{"$children", []interface{}{"test_connector"}}},
		},
		{
			Expression:        "children NOT CONTAINS \"test_connector\"",
			ExpectedQuery:     bson.M{"children": bson.M{"$nin": []interface{}{"test_connector"}}},
			ExpectedExprQuery: bson.M{"$not": bson.M{"$in": bson.A{"$children", []interface{}{"test_connector"}}}},
		},
		{
			Expression:        "NOT connector=\"test_connector1\"",
			ExpectedQuery:     bson.M{"connector": bson.M{"$not": bson.M{"$eq": "test_connector1"}}},
			ExpectedExprQuery: bson.M{"$not": bson.M{"$eq": bson.A{"$connector", "test_connector1"}}},
		},
		{
			Expression: "connector=\"test_connector1\" AND connector=\"test_connector2\"",
			ExpectedQuery: bson.M{"$and": []bson.M{
				{"connector": bson.M{"$eq": "test_connector1"}},
				{"connector": bson.M{"$eq": "test_connector2"}},
			}},
			ExpectedExprQuery: bson.M{"$and": []bson.M{
				{"$eq": bson.A{"$connector", "test_connector1"}},
				{"$eq": bson.A{"$connector", "test_connector2"}},
			}},
		},
		{
			Expression: "connector=\"test_connector1\" OR connector=\"test_connector2\"",
			ExpectedQuery: bson.M{"$or": []bson.M{
				{"connector": bson.M{"$eq": "test_connector1"}},
				{"connector": bson.M{"$eq": "test_connector2"}},
			}},
			ExpectedExprQuery: bson.M{"$or": []bson.M{
				{"$eq": bson.A{"$connector", "test_connector1"}},
				{"$eq": bson.A{"$connector", "test_connector2"}},
			}},
		},
		{
			Expression: "connector=\"test_connector\" AND resource=\"test_resource1\" OR resource=\"test_resource2\"",
			ExpectedQuery: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"connector": bson.M{"$eq": "test_connector"}},
					{"resource": bson.M{"$eq": "test_resource1"}},
				}},
				{"resource": bson.M{"$eq": "test_resource2"}},
			}},
			ExpectedExprQuery: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"$eq": bson.A{"$connector", "test_connector"}},
					{"$eq": bson.A{"$resource", "test_resource1"}},
				}},
				{"$eq": bson.A{"$resource", "test_resource2"}},
			}},
		},
		{
			Expression: "connector=\"test_connector\" OR resource=\"test_resource1\" AND resource=\"test_resource2\"",
			ExpectedQuery: bson.M{"$or": []bson.M{
				{"connector": bson.M{"$eq": "test_connector"}},
				{"$and": []bson.M{
					{"resource": bson.M{"$eq": "test_resource1"}},
					{"resource": bson.M{"$eq": "test_resource2"}},
				}},
			}},
			ExpectedExprQuery: bson.M{"$or": []bson.M{
				{"$eq": bson.A{"$connector", "test_connector"}},
				{"$and": []bson.M{
					{"$eq": bson.A{"$resource", "test_resource1"}},
					{"$eq": bson.A{"$resource", "test_resource2"}},
				}},
			}},
		},
		{
			Expression:        "v.connector=\"test_connector\"",
			ExpectedQuery:     bson.M{"v.connector": bson.M{"$eq": "test_connector"}},
			ExpectedExprQuery: bson.M{"$eq": bson.A{"$v.connector", "test_connector"}},
		},
		{
			Expression:        "v.state.val=0",
			ExpectedQuery:     bson.M{"v.state.val": bson.M{"$eq": 0}},
			ExpectedExprQuery: bson.M{"$eq": bson.A{"$v.state.val", 0}},
		},
		{
			Expression:        "v.state.val=1",
			ExpectedQuery:     bson.M{"v.state.val": bson.M{"$eq": 1}},
			ExpectedExprQuery: bson.M{"$eq": bson.A{"$v.state.val", 1}},
		},
		{
			Expression:        "v.state.val>1.5",
			ExpectedQuery:     bson.M{"v.state.val": bson.M{"$gt": 1.5}},
			ExpectedExprQuery: bson.M{"$gt": bson.A{"$v.state.val", 1.5}},
		},
		{
			Expression:        "v.state.val=v.status.val",
			ExpectedQuery:     bson.M{"v.state.val": bson.M{"$eq": "v.status.val"}},
			ExpectedExprQuery: bson.M{"$eq": bson.A{"$v.state.val", "v.status.val"}},
		},
		{
			Expression:        "v.state.val=NULL",
			ExpectedQuery:     bson.M{"v.state.val": bson.M{"$eq": nil}},
			ExpectedExprQuery: bson.M{"$eq": bson.A{"$v.state.val", nil}},
		},
		{
			Expression:        "v.meta=TRUE",
			ExpectedQuery:     bson.M{"v.meta": bson.M{"$eq": true}},
			ExpectedExprQuery: bson.M{"$eq": bson.A{"$v.meta", true}},
		},
		{
			Expression:        "v.meta=FALSE",
			ExpectedQuery:     bson.M{"v.meta": bson.M{"$eq": false}},
			ExpectedExprQuery: bson.M{"$eq": bson.A{"$v.meta", false}},
		},
		{
			Expression: "NOT connector=\"test_connector1\" OR connector=\"test_connector2\" AND connector LIKE \"test_connector3\"",
			ExpectedQuery: bson.M{"$or": []bson.M{
				{"connector": bson.M{"$not": bson.M{"$eq": "test_connector1"}}},
				{"$and": []bson.M{
					{"connector": bson.M{"$eq": "test_connector2"}},
					{"connector": bson.M{"$regex": "test_connector3"}},
				}},
			}},
			ExpectedExprQuery: bson.M{"$or": []bson.M{
				{"$not": bson.M{"$eq": bson.A{"$connector", "test_connector1"}}},
				{"$and": []bson.M{
					{"$eq": bson.A{"$connector", "test_connector2"}},
					{"$regexMatch": bson.M{"input": "$connector", "regex": "test_connector3"}},
				}},
			}},
		},
		{
			Expression: "ressete",
			Err:        fmt.Errorf("comparison not found"),
		},
	}
}
