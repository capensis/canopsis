package view

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewtab"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const (
	permissionPrefix = "Rights on view :"
	defaultTabTitle  = "Default"
)

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetOneBy(ctx context.Context, id string) (*Response, error)
	Insert(ctx context.Context, r EditRequest, withDefaultTab bool) (*Response, error)
	Update(ctx context.Context, r EditRequest) (*Response, error)
	// UpdatePositions receives some groups and views with updated positions and updates
	// positions for all groups and views in db and moves views to another groups if necessary.
	UpdatePositions(ctx context.Context, r EditPositionRequest) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	Copy(ctx context.Context, id string, r EditRequest) (*Response, error)
	Export(ctx context.Context, r ExportRequest) (ExportResponse, error)
	Import(ctx context.Context, r ImportRequest, userId string) error
}

func NewStore(dbClient mongo.DbClient, tabStore viewtab.Store) Store {
	return &store{
		client:                dbClient,
		collection:            dbClient.Collection(mongo.ViewMongoCollection),
		tabCollection:         dbClient.Collection(mongo.ViewTabMongoCollection),
		widgetCollection:      dbClient.Collection(mongo.WidgetMongoCollection),
		filterCollection:      dbClient.Collection(mongo.WidgetFiltersMongoCollection),
		groupCollection:       dbClient.Collection(mongo.ViewGroupMongoCollection),
		aclCollection:         dbClient.Collection(mongo.RightsMongoCollection),
		userPrefCollection:    dbClient.Collection(mongo.UserPreferencesMongoCollection),
		defaultSearchByFields: []string{"_id", "title", "description", "author"},
		defaultSortBy:         "position",

		tabStore: tabStore,
	}
}

type store struct {
	client                mongo.DbClient
	collection            mongo.DbCollection
	tabCollection         mongo.DbCollection
	widgetCollection      mongo.DbCollection
	filterCollection      mongo.DbCollection
	groupCollection       mongo.DbCollection
	aclCollection         mongo.DbCollection
	userPrefCollection    mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string

	tabStore viewtab.Store
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := make([]bson.M, 0)

	if len(r.Ids) > 0 {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"_id": bson.M{"$in": r.Ids}}})
	}

	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	pipeline = append(pipeline, []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.ViewGroupMongoCollection,
			"localField":   "group_id",
			"foreignField": "_id",
			"as":           "group",
		}},
		{"$unwind": bson.M{"path": "$group", "preserveNullAndEmptyArrays": true}},
	}...)
	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		common.GetSortQuery("position", common.SortAsc),
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

func (s *store) GetOneBy(ctx context.Context, id string) (*Response, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		model := &Response{}
		err := cursor.Decode(model)
		if err != nil {
			return nil, err
		}

		return model, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, r EditRequest, withDefaultTab bool) (*Response, error) {
	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		count, err := s.collection.CountDocuments(ctx, bson.M{})
		if err != nil {
			return err
		}

		now := types.CpsTime{Time: time.Now()}
		id := utils.NewID()
		_, err = s.collection.InsertOne(ctx, view.View{
			ID:              id,
			Enabled:         *r.Enabled,
			Title:           r.Title,
			Description:     r.Description,
			Group:           r.Group,
			Position:        count,
			Tags:            r.Tags,
			PeriodicRefresh: r.PeriodicRefresh,
			Author:          r.Author,
			Created:         now,
			Updated:         now,
		})
		if err != nil {
			return err
		}

		if withDefaultTab {
			_, err := s.tabCollection.InsertOne(ctx, view.Tab{
				ID:       utils.NewID(),
				Title:    defaultTabTitle,
				View:     id,
				Author:   r.Author,
				Position: 0,
				Created:  now,
				Updated:  now,
			})
			if err != nil {
				return err
			}
		}

		response, err = s.GetOneBy(ctx, id)
		if err != nil {
			return err
		}

		return s.createPermissions(ctx, r.Author, map[string]string{response.ID: response.Title})
	})

	return response, err
}

