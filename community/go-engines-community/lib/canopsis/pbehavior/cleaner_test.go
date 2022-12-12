package pbehavior_test

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCleaner_Clean_GivenPbehaviorsWithoutRrule_ShouldDeleteThem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	before := types.CpsTime{Time: time.Now().Add(-time.Hour * 24 * 7)}
	var limit int64 = 1000
	var expectedDeleted int64 = 1
	pbhId := "test-pbh"

	mockClient := mock_mongo.NewMockDbClient(ctrl)
	mockCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockCursor := mock_mongo.NewMockCursor(ctrl)
	mockCursor.EXPECT().Next(gomock.Any()).Return(true)
	mockCursor.EXPECT().Next(gomock.Any()).Return(false)
	mockCursor.EXPECT().Close(gomock.Any()).Return(nil)
	mockCursor.EXPECT().Decode(gomock.Any()).Do(func(v *pbehavior.PBehavior) {
		v.ID = pbhId
	})
	emptyMockCursor := mock_mongo.NewMockCursor(ctrl)
	emptyMockCursor.EXPECT().Next(gomock.Any()).Return(false)
	emptyMockCursor.EXPECT().Close(gomock.Any()).Return(nil)
	mockClient.EXPECT().Collection(mongo.PbehaviorMongoCollection).Return(mockCollection)
	mockCollection.EXPECT().
		DeleteMany(gomock.Any(), gomock.Eq(bson.M{"_id": bson.M{"$in": []string{pbhId}}})).
		Return(expectedDeleted, nil)
	mockCollection.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCursor, nil)
	mockCollection.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(emptyMockCursor, nil)
	cleaner := pbehavior.NewCleaner(mockClient, 1000, zerolog.Nop())

	deleted, err := cleaner.Clean(ctx, before, limit)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if deleted != expectedDeleted {
		t.Errorf("expected deleted %v but got %v", expectedDeleted, deleted)
	}
}

func TestCleaner_Clean_GivenPbehaviorsWithRrule_ShouldDeleteThem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	now := time.Now()
	before := types.CpsTime{Time: now.Add(-time.Hour * 24 * 7)}
	var limit int64 = 1000
	start := types.CpsTime{Time: now.AddDate(0, 0, -100)}
	var expectedDeleted int64 = 2
	pbehaviors := []pbehavior.PBehavior{
		{
			ID:    "test-pbehavior-1",
			RRule: "FREQ=DAILY;INTERVAL=7",
			Start: &start,
		},
		{
			ID:    "test-pbehavior-2",
			RRule: "FREQ=DAILY;INTERVAL=7;COUNT=5",
			Start: &start,
		},
		{
			ID:    "test-pbehavior-3",
			RRule: "FREQ=DAILY;INTERVAL=7;UNTIL=" + now.AddDate(0, 0, -8).Format("20060102T150405"),
			Start: &start,
		},
		{
			ID:    "test-pbehavior-4",
			RRule: "FREQ=DAILY;INTERVAL=7;COUNT=500",
			Start: &start,
		},
		{
			ID:    "test-pbehavior-5",
			RRule: "FREQ=DAILY;INTERVAL=7;UNTIL=" + now.AddDate(0, 0, 100).Format("20060102T150405"),
			Start: &start,
		},
	}

	mockClient := mock_mongo.NewMockDbClient(ctrl)
	mockCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockCursor := mock_mongo.NewMockCursor(ctrl)
	nextIndex := -1
	mockCursor.EXPECT().Next(gomock.Any()).DoAndReturn(func(_ context.Context) bool {
		nextIndex++
		return nextIndex < len(pbehaviors)
	}).Times(len(pbehaviors) + 1)
	decodeIndedx := -1
	mockCursor.EXPECT().Decode(gomock.Any()).DoAndReturn(func(v *pbehavior.PBehavior) error {
		decodeIndedx++
		*v = pbehaviors[decodeIndedx]
		return nil
	}).Times(len(pbehaviors))
	mockCursor.EXPECT().Close(gomock.Any()).Return(nil)
	emptyMockCursor := mock_mongo.NewMockCursor(ctrl)
	emptyMockCursor.EXPECT().Next(gomock.Any()).Return(false)
	emptyMockCursor.EXPECT().Close(gomock.Any()).Return(nil)
	mockClient.EXPECT().Collection(mongo.PbehaviorMongoCollection).Return(mockCollection)
	mockCollection.EXPECT().
		DeleteMany(gomock.Any(), gomock.Eq(bson.M{"_id": bson.M{"$in": []string{
			pbehaviors[1].ID,
			pbehaviors[2].ID,
		}}})).
		Return(expectedDeleted, nil)
	mockCollection.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(emptyMockCursor, nil)
	mockCollection.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCursor, nil)

	cleaner := pbehavior.NewCleaner(mockClient, 1000, zerolog.Nop())
	deleted, err := cleaner.Clean(ctx, before, limit)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	if deleted != expectedDeleted {
		t.Errorf("expected deleted %v but got %v", expectedDeleted, deleted)
	}
}
