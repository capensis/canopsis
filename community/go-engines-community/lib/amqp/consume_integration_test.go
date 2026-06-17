package amqp_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func TestIntegration_ConsumeWithReconnect_ShouldNotDeadlockOnShutdown(t *testing.T) {
	conn, err := amqp.New(3, 200*time.Millisecond, zerolog.Nop())
	if err != nil {
		t.Fatalf("New: %+v", err)
	}

	t.Cleanup(func() { _ = conn.Close() })

	pubCh, err := conn.Channel(t.Context())
	if err != nil {
		t.Fatalf("Channel: %+v", err)
	}

	t.Cleanup(func() { _ = pubCh.Close() })

	consumeCh, err := conn.Channel(t.Context())
	if err != nil {
		t.Fatalf("Channel: %+v", err)
	}

	q, err := pubCh.QueueDeclare(t.Context(), "", false, false, true, false, nil)
	if err != nil {
		t.Fatalf("QueueDeclare: %+v", err)
	}

	// Publish enough messages that the broker keeps pushing while we shut down.
	const total = 500
	for i := 0; i < total; i++ {
		err = pubCh.PublishWithContext(t.Context(), "", q.Name, false, false,
			amqp091.Publishing{Body: []byte("x")})
		if err != nil {
			t.Fatalf("Publish: %+v", err)
		}
	}

	// Slow processor so cancel() catches us with deliveries still queued in amqp091's consumers.buffer for this consumer.
	process := func(ctx context.Context, _ amqp091.Delivery) (amqp.AckAction, error) {
		select {
		case <-ctx.Done():
			return amqp.Nack, nil
		case <-time.After(50 * time.Millisecond):
			return amqp.Ack, nil
		}
	}

	ctx, cancel := context.WithCancel(t.Context())
	consumeDone := make(chan error, 1)
	go func() {
		consumeDone <- amqp.ConsumeWithReconnect(ctx, consumeCh, amqp.ConsumeOptions{Queue: q.Name}, 4, process, zerolog.Nop())
	}()

	// Let the broker fill the pipeline and the workers chew on a few.
	time.Sleep(300 * time.Millisecond)

	// Trigger shutdown.
	cancel()

	// ConsumeWithReconnect must return without blocking.
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("ConsumeWithReconnect did not return after cancel (drain missing?)")
	case err := <-consumeDone:
		if err != nil && !errors.Is(err, context.Canceled) {
			t.Fatalf("ConsumeWithReconnect: %+v", err)
		}
	}

	// Closing the channel must not hang either. Without the drain, this is where the deadlock manifests in production.
	closeDone := make(chan error, 1)
	go func() { closeDone <- consumeCh.Close() }()
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("consumeCh.Close() deadlocked")
	case err := <-closeDone:
		if err != nil {
			t.Fatalf("consumeCh.Close: %+v", err)
		}
	}
}
