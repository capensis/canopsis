package author

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
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
	Find(ctx context.Context, id string) (Author, error)
	Pipeline() []bson.M
	PipelineForField(field string) []bson.M
	GetDisplayNameQuery(field string) bson.M
}

func NewProvider(client mongo.DbClient, configProvider config.ApiConfigProvider) Provider {
	return &provider{
		collection:     client.Collection(mongo.UserCollection),
		configProvider: configProvider,
	}
}

type provider struct {
	collection     mongo.DbCollection
	configProvider config.ApiConfigProvider
}

func (p *provider) Find(ctx context.Context, id string) (Author, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$addFields": bson.M{
			"username": "$name",
		}},
		{"$addFields": bson.M{
			"display_name": p.GetDisplayNameQuery(""),
		}},
	}
	author := Author{}
	cursor, err := p.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return author, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&author)
		return author, err
	}

	return author, nil
}

func (p *provider) Pipeline() []bson.M {
	return p.PipelineForField("author")
}

func (p *provider) PipelineForField(field string) []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.UserCollection,
			"localField":   field,
			"foreignField": "_id",
			"as":           field,
		}},
		{"$unwind": bson.M{"path": "$" + field, "preserveNullAndEmptyArrays": true}},
		{"$addFields": bson.M{
			field + ".username": "$" + field + ".name",
		}},
		{"$addFields": bson.M{
			field + ".display_name": p.GetDisplayNameQuery(field),
		}},
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
