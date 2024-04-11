package parser

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics/schema"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
)

func TestParser_ParseMongo(t *testing.T) {
	p := NewParser()
	for i, data := range getParseMongoDataSets() {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			output, err := p.Parse(data.Expression, nil)

			if err != nil {
				if data.Err == nil || data.Err.Error() != err.Error() {
					t.Errorf("%s: expected no error but got %s %+v", data.Expression, err, data.Err)
				}
				return
			}

			if diff := pretty.Compare(output.MongoQuery(), data.ExpectedQuery); diff != "" {
				t.Errorf("unexpected query: %s", diff)
			}
			if diff := pretty.Compare(output.MongoExprQuery(), data.ExpectedExprQuery); diff != "" {
				t.Errorf("unexpected query: %s", diff)
			}
		})
	}
}

func TestParser_ParsePostgres(t *testing.T) {
	p := NewParser()
	for i, data := range getParsePostgresDataSets() {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			output, err := p.Parse(data.Expression, func(f string) bool {
				return schema.GetAllowedSearchEntityMetaFields()[f] || strings.HasPrefix(f, "infos.") || strings.HasPrefix(f, "component_infos.")
			})

			if err != nil {
				if data.Err == nil || data.Err.Error() != err.Error() {
					t.Errorf("%s: expected no error but got %s %+v", data.Expression, err, data.Err)
				}
				return
			}

			resQuery, resArgs := output.PostgresQuery(data.Prefix)
			for k, v := range resArgs {
				arg := "@" + k
				switch val := v.(type) {
				case string:
					resQuery = strings.Replace(resQuery, arg, "'"+val+"'", -1)
				case int:
					resQuery = strings.Replace(resQuery, arg, strconv.Itoa(val), -1)
				case float64:
					resQuery = strings.Replace(resQuery, arg, strconv.FormatFloat(val, 'f', -1, 64), -1)
				case bool:
					resQuery = strings.Replace(resQuery, arg, strconv.FormatBool(val), -1)
				case []any:
					strSlice := make([]string, len(val))
					for idx := range val {
						strSlice[idx] = "'" + val[idx].(string) + "'"
					}

					resQuery = strings.Replace(resQuery, arg, "ARRAY ["+strings.Join(strSlice, ",")+"]", -1)
				}
			}

			if resQuery != data.ExpectedQuery {
				t.Errorf("expected query = %q, but got %q\n", data.ExpectedQuery, resQuery)
			}
		})
	}
}

type mongoDataSet struct {
	Expression        string
	ExpectedQuery     bson.M
	ExpectedExprQuery bson.M
	Err               error
}

type postgresDataSet struct {
	Expression    string
	ExpectedQuery string
	Prefix        string
	Err           error
}

func getParseMongoDataSets() []mongoDataSet {
	return []mongoDataSet{
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
			ExpectedQuery:     bson.M{"children": bson.M{"$in": []any{"test_connector"}}},
			ExpectedExprQuery: bson.M{"$in": bson.A{"$children", []any{"test_connector"}}},
		},
		{
			Expression:        "children NOT CONTAINS \"test_connector\"",
			ExpectedQuery:     bson.M{"children": bson.M{"$nin": []any{"test_connector"}}},
			ExpectedExprQuery: bson.M{"$not": bson.M{"$in": bson.A{"$children", []any{"test_connector"}}}},
		},
		{
			Expression:        "children CONTAINS \"test_connector\" \"test_connector_2\"",
			ExpectedQuery:     bson.M{"children": bson.M{"$in": []any{"test_connector", "test_connector_2"}}},
			ExpectedExprQuery: bson.M{"$in": bson.A{"$children", []any{"test_connector", "test_connector_2"}}},
		},
		{
			Expression:        "children NOT CONTAINS \"test_connector\" \"test_connector_2\"",
			ExpectedQuery:     bson.M{"children": bson.M{"$nin": []any{"test_connector", "test_connector_2"}}},
			ExpectedExprQuery: bson.M{"$not": bson.M{"$in": bson.A{"$children", []any{"test_connector", "test_connector_2"}}}},
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
			Err:        errors.New("comparison not found"),
		},
	}
}