func (s *store) Update(ctx context.Context, r EditRequest) (*Response, error) {
	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil

		oldView := view.View{}
		err := s.collection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&oldView)
		if err != nil {
			return err
		}

		now := types.CpsTime{Time: time.Now()}
		_, err = s.collection.UpdateOne(ctx, bson.M{"_id": oldView.ID}, bson.M{"$set": bson.M{
			"enabled":          *r.Enabled,
			"title":            r.Title,
			"description":      r.Description,
			"group_id":         r.Group,
			"tags":             r.Tags,
			"periodic_refresh": r.PeriodicRefresh,
			"author":           r.Author,
			"updated":          now,
		}})
		if err != nil {
			return err
		}

		if r.Group != oldView.Group {
			err := s.normalizePositionsOnViewMove(ctx, r.ID, r.Group)
			if err != nil {
				return err
			}
		}

		response, err = s.GetOneBy(ctx, r.ID)
		if err != nil {
			return err
		}

		return s.updatePermissions(ctx, *response)
	})

	return response, err
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	res := false
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false
		delCount, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil || delCount == 0 {
			return err
		}

		err = s.deletePermissions(ctx, id)
		if err != nil {
			return err
		}

		err = s.deleteTabs(ctx, id)
		if err != nil {
			return err
		}

		res = true
		return nil
	})

	return res, err
}

func (s *store) Copy(ctx context.Context, id string, r EditRequest) (*Response, error) {
	var response *Response
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		response = nil
		v, err := s.GetOneBy(ctx, id)
		if err != nil || v == nil {
			return err
		}

		newView, err := s.Insert(ctx, r, false)
		if err != nil {
			return err
		}

		err = s.tabStore.CopyForView(ctx, v.ID, newView.ID, r.Author)
		if err != nil {
			return err
		}

		response = newView
		return nil
	})

	return response, err
}

func (s *store) UpdatePositions(ctx context.Context, r EditPositionRequest) (bool, error) {
	res := false
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		res = false
		var err error
		res, err = s.updatePositions(ctx, r)
		return err
	})

	return res, err
}

