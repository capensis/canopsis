package view

import (
	"context"
	"encoding/json"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewgroup"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const permissionPrefix = "Rights on view :"

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetOneBy(ctx context.Context, r string) (*viewgroup.View, error)
	Insert(ctx context.Context, userID string, r []EditRequest) ([]viewgroup.View, error)
	Update(ctx context.Context, r []BulkUpdateRequestItem) ([]viewgroup.View, error)
	// UpdatePositions receives some groups and views with updated positions and updates
	// positions for all groups and views in db and moves views to another groups if necessary.
	UpdatePositions(ctx context.Context, r EditPositionRequest) (bool, error)
	Delete(ctx context.Context, id []string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbCollection:          dbClient.Collection(mongo.ViewMongoCollection),
		dbGroupCollection:     dbClient.Collection(mongo.ViewGroupMongoCollection),
		aclDbCollection:       dbClient.Collection(mongo.RightsMongoCollection),
		userPrefDbCollection:  dbClient.Collection(mongo.UserPreferencesMongoCollection),
		defaultSearchByFields: []string{"_id", "title", "description", "author"},
		defaultSortBy:         "position",
	}
}

type store struct {
	dbCollection          mongo.DbCollection
	dbGroupCollection     mongo.DbCollection
	aclDbCollection       mongo.DbCollection
	userPrefDbCollection  mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
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

	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
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

func (s *store) GetOneBy(ctx context.Context, id string) (*viewgroup.View, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		model := &viewgroup.View{}
		err := cursor.Decode(model)
		if err != nil {
			return nil, err
		}

		return model, nil
	}

	return nil, nil
}

func (s *store) Insert(ctx context.Context, userID string, r []EditRequest) ([]viewgroup.View, error) {
	count, err := s.dbCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	now := types.CpsTime{Time: time.Now()}
	ids := make([]string, len(r))
	docs := make([]interface{}, len(r))
	for i, item := range r {
		tabs, err := transformTabRequestToModel(item.Tabs, nil)
		if err != nil {
			return nil, err
		}

		ids[i] = utils.NewID()
		docs[i] = view.View{
			ID:              ids[i],
			Enabled:         *item.Enabled,
			Title:           item.Title,
			Description:     item.Description,
			Group:           item.Group,
			Tabs:            tabs,
			Position:        count + int64(i),
			Tags:            item.Tags,
			PeriodicRefresh: item.PeriodicRefresh,
			Author:          item.Author,
			Created:         now,
			Updated:         now,
		}
	}

	_, err = s.dbCollection.InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}

	views, err := s.findByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	err = s.createPermissions(ctx, userID, views)
	if err != nil {
		return nil, err
	}

	return views, nil
}

func (s *store) Update(ctx context.Context, r []BulkUpdateRequestItem) ([]viewgroup.View, error) {
	ids := make([]string, len(r))
	rByID := make(map[string]BulkUpdateRequestItem, len(r))
	for i, item := range r {
		ids[i] = item.ID
		rByID[item.ID] = item
	}

	cursor, err := s.dbCollection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	oldViews := make([]view.View, 0)
	err = cursor.All(ctx, &oldViews)
	if err != nil || len(oldViews) < len(ids) {
		return nil, err
	}

	now := types.CpsTime{Time: time.Now()}
	models := make([]mongodriver.WriteModel, len(oldViews))
	changedGroup := make(map[string]string)
	oldTabs := make([]view.Tab, 0)
	newTabs := make([]view.Tab, 0)

	for i, oldModel := range oldViews {
		item := rByID[oldModel.ID]
		tabs, err := transformTabRequestToModel(item.Tabs, oldModel.Tabs)
		if err != nil {
			return nil, err
		}

		oldTabs = append(oldTabs, oldModel.Tabs...)
		newTabs = append(newTabs, tabs...)

		if item.Group != oldModel.Group {
			changedGroup[item.Group] = item.ID
		}

		models[i] = mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": oldModel.ID}).
			SetUpdate(bson.M{"$set": bson.M{
				"enabled":          *item.Enabled,
				"title":            item.Title,
				"description":      item.Description,
				"group_id":         item.Group,
				"tabs":             tabs,
				"tags":             item.Tags,
				"periodic_refresh": item.PeriodicRefresh,
				"author":           item.Author,
				"updated":          now,
			}})
	}

	_, err = s.dbCollection.BulkWrite(ctx, models)
	if err != nil {
		return nil, err
	}

	err = s.normalizePositionsOnViewMove(ctx, changedGroup)
	if err != nil {
		return nil, err
	}

	newViews, err := s.findByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	err = s.updatePermissions(ctx, newViews)
	if err != nil {
		return nil, err
	}

	err = s.handleWidgets(ctx, oldTabs, newTabs)
	if err != nil {
		return nil, err
	}

	return newViews, nil
}

