package amqp

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
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

	apiClient := newRabbitMQAPIClient(t)
	username := setupTempAMQPUser(t, apiClient)

	c, err := NewConfig(20, 200*time.Millisecond, config, zerolog.Nop())
	if err != nil {
		t.Fatalf("New: %+v", err)
	}

	conn := c.(*connection)

	initial := conn.waitReconnect(t.Context().Done(), nil)
	if initial == nil {
		t.Fatal("initial conn nil")
	}

	killAllConnectionsOfUser(t, apiClient, username)

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

func killAllConnectionsOfUser(t *testing.T, c *rabbithole.Client, username string) {
	t.Helper()

	resp, err := c.CloseAllConnectionsOfUser(username)
	if err != nil {
		t.Fatalf("cannot close all user connections: %+v", err)
	}

	checkAMQPResponseCode(t, resp, http.StatusNoContent)
}

func setupTempAMQPUser(t *testing.T, c *rabbithole.Client) string {
	username, pwd := createTempAMQPUser(t, c)
	t.Cleanup(func() {
		deleteTempAMQPUser(t, c, username)
	})

	overrideEnvURLUser(t, username, pwd)

	return username
}

func createTempAMQPUser(t *testing.T, c *rabbithole.Client) (username string, pwd string) {
	t.Helper()
	vhost := getAMQPVhost(t)
	username = t.Name() + strconv.FormatInt(time.Now().Unix(), 10)
	pwd = "test"

	resp, err := c.PutUser(username, rabbithole.UserSettings{
		Password: pwd,
	})
	if err != nil {
		t.Fatalf("cannot create temporary rabbitmq user: %+v", err)
	}

	checkAMQPResponseCode(t, resp, http.StatusCreated)

	resp, err = c.UpdatePermissionsIn(vhost, username, rabbithole.Permissions{
		Configure: ".*",
		Write:     ".*",
		Read:      ".*",
	})
	if err != nil {
		_, _ = c.DeleteUser(username)
		t.Fatalf("cannot grant temporary rabbitmq user permissions: %+v", err)
	}

	checkAMQPResponseCode(t, resp, http.StatusCreated)

	return username, pwd
}

func deleteTempAMQPUser(t *testing.T, c *rabbithole.Client, username string) {
	t.Helper()

	resp, err := c.DeleteUser(username)
	if err != nil {
		t.Fatalf("cannot delete temporary rabbitmq user: %+v", err)
	}

	checkAMQPResponseCode(t, resp, http.StatusNoContent)
}

func overrideEnvURLUser(t *testing.T, username, pwd string) {
	t.Helper()
	u, err := url.Parse(os.Getenv(EnvURL))
	if err != nil {
		t.Fatalf("cannot parse amqp url: %+v", err)
	}

	u.User = url.UserPassword(username, pwd)
	t.Setenv(EnvURL, u.String())
}

func getAMQPVhost(t *testing.T) string {
	t.Helper()
	u, err := url.Parse(os.Getenv(EnvURL))
	if err != nil {
		t.Fatalf("cannot parse amqp url: %+v", err)
	}

	vhost := strings.TrimPrefix(u.Path, "/")
	if vhost == "" {
		return "/"
	}

	decoded, err := url.PathUnescape(vhost)
	if err != nil {
		t.Fatalf("cannot decode amqp vhost: %+v", err)
	}

	return decoded
}

func checkAMQPResponseCode(t *testing.T, resp *http.Response, expectedStatusCode int) {
	t.Helper()

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("cannot read response body: %+v", err)
	}

	if resp.StatusCode != expectedStatusCode {
		t.Fatalf("unexpected status code: %d, expected %d, request url: %s, response body: %s",
			resp.StatusCode, expectedStatusCode, resp.Request.URL.String(), string(b))
	}
}

func newRabbitMQAPIClient(t *testing.T) *rabbithole.Client {
	t.Helper()
	u := os.Getenv(EnvHTTPURL)
	if u == "" {
		t.Fatalf("environment variable %s not set", EnvHTTPURL)
	}

	user := os.Getenv(EnvHTTPUser)
	pass := os.Getenv(EnvHTTPPassword)

	c, err := rabbithole.NewClient(u, user, pass)
	if err != nil {
		t.Fatalf("cannot create rabbitmq api client: %+v", err)
	}

	return c
}
