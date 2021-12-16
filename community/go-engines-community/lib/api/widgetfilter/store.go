package widgetfilter

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	FindViewId(ctx context.Context, id string) (string, error)
	GetOneBy(ctx context.Context, id, userId string) (*view.Filter, error)
	Insert(ctx context.Context, r CreateRequest, userId string) (*view.Filter, error)
	Update(ctx context.Context, r EditRequest, userId string) (*view.Filter, error)
	Delete(ctx context.Context, id, userId string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		collection:         dbClient.Collection(mongo.WidgetFiltersMongoCollection),
		widgetCollection:   dbClient.Collection(mongo.WidgetMongoCollection),
		userPrefCollection: dbClient.Collection(mongo.UserPreferencesMongoCollection),
	}
}

type store struct {
	collection         mongo.DbCollection
	widgetCollection   mongo.DbCollection
	userPrefCollection mongo.DbCollection
}

func (s *store) FindViewId(ctx context.Context, id string) (string, error) {
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$lookup": bson.M{
			"from":         mongo.WidgetMongoCollection,
			"localField":   "widget",
			"foreignField": "_id",
			"as":           "widget",
		}},
		{"$unwind": bson.M{"path": "$widget"}},
		{"$lookup": bson.M{
			"from":         mongo.ViewTabMongoCollection,
			"localField":   "widget.tab",
			"foreignField": "_id",
			"as":           "tab",
		}},
		{"$unwind": bson.M{"path": "$tab"}},
		{"$project": bson.M{
			"view": "$tab.view",
		}},
	})
	if err != nil {
		return "", err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		doc := struct {
			View string `bson:"view"`
		}{}
		err = cursor.Decode(&doc)
		if err != nil {
			return "", err
		}

		return doc.View, nil
	}

	return "", nil
}

func (s *store) GetOneBy(ctx context.Context, id, userId string) (*view.Filter, error) {
	model := &view.Filter{}
	err := s.collection.FindOne(ctx, bson.M{"_id": id, "$or": bson.A{
		bson.M{"user": userId},
		bson.M{"user": nil},
	}}).Decode(&model)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return model, nil
}

func (s *store) Insert(ctx context.Context, r CreateRequest, userId string) (*view.Filter, error) {
	now := types.CpsTime{Time: time.Now()}
	filter := transformEditRequestToModel(r.EditRequest)
	filter.ID = utils.NewID()
	filter.Widget = r.Widget
	filter.Created = &now
	filter.Updated = &now

	if r.Personal != nil && *r.Personal {
		filter.User = userId
	}

	_, err := s.collection.InsertOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &filter, nil
}

func (s *store) Update(ctx context.Context, r EditRequest, userId string) (*view.Filter, error) {
	oldFilter, err := s.GetOneBy(ctx, r.ID, userId)
	if err != nil || oldFilter == nil {
		return nil, err
	}

	now := types.CpsTime{Time: time.Now()}
	filter := transformEditRequestToModel(r)
	filter.ID = oldFilter.ID
	filter.Widget = oldFilter.Widget
	filter.User = oldFilter.User
	filter.Created = oldFilter.Created
	filter.Updated = &now

	_, err = s.collection.UpdateOne(ctx, bson.M{
		"_id": filter.ID,
		"$or": bson.A{
			bson.M{"user": userId},
			bson.M{"user": nil},
		},
	}, bson.M{"$set": filter})
	if err != nil {
		return nil, err
	}

	return &filter, nil
}

func (s *store) Delete(ctx context.Context, id, userId string) (bool, error) {
	delCount, err := s.collection.DeleteOne(ctx, bson.M{"_id": id, "$or": bson.A{
		bson.M{"user": userId},
		bson.M{"user": nil},
	}})
	if err != nil {
		return false, err
	}

	if delCount == 0 {
		return false, nil
	}

	err = s.updateWidgets(ctx, id)
	if err != nil {
		return false, err
	}

	err = s.updateUserPreferences(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) updateWidgets(ctx context.Context, filterId string) error {
	_, err := s.widgetCollection.UpdateMany(ctx, bson.M{
		"parameters.main_filter": filterId,
	}, bson.M{
		"$unset": bson.M{"parameters.main_filter": ""},
	})

	return err
}

func (s *store) updateUserPreferences(ctx context.Context, filterId string) error {
	_, err := s.userPrefCollection.UpdateMany(ctx, bson.M{
		"content.main_filter": filterId,
	}, bson.M{
		"$unset": bson.M{"content.main_filter": ""},
	})

	return err
}

func transformEditRequestToModel(request EditRequest) view.Filter {
	return view.Filter{
		Title:  request.Title,
		Query:  request.Query,
		Author: request.Author,
	}
}
