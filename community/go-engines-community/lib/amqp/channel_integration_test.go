package amqp

import (
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func TestIntegration_Channel_GivenLiveBroker_ShouldPublishAndConsumeMessage(t *testing.T) {
	conn, err := New(3, 200*time.Millisecond, zerolog.Nop())
	if err != nil {
		t.Fatalf("New: %+v", err)
	}

	t.Cleanup(func() {
		_ = conn.Close()
	})

	publishCh, err := conn.Channel(t.Context())
	if err != nil {
		t.Fatalf("Channel: %+v", err)
	}

	t.Cleanup(func() {
		_ = publishCh.Close()
	})

	consumeCh, err := conn.Channel(t.Context())
	if err != nil {
		t.Fatalf("Channel: %+v", err)
	}

	t.Cleanup(func() {
		_ = consumeCh.Close()
	})

	q, err := publishCh.QueueDeclare(t.Context(), "", false, false, true, false, nil)
	if err != nil {
		t.Fatalf("QueueDeclare: %+v", err)
	}

	msg := []byte("test")
	err = publishCh.PublishWithContext(t.Context(), "", q.Name, false, false, amqp.Publishing{
		Body:         msg,
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		t.Fatalf("PublishWithContext: %+v", err)
	}

	d, err := consumeCh.ConsumeWithContext(t.Context(), q.Name, "", false, false, false, false, nil)
	if err != nil {
		t.Fatalf("ConsumeWithContext: %+v", err)
	}

	select {
	case <-time.After(5 * time.Second):
		t.Fatalf("timeout exceeded")
	case got, ok := <-d:
		if !ok {
			t.Fatalf("channel closed unexpectedly")
		}

		if string(got.Body) != string(msg) {
			t.Fatalf("expected %s, got %s", string(msg), string(got.Body))
		}

		err = got.Ack(false)
		if err != nil {
			t.Fatalf("Ack: %+v", err)
		}
	}

	if err = publishCh.Close(); err != nil {
		t.Fatalf("publishCh.Close: %+v", err)
	}

	if err = consumeCh.Close(); err != nil {
		t.Fatalf("consumeCh.Close: %+v", err)
	}

	if err = conn.Close(); err != nil {
		t.Fatalf("Close: %+v", err)
	}
}

func TestIntegration_Channel_GivenBrokerDrop_ShouldReconnect(t *testing.T) {
	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	connName := t.Name() + utils.NewID()
	config.Properties.SetClientConnectionName(connName)

	conn, err := NewConfig(20, 200*time.Millisecond, config, zerolog.Nop())
	if err != nil {
		t.Fatalf("New: %+v", err)
	}

	t.Cleanup(func() {
		_ = conn.Close()
	})

	ch, err := conn.Channel(t.Context())
	if err != nil {
		t.Fatalf("Channel: %+v", err)
	}

	t.Cleanup(func() {
		_ = ch.Close()
	})

	connection := conn.(*connection)
	channel := ch.(*channel)

	initialConn := connection.waitReconnect(t.Context().Done(), nil)
	if initialConn == nil {
		t.Fatal("initial connection nil")
	}

	initial := channel.waitReconnect(t.Context(), nil)
	if initial == nil {
		t.Fatal("initial channel nil")
	}

	apiClient := newTestAPIClient(t)
	realConnName := findRealConnName(t, apiClient, connName)
	killConnection(t, apiClient, realConnName)

	gotConn := connection.waitReconnect(t.Context().Done(), initialConn)
	if gotConn == nil {
		t.Fatal("waitReconnect returned nil after drop")
	}

	got := channel.waitReconnect(t.Context(), initial)
	if got == nil {
		t.Fatal("waitReconnect returned nil after drop")
	}

	if gotConn == initialConn {
		t.Fatal("waitReconnect returned the dropped connection")
	}

	if got == initial {
		t.Fatal("waitReconnect returned the dropped channel")
	}

	if err = ch.Close(); err != nil {
		t.Fatalf("ch.Close: %+v", err)
	}

	if err = conn.Close(); err != nil {
		t.Fatalf("Close: %+v", err)
	}
}

func TestIntegration_Channel_GivenChannelDrop_ShouldReconnect(t *testing.T) {
	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	connName := t.Name() + utils.NewID()
	config.Properties.SetClientConnectionName(connName)

	conn, err := NewConfig(20, 200*time.Millisecond, config, zerolog.Nop())
	if err != nil {
		t.Fatalf("New: %+v", err)
	}

	t.Cleanup(func() {
		_ = conn.Close()
	})

	ch, err := conn.Channel(t.Context())
	if err != nil {
		t.Fatalf("Channel: %+v", err)
	}

	t.Cleanup(func() {
		_ = ch.Close()
	})

	connection := conn.(*connection)
	channel := ch.(*channel)

	initialConn := connection.waitReconnect(t.Context().Done(), nil)
	if initialConn == nil {
		t.Fatal("initial connection nil")
	}

	initial := channel.waitReconnect(t.Context(), nil)
	if initial == nil {
		t.Fatal("initial channel nil")
	}

	if err = initial.Close(); err != nil {
		t.Fatalf("initial.Close: %+v", err)
	}

	got := channel.waitReconnect(t.Context(), initial)
	if got == nil {
		t.Fatal("waitReconnect returned nil after drop")
	}

	gotConn := connection.waitReconnect(t.Context().Done(), nil)
	if gotConn == nil {
		t.Fatal("waitReconnect returned nil after drop")
	}

	if got == initial {
		t.Fatal("waitReconnect returned the dropped channel")
	}

	if gotConn != initialConn {
		t.Fatal("waitReconnect returned a new connection")
	}

	if err = ch.Close(); err != nil {
		t.Fatalf("ch.Close: %+v", err)
	}

	if err = conn.Close(); err != nil {
		t.Fatalf("Close: %+v", err)
	}
}

func TestIntegration_Channel_GivenBrokerDrop_ShouldResumeMessageFlow(t *testing.T) {
	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	connName := t.Name() + utils.NewID()
	config.Properties.SetClientConnectionName(connName)

	conn, err := NewConfig(20, 200*time.Millisecond, config, zerolog.Nop())
	if err != nil {
		t.Fatalf("New: %+v", err)
	}

	t.Cleanup(func() {
		_ = conn.Close()
	})

	publishCh, err := conn.Channel(t.Context())
	if err != nil {
		t.Fatalf("Channel: %+v", err)
	}

	t.Cleanup(func() {
		_ = publishCh.Close()
	})

	consumeCh, err := conn.Channel(t.Context())
	if err != nil {
		t.Fatalf("Channel: %+v", err)
	}

	t.Cleanup(func() {
		_ = consumeCh.Close()
	})

	q, err := publishCh.QueueDeclare(t.Context(), "", true, false, false, false, nil)
	if err != nil {
		t.Fatalf("QueueDeclare: %+v", err)
	}

	t.Cleanup(func() {
		deleteTestQueue(t, q.Name)
	})

	msgs := [][]byte{
		[]byte("test1"),
		[]byte("test2"),
	}

	consumeMsgCh := make(chan []byte, 1)
	consumeErrCh := make(chan error, 1)
	var consumeWithCtxCalled atomic.Int32

	go func() {
		consumedMsgCount := 0
		for {
			consumeWithCtxCalled.Add(1)
			d, notifyClose, err := consumeCh.ConsumeWithCloseNotify(t.Context(), q.Name, "", true, false, false, false, nil)
			if err != nil {
				select {
				case <-t.Context().Done():
				case consumeErrCh <- err:
				}

				return
			}

		loop:
			for {
				select {
				case <-t.Context().Done():
					return
				case <-notifyClose:
					break loop
				case got, ok := <-d:
					if !ok {
						select {
						case <-t.Context().Done():
						case consumeErrCh <- errors.New("channel closed unexpectedly"):
						}

						return
					}

					consumedMsgCount++

					select {
					case <-t.Context().Done():
					case consumeMsgCh <- got.Body:
					}

					if consumedMsgCount == len(msgs) {
						return
					}
				}
			}
		}
	}()

	err = publishCh.PublishWithContext(t.Context(), "", q.Name, false, false, amqp.Publishing{
		Body:         msgs[0],
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		t.Fatalf("PublishWithContext: %+v", err)
	}

	select {
	case <-time.After(5 * time.Second):
		t.Fatal("timeout exceeded")
	case gotErr := <-consumeErrCh:
		t.Fatalf("consume error: %+v", gotErr)
	case got := <-consumeMsgCh:
		if string(got) != string(msgs[0]) {
			t.Fatalf("expected: %s, got: %s", string(msgs[0]), string(got))
		}
	}

	apiClient := newTestAPIClient(t)
	realConnName := findRealConnName(t, apiClient, connName)
	killConnection(t, apiClient, realConnName)

	err = publishCh.PublishWithContext(t.Context(), "", q.Name, false, false, amqp.Publishing{
		Body:         msgs[1],
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		t.Fatalf("PublishWithContext: %+v", err)
	}

	select {
	case <-time.After(5 * time.Second):
		t.Fatal("timeout exceeded")
	case gotErr := <-consumeErrCh:
		t.Fatalf("consume error: %+v", gotErr)
	case got := <-consumeMsgCh:
		if string(got) != string(msgs[1]) {
			t.Fatalf("expected: %s, got: %s", string(msgs[1]), string(got))
		}
	}

	if got := consumeWithCtxCalled.Load(); got != 2 {
		t.Fatalf("expected %d channel reconnects, got %d", 2, got)
	}

	if err = publishCh.Close(); err != nil {
		t.Fatalf("publishCh.Close: %+v", err)
	}

	if err = consumeCh.Close(); err != nil {
		t.Fatalf("consumeCh.Close: %+v", err)
	}

	if err = conn.Close(); err != nil {
		t.Fatalf("Close: %+v", err)
	}
}

func deleteTestQueue(t *testing.T, queueName string) {
	conn, err := amqp.Dial(os.Getenv(EnvURL))
	if err != nil {
		t.Logf("Dial: %+v", err)

		return
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Logf("Channel: %+v", err)

		return
	}

	defer ch.Close()

	_, err = ch.QueueDelete(queueName, false, false, false)
	if err != nil {
		t.Logf("QueueDelete: %+v", err)
	}
}