func (s *store) updatePositions(ctx context.Context, r EditPositionRequest) (bool, error) {
	groupPositions, viewPositionsByGroup, err := s.findViewPositions(ctx)
	if err != nil || len(groupPositions) == 0 {
		return false, err
	}

	newGroupPositions, newViewPositions, changedViewGroup := computePositions(r,
		groupPositions, viewPositionsByGroup)
	if len(newGroupPositions) == 0 {
		return false, nil
	}

	err = s.writePositions(ctx, newGroupPositions, newViewPositions, changedViewGroup)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) Export(ctx context.Context, r ExportRequest) (ExportResponse, error) {
	groups := make([]ExportViewGroupResponse, 0)
	views := make([]Response, 0)

	nestedObjectsPipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.ViewTabMongoCollection,
			"localField":   "_id",
			"foreignField": "view",
			"as":           "tabs",
		}},
		{"$unwind": bson.M{"path": "$tabs", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.WidgetMongoCollection,
			"localField":   "tabs._id",
			"foreignField": "tab",
			"as":           "widgets",
		}},
		{"$unwind": bson.M{"path": "$widgets", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from": mongo.WidgetFiltersMongoCollection,
			"let":  bson.M{"id": "$widgets._id"},
			"pipeline": []bson.M{
				{"$match": bson.M{"$and": []bson.M{
					{"$expr": bson.M{"$eq": bson.A{"$widget", "$$id"}}},
					{"is_private": false},
					{"old_mongo_query": nil}, //do not import old filters
				}}},
			},
			"as": "filters",
		}},
		{"$unwind": bson.M{"path": "$filters", "preserveNullAndEmptyArrays": true}},
		{"$project": bson.M{
			"filters.corporate_alarm_pattern":           0,
			"filters.corporate_alarm_pattern_title":     0,
			"filters.corporate_entity_pattern":          0,
			"filters.corporate_entity_pattern_title":    0,
			"filters.corporate_pbehavior_pattern":       0,
			"filters.corporate_pbehavior_pattern_title": 0,
			"filters.is_private":                        0,
			"filters.author":                            0,
			"filters.updated":                           0,
			"filters.created":                           0,
		}},
		{"$sort": bson.M{"filters.title": 1}},
		{"$group": bson.M{
			"_id": bson.M{
				"_id":    "$_id",
				"tab":    "$tabs._id",
				"widget": "$widgets._id",
			},
			"data":    bson.M{"$first": "$$ROOT"},
			"tabs":    bson.M{"$first": "$tabs"},
			"widgets": bson.M{"$first": "$widgets"},
			"filters": bson.M{"$push": "$filters"},
		}},
		{"$addFields": bson.M{
			"_id":             "$_id._id",
			"widgets.filters": "$filters",
		}},
		{"$project": bson.M{
			"widgets._id":     0,
			"widgets.author":  0,
			"widgets.updated": 0,
			"widgets.created": 0,
		}},
		{"$sort": bson.D{
			{Key: "widgets.grid_parameters.desktop.y", Value: 1},
			{Key: "widgets.grid_parameters.desktop.x", Value: 1},
		}},
		{"$group": bson.M{
			"_id": bson.M{
				"_id": "$_id",
				"tab": "$tabs._id",
			},
			"data":    bson.M{"$first": "$data"},
			"tabs":    bson.M{"$first": "$tabs"},
			"widgets": bson.M{"$push": "$widgets"},
		}},
		{"$addFields": bson.M{
			"tabs.widgets": bson.M{"$filter": bson.M{
				"input": "$widgets",
				"cond":  "$$this.type",
			}},
		}},
		{"$sort": bson.M{"tabs.position": 1}},
		{"$project": bson.M{
			"tabs._id":      0,
			"tabs.author":   0,
			"tabs.updated":  0,
			"tabs.created":  0,
			"tabs.position": 0,
		}},
		{"$group": bson.M{
			"_id":  "$_id._id",
			"data": bson.M{"$first": "$data"},
			"tabs": bson.M{"$push": "$tabs"},
		}},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{"tabs": bson.M{"$filter": bson.M{
				"input": "$tabs",
				"cond":  "$$this.title",
			}}},
		}}}},
	}

	if len(r.Views) > 0 {
		pipeline := []bson.M{{"$match": bson.M{"_id": bson.M{"$in": r.Views}}}}
		pipeline = append(pipeline, nestedObjectsPipeline...)
		pipeline = append(pipeline, common.GetSortQuery("position", common.SortAsc))
		pipeline = append(pipeline, bson.M{"$project": bson.M{
			"_id":      0,
			"author":   0,
			"updated":  0,
			"created":  0,
			"position": 0,
		}})
		cursor, err := s.collection.Aggregate(ctx, pipeline)
		if err != nil {
			return ExportResponse{}, err
		}

		err = cursor.All(ctx, &views)
		if err != nil {
			return ExportResponse{}, err
		}

		if len(views) != len(r.Views) {
			return ExportResponse{}, ValidationError{field: "views", error: fmt.Errorf("views not found")}
		}
	}
	if len(r.Groups) > 0 {
		groupIds := make([]string, len(r.Groups))
		viewsByGroup := make(map[string]map[string]bool, len(r.Groups))
		for i, group := range r.Groups {
			groupIds[i] = group.ID
			viewsByGroup[group.ID] = make(map[string]bool, len(group.Views))
			for _, v := range group.Views {
				viewsByGroup[group.ID][v] = true
			}
		}

		pipeline := []bson.M{{"$match": bson.M{"_id": bson.M{"$in": groupIds}}}}
		pipeline = append(pipeline, []bson.M{
			{"$lookup": bson.M{
				"from":         mongo.ViewMongoCollection,
				"localField":   "_id",
				"foreignField": "group_id",
				"as":           "views",
			}},
			{"$unwind": bson.M{"path": "$views", "preserveNullAndEmptyArrays": true}},
			{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
				"$views",
				bson.M{"group_obj": "$$ROOT"},
			}}}},
		}...)
		pipeline = append(pipeline, nestedObjectsPipeline...)
		pipeline = append(pipeline, common.GetSortQuery("position", common.SortAsc))
		pipeline = append(pipeline, []bson.M{
			{"$project": bson.M{
				"author":   0,
				"updated":  0,
				"created":  0,
				"position": 0,
			}},
			{"$group": bson.M{
				"_id":   "$group_obj._id",
				"group": bson.M{"$first": "$group_obj"},
				"views": bson.M{"$push": "$$ROOT"},
			}},
			{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
				"$group",
				bson.M{"views": bson.M{"$filter": bson.M{
					"input": "$views",
					"cond":  "$$this.title",
				}}},
			}}}},
			{"$project": bson.M{
				"author":   0,
				"updated":  0,
				"created":  0,
				"position": 0,
			}},
		}...)
		pipeline = append(pipeline, common.GetSortQuery("position", common.SortAsc))
		cursor, err := s.groupCollection.Aggregate(ctx, pipeline)
		if err != nil {
			return ExportResponse{}, err
		}

		err = cursor.All(ctx, &groups)
		if err != nil {
			return ExportResponse{}, err
		}

		if len(groups) != len(r.Groups) {
			return ExportResponse{}, ValidationError{field: "groups", error: fmt.Errorf("groups not found")}
		}

		for i, group := range groups {
			foundViews := make([]Response, 0, len(viewsByGroup[group.ID]))
			for _, v := range group.Views {
				if viewsByGroup[group.ID][v.ID] {
					v.ID = ""
					foundViews = append(foundViews, v)
				}
			}
			groups[i].Views = foundViews
			if len(groups[i].Views) != len(viewsByGroup[group.ID]) {
				return ExportResponse{}, ValidationError{field: fmt.Sprintf("groups.%d", i), error: fmt.Errorf("views not found")}
			}

			groups[i].ID = ""
		}
	}

	return ExportResponse{
		Groups: groups,
		Views:  views,
	}, nil
}

