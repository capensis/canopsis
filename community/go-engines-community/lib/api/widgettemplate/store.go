package widgettemplate

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetOneById(ctx context.Context, id string) (*Response, error)
	Insert(ctx context.Context, r EditRequest) (*Response, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	Delete(ctx context.Context, id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient, authorProvider author.Provider) Store {
	return &store{
		client:           dbClient,
		collection:       dbClient.Collection(mongo.WidgetTemplateMongoCollection),
		widgetCollection: dbClient.Collection(mongo.WidgetMongoCollection),
		authorProvider:   authorProvider,

		widgetParameters: view.GetWidgetTemplateParameters(),

		defaultSearchByFields: []string{"_id", "title", "type", "author.name"},
		defaultSortBy:         "created",
	}
}

type store struct {
	client           mongo.DbClient
	collection       mongo.DbCollection
	widgetCollection mongo.DbCollection
	authorProvider   author.Provider

	widgetParameters map[string]map[string][]string

	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)
	if r.Type != "" {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"type": r.Type}})
	}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)
	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := s.defaultSortBy
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		common.GetSortQuery(sortBy, r.Sort),
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

func (s *store) GetOneById(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, s.authorProvider.Pipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		res := Response{}
		err = cursor.Decode(&res)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r EditRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	model := transformEditRequestToModel(r)
	model.ID = utils.NewID()
	model.Created = now
	model.Updated = now

	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		_, err := s.collection.InsertOne(ctx, model)
		if err != nil {
			return err
		}

		response, err = s.GetOneById(ctx, model.ID)
		return err
	})

	return response, err
}

func (s *store) Update(ctx context.Context, r EditRequest) (*Response, error) {
	now := datetime.NewCpsTime()
	model := transformEditRequestToModel(r)
	model.ID = r.ID
	model.Updated = now

	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		_, err := s.collection.UpdateOne(ctx,
			bson.M{"_id": model.ID},
			bson.M{"$set": model},
		)
		if err != nil {
			return err
		}

		response, err = s.GetOneById(ctx, r.ID)
		if err != nil || response == nil {
			return err
		}

		return s.updateLinkedWidgets(ctx, *response)
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	res := false
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false

		model := view.WidgetTemplate{}
		err := s.collection.FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(&model)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}
			return err
		}

		err = s.cleanLinkedWidgets(ctx, model)
		if err != nil {
			return err
		}

		res = true
		return nil
	})

	return res, err
}

func (s *store) updateLinkedWidgets(ctx context.Context, tpl Response) error {
	for widgetType, parametersByType := range s.widgetParameters {
		parameters := parametersByType[tpl.Type]
		for _, parameter := range parameters {
			var val any
			switch tpl.Type {
			case view.WidgetTemplateTypeAlarmColumns,
				view.WidgetTemplateTypeEntityColumns:
				val = tpl.Columns
			case view.WidgetTemplateTypeAlarmMoreInfos,
				view.WidgetTemplateTypeAlarmExportToPDF,
				view.WidgetTemplateTypeServiceWeatherItem,
				view.WidgetTemplateTypeServiceWeatherModal,
				view.WidgetTemplateTypeServiceWeatherEntity:
				val = tpl.Content
			}

			_, err := s.widgetCollection.UpdateMany(
				ctx,
				bson.M{
					"type":                                 widgetType,
					"parameters." + parameter + "Template": tpl.ID,
				},
				bson.M{"$set": bson.M{
					"parameters." + parameter:                   val,
					"parameters." + parameter + "TemplateTitle": tpl.Title,
				}},
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *store) cleanLinkedWidgets(ctx context.Context, model view.WidgetTemplate) error {
	for widgetType, parametersByType := range s.widgetParameters {
		parameters := parametersByType[model.Type]
		for _, parameter := range parameters {
			_, err := s.widgetCollection.UpdateMany(
				ctx,
				bson.M{
					"type":                                 widgetType,
					"parameters." + parameter + "Template": model.ID,
				},
				bson.M{"$unset": bson.M{
					"parameters." + parameter + "Template":      "",
					"parameters." + parameter + "TemplateTitle": "",
				}},
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func transformEditRequestToModel(r EditRequest) view.WidgetTemplate {
	return view.WidgetTemplate{
		Title:   r.Title,
		Type:    r.Type,
		Columns: r.Columns,
		Content: r.Content,
		Author:  r.Author,
	}
}
