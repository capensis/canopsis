package mongo

import "time"

// Globals
const (
	EnvURL = "CPS_MONGO_URL"

	// Timeout is aimed to be used for first connection with mongo.NewSession
	Timeout = time.Second * 3

	// TimeoutWork is used by mongo.NewSession to set session.SetSocketTimeout()
	// after connection was successful.
	TimeoutWork = time.Second * 30
)