func (s *store) Import(ctx context.Context, r ImportRequest, userId string) error {
	err := s.client.WithTransaction(ctx, func(ctx context.Context) error {
		maxViewPosition, err := s.collection.CountDocuments(ctx, bson.M{})
		if err != nil {
			return err
		}
		maxGroupPosition, err := s.groupCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			return err
		}
		groupIds := make([]string, 0, len(r.Items))
		viewIds := make([]string, 0, len(r.Items))

		for _, g := range r.Items {
			if g.ID != "" {
				groupIds = append(groupIds, g.ID)
			}
			if g.Views != nil {
				for _, v := range g.Views {
					if v.ID != "" {
						viewIds = append(viewIds, v.ID)
					}
				}
			}
		}

		existedViewIds := make(map[string]bool)
		existedGroupIds := make(map[string]bool)
		if len(viewIds) > 0 {
			cursor, err := s.collection.Find(ctx, bson.M{"_id": bson.M{"$in": viewIds}})
			if err != nil {
				return err
			}
			defer cursor.Close(ctx)
			for cursor.Next(ctx) {
				model := struct {
					ID string `bson:"_id"`
				}{}
				err := cursor.Decode(&model)
				if err != nil {
					return err
				}
				existedViewIds[model.ID] = true
			}
		}
		if len(groupIds) > 0 {
			cursor, err := s.groupCollection.Find(ctx, bson.M{"_id": bson.M{"$in": groupIds}})
			if err != nil {
				return err
			}
			defer cursor.Close(ctx)
			for cursor.Next(ctx) {
				model := struct {
					ID string `bson:"_id"`
				}{}
				err := cursor.Decode(&model)
				if err != nil {
					return err
				}
				existedGroupIds[model.ID] = true
			}
		}

		newGroups := make([]interface{}, 0, len(r.Items))
		newViews := make([]interface{}, 0, len(r.Items))
		newTabs := make([]interface{}, 0, len(r.Items))
		newWidgets := make([]interface{}, 0, len(r.Items))
		newWidgetFilters := make([]interface{}, 0, len(r.Items))
		newViewTitles := make(map[string]string, len(r.Items))
		positionItems := make([]EditPositionItemRequest, 0, len(r.Items))
		now := types.NewCpsTime()
		for gi, g := range r.Items {
			groupId := g.ID

			if g.ID == "" || !existedGroupIds[g.ID] {
				groupId = utils.NewID()
				if g.Title == "" {
					return ValidationError{
						field: fmt.Sprintf("%d.title", gi),
						error: fmt.Errorf("value is missing"),
					}
				}
				newGroups = append(newGroups, view.Group{
					ID:       groupId,
					Title:    g.Title,
					Position: maxGroupPosition,
					Author:   userId,
					Created:  now,
					Updated:  now,
				})
				maxGroupPosition++
			}

			groupViewIds := make([]string, 0)

			if g.Views != nil {
				for vi, v := range g.Views {
					if v.ID != "" && existedViewIds[v.ID] {
						groupViewIds = append(groupViewIds, v.ID)
						continue
					}

					if v.Title == "" {
						return ValidationError{
							field: fmt.Sprintf("%d.views.%d.title", gi, vi),
							error: fmt.Errorf("value is missing"),
						}
					}

					viewId := utils.NewID()
					groupViewIds = append(groupViewIds, viewId)
					newViews = append(newViews, view.View{
						ID:              viewId,
						Enabled:         v.Enabled,
						Title:           v.Title,
						Description:     v.Description,
						Position:        maxViewPosition,
						Group:           groupId,
						Tags:            v.Tags,
						PeriodicRefresh: v.PeriodicRefresh,
						Author:          userId,
						Created:         now,
						Updated:         now,
					})
					maxViewPosition++
					newViewTitles[viewId] = v.Title

					if v.Tabs != nil {
						for ti, tab := range *v.Tabs {
							if tab.Title == "" {
								return ValidationError{
									field: fmt.Sprintf("%d.views.%d.tabs.%d.title", gi, vi, ti),
									error: fmt.Errorf("value is missing"),
								}
							}

							tabId := utils.NewID()
							newTabs = append(newTabs, view.Tab{
								ID:       tabId,
								Title:    tab.Title,
								View:     viewId,
								Author:   userId,
								Position: int64(ti),
								Created:  now,
								Updated:  now,
							})

							if tab.Widgets != nil {
								for wi, widget := range *tab.Widgets {
									if widget.Type == "" {
										return ValidationError{
											field: fmt.Sprintf("%d.views.%d.tabs.%d.widgets.%d.type", gi, vi, ti, wi),
											error: fmt.Errorf("value is missing"),
										}
									}

									widgetId := utils.NewID()
									mainFilterId := ""

									for fi, filter := range widget.Filters {
										if filter.Title == "" {
											return ValidationError{
												field: fmt.Sprintf("%d.views.%d.tabs.%d.widgets.%d.filters.%d.title", gi, vi, ti, wi, fi),
												error: fmt.Errorf("value is missing"),
											}
										}
										if len(filter.AlarmPattern) == 0 && len(filter.EntityPattern) == 0 && len(filter.PbehaviorPattern) == 0 {
											return ValidationError{
												field: fmt.Sprintf("%d.views.%d.tabs.%d.widgets.%d.filters.%d.alarm_pattern", gi, vi, ti, wi, fi),
												error: fmt.Errorf("value is missing"),
											}
										}

										filterId := utils.NewID()
										newWidgetFilters = append(newWidgetFilters, view.WidgetFilter{
											ID:        filterId,
											Title:     filter.Title,
											Widget:    widgetId,
											IsPrivate: false,
											AlarmPatternFields: savedpattern.AlarmPatternFields{
												AlarmPattern: filter.AlarmPattern,
											},
											EntityPatternFields: savedpattern.EntityPatternFields{
												EntityPattern: filter.EntityPattern,
											},
											PbehaviorPatternFields: savedpattern.PbehaviorPatternFields{
												PbehaviorPattern: filter.PbehaviorPattern,
											},
											Author:  userId,
											Created: now,
											Updated: now,
										})

										if widget.Parameters.MainFilter != "" && filter.ID == widget.Parameters.MainFilter {
											mainFilterId = filterId
										}
									}

									widget.Parameters.MainFilter = mainFilterId
									newWidgets = append(newWidgets, view.Widget{
										ID:             widgetId,
										Tab:            tabId,
										Title:          widget.Title,
										Type:           widget.Type,
										GridParameters: widget.GridParameters,
										Parameters:     widget.Parameters,
										Author:         userId,
										Created:        now,
										Updated:        now,
									})
								}
							}
						}
					}
				}
			}

			positionItems = append(positionItems, EditPositionItemRequest{
				ID:    groupId,
				Views: groupViewIds,
			})
		}
		if len(newGroups) > 0 {
			_, err := s.groupCollection.InsertMany(ctx, newGroups)
			if err != nil {
				return err
			}
		}
		if len(newViews) > 0 {
			_, err := s.collection.InsertMany(ctx, newViews)
			if err != nil {
				return err
			}
		}
		if len(newTabs) > 0 {
			_, err := s.tabCollection.InsertMany(ctx, newTabs)
			if err != nil {
				return err
			}
		}
		if len(newWidgets) > 0 {
			_, err := s.widgetCollection.InsertMany(ctx, newWidgets)
			if err != nil {
				return err
			}
		}
		if len(newWidgetFilters) > 0 {
			_, err := s.filterCollection.InsertMany(ctx, newWidgetFilters)
			if err != nil {
				return err
			}
		}

		err = s.createPermissions(ctx, userId, newViewTitles)
		if err != nil {
			return err
		}

		_, err = s.updatePositions(ctx, EditPositionRequest{Items: positionItems})
		return err
	})

	return err
}

