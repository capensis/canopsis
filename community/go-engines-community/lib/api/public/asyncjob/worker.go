package asyncjob

import (
	"context"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

// Worker consumes job signals and dispatches each of them to the Handler registered for its type.
type Worker struct {
	chPool   libamqp.ChannelPool
	decoder  encoding.Decoder
	handlers map[string]Handler
	workers  int
	logger   zerolog.Logger
}

// NewWorker returns a Worker which runs the given handlers, keyed by job type.
// workers is the number of jobs processed in parallel: all types share the same queue,
// so a single worker running a long job delays the other ones.
func NewWorker(
	chPool libamqp.ChannelPool,
	decoder encoding.Decoder,
	handlers map[string]Handler,
	workers int,
	logger zerolog.Logger,
) *Worker {
	if workers < 1 {
		workers = 1
	}

	return &Worker{
		chPool:   chPool,
		decoder:  decoder,
		handlers: handlers,
		workers:  workers,
		logger:   logger,
	}
}

// Run opens its own AMQP channel, declares the job queue and consumes it until ctx is done.
// It returns nil on graceful shutdown.
func (w *Worker) Run(ctx context.Context) error {
	queue := canopsis.PublicApiJobQueueName

	ch, err := w.chPool.Get(ctx)
	if err != nil {
		return err
	}

	defer w.chPool.Put(ch)

	// check if queue exists because Consume method doesn't return appropriate error
	_, err = ch.QueueDeclarePassive(ctx, queue, false, false, false, false, nil)
	if err != nil {
		return err
	}

	types := make([]string, 0, len(w.handlers))
	for jobType := range w.handlers {
		types = append(types, jobType)
	}

	w.logger.Info().
		Str("queue", queue).
		Strs("types", types).
		Int("workers", w.workers).
		Msg("async job worker started")

	return libamqp.ConsumeWithReconnect(
		ctx,
		ch,
		libamqp.ConsumeOptions{
			Queue:         queue,
			Consumer:      canopsis.PublicApiJobConsumerName,
			PrefetchCount: w.workers,
		},
		w.workers,
		w.process,
		w.logger,
	)
}

// process never returns an error: a broken message or a failed job must not tear down the consume loop.
// Malformed messages are dropped, and a failed job is reported by its Handler on the job itself,
// so the message is acknowledged in both cases.
func (w *Worker) process(ctx context.Context, d amqp.Delivery) (libamqp.AckAction, error) {
	var msg Message
	err := w.decoder.Decode(d.Body, &msg)
	if err != nil {
		w.logger.Err(err).Str("body", string(d.Body)).Msg("cannot decode job message, message is dropped")

		return libamqp.Ack, nil
	}

	if msg.JobID == "" {
		w.logger.Error().Str("body", string(d.Body)).Msg("job message has no job id, message is dropped")

		return libamqp.Ack, nil
	}

	handler, ok := w.handlers[msg.Type]
	if !ok {
		w.logger.Error().
			Str("type", msg.Type).
			Str("job_id", msg.JobID).
			Msg("no handler for job type, message is dropped")

		return libamqp.Ack, nil
	}

	w.logger.Debug().Str("type", msg.Type).Str("job_id", msg.JobID).Msg("job started")

	start := time.Now()
	err = handler.Handle(ctx, msg.JobID)
	if err != nil {
		if ctx.Err() != nil {
			// Shutdown in the middle of the job: requeue so it's picked up on the next start.
			return libamqp.Nack, nil
		}

		w.logger.Err(err).Str("type", msg.Type).Str("job_id", msg.JobID).Msg("cannot process job")

		return libamqp.Ack, nil
	}

	w.logger.Info().
		Str("type", msg.Type).
		Str("job_id", msg.JobID).
		Str("duration", time.Since(start).String()).
		Msg("job done")

	return libamqp.Ack, nil
}
