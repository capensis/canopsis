package che

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	mock_contextgraph "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/contextgraph"
	mock_encoding "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/encoding"
	mock_eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"reflect"
	"testing"
)

func TestMessageProcessor_Process_GivenRecomputeEntityServiceEvent_ShouldPassItToNextQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	body := []byte("{\"event_type\":\"recomputeentityservice\"}")
	event := types.Event{
		EventType:     types.EventTypeRecomputeEntityService,
		SourceType:    types.SourceTypeComponent,
		Connector:     "test-connector",
		ConnectorName: "test-connector-name",
		Component:     "test-component",
	}
	expectedBody := []byte("test-next-body")
	mockAlarmConfigProvider := mock_config.NewMockAlarmConfigProvider(ctrl)
	mockAlarmConfigProvider.EXPECT().Get().Return(config.AlarmConfig{})
	mockEventFilterService := mock_eventfilter.NewMockService(ctrl)

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().WithTransaction(gomock.Any(), gomock.Any()).Do(func(ctx context.Context, f func(context.Context) error) error {
		return f(ctx)
	})

	manager := mock_contextgraph.NewMockManager(ctrl)
	manager.EXPECT().RecomputeService(gomock.Any(), gomock.Eq("test-component")).Return([]types.Entity{}, nil)
	manager.EXPECT().UpdateEntities(gomock.Any(), gomock.Eq([]types.Entity{})).Return(nil)
	mockDecoder := mock_encoding.NewMockDecoder(ctrl)
	mockDecoder.EXPECT().Decode(gomock.Eq(body), gomock.Any()).Do(func(_ []byte, e *types.Event) {
		*e = event
	}).Return(nil)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockEncoder.EXPECT().Encode(gomock.Any()).Do(func(event types.Event) {
		if event.EventType != types.EventTypeRecomputeEntityService {
			t.Errorf("expected event %s but got %s", types.EventTypeRecomputeEntityService, event.EventType)
		}
	}).Return(expectedBody, nil)
	processor := &messageProcessor{
		DbClient:            dbClient,
		AlarmConfigProvider: mockAlarmConfigProvider,
		EventFilterService:  mockEventFilterService,
		ContextGraphManager: manager,
		Encoder:             mockEncoder,
		Decoder:             mockDecoder,
		Logger:              zerolog.Logger{},
	}

	resBody, err := processor.Process(context.Background(), amqp.Delivery{
		Body: body,
	})
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if !reflect.DeepEqual(expectedBody, resBody) {
		t.Errorf("expected result %s but got %s", expectedBody, resBody)
	}
}