func (s *store) createPermissions(ctx context.Context, userID string, views map[string]string) error {
	if len(views) == 0 {
		return nil
	}

	newPermissions := make([]interface{}, 0, len(views))
	setRole := bson.M{}
	for viewId, viewTitle := range views {
		newPermissions = append(newPermissions, bson.M{
			"_id":          viewId,
			"crecord_name": viewId,
			"crecord_type": securitymodel.LineTypeObject,
			"desc":         fmt.Sprintf("%s %s", permissionPrefix, viewTitle),
			"type":         securitymodel.LineObjectTypeRW,
		})
		setRole["rights."+viewId] = bson.M{
			"checksum": securitymodel.PermissionBitmaskRead |
				securitymodel.PermissionBitmaskUpdate |
				securitymodel.PermissionBitmaskDelete,
		}
	}
	_, err := s.aclCollection.InsertMany(ctx, newPermissions)
	if err != nil {
		return err
	}

	res := s.aclCollection.FindOne(ctx, bson.M{
		"_id":          userID,
		"crecord_type": securitymodel.LineTypeSubject,
	})
	if err := res.Err(); err != nil {
		return err
	}

	user := struct {
		Role string `json:"role"`
	}{}
	err = res.Decode(&user)
	if err != nil {
		return err
	}

	_, err = s.aclCollection.UpdateOne(ctx,
		bson.M{
			"_id":          user.Role,
			"crecord_type": securitymodel.LineTypeRole,
		},
		bson.M{"$set": setRole},
	)
	if err != nil {
		return err
	}

	return err
}

