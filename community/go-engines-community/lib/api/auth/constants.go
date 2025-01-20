package auth

import "time"

const (
	// UserKey is the context name for user credential.
	UserKey  = "user_id"
	Username = "username"
	// ApiKey is the context name for user api key.
	ApiKey = "api_key"
	// RolesKey is the context name for user's roles
	RolesKey = "roles"
)

const DefaultMetaValidDuration = time.Hour
