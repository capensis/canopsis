package axe

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
)

// ProcessResumeJobForTest lets the existing processor tests deliver durable work
// synchronously while preserving their per-event dependency ordering hooks.
func ProcessResumeJobForTest(ctx context.Context, client mongo.DbClient, publisher amqp.Publisher, id string) error {
	w := newResumeWorker(client, client, nil, publisher, json.NewEncoder(),
		event.NewGenerator(canopsis.AxeConnector, canopsis.AxeConnector), zerolog.Nop())

	return w.processJob(ctx, id)
}

// ProcessPauseJobForTest delivers the pause worker's durable work through the
// existing processor tests without starting a background engine routine.
func ProcessPauseJobForTest(ctx context.Context, client mongo.DbClient, publisher amqp.Publisher, id string) error {
	w := newPauseWorker(client, client, nil, publisher, json.NewEncoder(),
		event.NewGenerator(canopsis.AxeConnector, canopsis.AxeConnector), zerolog.Nop())
	return w.processJob(ctx, id)
}