func (s *store) updatePermissions(ctx context.Context, view Response) error {
	_, err := s.aclCollection.UpdateOne(ctx,
		bson.M{
			"_id":          view.ID,
			"crecord_type": securitymodel.LineTypeObject,
		},
		bson.M{"$set": bson.M{
			"desc": fmt.Sprintf("%s %s", permissionPrefix, view.Title),
		}},
	)

	return err
}

func (s *store) deletePermissions(ctx context.Context, viewID string) error {
	_, err := s.aclCollection.UpdateMany(ctx,
		bson.M{"crecord_type": securitymodel.LineTypeRole},
		bson.M{"$unset": bson.M{
			"rights." + viewID: "",
		}},
	)
	if err != nil {
		return err
	}

	_, err = s.aclCollection.DeleteOne(ctx, bson.M{
		"_id":          viewID,
		"crecord_type": securitymodel.LineTypeObject,
	})

	return err
}

func (s *store) normalizePositionsOnViewMove(ctx context.Context, viewID, groupID string) error {
	groupPositions, viewPositionsByGroup, err := s.findViewPositions(ctx)
	if err != nil {
		return err
	}

	index := -1
	for i, v := range viewPositionsByGroup[groupID] {
		if v == viewID {
			index = i
		}
	}

	viewPositionsByGroup[groupID] = append(viewPositionsByGroup[groupID][:index], viewPositionsByGroup[groupID][index+1:]...)
	viewPositionsByGroup[groupID] = append(viewPositionsByGroup[groupID], viewID)

	viewPositions := make([]string, 0)
	for id := range viewPositionsByGroup {
		viewPositions = append(viewPositions, viewPositionsByGroup[id]...)
	}

	return s.writePositions(ctx, groupPositions, viewPositions, nil)
}

func (s *store) findViewPositions(ctx context.Context) ([]string, map[string][]string, error) {
	cursor, err := s.groupCollection.Aggregate(ctx, []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "_id",
			"foreignField": "group_id",
			"as":           "views",
		}},
		{"$unwind": bson.M{"path": "$views", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"views.position": 1}},
		{"$group": bson.M{
			"_id":      "$_id",
			"position": bson.M{"$first": "$position"},
			"views":    bson.M{"$push": "$views._id"},
		}},
		{"$sort": bson.M{"position": 1}},
	})
	if err != nil {
		return nil, nil, err
	}

	res := make([]struct {
		ID    string   `bson:"_id"`
		Views []string `bson:"views"`
	}, 0)

	if err := cursor.All(ctx, &res); err != nil {
		return nil, nil, err
	}

	if len(res) == 0 {
		return nil, nil, nil
	}

	groupPositions := make([]string, len(res))
	viewPositionsByGroup := make(map[string][]string, len(res))
	for i, group := range res {
		groupPositions[i] = group.ID
		viewPositionsByGroup[group.ID] = make([]string, len(group.Views))
		copy(viewPositionsByGroup[group.ID], group.Views)
	}

	return groupPositions, viewPositionsByGroup, nil
}

