package author

import (
	"encoding/json"
	"regexp"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Author struct {
	ID          string `bson:"_id" json:"_id"`
	Name        string `bson:"name" json:"name"`
	DisplayName string `bson:"display_name" json:"display_name"`
}

type Role struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
}

type Provider interface {
	Pipeline() []bson.M
	PipelineForField(field string) []bson.M
	GetDisplayNameQuery(field string) bson.M
}

func NewProvider(configProvider config.ApiConfigProvider) Provider {
	return &provider{
		configProvider: configProvider,
	}
}

type provider struct {
	configProvider config.ApiConfigProvider
}

func (p *provider) Pipeline() []bson.M {
	return p.PipelineForField("author")
}

func (p *provider) PipelineForField(field string) []bson.M {
	tmpField := getTempField(field)
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.UserCollection,
			"localField":   field,
			"foreignField": "_id",
			"as":           tmpField,
		}},
		{"$unwind": bson.M{"path": "$" + tmpField, "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			field + "._id":      "$" + tmpField + "._id",
			field + ".name":     "$" + tmpField + ".name",
			field + ".username": "$" + tmpField + ".name",
		}},
		// mirror field "username" from AuthorScheme to tmpField, so GetDisplayNameQuery can access it
		{"$addFields": bson.M{
			tmpField + ".username": "$" + tmpField + ".name",
		}},
		{"$addFields": bson.M{
			field + ".display_name": p.GetDisplayNameQuery(tmpField),
		}},
		{"$project": bson.M{tmpField: 0}},
		{"$addFields": bson.M{
			field: bson.M{"$cond": bson.M{
				"if":   "$" + field + "._id",
				"then": "$" + field,
				"else": "$$REMOVE",
			}},
		}},
	}
}

func (p *provider) GetDisplayNameQuery(field string) bson.M {
	authorScheme := p.configProvider.Get().AuthorScheme
	concat := make([]any, len(authorScheme))
	for i, v := range authorScheme {
		if len(v) > 0 && v[0] == '$' {
			f := v
			if field != "" {
				f = "$" + field + "." + v[1:]
			}

			concat[i] = bson.M{"$ifNull": bson.A{f, ""}}
		} else {
			concat[i] = v
		}
	}

	return bson.M{"$concat": concat}
}

// getTempField returns a field name with "_" prefix and 3 random characters
func getTempField(s string) string {
	prefix := "_" + utils.RandString(3)

	if fs := strings.Split(s, "."); len(fs) > 1 {
		return strings.Join(fs[:len(fs)-1], ".") + "." + prefix + fs[len(fs)-1]
	}
	return prefix + s
}

// StripAuthorRandomPrefix removes random prefix from author field in pipeline
// This is required to make unit test's comparison
func StripAuthorRandomPrefix(pipeline []bson.M) []bson.M {
	b, err := json.Marshal(pipeline)
	if err != nil {
		return pipeline
	}
	re := regexp.MustCompile(`\._[0-9a-z]{3}author\b`)
	b = re.ReplaceAll(b, []byte(`.author`))
	var pipeline2 []bson.M
	err = json.Unmarshal(b, &pipeline2)
	if err != nil {
		return pipeline
	}
	return pipeline2
}
