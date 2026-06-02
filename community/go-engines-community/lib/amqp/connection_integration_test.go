package amqp

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
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

	apiClient := newTestAPIClient(t)
	testUser := createTemporaryAMQPUser(t, apiClient)
	setTemporaryAMQPUserURL(t, testUser)

	c, err := NewConfig(20, 200*time.Millisecond, config, zerolog.Nop())
	if err != nil {
		t.Fatalf("New: %+v", err)
	}

	conn := c.(*connection)

	initial := conn.waitReconnect(t.Context().Done(), nil)
	if initial == nil {
		t.Fatal("initial conn nil")
	}

	killAllConnectionsOfUser(t, apiClient, testUser.Username)

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

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("cannot read response body: %+v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("cannot close all user connections: status code: %d, response body: %s", resp.StatusCode, string(b))
	}
}

type temporaryAMQPUser struct {
	Username string
	Password string
}

func createTemporaryAMQPUser(t *testing.T, c *rabbithole.Client) temporaryAMQPUser {
	t.Helper()
	vhost := getAMQPVhost(t)
	username := "amqp-test-" + strings.ToLower(strings.ReplaceAll(utils.NewID(), "_", "-"))
	password := "pwd-" + utils.NewID()

	resp, err := c.PutUser(username, rabbithole.UserSettings{
		Password: password,
	})
	if err != nil {
		t.Fatalf("cannot create temporary rabbitmq user: %+v", err)
	}
	closeManagementResponse(t, resp)
	assertManagementStatus(t, resp, http.StatusCreated, http.StatusNoContent)

	resp, err = c.UpdatePermissionsIn(vhost, username, rabbithole.Permissions{
		Configure: ".*",
		Write:     ".*",
		Read:      ".*",
	})
	if err != nil {
		_, _ = c.DeleteUser(username)
		t.Fatalf("cannot grant temporary rabbitmq user permissions: %+v", err)
	}
	closeManagementResponse(t, resp)
	assertManagementStatus(t, resp, http.StatusCreated, http.StatusNoContent)

	u := temporaryAMQPUser{Username: username, Password: password}
	t.Cleanup(func() {
		deleteTemporaryAMQPUser(t, c, u.Username)
	})
	return u
}

func deleteTemporaryAMQPUser(t *testing.T, c *rabbithole.Client, username string) {
	t.Helper()

	resp, err := c.DeleteUser(username)
	if err != nil {
		t.Fatalf("cannot delete temporary rabbitmq user: %+v", err)
	}

	closeManagementResponse(t, resp)
	assertManagementStatus(t, resp, http.StatusNoContent)
}

func setTemporaryAMQPUserURL(t *testing.T, user temporaryAMQPUser) {
	t.Helper()
	u, err := url.Parse(os.Getenv(EnvURL))
	if err != nil {
		t.Fatalf("cannot parse amqp url: %+v", err)
	}

	u.User = url.UserPassword(user.Username, user.Password)
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

func closeManagementResponse(t *testing.T, resp *http.Response) {
	t.Helper()
	if resp == nil || resp.Body == nil {
		return
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
}

func assertManagementStatus(t *testing.T, resp *http.Response, allowedStatuses ...int) {
	t.Helper()
	if resp == nil {
		t.Fatal("management response is nil")
	}

	if slices.Contains(allowedStatuses, resp.StatusCode) {
		return
	}

	t.Fatalf("unexpected management response status: %d", resp.StatusCode)
}

func newTestAPIClient(t *testing.T) *rabbithole.Client {
	t.Helper()
	url := os.Getenv(EnvHTTPURL)
	user := os.Getenv(EnvHTTPUser)
	pass := os.Getenv(EnvHTTPPassword)

	c, err := rabbithole.NewClient(url, user, pass)
	if err != nil {
		t.Fatalf("cannot create rabbitmq api client: %+v", err)
	}

	return c
}
