package entity_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	libapi "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

var trans *ut.UniversalTranslator

func init() {
	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}, zerolog.Nop()), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))
	var err error
	trans, err = libapi.RegisterValidators(security.Config{}, tplExecutor)
	if err != nil {
		panic(err)
	}
}

func BenchmarkBulkPause_Given100CreateItems(b *testing.B) {
	benchmarkBulkPause_givenNCreateItems(b, 100)
}

func BenchmarkBulkPause_Given500CreateItems(b *testing.B) {
	benchmarkBulkPause_givenNCreateItems(b, 500)
}

func BenchmarkBulkPause_Given1000CreateItems(b *testing.B) {
	benchmarkBulkPause_givenNCreateItems(b, 1000)
}

func benchmarkBulkPause_givenNCreateItems(b *testing.B, itemCount int) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()
	dbClient, err := mongo.NewClient(b.Context())
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}

	defer func() {
		err := dbClient.Disconnect(context.WithoutCancel(b.Context()))
		if err != nil {
			b.Errorf("unexpected error %v", err)
		}
	}()

	ctx := b.Context()
	loader := fixtures.NewLoader(dbClient, []string{"./testdata/fixtures/bulk_pause.yml"},
		fixtures.NewParser(fixtures.NewFaker(password.NewBcryptEncoder())), zerolog.Nop())
	err = loader.Load(ctx)
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}

	defer func() {
		err := loader.Clean(context.WithoutCancel(ctx))
		if err != nil {
			b.Errorf("unexpected error %v", err)
		}
	}()

	ch := make(chan rpc.PbehaviorRecomputeEvent, b.N)
	defer close(ch)
	go func() {
		for {
			select {
			case <-ctx.Done():
			case _, ok := <-ch:
				if !ok {
					return
				}
			}
		}
	}()
	authorProvider := author.NewProvider(&config.BaseApiConfigProvider{})
	amqpPublisher := mock_amqp.NewMockPublisher(ctrl)
	amqpPublisher.EXPECT().PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	store := entity.NewStore(dbClient, nil, nil, nil, authorProvider, nil, json.NewEncoder(), json.NewDecoder(), amqpPublisher, event.NewGenerator("", ""))
	errorResponder := response.NewErrorResponder(validation.NewErrorTranslator(trans, zerolog.Nop()), zerolog.Nop())
	api := entity.NewApi(store, nil, nil, nil, nil, nil, errorResponder, zerolog.Nop())
	encoder := json.NewEncoder()
	reqItems := make([]map[string]any, itemCount)
	for i := 0; i < itemCount; i++ {
		reqItems[i] = map[string]any{
			"_id":     "test-resource-" + strconv.Itoa(i) + "/test-component",
			"comment": "pause comment",
			"type":    "test-pause-type",
		}
	}

	buf, err := encoder.Encode(reqItems)
	if err != nil {
		b.Errorf("unexpected error %v", err)
	}

	reqBody := io.NopCloser(bytes.NewBuffer(buf))
	r, err := http.NewRequestWithContext(ctx, http.MethodPut, "/api/v4/bulk/entities/pause", nil)
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = r
	authctx.SetUserKey(ginCtx, "root")
	authctx.SetUsername(ginCtx, "root")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Body = reqBody
		api.BulkPause(ginCtx)
		if w.Code != http.StatusMultiStatus {
			b.Fatalf("unexpected status code %v", w.Code)
		}
	}
}
