package che

import (
	"context"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	libcontext "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	mock_context "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/context"
	mock_encoding "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/encoding"
	mock_eventfilter "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/eventfilter"
	"github.com/golang/mock/gomock"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
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
	mockMetricsConfigProvider := mock_config.NewMockMetricsConfigProvider(ctrl)
	mockMetricsConfigProvider.EXPECT().Get().Return(config.MetricsConfig{EnableTechMetrics: false})
	mockEventFilterService := mock_eventfilter.NewMockService(ctrl)
	mockEventFilterService.EXPECT().ProcessEvent(gomock.Any(), gomock.Any()).Return(event, nil)
	mockEnrichmentCenter := mock_context.NewMockEnrichmentCenter(ctrl)
	mockEnrichmentCenter.EXPECT().HandleEntityServiceUpdate(gomock.Any(), gomock.Eq("test-component")).
		Return(&libcontext.UpdatedEntityServices{}, nil)
	mockEnrichmentCenter.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)
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
		MetricsConfigProvider:  mockMetricsConfigProvider,
		FeatureEventProcessing: true,
		FeatureContextCreation: true,
		AlarmConfigProvider:    mockAlarmConfigProvider,
		EventFilterService:     mockEventFilterService,
		EnrichmentCenter:       mockEnrichmentCenter,
		Encoder:                mockEncoder,
		Decoder:                mockDecoder,
		Logger:                 zerolog.Logger{},
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
