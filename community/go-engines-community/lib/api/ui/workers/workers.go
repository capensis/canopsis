package workers

import (
	"context"
	"errors"
	"fmt"
	"sync"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type JobPublisher interface {
	Publish(ctx context.Context, id string) error
}

type JobExecutor func(ctx context.Context, id string) error

type Job struct {
	ID   string `json:"_id"`
	Type string `json:"type"`
}

func NewRunner(amqpConn libamqp.Connection, logger zerolog.Logger) *Runner {
	return &Runner{
		amqpConn:     amqpConn,
		queue:        canopsis.ApiWorkersQueueName,
		decoder:      json.NewDecoder(),
		workers:      10,
		jobExecutors: make(map[string]JobExecutor),
		logger:       logger,
	}
}

type Runner struct {
	amqpConn       libamqp.Connection
	queue          string
	decoder        encoding.Decoder
	workers        int
	jobExecutorsMx sync.RWMutex
	jobExecutors   map[string]JobExecutor
	logger         zerolog.Logger
}

func (r *Runner) Run(ctx context.Context) error {
	amqpConsumer, err := r.amqpConn.Channel()
	if err != nil {
		return fmt.Errorf("cannot create rmq channel: %w", err)
	}

	defer amqpConsumer.Close()
	// check if queue exists because Consume method doesn't return appropriate error
	_, err = amqpConsumer.QueueInspect(r.queue)
	if err != nil {
		return err
	}

	ch, err := amqpConsumer.Consume(r.queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to consume jobs: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < r.workers; i++ {
		g.Go(func() error {
			var err error
			for {
				select {
				case <-ctx.Done():
					return nil
				case msg, ok := <-ch:
					if !ok {
						return errors.New("channel closed")
					}

					r.logger.Debug().Str("msg", string(msg.Body)).Msg("worker job message")
					job := Job{}
					err = r.decoder.Decode(msg.Body, &job)
					if err != nil {
						r.logger.Err(err).Msg("failed to decode job")
						continue
					}

					r.jobExecutorsMx.RLock()
					e, ok := r.jobExecutors[job.Type]
					r.jobExecutorsMx.RUnlock()
					if !ok {
						r.logger.Warn().Str("type", job.Type).Str("id", job.ID).Msg("unknown job type")
						continue
					}

					err = e(ctx, job.ID)
					if err != nil {
						r.logger.Err(err).Str("type", job.Type).Str("id", job.ID).Msg("failed to execute job")
						if mongo.IsConnectionError(err) {
							err = amqpConsumer.Nack(msg.DeliveryTag, false, true)
							if err != nil {
								r.logger.Err(err).Msg("failed to negatively acknowledge message")
							}

							continue
						}
					}

					err = amqpConsumer.Ack(msg.DeliveryTag, false)
					if err != nil {
						r.logger.Err(err).Msg("failed to acknowledge message")
					}
				}
			}
		})
	}

	return g.Wait()
}

func (r *Runner) AddJobExecutor(t string, e JobExecutor) {
	r.jobExecutorsMx.Lock()
	defer r.jobExecutorsMx.Unlock()
	if _, ok := r.jobExecutors[t]; ok {
		panic("duplicate job executor")
	}

	r.jobExecutors[t] = e
}

func NewJobPublisher(t string, amqpPublisher libamqp.Channel) JobPublisher {
	return &jobPublisher{
		jobType:       t,
		amqpPublisher: amqpPublisher,
		exchange:      canopsis.DefaultExchangeName,
		queue:         canopsis.ApiWorkersQueueName,
		encoder:       json.NewEncoder(),
		contentType:   canopsis.JsonContentType,
	}
}

type jobPublisher struct {
	jobType         string
	amqpPublisher   libamqp.Channel
	exchange, queue string
	encoder         encoding.Encoder
	contentType     string
}

func (r *jobPublisher) Publish(ctx context.Context, id string) error {
	job := Job{ID: id, Type: r.jobType}
	b, err := r.encoder.Encode(job)
	if err != nil {
		return fmt.Errorf("failed to encode job: %w", err)
	}

	err = r.amqpPublisher.PublishWithContext(ctx, r.exchange, r.queue, false, false, amqp.Publishing{
		ContentType:  r.contentType,
		DeliveryMode: amqp.Persistent,
		Body:         b,
	})
	if err != nil {
		return fmt.Errorf("failed to send %q job %q to amqp: %w", job.Type, job.ID, err)
	}

	return nil
}