func (s *store) Delete(ctx context.Context, ids []string) (bool, error) {
	if len(ids) == 0 {
		return false, nil
	}

	cursor, err := s.dbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"_id": bson.M{"$in": ids}}},
		{"$unwind": bson.M{"path": "$tabs", "preserveNullAndEmptyArrays": true}},
		{"$unwind": bson.M{"path": "$tabs.widgets", "preserveNullAndEmptyArrays": true}},
		{"$project": bson.M{
			"widget": "$tabs.widgets._id",
		}},
	})
	if err != nil {
		return false, err
	}
	defer cursor.Close(ctx)

	widgetIds := make([]string, 0)
	foundViews := make(map[string]bool)
	for cursor.Next(ctx) {
		data := make(map[string]string)
		err := cursor.Decode(&data)
		if err != nil {
			return false, err
		}
		foundViews[data["_id"]] = true
		widgetIds = append(widgetIds, data["widget"])
	}

	if len(foundViews) < len(ids) {
		return false, nil
	}

	delCount, err := s.dbCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return false, err
	}

	if delCount == 0 {
		return false, nil
	}

	err = s.deletePermissions(ctx, ids)
	if err != nil {
		return false, err
	}

	err = s.deleteUserPreferences(ctx, widgetIds)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) UpdatePositions(ctx context.Context, r EditPositionRequest) (bool, error) {
	groupPositions, viewPositionsByGroup, err := s.findViewPositions(ctx)
	if err != nil || len(groupPositions) == 0 {
		return false, err
	}

	newGroupPositions, newViewPositions, changedViewGroup := computePositions(r,
		groupPositions, viewPositionsByGroup)
	if len(newGroupPositions) == 0 {
		return false, nil
	}

	err = s.updatePositions(ctx, newGroupPositions, newViewPositions, changedViewGroup)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *store) findByIDs(ctx context.Context, ids []string) ([]viewgroup.View, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": bson.M{"$in": ids}}}}
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	pipeline = append(pipeline, common.GetSortQuery("position", common.SortAsc))
	cursor, err := s.dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	views := make([]viewgroup.View, 0)
	err = cursor.All(ctx, &views)
	if err != nil {
		return nil, err
	}

	return views, nil
}

func (s *store) createPermissions(ctx context.Context, userID string, views []viewgroup.View) error {

	docs := make([]interface{}, len(views))
	set := bson.M{}

	for i, v := range views {
		docs[i] = bson.M{
			"_id":          v.ID,
			"crecord_name": v.ID,
			"crecord_type": securitymodel.LineTypeObject,
			"desc":         fmt.Sprintf("%s %s", permissionPrefix, v.Title),
			"type":         securitymodel.LineObjectTypeRW,
		}
		set["rights."+v.ID] = bson.M{
			"checksum": securitymodel.PermissionBitmaskRead |
				securitymodel.PermissionBitmaskUpdate |
				securitymodel.PermissionBitmaskDelete,
		}
	}

	_, err := s.aclDbCollection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	res := s.aclDbCollection.FindOne(ctx, bson.M{
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

	_, err = s.aclDbCollection.UpdateOne(ctx,
		bson.M{
			"_id":          user.Role,
			"crecord_type": securitymodel.LineTypeRole,
		},
		bson.M{"$set": set},
	)
	if err != nil {
		return err
	}

	return err
}

func (s *store) updatePermissions(ctx context.Context, views []viewgroup.View) error {

	models := make([]mongodriver.WriteModel, len(views))
	for i, v := range views {
		models[i] = mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{
				"_id":          v.ID,
				"crecord_type": securitymodel.LineTypeObject,
			}).
			SetUpdate(bson.M{"$set": bson.M{
				"desc": fmt.Sprintf("%s %s", permissionPrefix, v.Title),
			}})
	}

	_, err := s.aclDbCollection.BulkWrite(ctx, models)

	return err
}

func (s *store) deletePermissions(ctx context.Context, viewIDs []string) error {
	unset := bson.M{}
	for _, viewID := range viewIDs {
		unset["rights."+viewID] = ""
	}

	_, err := s.aclDbCollection.UpdateMany(ctx,
		bson.M{"crecord_type": securitymodel.LineTypeRole},
		bson.M{"$unset": unset},
	)
	if err != nil {
		return err
	}

	_, err = s.aclDbCollection.DeleteMany(ctx, bson.M{
		"_id":          bson.M{"$in": viewIDs},
		"crecord_type": securitymodel.LineTypeObject,
	})

	return err
}

func (s *store) handleWidgets(ctx context.Context, oldTabs, newTabs []view.Tab) error {
	widgetsByID := make(map[string]map[string]view.Widget)
	for _, tab := range newTabs {
		widgetsByID[tab.ID] = make(map[string]view.Widget, len(tab.Widgets))
		for _, widget := range tab.Widgets {
			widgetsByID[tab.ID][widget.ID] = widget
		}
	}

	removedWidgets := make([]string, 0)

	for _, tab := range oldTabs {
		for _, oldWidget := range tab.Widgets {
			if newWidget, ok := widgetsByID[tab.ID][oldWidget.ID]; ok {
				switch oldWidget.Type {
				case view.WidgetTypeJunit:
					err := s.handleJunitWidget(ctx, oldWidget, newWidget)
					if err != nil {
						return err
					}
				}
			} else {
				removedWidgets = append(removedWidgets, oldWidget.ID)
			}
		}
	}

	err := s.deleteUserPreferences(ctx, removedWidgets)
	if err != nil {
		return err
	}

	return nil
}

func (s *store) handleJunitWidget(ctx context.Context, oldModel, newModel view.Widget) error {
	oldDir := oldModel.GetStringParameter(view.WidgetParamJunitDir, "")
	newDir := newModel.GetStringParameter(view.WidgetParamJunitDir, "")
	if oldDir == newDir {
		return nil
	}

	_, err := s.dbCollection.UpdateMany(ctx,
		bson.M{
			"tabs.widgets._id": oldModel.ID,
		},
		bson.M{"$unset": bson.M{
			"tabs.$[tab].widgets.$[widget].internal_parameters." + view.WidgetInternalParamJunitTestSuites: "",
		}},
		options.Update().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"tab.widgets._id": oldModel.ID},
				bson.M{"widget._id": oldModel.ID},
			},
		}),
	)

	return err
}

