package amqp

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Environment variables linked to session parameters
const (
	EnvURL                = "CPS_AMQP_URL"
	EnvInsecureSkipVerify = "CPS_AMQP_INSECURE_SKIP_VERIFY"
	EnvCaCertFile         = "CPS_AMQP_CA_CERT_FILE"
)

// NewSession creates a new connection to an AMQP bus,
// using env var EnvCpsAmqpURL as configuration.
// New connection tries to reconnect to AMQP if connection is lost.
func NewConnection(logger zerolog.Logger, reconnectCount int, minReconnectTimeout time.Duration) (Connection, error) {
	url := os.Getenv(EnvURL)
	if url == "" {
		return nil, fmt.Errorf("environment variable %s empty", EnvURL)
	}
	return Dial(url, logger, reconnectCount, minReconnectTimeout)
}

func BuildRoutingKey(prefix string, initiator string) string {
	return prefix + "_" + initiator
}
