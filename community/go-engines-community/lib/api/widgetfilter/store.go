package widgetfilter

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(ctx context.Context, r ListRequest, userId string) (*AggregationResult, error)
	FindViewId(ctx context.Context, id string) (string, error)
	FindViewIdByWidget(ctx context.Context, widgetId string) (string, error)
	GetOneBy(ctx context.Context, id, userId string) (*Response, error)
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, id, userId string) (bool, error)
	UpdatePositions(ctx context.Context, filters []string, widgetId, userId string, isPrivate bool) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		client:             dbClient,
		collection:         dbClient.Collection(mongo.WidgetFiltersMongoCollection),
		widgetCollection:   dbClient.Collection(mongo.WidgetMongoCollection),
		userPrefCollection: dbClient.Collection(mongo.UserPreferencesMongoCollection),
	}
}

type store struct {
	client             mongo.DbClient
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

func (s *store) FindViewIdByWidget(ctx context.Context, widgetId string) (string, error) {
	cursor, err := s.widgetCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": widgetId}},
		{"$lookup": bson.M{
			"from":         mongo.ViewTabMongoCollection,
			"localField":   "tab",
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

func (s *store) Find(ctx context.Context, r ListRequest, userId string) (*AggregationResult, error) {
	match := bson.M{"widget": r.Widget}

	if r.Private == nil {
		match["$or"] = []bson.M{
			{"author": userId},
			{"is_private": false},
		}
	} else if *r.Private {
		match["author"] = userId
		match["is_private"] = true
	} else {
		match["is_private"] = false
	}

	pipeline := []bson.M{
		{"$match": match},
	}

	var sort bson.M
	if r.Private == nil {
		sort = bson.M{"$sort": bson.D{
			{Key: "is_private", Value: 1},
			{Key: "position", Value: 1},
			{Key: "_id", Value: 1},
		}}
	} else {
		sort = bson.M{"$sort": bson.D{
			{Key: "position", Value: 1},
			{Key: "_id", Value: 1},
		}}
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		sort,
		author.Pipeline(),
	))

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	res := AggregationResult{}

	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (s *store) GetOneBy(ctx context.Context, id, userId string) (*Response, error) {
	pipeline := []bson.M{
		{"$match": bson.M{
			"_id": id,
			"$or": bson.A{
				bson.M{"author": userId},
				bson.M{"is_private": false},
			}},
		},
	}
	pipeline = append(pipeline, author.Pipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		model := Response{}
		err = cursor.Decode(&model)
		if err != nil {
			return nil, err
		}

		return &model, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r EditRequest) (*Response, error) {
	now := types.NewCpsTime()
	filter := transformEditRequestToModel(r)
	filter.ID = utils.NewID()
	filter.Widget = r.Widget
	filter.Created = now
	filter.Updated = now

	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		position, err := s.getNextPosition(ctx, r.Widget, *r.IsPrivate, r.Author)
		if err != nil {
			return err
		}
		filter.Position = position

		_, err = s.collection.InsertOne(ctx, filter)
		if err != nil {
			return err
		}

		response, err = s.GetOneBy(ctx, filter.ID, r.Author)
		return err
	})

	return response, err
}

func (s *store) Update(ctx context.Context, r EditRequest) (*Response, error) {
	now := types.NewCpsTime()
	filter := transformEditRequestToModel(r)
	filter.ID = r.ID
	filter.Widget = r.Widget
	filter.Updated = now

	update := bson.M{
		"$set": filter,
	}
	if len(filter.EntityPattern) > 0 || len(filter.AlarmPattern) > 0 || len(filter.PbehaviorPattern) > 0 {
		update["$unset"] = bson.M{"old_mongo_query": ""}
	}

	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		oldFilter := view.WidgetFilter{}
		err := s.collection.
			FindOne(ctx, bson.M{"_id": filter.ID}, options.FindOne().SetProjection(bson.M{"position": 1})).
			Decode(&oldFilter)

		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}
			return err
		}
		filter.Position = oldFilter.Position

		_, err = s.collection.UpdateOne(ctx,
			bson.M{"_id": filter.ID},
			update,
		)
		if err != nil {
			return err
		}

		response, err = s.GetOneBy(ctx, r.ID, r.Author)
		return err
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, id, userId string) (bool, error) {
	res := false
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false
		delCount, err := s.collection.DeleteOne(ctx, bson.M{"_id": id, "$or": bson.A{
			bson.M{"author": userId},
			bson.M{"is_private": false},
		}})
		if err != nil {
			return err
		}

		if delCount == 0 {
			return nil
		}

		err = s.updateWidgets(ctx, id)
		if err != nil {
			return err
		}

		err = s.updateUserPreferences(ctx, id)
		if err != nil {
			return err
		}

		res = true
		return nil
	})

	return res, err
}