func (s *store) writePositions(
	ctx context.Context,
	groups, views []string,
	changedViewGroup map[string]string,
) error {
	if len(groups) > 0 {
		groupModels := make([]mongodriver.WriteModel, len(groups))
		for position, id := range groups {
			groupModels[position] = mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": id}).
				SetUpdate(bson.M{"$set": bson.M{"position": position}})
		}

		_, err := s.groupCollection.BulkWrite(ctx, groupModels)
		if err != nil {
			return err
		}
	}

	if len(views) > 0 {
		models := make([]mongodriver.WriteModel, len(views))
		for position, id := range views {
			update := bson.M{"position": position}

			if groupID, ok := changedViewGroup[id]; ok {
				update["group_id"] = groupID
			}

			models[position] = mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": id}).
				SetUpdate(bson.M{"$set": update})
		}

		_, err := s.collection.BulkWrite(ctx, models)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) deleteTabs(ctx context.Context, id string) error {
	tabCursor, err := s.tabCollection.Find(ctx, bson.M{"view": id})
	if err != nil {
		return err
	}
	tabs := make([]view.Tab, 0)
	err = tabCursor.All(ctx, &tabs)
	if err != nil || len(tabs) == 0 {
		return err
	}

	tabIds := make([]string, len(tabs))
	for i, tab := range tabs {
		tabIds[i] = tab.ID
	}

	_, err = s.tabCollection.DeleteMany(ctx, bson.M{"view": id})
	if err != nil {
		return err
	}

	widgetCursor, err := s.widgetCollection.Find(ctx, bson.M{"tab": bson.M{"$in": tabIds}})
	if err != nil {
		return err
	}
	widgets := make([]view.Widget, 0)
	err = widgetCursor.All(ctx, &widgets)
	if err != nil || len(widgets) == 0 {
		return err
	}

	widgetIds := make([]string, len(widgets))
	for i, w := range widgets {
		widgetIds[i] = w.ID
	}

	_, err = s.widgetCollection.DeleteMany(ctx, bson.M{"tab": bson.M{"$in": tabIds}})
	if err != nil {
		return err
	}

	_, err = s.userPrefCollection.DeleteMany(ctx, bson.M{"widget": bson.M{"$in": widgetIds}})
	if err != nil {
		return err
	}

	return nil
}

func getNestedObjectsPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.ViewTabMongoCollection,
			"localField":   "_id",
			"foreignField": "view",
			"as":           "tabs",
		}},
		{"$unwind": bson.M{"path": "$tabs", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.WidgetMongoCollection,
			"localField":   "tabs._id",
			"foreignField": "tab",
			"as":           "widgets",
		}},
		{"$unwind": bson.M{"path": "$widgets", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.WidgetFiltersMongoCollection,
			"localField":   "widgets._id",
			"foreignField": "widget",
			"as":           "filters",
		}},
		{"$unwind": bson.M{"path": "$filters", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"filters.title": 1}},
		{"$group": bson.M{
			"_id": bson.M{
				"_id":    "$_id",
				"tab":    "$tabs._id",
				"widget": "$widgets._id",
			},
			"data":    bson.M{"$first": "$$ROOT"},
			"tabs":    bson.M{"$first": "$tabs"},
			"widgets": bson.M{"$first": "$widgets"},
			"filters": bson.M{"$push": "$filters"},
		}},
		{"$addFields": bson.M{
			"_id": "$_id._id",
			"widgets.filters": bson.M{"$filter": bson.M{
				"input": "$filters",
				"cond":  bson.M{"$eq": bson.A{"$$this.is_private", false}},
			}},
		}},
		{"$sort": bson.D{
			{Key: "widgets.grid_parameters.desktop.y", Value: 1},
			{Key: "widgets.grid_parameters.desktop.x", Value: 1},
		}},
		{"$group": bson.M{
			"_id": bson.M{
				"_id": "$_id",
				"tab": "$tabs._id",
			},
			"data":    bson.M{"$first": "$data"},
			"tabs":    bson.M{"$first": "$tabs"},
			"widgets": bson.M{"$push": "$widgets"},
		}},
		{"$addFields": bson.M{
			"tabs.widgets": bson.M{"$filter": bson.M{
				"input": "$widgets",
				"cond":  "$$this._id",
			}},
		}},
		{"$sort": bson.M{"tabs.position": 1}},
		{"$group": bson.M{
			"_id":  "$_id._id",
			"data": bson.M{"$first": "$data"},
			"tabs": bson.M{"$push": "$tabs"},
		}},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$data",
			bson.M{"tabs": bson.M{"$filter": bson.M{
				"input": "$tabs",
				"cond":  "$$this._id",
			}}},
		}}}},
		{"$lookup": bson.M{
			"from":         mongo.ViewGroupMongoCollection,
			"localField":   "group_id",
			"foreignField": "_id",
			"as":           "group",
		}},
		{"$unwind": bson.M{"path": "$group", "preserveNullAndEmptyArrays": true}},
	}
}

