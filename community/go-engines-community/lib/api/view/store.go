package view

import (
	"context"
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
	"time"
)

const permissionPrefix = "Rights on view :"

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
	GetOneBy(ctx context.Context, r string) (*viewgroup.View, error)
	Insert(ctx context.Context, userID string, r EditRequest) (*viewgroup.View, error)
	Update(ctx context.Context, r EditRequest) (*viewgroup.View, error)
	// UpdatePositions receives some groups and views with updated positions and updates
	// positions for all groups and views in db and moves views to another groups if necessary.
	UpdatePositions(ctx context.Context, r EditPositionRequest) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		collection:            dbClient.Collection(mongo.ViewMongoCollection),
		tabCollection:         dbClient.Collection(mongo.ViewTabMongoCollection),
		widgetCollection:      dbClient.Collection(mongo.WidgetMongoCollection),
		groupCollection:       dbClient.Collection(mongo.ViewGroupMongoCollection),
		aclCollection:         dbClient.Collection(mongo.RightsMongoCollection),
		userPrefCollection:    dbClient.Collection(mongo.UserPreferencesMongoCollection),
		defaultSearchByFields: []string{"_id", "title", "description", "author"},
		defaultSortBy:         "position",
	}
}

type store struct {
	collection            mongo.DbCollection
	tabCollection         mongo.DbCollection
	widgetCollection      mongo.DbCollection
	groupCollection       mongo.DbCollection
	aclCollection         mongo.DbCollection
	userPrefCollection    mongo.DbCollection
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

func (s *store) GetOneBy(ctx context.Context, id string) (*viewgroup.View, error) {
	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	pipeline = append(pipeline, []bson.M{
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
		{"$sort": bson.M{"widgets.position": 1}},
		{"$group": bson.M{
			"_id": bson.M{
				"_id": "$_id",
				"tab": "$tabs._id",
			},
			"data":    bson.M{"$first": "$$ROOT"},
			"tabs":    bson.M{"$first": "$tabs"},
			"widgets": bson.M{"$push": "$widgets"},
		}},
		{"$addFields": bson.M{
			"tabs.widgets": "$widgets",
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
	}...)
	pipeline = append(pipeline, getNestedObjectsPipeline()...)
	cursor, err := s.collection.Aggregate(ctx, pipeline)
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

func (s *store) Insert(ctx context.Context, userID string, r EditRequest) (*viewgroup.View, error) {
	count, err := s.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
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
		return nil, err
	}

	newView, err := s.GetOneBy(ctx, id)
	if err != nil {
		return nil, err
	}

	err = s.createPermissions(ctx, userID, *newView)
	if err != nil {
		return nil, err
	}

	return newView, nil
}

func (s *store) Update(ctx context.Context, r EditRequest) (*viewgroup.View, error) {
	oldView := view.View{}
	err := s.collection.FindOne(ctx, bson.M{"_id": r.ID}).Decode(&oldView)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	if r.Group != oldView.Group {
		err := s.normalizePositionsOnViewMove(ctx, r.ID, r.Group)
		if err != nil {
			return nil, err
		}
	}

	newView, err := s.GetOneBy(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	err = s.updatePermissions(ctx, *newView)
	if err != nil {
		return nil, err
	}

	return newView, nil
}

func (s *store) Delete(ctx context.Context, id string) (bool, error) {
	delCount, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	if delCount == 0 {
		return false, nil
	}

	err = s.deletePermissions(ctx, id)
	if err != nil {
		return false, err
	}

	err = s.deleteTabs(ctx, id)
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

func (s *store) createPermissions(ctx context.Context, userID string, view viewgroup.View) error {
	_, err := s.aclCollection.InsertOne(ctx, bson.M{
		"_id":          view.ID,
		"crecord_name": view.ID,
		"crecord_type": securitymodel.LineTypeObject,
		"desc":         fmt.Sprintf("%s %s", permissionPrefix, view.Title),
		"type":         securitymodel.LineObjectTypeRW,
	})
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
		bson.M{"$set": bson.M{
			"rights." + view.ID: bson.M{
				"checksum": securitymodel.PermissionBitmaskRead |
					securitymodel.PermissionBitmaskUpdate |
					securitymodel.PermissionBitmaskDelete,
			},
		}},
	)
	if err != nil {
		return err
	}

	return err
}

func (s *store) updatePermissions(ctx context.Context, view viewgroup.View) error {
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

	return s.updatePositions(ctx, groupPositions, viewPositions, nil)
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
