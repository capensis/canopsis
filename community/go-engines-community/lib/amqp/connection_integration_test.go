package amqp

import (
	"io"
	"net/http"
	"os"
	"slices"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	rabbithole "github.com/michaelklishin/rabbit-hole/v3"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func TestIntegration_Connection_Channel_GivenLiveBroker_ShouldReturnChannel(t *testing.T) {
	c, err := New(3, 200*time.Millisecond, zerolog.Nop())
	if err != nil {
		t.Fatalf("New: %+v", err)
	}

	ch, err := c.Channel(t.Context())
	if err != nil {
		t.Fatalf("Channel: %+v", err)
	}

	if err = ch.Close(); err != nil {
		t.Fatalf("ch.Close: %+v", err)
	}

	if err = c.Close(); err != nil {
		t.Fatalf("Close: %+v", err)
	}
}

func TestIntegration_Connection_GivenBrokerDrop_ShouldReconnect(t *testing.T) {
	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	connName := t.Name() + utils.NewID()
	config.Properties.SetClientConnectionName(connName)

	c, err := NewConfig(20, 200*time.Millisecond, config, zerolog.Nop())
	if err != nil {
		t.Fatalf("New: %+v", err)
	}

	conn := c.(*connection)

	initial := conn.waitReconnect(t.Context().Done(), nil)
	if initial == nil {
		t.Fatal("initial conn nil")
	}

	apiClient := newTestAPIClient(t)
	realConnName := findRealConnName(t, apiClient, connName)
	killConnection(t, apiClient, realConnName)

	got := conn.waitReconnect(t.Context().Done(), initial)
	if got == nil {
		t.Fatal("waitReconnect returned nil after drop")
	}

	if got == initial {
		t.Fatal("waitReconnect returned the dropped conn")
	}

	if err = c.Close(); err != nil {
		t.Fatalf("Close: %+v", err)
	}
}

func findRealConnName(t *testing.T, c *rabbithole.Client, connName string) string {
	t.Helper()

	const retries = 50
	const timeout = 100 * time.Millisecond
	for i := 0; i < retries; i++ {
		conns, err := c.ListConnections()
		if err != nil {
			t.Fatalf("cannot list rabbitmq connections: %+v", err)
		}

		idx := slices.IndexFunc(conns, func(info rabbithole.ConnectionInfo) bool {
			return info.ClientProperties["connection_name"] == connName
		})
		if idx >= 0 {
			return conns[idx].Name
		}

		if i < retries-1 {
			select {
			case <-t.Context().Done():
			case <-time.After(timeout):
			}
		}
	}

	t.Fatalf("cannot find rabbitmq connection: %s", connName)

	return ""
}

func killConnection(t *testing.T, c *rabbithole.Client, realConnName string) {
	t.Helper()

	resp, err := c.CloseConnection(realConnName)
	if err != nil {
		t.Fatalf("cannot close connection: %+v", err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("cannot read response body: %+v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("cannot close connection: %s, status code: %d, response body: %s", realConnName, resp.StatusCode, string(b))
	}
}

func newTestAPIClient(t *testing.T) *rabbithole.Client {
	t.Helper()
	url := os.Getenv(EnvHTTPURL)
	if url == "" {
		t.Fatalf("environment variable %s not set", EnvHTTPURL)
	}

	user := os.Getenv(EnvHTTPUser)
	pass := os.Getenv(EnvHTTPPassword)

	c, err := rabbithole.NewClient(url, user, pass)
	if err != nil {
		t.Fatalf("cannot create rabbitmq api client: %+v", err)
	}

	return c
}