func computePositions(
	r EditPositionRequest,
	oldGroupPositions []string,
	oldViewPositionsByGroup map[string][]string,
) ([]string, []string, map[string]string) {
	newGroupPositions := make([]string, 0, len(oldGroupPositions))
	newViewPositionsByGroup := make(map[string][]string, len(oldGroupPositions))
	minUpdatedGroupOldPos := -1
	minUpdatedViewOldPosByGroup := make(map[string]int, len(oldGroupPositions))
	notUpdatedGroups := inverseStrSlice(oldGroupPositions)
	notUpdatedViews := make(map[string]map[string]int, len(oldGroupPositions))
	oldViewGroup := make(map[string]string)
	for group, positions := range oldViewPositionsByGroup {
		notUpdatedViews[group] = inverseStrSlice(positions)
		for _, viewID := range positions {
			oldViewGroup[viewID] = group
		}
	}
	changedViewGroup := make(map[string]string)

	for _, item := range r.Items {
		if oldPos, ok := notUpdatedGroups[item.ID]; ok {
			delete(notUpdatedGroups, item.ID)

			if minUpdatedGroupOldPos < 0 || minUpdatedGroupOldPos > oldPos {
				minUpdatedGroupOldPos = oldPos
			}

			newGroupPositions = append(newGroupPositions, item.ID)
		} else {
			return nil, nil, nil
		}

		newViewPositionsByGroup[item.ID] = make([]string, 0, len(item.Views))
		minUpdatedViewOldPosByGroup[item.ID] = -1

		for _, viewID := range item.Views {
			if oldGroup, ok := oldViewGroup[viewID]; ok {
				oldViewPos := notUpdatedViews[oldGroup][viewID]
				delete(notUpdatedViews[oldGroup], viewID)

				if oldGroup == item.ID {
					if minUpdatedViewOldPosByGroup[item.ID] < 0 || minUpdatedViewOldPosByGroup[item.ID] > oldViewPos {
						minUpdatedViewOldPosByGroup[item.ID] = oldViewPos
					}
				} else {
					changedViewGroup[viewID] = item.ID
				}

				newViewPositionsByGroup[item.ID] = append(newViewPositionsByGroup[item.ID], viewID)
			} else {
				return nil, nil, nil
			}
		}
	}

	newGroupPositions = mergePositions(
		newGroupPositions,
		oldGroupPositions,
		notUpdatedGroups,
		minUpdatedGroupOldPos,
	)

	for group := range oldViewPositionsByGroup {
		index, ok := minUpdatedViewOldPosByGroup[group]
		if !ok {
			index = -1
		}
		newViewPositionsByGroup[group] = mergePositions(
			newViewPositionsByGroup[group],
			oldViewPositionsByGroup[group],
			notUpdatedViews[group],
			index,
		)
	}

	newViewPositions := make([]string, 0)
	for _, group := range newGroupPositions {
		newViewPositions = append(newViewPositions, newViewPositionsByGroup[group]...)
	}

	return newGroupPositions, newViewPositions, changedViewGroup
}

func prependStr(slice []string, elem string) []string {
	slice = append(slice, "")
	copy(slice[1:], slice)
	slice[0] = elem
	return slice
}

func mergePositions(
	newPositions []string,
	oldPositions []string,
	notUpdated map[string]int,
	minUpdatedOldPosition int,
) []string {
	// Add not updated items to the begin if they were before first updated item.
	for oldPosition := len(oldPositions) - 1; oldPosition >= 0; oldPosition-- {
		id := oldPositions[oldPosition]
		if _, ok := notUpdated[id]; !ok {
			continue
		}

		if oldPosition < minUpdatedOldPosition {
			delete(notUpdated, id)
			newPositions = prependStr(newPositions, id)
		}
	}

	// Add not updated items to the end if they were after first updated item.
	for _, id := range oldPositions {
		if _, ok := notUpdated[id]; !ok {
			continue
		}

		delete(notUpdated, id)
		newPositions = append(newPositions, id)
	}

	return newPositions
}

func inverseStrSlice(slice []string) map[string]int {
	m := make(map[string]int, len(slice))
	for k, v := range slice {
		m[v] = k
	}

	return m
}