func (s *store) UpdatePositions(ctx context.Context, ids []string, widgetId, userId string, isPrivate bool) (bool, error) {
	res := false
	notFoundIds := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		notFoundIds[id] = struct{}{}
	}
	if len(ids) != len(notFoundIds) {
		return false, nil
	}

	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false

		match := bson.M{
			"widget":     widgetId,
			"is_private": isPrivate,
		}
		if isPrivate {
			match["author"] = userId
		}

		cursor, err := s.collection.Find(ctx, match)
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			filter := view.WidgetFilter{}
			err = cursor.Decode(&filter)
			if err != nil {
				return err
			}

			if _, ok := notFoundIds[filter.ID]; ok {
				delete(notFoundIds, filter.ID)
			} else {
				return ValidationErr{error: errors.New("filters are related to different widgets or users")}
			}
		}

		if len(notFoundIds) > 0 {
			return ValidationErr{error: errors.New("filters are related to different widgets or users")}
		}

		writeModels := make([]mongodriver.WriteModel, len(ids))
		for i, id := range ids {
			writeModels[i] = mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": id}).
				SetUpdate(bson.M{"$set": bson.M{"position": i}})
		}

		writeRes, err := s.collection.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}

		res = writeRes.MatchedCount > 0
		return nil
	})

	return res, err
}

func (s *store) updateWidgets(ctx context.Context, filterId string) error {
	_, err := s.widgetCollection.UpdateMany(ctx, bson.M{
		"parameters.mainFilter": filterId,
	}, bson.M{
		"$unset": bson.M{"parameters.mainFilter": ""},
	})

	return err
}

func (s *store) updateUserPreferences(ctx context.Context, filterId string) error {
	_, err := s.userPrefCollection.UpdateMany(ctx, bson.M{
		"content.mainFilter": filterId,
	}, bson.M{
		"$unset": bson.M{"content.mainFilter": ""},
	})

	return err
}

func (s *store) getNextPosition(ctx context.Context, widget string, isPrivate bool, user string) (int64, error) {
	match := bson.M{"widget": widget, "is_private": isPrivate}
	if isPrivate {
		match["author"] = user
	}
	cursor, err := s.collection.Aggregate(ctx, []bson.M{
		{"$match": match},
		{"$group": bson.M{
			"_id":      nil,
			"position": bson.M{"$max": "$position"},
		}},
	})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		data := struct {
			Position int64 `bson:"position"`
		}{}
		err = cursor.Decode(&data)
		return data.Position + 1, err
	}

	return 0, nil
}

func transformEditRequestToModel(request EditRequest) view.WidgetFilter {
	return view.WidgetFilter{
		Title:     request.Title,
		IsPrivate: *request.IsPrivate,
		Author:    request.Author,

		AlarmPatternFields:     request.AlarmPatternFieldsRequest.ToModel(),
		EntityPatternFields:    request.EntityPatternFieldsRequest.ToModel(),
		PbehaviorPatternFields: request.PbehaviorPatternFieldsRequest.ToModel(),

		WeatherServicePattern: request.WeatherServicePattern,
	}
}
