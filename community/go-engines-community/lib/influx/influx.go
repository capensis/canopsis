package influx

import (
	"net/url"
	"os"

	influxmod "github.com/influxdata/influxdb/client/v2"
)

// Consts for InfluxDB sessions
const (
	EnvURL = "CPS_INFLUX_URL"
)

// NewSession creates a new connection to the InfluxDB database.
// It uses the EnvCpsInfluxUrl as configuration.
func NewSession() (influxmod.Client, error) {
	influxRawURL := os.Getenv(EnvURL)
	influxURL, err := url.ParseRequestURI(influxRawURL)
	if err != nil {
		return nil, err
	}

	influxUsername := ""
	influxPassword := ""
	influxPasswordSet := false
	if influxURL.User != nil {
		influxUsername = influxURL.User.Username()
		influxPassword, influxPasswordSet = influxURL.User.Password()
		if !influxPasswordSet {
			influxPassword = ""
		}
	}

	session, err := influxmod.NewHTTPClient(influxmod.HTTPConfig{
		Addr:     "http://" + influxURL.Host,
		Username: influxUsername,
		Password: influxPassword,
	})
	if err != nil {
		return nil, err
	}

	return session, nil
}
