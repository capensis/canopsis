package statistics_test

import (
	"context"
	"reflect"
	"strconv"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statistics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_v8 "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/github.com/go-redis/redis/v8"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestStatsListener_Listen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	flushInterval := time.Millisecond * 100
	ctx, cancel := context.WithTimeout(context.Background(), flushInterval+time.Millisecond*10)
	defer cancel()

	now := time.Now()
	y, m, d := now.Date()
	nowMinute := time.Date(y, m, d, now.Hour(), now.Minute(), 0, 0, now.Location())
	key1 := nowMinute.Unix()
	key2 := nowMinute.Add(3 * time.Minute).Unix()
	keys := []string{strconv.FormatInt(key1, 10), strconv.FormatInt(key2, 10)}
	expectedModelsCount := 4
	expectedIds := []int64{
		nowMinute.Unix(),
		nowMinute.Add(1 * time.Minute).Unix(),
		nowMinute.Add(2 * time.Minute).Unix(),
		nowMinute.Add(3 * time.Minute).Unix(),
	}

	storeIntervals := map[string]int64{
		mongo.MessageRateStatsMinuteCollectionName: 1,
	}
	mockMongoClient := mock_mongo.NewMockDbClient(ctrl)
	mockCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockMongoClient.EXPECT().Collection(gomock.Eq(mongo.MessageRateStatsMinuteCollectionName)).
		Return(mockCollection)
	mockCursor := mock_mongo.NewMockCursor(ctrl)
	mockCursor.EXPECT().Next(gomock.Any()).Return(true)
	mockCursor.EXPECT().Decode(gomock.Any()).Return(nil)
	mockCursor.EXPECT().Close(gomock.Any()).Return(nil)
	mockCollection.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockCursor, nil)
	mockCollection.EXPECT().BulkWrite(gomock.Any(), gomock.Any()).Do(func(_ context.Context, writeModels []mongodriver.WriteModel, _ ...options.BulkWriteOptions) {
		if len(writeModels) != expectedModelsCount {
			t.Errorf("expected %v write models but got %v", expectedModelsCount, writeModels)
			return
		}

		ids := make([]int64, expectedModelsCount)
		for i, model := range writeModels {
			updateModel := model.(*mongodriver.UpdateOneModel)
			filter := updateModel.Filter.(bson.M)
			ids[i] = filter["_id"].(int64)
		}

		if !reflect.DeepEqual(ids, expectedIds) {
			t.Errorf("expected %v ids but got %v", expectedIds, ids)
			return
		}
	}).Return(&mongodriver.BulkWriteResult{}, nil)
	mockRedisClient := mock_v8.NewMockCmdable(ctrl)
	mockRedisClient.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewScanCmdResult(keys, 0, nil))
	mockRedisClient.EXPECT().HGet(gomock.Any(), gomock.Eq(keys[0]), gomock.Any()).
		Return(redis.NewStringResult("2", nil))
	mockRedisClient.EXPECT().HGet(gomock.Any(), gomock.Eq(keys[0]), gomock.Any()).
		Return(redis.NewStringResult("0", nil))
	mockRedisClient.EXPECT().HGet(gomock.Any(), gomock.Eq(keys[1]), gomock.Any()).
		Return(redis.NewStringResult("3", nil))
	mockRedisClient.EXPECT().HGet(gomock.Any(), gomock.Eq(keys[1]), gomock.Any()).
		Return(redis.NewStringResult("0", nil))
	mockRedisClient.EXPECT().FlushDB(gomock.Any()).Return(redis.NewStatusResult("", nil))

	listener := statistics.NewStatsListener(mockMongoClient, mockRedisClient, flushInterval,
		storeIntervals, zerolog.Nop())
	ch := make(chan statistics.Message)
	defer close(ch)
	listener.Listen(ctx, ch)
}
