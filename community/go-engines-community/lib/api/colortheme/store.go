package colortheme

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	Canopsis       = "canopsis"
	CanopsisDark   = "canopsis_dark"
	ColorBlind     = "color_blind"
	ColorBlindDark = "color_blind_dark"
)

type Store interface {
	Insert(ctx context.Context, r EditRequest) (*Theme, error)
	GetById(ctx context.Context, id string) (*Theme, error)
	Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error)
	Update(ctx context.Context, r EditRequest) (*Theme, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type store struct {
	dbColorCollection libmongo.DbCollection
	dbUserCollection  libmongo.DbCollection

	defaultSearchByFields []string
	defaultThemeIDs       map[string]struct{}
}

func NewStore(
	dbClient libmongo.DbClient,
) Store {
	return &store{
		dbColorCollection:     dbClient.Collection(libmongo.ColorThemeCollection),
		dbUserCollection:      dbClient.Collection(libmongo.UserCollection),
		defaultSearchByFields: []string{"_id", "name"},
		defaultThemeIDs: map[string]struct{}{
			Canopsis:       {},
			CanopsisDark:   {},
			ColorBlind:     {},
			ColorBlindDark: {},
		},
	}
}

func (s *store) Insert(ctx context.Context, r EditRequest) (*Theme, error) {
	theme := s.transformRequestToDocument(r)
	theme.ID = utils.NewID()

	_, err := s.dbColorCollection.InsertOne(ctx, theme)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, common.NewValidationError("name", "Name already exists.")
		}

		return nil, err
	}

	return &theme, nil
}

func (s *store) GetById(ctx context.Context, id string) (*Theme, error) {
	var res Theme
	err := s.dbColorCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &res, nil
}

func (s *store) Find(ctx context.Context, query FilteredQuery) (*AggregationResult, error) {
	var pipeline []bson.M

	filter := common.GetSearchQuery(query.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := "name"
	if query.SortBy != "" {
		sortBy = query.SortBy
	}

	cursor, err := s.dbColorCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		query.Query,
		pipeline,
		common.GetSortQuery(sortBy, query.Sort),
	))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var result AggregationResult
	if cursor.Next(ctx) {
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (s *store) Update(ctx context.Context, r EditRequest) (*Theme, error) {
	if s.isDefaultTheme(r.ID) {
		return nil, ErrDefaultTheme
	}

	theme := s.transformRequestToDocument(r)
	theme.ID = r.ID

	res, err := s.dbColorCollection.UpdateOne(ctx, bson.M{"_id": theme.ID}, bson.M{"$set": theme})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, common.NewValidationError("name", "Name already exists.")
		}

		return nil, err
	}

	if res.MatchedCount == 0 {
		return nil, nil
	}

	return &theme, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	if s.isDefaultTheme(id) {
		return false, ErrDefaultTheme
	}

	deleted, err := s.dbColorCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil || deleted == 0 {
		return false, err
	}

	_, err = s.dbUserCollection.UpdateMany(ctx, bson.M{"ui_theme": id}, bson.M{"$set": bson.M{"ui_theme": Canopsis}})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) transformRequestToDocument(r EditRequest) Theme {
	return Theme{
		Name:      r.Name,
		Colors:    r.Colors,
		FontSize:  r.FontSize,
		Updated:   datetime.NewCpsTime(),
		Deletable: true,
	}
}

func (s *store) isDefaultTheme(id string) bool {
	_, ok := s.defaultThemeIDs[id]
	return ok
}
