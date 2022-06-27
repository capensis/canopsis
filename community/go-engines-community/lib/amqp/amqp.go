package amqp

import (
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

// Environment variables linked to session parameters
const (
	EnvURL = "CPS_AMQP_URL"
)

// NewSession creates a new connection to an AMQP bus,
// using env var EnvCpsAmqpURL as configuration.
// Use NewConnection for reconnection feature.
func NewSession() (*amqp.Connection, error) {
	url := os.Getenv(EnvURL)
	if url == "" {
		return nil, fmt.Errorf("environment variable %s empty", EnvURL)
	}
	return amqp.Dial(url)
}

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
