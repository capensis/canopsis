package authctx

type contextKey int

const (
	// UserKey is the context name for user id.
	UserKey contextKey = iota
	Username
	// APIKey is the context name for user api key.
	APIKey
	// Roles is the context name for user's roles
	Roles
	// Locale is the context name for user's preferred locale.
	Locale
)
