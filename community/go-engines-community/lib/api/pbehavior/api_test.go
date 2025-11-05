package pbehavior_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	libapi "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/fixtures"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func init() {
	dbClient, err := mongo.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}, zerolog.Nop()), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))

	err = libapi.RegisterValidators(dbClient, security.Config{}, nil, tplExecutor)
	if err != nil {
		panic(err)
	}
}

func BenchmarkBulkConnectorEdit_Given100CreateItems(b *testing.B) {
	benchmarkBulkConnectorEdit_givenNCreateItems(b, 100)
}

func BenchmarkBulkConnectorEdit_Given500CreateItems(b *testing.B) {
	benchmarkBulkConnectorEdit_givenNCreateItems(b, 500)
}

func BenchmarkBulkConnectorEdit_Given1000CreateItems(b *testing.B) {
	benchmarkBulkConnectorEdit_givenNCreateItems(b, 1000)
}

func benchmarkBulkConnectorEdit_givenNCreateItems(b *testing.B, itemCount int) {
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
	loader := fixtures.NewLoader(dbClient, []string{"./testdata/fixtures/bulk_connector_edit.yml"},
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
	store := pbehavior.NewStore(dbClient, nil, nil, nil, nil, nil, authorProvider, nil, nil, nil)
	api := pbehavior.NewApi(store, nil, ch, nil, zerolog.Nop())
	reqBodies := make([]io.ReadCloser, b.N)
	now := time.Now().Unix()
	tomorrow := time.Now().AddDate(0, 0, 1).Unix()
	itemTpl := map[string]any{
		"action":  "create",
		"origin":  "centreon/centreon_test",
		"comment": "Downtime set by test",
		"name":    "centreon/centreon_test downtime " + strconv.FormatInt(now, 10) + "-" + strconv.FormatInt(tomorrow, 10) + " Downtime set by test",
		"color":   "#73D8FF",
		"reason":  "default_reason",
		"type":    "default_maintenance",
		"tstart":  now,
		"tstop":   tomorrow,
	}
	for i := 0; i < b.N; i++ {
		reqItems := make([]map[string]any, itemCount)
		for j := 0; j < itemCount; j++ {
			reqItems[j] = make(map[string]any, len(itemTpl))
			maps.Copy(reqItems[j], itemTpl)
			reqItems[j]["entities"] = []string{"cps-resource-" + strconv.Itoa(i) + "-" + strconv.Itoa(j) + "/cps-component"}
		}

		buf := bytes.NewBuffer(nil)
		encoder := json.NewEncoder(buf)
		err = encoder.Encode(reqItems)
		if err != nil {
			b.Errorf("unexpected error %v", err)
		}

		reqBodies[i] = io.NopCloser(buf)
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPut, "/api/v4/bulk/connector-pbehaviors", nil)
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = r
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Body = reqBodies[i]
		api.BulkConnectorEdit(ginCtx)
		if w.Code != http.StatusMultiStatus {
			b.Fatalf("unexpected status code %v", w.Code)
		}
	}
}