func (s *store) normalizePositionsOnViewMove(ctx context.Context, changedGroup map[string]string) error {
	if len(changedGroup) == 0 {
		return nil
	}

	groupPositions, viewPositionsByGroup, err := s.findViewPositions(ctx)
	if err != nil {
		return err
	}

	for groupID, viewID := range changedGroup {
		index := -1
		for i, v := range viewPositionsByGroup[groupID] {
			if v == viewID {
				index = i
			}
		}

		viewPositionsByGroup[groupID] = append(viewPositionsByGroup[groupID][:index], viewPositionsByGroup[groupID][index+1:]...)
		viewPositionsByGroup[groupID] = append(viewPositionsByGroup[groupID], viewID)
	}

	viewPositions := make([]string, 0)
	for id := range viewPositionsByGroup {
		viewPositions = append(viewPositions, viewPositionsByGroup[id]...)
	}

	return s.updatePositions(ctx, groupPositions, viewPositions, nil)
}

func (s *store) findViewPositions(ctx context.Context) ([]string, map[string][]string, error) {
	cursor, err := s.dbGroupCollection.Aggregate(ctx, []bson.M{
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

		for j, viewID := range group.Views {
			viewPositionsByGroup[group.ID][j] = viewID
		}
	}

	return groupPositions, viewPositionsByGroup, nil
}

func (s *store) updatePositions(
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

		_, err := s.dbGroupCollection.BulkWrite(ctx, groupModels)
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

		_, err := s.dbCollection.BulkWrite(ctx, models)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *store) deleteUserPreferences(ctx context.Context, widgetIDs []string) error {
	if len(widgetIDs) == 0 {
		return nil
	}

	_, err := s.userPrefDbCollection.DeleteMany(ctx, bson.M{
		"widget": bson.M{"$in": widgetIDs},
	})

	return err
}

func getNestedObjectsPipeline() []bson.M {
	return []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.ViewGroupMongoCollection,
			"localField":   "group_id",
			"foreignField": "_id",
			"as":           "group",
		}},
		{"$unwind": bson.M{"path": "$group", "preserveNullAndEmptyArrays": true}},
	}
}

func transformTabRequestToModel(r []TabRequest, old []view.Tab) ([]view.Tab, error) {
	oldWidgetByID := make(map[string]view.Widget, len(old))
	for _, tab := range old {
		for _, widget := range tab.Widgets {
			oldWidgetByID[widget.ID] = widget
		}
	}

	tabs := make([]view.Tab, len(r))
	for i, tab := range r {
		widgets := make([]view.Widget, len(tab.Widgets))
		for j, widget := range tab.Widgets {
			var params map[string]interface{}
			b, err := json.Marshal(widget.Parameters)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(b, &params)
			if err != nil {
				return nil, err
			}

			var internalParameters map[string]interface{}
			if oldWidget, ok := oldWidgetByID[widget.ID]; ok {
				internalParameters = oldWidget.InternalParameters
			}

			widgets[j] = view.Widget{
				ID:                 widget.ID,
				Title:              widget.Title,
				Type:               widget.Type,
				GridParameters:     widget.GridParameters,
				Parameters:         params,
				InternalParameters: internalParameters,
			}
		}
		tabs[i] = view.Tab{
			ID:      tab.ID,
			Title:   tab.Title,
			Widgets: widgets,
		}
	}

	return tabs, nil
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