func getParsePostgresDataSets() []postgresDataSet {
	return []postgresDataSet{
		{
			Expression:    "connector=\"test_connector\"",
			ExpectedQuery: "WHERE connector = 'test_connector'",
		},
		{
			Expression:    "NOT connector=\"test_connector\"",
			ExpectedQuery: "WHERE (connector IS NULL OR connector != 'test_connector')",
		},
		{
			Expression:    "connector!=\"test_connector\"",
			ExpectedQuery: "WHERE (connector IS NULL OR connector != 'test_connector')",
		},
		{
			Expression:    "NOT connector!=\"test_connector\"",
			ExpectedQuery: "WHERE connector = 'test_connector'",
		},
		{
			Expression:    "connector<=\"test_connector\"",
			ExpectedQuery: "WHERE connector <= 'test_connector'",
		},
		{
			Expression:    "NOT connector<=\"test_connector\"",
			ExpectedQuery: "WHERE connector > 'test_connector'",
		},
		{
			Expression:    "connector>=\"test_connector\"",
			ExpectedQuery: "WHERE connector >= 'test_connector'",
		},
		{
			Expression:    "NOT connector>=\"test_connector\"",
			ExpectedQuery: "WHERE connector < 'test_connector'",
		},
		{
			Expression:    "connector<\"test_connector\"",
			ExpectedQuery: "WHERE connector < 'test_connector'",
		},
		{
			Expression:    "NOT connector<\"test_connector\"",
			ExpectedQuery: "WHERE connector >= 'test_connector'",
		},
		{
			Expression:    "connector>\"test_connector\"",
			ExpectedQuery: "WHERE connector > 'test_connector'",
		},
		{
			Expression:    "NOT connector>\"test_connector\"",
			ExpectedQuery: "WHERE connector <= 'test_connector'",
		},
		{
			Expression:    "connector LIKE \"test_connector\"",
			ExpectedQuery: "WHERE connector ~ 'test_connector'",
		},
		{
			Expression:    "name LIKE \"criticité\"",
			ExpectedQuery: "WHERE name ~ 'criticité'",
		},
		{
			Expression:    "name LIKE 10",
			ExpectedQuery: "WHERE name ~ '10'",
		},
		{
			Expression:    "infos.name LIKE \"criticité\"",
			ExpectedQuery: "WHERE infos->>'name' ~ 'criticité'",
		},
		{
			Expression:    "component_infos.name = \"criticité\"",
			ExpectedQuery: "WHERE component_infos->>'name' = 'criticité'",
		},
		{
			Expression:    "NOT connector LIKE \"test_connector\"",
			ExpectedQuery: "WHERE (connector IS NULL OR connector !~ 'test_connector')",
		},
		{
			Expression:    "connector NOT LIKE \"test_connector\"",
			ExpectedQuery: "WHERE (connector IS NULL OR connector !~ 'test_connector')",
		},
		{
			Expression:    "connector NOT LIKE 10",
			ExpectedQuery: "WHERE (connector IS NULL OR connector !~ '10')",
		},
		{
			Expression:    "NOT connector NOT LIKE \"test_connector\"",
			ExpectedQuery: "WHERE connector ~ 'test_connector'",
		},
		{
			Expression:    "connector CONTAINS \"test_connector\"",
			ExpectedQuery: "WHERE connector = ANY (ARRAY ['test_connector'])",
		},
		{
			Expression:    "connector CONTAINS 123",
			ExpectedQuery: "WHERE connector = ANY (ARRAY ['123'])",
		},
		{
			Expression:    "NOT connector CONTAINS \"test_connector\"",
			ExpectedQuery: "WHERE (connector IS NULL OR NOT connector = ANY (ARRAY ['test_connector']))",
		},
		{
			Expression:    "connector NOT CONTAINS \"test_connector\"",
			ExpectedQuery: "WHERE (connector IS NULL OR NOT connector = ANY (ARRAY ['test_connector']))",
		},
		{
			Expression:    "connector NOT CONTAINS 123",
			ExpectedQuery: "WHERE (connector IS NULL OR NOT connector = ANY (ARRAY ['123']))",
		},
		{
			Expression:    "NOT connector NOT CONTAINS \"test_connector\"",
			ExpectedQuery: "WHERE connector = ANY (ARRAY ['test_connector'])",
		},
		{
			Expression:    "connector CONTAINS \"test_connector\" \"test_connector_2\"",
			ExpectedQuery: "WHERE connector = ANY (ARRAY ['test_connector','test_connector_2'])",
		},
		{
			Expression:    "connector CONTAINS \"test_connector\" 123",
			ExpectedQuery: "WHERE connector = ANY (ARRAY ['test_connector','123'])",
		},
		{
			Expression:    "NOT connector CONTAINS \"test_connector\" 123",
			ExpectedQuery: "WHERE (connector IS NULL OR NOT connector = ANY (ARRAY ['test_connector','123']))",
		},
		{
			Expression:    "connector NOT CONTAINS \"test_connector\" \"test_connector_2\"",
			ExpectedQuery: "WHERE (connector IS NULL OR NOT connector = ANY (ARRAY ['test_connector','test_connector_2']))",
		},
		{
			Expression:    "connector NOT CONTAINS \"test_connector\" 123",
			ExpectedQuery: "WHERE (connector IS NULL OR NOT connector = ANY (ARRAY ['test_connector','123']))",
		},
		{
			Expression:    "NOT connector NOT CONTAINS \"test_connector\" 123",
			ExpectedQuery: "WHERE connector = ANY (ARRAY ['test_connector','123'])",
		},
		{
			Expression:    "connector=\"test_connector1\" AND connector=\"test_connector2\"",
			ExpectedQuery: "WHERE connector = 'test_connector1' AND connector = 'test_connector2'",
		},
		{
			Expression:    "connector=\"test_connector1\" OR connector=\"test_connector2\"",
			ExpectedQuery: "WHERE (connector = 'test_connector1' OR connector = 'test_connector2')",
		},
		{
			Expression:    "connector=\"test_connector\" AND name=\"test_resource1\" OR name=\"test_resource2\"",
			ExpectedQuery: "WHERE (connector = 'test_connector' AND name = 'test_resource1' OR name = 'test_resource2')",
		},
		{
			Expression:    "connector=\"test_connector\" OR name=\"test_resource1\" AND name=\"test_resource2\"",
			ExpectedQuery: "WHERE (connector = 'test_connector' OR name = 'test_resource1' AND name = 'test_resource2')",
		},
		{
			Expression:    "infos.val>1.5",
			ExpectedQuery: "WHERE infos->>'val' > 1.5",
		},
		{
			Expression:    "NOT infos.val>1.5",
			ExpectedQuery: "WHERE infos->>'val' <= 1.5",
		},
		{
			Expression:    "infos.val<1",
			ExpectedQuery: "WHERE infos->>'val' < 1",
		},
		{
			Expression:    "NOT infos.val<1",
			ExpectedQuery: "WHERE infos->>'val' >= 1",
		},
		{
			Expression:    "infos.val>=1.5",
			ExpectedQuery: "WHERE infos->>'val' >= 1.5",
		},
		{
			Expression:    "NOT infos.val>=1.5",
			ExpectedQuery: "WHERE infos->>'val' < 1.5",
		},
		{
			Expression:    "infos.val<=1",
			ExpectedQuery: "WHERE infos->>'val' <= 1",
		},
		{
			Expression:    "NOT infos.val<=1",
			ExpectedQuery: "WHERE infos->>'val' > 1",
		},
		{
			Expression:    "category=NULL",
			ExpectedQuery: "WHERE category IS NULL",
		},
		{
			Expression:    "category!=NULL",
			ExpectedQuery: "WHERE category IS NOT NULL",
		},
		{
			Expression:    "NOT category=NULL",
			ExpectedQuery: "WHERE category IS NOT NULL",
		},
		{
			Expression:    "NOT category!=NULL",
			ExpectedQuery: "WHERE category IS NULL",
		},
		{
			Expression:    "name=TRUE",
			ExpectedQuery: "WHERE name = true",
		},
		{
			Expression:    "name=FALSE",
			ExpectedQuery: "WHERE name = false",
		},
		{
			Expression:    "NOT connector=\"test_connector1\" OR connector=\"test_connector2\" AND connector LIKE \"test_connector3\"",
			ExpectedQuery: "WHERE ((connector IS NULL OR connector != 'test_connector1') OR connector = 'test_connector2' AND connector ~ 'test_connector3')",
		},
		{
			Expression:    "NOT connector=\"test_connector1\" OR NOT connector=\"test_connector2\" AND NOT connector NOT LIKE \"test_connector3\"",
			ExpectedQuery: "WHERE ((connector IS NULL OR connector != 'test_connector1') OR (connector IS NULL OR connector != 'test_connector2') AND connector ~ 'test_connector3')",
		},
		{
			Expression:    "connector=\"test_connector\"",
			ExpectedQuery: "WHERE e.connector = 'test_connector'",
			Prefix:        "e",
		},
		{
			Expression:    "NOT connector=\"test_connector1\" OR connector=\"test_connector2\" AND connector LIKE \"test_connector3\"",
			ExpectedQuery: "WHERE ((e.connector IS NULL OR e.connector != 'test_connector1') OR e.connector = 'test_connector2' AND e.connector ~ 'test_connector3')",
			Prefix:        "e",
		},
		{
			Expression: "ressete",
			Err:        errors.New("comparison not found"),
		},
		{
			Expression: "NOT connector=\"test_connector1\" OR links=\"test_connector2\" AND connector LIKE \"test_connector3\"",
			Err:        errors.New("field links is not allowed"),
		},
	}
}
